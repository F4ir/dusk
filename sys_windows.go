//go:build windows

package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
}

func isAdmin() bool { return windows.GetCurrentProcessToken().IsElevated() }

func relaunchAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(exe)
	args, _ := windows.UTF16PtrFromString("--no-intro")
	return windows.ShellExecute(0, verb, file, args, nil, windows.SW_NORMAL)
}

func restartExplorer() error {
	c := exec.Command("taskkill", "/f", "/im", "explorer.exe")
	hideWindow(c)
	c.Run()
	return exec.Command("explorer.exe").Start()
}

func wingetAvailable() bool {
	if _, err := exec.LookPath("winget"); err == nil {
		return true
	}
	// net geïnstalleerd: staat er al, maar PATH van dit proces is nog oud
	_, err := os.Stat(os.Getenv("LOCALAPPDATA") + `\Microsoft\WindowsApps\winget.exe`)
	return err == nil
}

func wingetList() (string, error) {
	c := exec.Command("winget", "list", "--accept-source-agreements", "--disable-interactivity")
	hideWindow(c)
	out, err := c.Output()
	if len(out) > 0 {
		err = nil
	}
	return string(out), err
}

func wingetShow(id string) string {
	c := exec.Command("winget", "show", "--id", id, "--exact", "--accept-source-agreements", "--disable-interactivity")
	hideWindow(c)
	out, _ := c.Output()
	return string(out)
}

// runOut draait een programma zonder venster en geeft de uitvoer terug.
func runOut(args []string) (string, error) {
	c := exec.Command(args[0], args[1:]...)
	hideWindow(c)
	out, err := c.Output()
	return string(out), err
}

// openDeviceManager via ShellExecute: mmc wil admin, en exec.Command kan
// geen programma starten dat om admin vraagt (daarom werkte de knop niet).
func openDeviceManager() error {
	verb, _ := windows.UTF16PtrFromString("open")
	file, _ := windows.UTF16PtrFromString("devmgmt.msc")
	return windows.ShellExecute(0, verb, file, nil, nil, windows.SW_NORMAL)
}

// runElevated draait een opdracht als administrator (Windows vraagt om
// toestemming). De uitvoer gaat naar een tijdelijk bestand dat we volgen.
func runElevated(args []string, line func(string)) (int, error) {
	f, err := os.CreateTemp("", "dusk-*.log")
	if err != nil {
		return -1, err
	}
	logPath := f.Name()
	f.Close()
	defer os.Remove(logPath)

	quoted := make([]string, len(args))
	for i, a := range args {
		if strings.ContainsAny(a, " &()^%!") {
			a = `"` + a + `"`
		}
		quoted[i] = a
	}
	params := `/c ` + strings.Join(quoted, " ") + ` > "` + logPath + `" 2>&1`
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString("cmd.exe")
	par, _ := windows.UTF16PtrFromString(params)
	info := &shellExecInfo{
		Mask:       0x00000040 | 0x00000400, // SEE_MASK_NOCLOSEPROCESS | SEE_MASK_FLAG_NO_UI
		Verb:       verb,
		File:       file,
		Parameters: par,
		Show:       windows.SW_HIDE,
	}
	if err := shellExecuteEx(info); err != nil {
		if errors.Is(err, windows.ERROR_CANCELLED) {
			return -1, errCancelled
		}
		return -1, err
	}
	defer windows.CloseHandle(info.Process)

	var off int64
	var rest string
	read := func() {
		lf, err := os.Open(logPath)
		if err != nil {
			return
		}
		defer lf.Close()
		lf.Seek(off, 0)
		buf, _ := io.ReadAll(lf)
		off += int64(len(buf))
		rest += string(buf)
		for {
			i := strings.IndexAny(rest, "\r\n")
			if i < 0 {
				break
			}
			l := strings.TrimSpace(rest[:i])
			rest = rest[i+1:]
			if strings.Trim(l, `-\|/ `) != "" {
				line(l)
			}
		}
	}
	for {
		ev, _ := windows.WaitForSingleObject(info.Process, 300)
		read()
		if ev == windows.WAIT_OBJECT_0 {
			break
		}
	}
	if l := strings.TrimSpace(rest); l != "" {
		line(l)
	}
	var code uint32
	windows.GetExitCodeProcess(info.Process, &code)
	return int(int32(code)), nil
}

func listAppx() ([]string, error) {
	c := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Get-AppxPackage | ForEach-Object { $_.Name }")
	hideWindow(c)
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, l := range strings.Split(string(out), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			names = append(names, l)
		}
	}
	return names, nil
}

func removeAppxArgs(pattern string) []string {
	return []string{"powershell", "-NoProfile", "-NonInteractive", "-Command",
		"Get-AppxPackage -Name '" + pattern + "' | Remove-AppxPackage -ErrorAction Stop"}
}

func rootKey(s string) registry.Key {
	if s == "HKLM" {
		return registry.LOCAL_MACHINE
	}
	return registry.CURRENT_USER
}

const (
	ultimateGUID = "e9a42b02-d5df-448d-aa00-03f14749eb61"
	balancedGUID = "381b4222-f694-41f0-9685-ff5bb260df2e"
)

// speciale tweaks die niet (alleen) in het register zitten

// Ultieme prestaties: het origineel (e9a42b02...) is een verborgen sjabloon dat je niet direct
// kunt aanzetten. Windows maakt er een kopie van met een eigen GUID. Die herkennen we aan de naam.
var reUltimate = regexp.MustCompile(`(?i)ultimate|ultieme|ultimative|höchstleistung|ultime|definitiv|máximo|massime|maksymalna|ultimativ`)
var reScheme = regexp.MustCompile(`([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})\s+\((.*?)\)`)

// ultimateScheme geeft de GUID van een bestaand plan "Ultieme prestaties" (leeg = geen).
func ultimateScheme() string {
	out, _ := runOut([]string{"powercfg", "/list"})
	for _, m := range reScheme.FindAllStringSubmatch(out, -1) {
		if strings.EqualFold(m[1], ultimateGUID) || reUltimate.MatchString(m[2]) {
			return m[1]
		}
	}
	return ""
}

func specialState(t Tweak) bool {
	switch t.Fn {
	case "power":
		out, _ := runOut([]string{"powercfg", "/getactivescheme"})
		m := reScheme.FindStringSubmatch(out)
		return m != nil && (strings.EqualFold(m[1], ultimateGUID) || reUltimate.MatchString(m[2]))
	case "hibernate":
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Power`, registry.QUERY_VALUE)
		if err != nil {
			return false
		}
		defer k.Close()
		n, _, err := k.GetIntegerValue("HibernateEnabled")
		return err == nil && n == 0
	}
	return false
}

func specialSet(t Tweak, on bool) error {
	switch t.Fn {
	case "power":
		if !on {
			_, err := runOut([]string{"powercfg", "/setactive", balancedGUID})
			return err
		}
		guid := ultimateScheme()
		if guid == "" {
			// nog niet aanwezig: kopie maken van het sjabloon; powercfg noemt de nieuwe GUID
			out, err := runOut([]string{"powercfg", "/duplicatescheme", ultimateGUID})
			m := reScheme.FindStringSubmatch(out)
			if m == nil {
				if err == nil {
					err = errors.New(strings.TrimSpace(out))
				}
				return err
			}
			guid = m[1]
		}
		_, err := runOut([]string{"powercfg", "/setactive", guid})
		return err
	case "hibernate":
		arg := "on"
		if on {
			arg = "off"
		}
		_, err := runOut([]string{"powercfg", "/hibernate", arg})
		return err
	}
	return nil
}

func tweakState(t Tweak) bool {
	if t.Fn != "" {
		return specialState(t)
	}
	for _, v := range t.Vals {
		k, err := registry.OpenKey(rootKey(v.Root), v.Path, registry.QUERY_VALUE)
		if err != nil {
			return false
		}
		ok := true
		switch v.Kind {
		case "dword":
			n, _, err := k.GetIntegerValue(v.Name)
			ok = err == nil && n == uint64(v.On.(int))
		case "sz":
			s, _, err := k.GetStringValue(v.Name)
			ok = err == nil && s == v.On.(string)
		}
		k.Close()
		if !ok {
			return false
		}
	}
	return true
}

func setTweak(t Tweak, on bool) error {
	if t.Fn != "" {
		return specialSet(t, on)
	}
	var errs []error
	for _, v := range t.Vals {
		root := rootKey(v.Root)
		if v.Kind == "key" {
			if on {
				k, _, err := registry.CreateKey(root, v.Path, registry.SET_VALUE)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				errs = append(errs, k.SetStringValue("", ""))
				k.Close()
			} else {
				registry.DeleteKey(root, v.Path)
				if i := strings.LastIndex(v.Path, `\`); i > 0 {
					registry.DeleteKey(root, v.Path[:i]) // lege oudersleutel ook weg
				}
			}
			continue
		}
		val := v.On
		if !on {
			val = v.Off
		}
		if val == nil {
			if k, err := registry.OpenKey(root, v.Path, registry.SET_VALUE); err == nil {
				if err := k.DeleteValue(v.Name); err != nil && !errors.Is(err, registry.ErrNotExist) {
					errs = append(errs, err)
				}
				k.Close()
			}
			continue
		}
		k, _, err := registry.CreateKey(root, v.Path, registry.SET_VALUE)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		switch v.Kind {
		case "dword":
			errs = append(errs, k.SetDWordValue(v.Name, uint32(val.(int))))
		case "sz":
			errs = append(errs, k.SetStringValue(v.Name, val.(string)))
		}
		k.Close()
	}
	return errors.Join(errs...)
}

func sysInfo() SysInfo {
	var s SysInfo
	if k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE); err == nil {
		name, _, _ := k.GetStringValue("ProductName")
		disp, _, _ := k.GetStringValue("DisplayVersion")
		build, _, _ := k.GetStringValue("CurrentBuild")
		k.Close()
		if b, _ := strconv.Atoi(build); b >= 22000 { // register zegt nog "Windows 10"
			name = strings.Replace(name, "Windows 10", "Windows 11", 1)
		}
		s.Windows = strings.TrimSpace(name + " " + disp)
		s.Build = build
	}
	if k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.QUERY_VALUE); err == nil {
		cpu, _, _ := k.GetStringValue("ProcessorNameString")
		k.Close()
		s.CPU = strings.TrimSpace(cpu)
	}
	var ms memStatus
	ms.Length = 64
	if r, _, _ := procMem.Call(uintptrOf(&ms)); r != 0 {
		s.RAM = int((ms.TotalPhys + (1 << 29)) >> 30)
	}
	s.PC, _ = os.Hostname()
	s.User = os.Getenv("USERNAME")
	return s
}

var procUILang = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultUILanguage")

func windowsIsDutch() bool {
	r, _, _ := procUILang.Call()
	return r&0x3ff == 0x13 // LANG_DUTCH
}
