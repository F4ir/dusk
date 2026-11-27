//go:build windows

package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// ---------- opstart-apps ----------

type StartupItem struct {
	Key, Name, Desc, Maker, Command, Source string
	Enabled, Admin                          bool
	Essential                               bool // driver- of systeemonderdeel: niet tonen, moet aan blijven
}

// Opstart-items van drivers, hardware-makers en Windows-beveiliging moeten aan blijven
// (audio-instellingen, sneltoetsen van laptops, touchpad, videokaart). Die laat Dusk niet zien.
var reEssentialMaker = regexp.MustCompile(`(?i)realtek|advanced micro|^amd\b|\bamd\b|nvidia|intel|synaptics|elan micro|elantech|conexant|dolby|waves|cirrus|alps|qualcomm|mediatek|broadcom|hewlett|^hp\b|hp inc|dell|lenovo|asustek|^asus|acer|micro-star|^msi\b|nahimic|a-volute|sonic studio|thx`)
var reEssentialName = regexp.MustCompile(`(?i)securityhealth|windows ?defender|rtkaud|rthdvcpl|realtek|igfx|nvbackend|nvcpl|radeon|amd|maxxaudio|hotkey|touchpad|synaptics|smartaudio|iastor|sonicstudio|nahimic`)

func essentialStartup(it StartupItem) bool {
	if reEssentialMaker.MatchString(it.Maker) || reEssentialName.MatchString(it.Name) {
		return true
	}
	// Windows zelf (SecurityHealth e.d.), maar niet apps zoals OneDrive of Teams
	if strings.HasPrefix(strings.ToLower(it.Maker), "microsoft") && !regexp.MustCompile(`(?i)onedrive|teams|edge|skype|copilot`).MatchString(it.Name+" "+it.Desc) {
		return true
	}
	return false
}

const (
	runKey      = `Software\Microsoft\Windows\CurrentVersion\Run`
	run32Key    = `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`
	approvedKey = `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\`
)

type startupSrc struct {
	id       string
	root     registry.Key
	runPath  string // leeg = map
	folder   string
	approved string // subsleutel onder StartupApproved
	admin    bool
}

func startupSources() []startupSrc {
	return []startupSrc{
		{"hkcu", registry.CURRENT_USER, runKey, "", "Run", false},
		{"hklm", registry.LOCAL_MACHINE, runKey, "", "Run", true},
		{"hklm32", registry.LOCAL_MACHINE, run32Key, "", "Run32", true},
		{"folder", registry.CURRENT_USER, "", os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup`, "StartupFolder", false},
		{"common", registry.LOCAL_MACHINE, "", os.Getenv("ProgramData") + `\Microsoft\Windows\Start Menu\Programs\StartUp`, "StartupFolder", true},
	}
}

// approvedState: geen waarde = aan; eerste byte even (02/06) = aan, oneven (03/07) = uit.
func approvedState(root registry.Key, sub, name string) bool {
	k, err := registry.OpenKey(root, approvedKey+sub, registry.QUERY_VALUE)
	if err != nil {
		return true
	}
	defer k.Close()
	b, _, err := k.GetBinaryValue(name)
	if err != nil || len(b) == 0 {
		return true
	}
	return b[0]%2 == 0
}

func ListStartup() []StartupItem {
	var out []StartupItem
	for _, s := range startupSources() {
		if s.runPath != "" {
			k, err := registry.OpenKey(s.root, s.runPath, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			names, _ := k.ReadValueNames(0)
			for _, n := range names {
				cmd, _, _ := k.GetStringValue(n)
				if strings.TrimSpace(cmd) == "" {
					continue
				}
				it := StartupItem{Key: s.id + "|" + n, Name: n, Command: cmd, Source: s.id, Admin: s.admin, Enabled: approvedState(s.root, s.approved, n)}
				it.Desc, it.Maker = fileInfo(exeFromCommand(cmd))
				it.Essential = essentialStartup(it)
				out = append(out, it)
			}
			k.Close()
			continue
		}
		files, _ := os.ReadDir(s.folder)
		for _, f := range files {
			if f.IsDir() || strings.EqualFold(f.Name(), "desktop.ini") {
				continue
			}
			name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			it := StartupItem{Key: s.id + "|" + f.Name(), Name: name, Command: filepath.Join(s.folder, f.Name()), Source: s.id, Admin: s.admin,
				Enabled: approvedState(s.root, s.approved, f.Name())}
			it.Essential = essentialStartup(it)
			out = append(out, it)
		}
	}
	return out
}

func SetStartup(key string, on bool) error {
	p := strings.SplitN(key, "|", 2)
	if len(p) != 2 {
		return nil
	}
	for _, s := range startupSources() {
		if s.id != p[0] {
			continue
		}
		k, _, err := registry.CreateKey(s.root, approvedKey+s.approved, registry.SET_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		b := make([]byte, 12)
		if on {
			b[0] = 2
		} else {
			b[0] = 3
			ft := windows.NsecToFiletime(time.Now().UnixNano())
			binary.LittleEndian.PutUint32(b[4:], ft.LowDateTime)
			binary.LittleEndian.PutUint32(b[8:], ft.HighDateTime)
		}
		return k.SetBinaryValue(p[1], b)
	}
	return nil
}

var reEnv = regexp.MustCompile(`%[A-Za-z0-9_()]+%`)

// exeFromCommand: "C:\x\app.exe" --min  ->  C:\x\app.exe
func exeFromCommand(c string) string {
	c = strings.TrimSpace(reEnv.ReplaceAllStringFunc(c, func(m string) string {
		if v := os.Getenv(strings.Trim(m, "%")); v != "" {
			return v
		}
		return m
	}))
	if strings.HasPrefix(c, `"`) {
		if i := strings.Index(c[1:], `"`); i >= 0 {
			return c[1 : i+1]
		}
	}
	if i := strings.Index(strings.ToLower(c), ".exe"); i >= 0 {
		return c[:i+4]
	}
	return c
}

// fileInfo leest FileDescription en CompanyName uit een exe.
func fileInfo(path string) (desc, maker string) {
	size, err := windows.GetFileVersionInfoSize(path, nil)
	if err != nil || size == 0 {
		return "", ""
	}
	buf := make([]byte, size)
	if windows.GetFileVersionInfo(path, 0, size, unsafe.Pointer(&buf[0])) != nil {
		return "", ""
	}
	var tr unsafe.Pointer
	var trLen uint32
	if windows.VerQueryValue(unsafe.Pointer(&buf[0]), `\VarFileInfo\Translation`, unsafe.Pointer(&tr), &trLen) != nil || trLen < 4 {
		return "", ""
	}
	t := (*[2]uint16)(tr)
	get := func(field string) string {
		var p unsafe.Pointer
		var n uint32
		q := `\StringFileInfo\` + hex4(t[0]) + hex4(t[1]) + `\` + field
		if windows.VerQueryValue(unsafe.Pointer(&buf[0]), q, unsafe.Pointer(&p), &n) != nil || n == 0 {
			return ""
		}
		return strings.TrimSpace(windows.UTF16PtrToString((*uint16)(p)))
	}
	return get("FileDescription"), get("CompanyName")
}

func hex4(v uint16) string {
	const d = "0123456789abcdef"
	return string([]byte{d[v>>12], d[v>>8&15], d[v>>4&15], d[v&15]})
}

// ---------- rommelbestanden ----------

type JunkCat struct {
	Key   string
	Size  int64
	Files int
	Admin bool
}

func junkPaths() map[string][]string {
	la, ra, pd, win := os.Getenv("LOCALAPPDATA"), os.Getenv("APPDATA"), os.Getenv("ProgramData"), os.Getenv("SystemRoot")
	glob := func(pats ...string) []string {
		var out []string
		for _, p := range pats {
			m, _ := filepath.Glob(p)
			out = append(out, m...)
		}
		return out
	}
	chromium := func(base string) []string {
		return glob(base+`\User Data\*\Cache`, base+`\User Data\*\Code Cache`, base+`\User Data\*\GPUCache`, base+`\User Data\ShaderCache`, base+`\User Data\GrShaderCache`)
	}
	browsers := append(chromium(la+`\Google\Chrome`), chromium(la+`\Microsoft\Edge`)...)
	browsers = append(browsers, chromium(la+`\BraveSoftware\Brave-Browser`)...)
	browsers = append(browsers, glob(la+`\Mozilla\Firefox\Profiles\*\cache2`, la+`\Opera Software\Opera GX Stable\Cache`, la+`\Opera Software\Opera Stable\Cache`)...)
	return map[string][]string{
		"temp":    {os.TempDir()},
		"wintemp": {win + `\Temp`},
		"wucache": {win + `\SoftwareDistribution\Download`},
		"shader":  glob(la+`\D3DSCache`, la+`\NVIDIA\DXCache`, la+`\NVIDIA\GLCache`, la+`\AMD\DxCache`, la+`\AMD\DxcCache`, la+`\AMD\VkCache`, la+`\AMD\GLCache`, la+`\Intel\ShaderCache`),
		"browser": browsers,
		"discord": glob(ra+`\discord\Cache`, ra+`\discord\Code Cache`, ra+`\discord\GPUCache`),
		"crash":   glob(la+`\CrashDumps`, pd+`\Microsoft\Windows\WER\ReportArchive`, pd+`\Microsoft\Windows\WER\ReportQueue`, win+`\Minidump`),
		"recycle": nil,
	}
}

var junkAdmin = map[string]bool{"wintemp": true, "wucache": true}

func dirSize(dirs []string) (int64, int) {
	var size int64
	n := 0
	for _, d := range dirs {
		filepath.WalkDir(d, func(p string, e os.DirEntry, err error) error {
			if err != nil || e.IsDir() {
				return nil
			}
			if fi, err := e.Info(); err == nil {
				size += fi.Size()
				n++
			}
			return nil
		})
	}
	return size, n
}

func ScanJunk() []JunkCat {
	var out []JunkCat
	for _, k := range []string{"temp", "shader", "browser", "discord", "crash", "recycle", "wintemp", "wucache"} {
		var size int64
		var n int
		if k == "recycle" {
			size, n = recycleInfo()
		} else {
			size, n = dirSize(junkPaths()[k])
		}
		out = append(out, JunkCat{k, size, n, junkAdmin[k]})
	}
	return out
}

// CleanJunk verwijdert bestanden; wat in gebruik is wordt overgeslagen.
func CleanJunk(keys []string) int64 {
	var freed int64
	cutoff := time.Now().Add(-time.Hour)
	for _, k := range keys {
		if k == "recycle" {
			s, _ := recycleInfo()
			if emptyRecycle() == nil {
				freed += s
			}
			continue
		}
		for _, d := range junkPaths()[k] {
			var dirs []string
			filepath.WalkDir(d, func(p string, e os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if e.IsDir() {
					if p != d {
						dirs = append(dirs, p)
					}
					return nil
				}
				fi, err := e.Info()
				if err != nil {
					return nil
				}
				if k == "temp" && fi.ModTime().After(cutoff) { // net gemaakte tijdelijke bestanden laten staan
					return nil
				}
				if os.Remove(p) == nil {
					freed += fi.Size()
				}
				return nil
			})
			for i := len(dirs) - 1; i >= 0; i-- {
				os.Remove(dirs[i]) // alleen lege mappen gaan weg
			}
		}
	}
	return freed
}

type shqueryrbinfo struct {
	Size     uint32
	_        uint32
	BinSize  int64
	NumItems int64
}

var (
	shell32         = windows.NewLazySystemDLL("shell32.dll")
	procQueryRecBin = shell32.NewProc("SHQueryRecycleBinW")
	procEmptyRecBin = shell32.NewProc("SHEmptyRecycleBinW")
)

func recycleInfo() (int64, int) {
	var i shqueryrbinfo
	i.Size = uint32(unsafe.Sizeof(i))
	r, _, _ := procQueryRecBin.Call(0, uintptr(unsafe.Pointer(&i)))
	if r != 0 {
		return 0, 0
	}
	return i.BinSize, int(i.NumItems)
}

func emptyRecycle() error {
	// SHERB_NOCONFIRMATION | SHERB_NOPROGRESSUI | SHERB_NOSOUND
	r, _, err := procEmptyRecBin.Call(0, 0, 0x1|0x2|0x4)
	if r != 0 && r != 0x8000FFFF {
		return err
	}
	return nil
}

// fixPinnedIcons: na een update kan een vastgemaakte taakbalk-snelkoppeling een wit icoon
// krijgen, als zijn icoon naar een bestand wijst dat de update heeft weggehaald (Discord e.d.
// wisselen van map). Alleen die snelkoppelingen krijgen het icoon van het programma terug.
// Daarna krijgt Windows een seintje om iconen te verversen (SHChangeNotify): geen
// herstart van Verkenner, dus de taakbalk knippert niet.
func fixPinnedIcons() []string {
	script := `$sh=New-Object -ComObject WScript.Shell
$cut=(Get-Date).AddHours(-2)
foreach($f in Get-ChildItem "$env:APPDATA\Microsoft\Internet Explorer\Quick Launch\User Pinned\TaskBar\*.lnk" -ErrorAction SilentlyContinue){
  $l=$sh.CreateShortcut($f.FullName); $ico=(($l.IconLocation -split ',')[0]).Trim()
  if(-not $l.TargetPath -or -not (Test-Path $l.TargetPath)){ continue }
  if((Get-Item $l.TargetPath).LastWriteTime -lt $cut){ continue }
  if(-not $ico -or (Test-Path ([Environment]::ExpandEnvironmentVariables($ico)))){ continue }
  $appico=Join-Path (Split-Path $l.TargetPath) 'app.ico'
  if(Test-Path $appico){ $l.IconLocation="$appico,0" } else { $l.IconLocation="$($l.TargetPath),0" }
  $l.Save(); $f.FullName
}`
	out, _ := runOut(psArgs(script))
	var fixed []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			p, _ := windows.UTF16PtrFromString(l)
			procSHChangeNotify.Call(0x00002000, 0x0005, uintptr(unsafe.Pointer(p)), 0) // SHCNE_UPDATEITEM, SHCNF_PATHW
			fixed = append(fixed, strings.TrimSuffix(filepath.Base(l), ".lnk"))
		}
	}
	if len(fixed) > 0 {
		procSHChangeNotify.Call(0x08000000, 0, 0, 0) // SHCNE_ASSOCCHANGED: iconen opnieuw laden
	}
	return fixed
}

var procSHChangeNotify = shell32.NewProc("SHChangeNotify")

// oemToolArgs: opdrachtregel van de officiële updatetool (nil = niet geïnstalleerd).
func oemToolArgs(brand string) []string {
	pf, pf86 := os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)")
	exists := func(p string) bool { _, err := os.Stat(p); return err == nil }
	switch brand {
	case "dell":
		for _, p := range []string{pf + `\Dell\CommandUpdate\dcu-cli.exe`, pf86 + `\Dell\CommandUpdate\dcu-cli.exe`} {
			if exists(p) {
				return []string{p, "/applyUpdates", "-reboot=disable", "-silent"}
			}
		}
	case "lenovo":
		for _, p := range []string{pf86 + `\Lenovo\System Update\tvsu.exe`, pf + `\Lenovo\System Update\tvsu.exe`} {
			if exists(p) {
				return []string{p, "/CM", "-search", "A", "-action", "INSTALL", "-includerebootpackages", "1,3,4", "-noicon", "-noreboot", "-nolicense"}
			}
		}
	}
	return nil
}

func openVendorApp(vendor string) bool {
	pf, pf86 := os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)")
	var cands []string
	switch vendor {
	case "amd":
		cands = []string{pf + `\AMD\CNext\CNext\RadeonSoftware.exe`}
	case "intel":
		cands = []string{pf86 + `\Intel\Driver and Support Assistant\DSATray.exe`, pf + `\Intel\Driver and Support Assistant\DSATray.exe`}
	}
	for _, p := range cands {
		if _, err := os.Stat(p); err == nil {
			verb, _ := windows.UTF16PtrFromString("open")
			file, _ := windows.UTF16PtrFromString(p)
			return windows.ShellExecute(0, verb, file, nil, nil, windows.SW_NORMAL) == nil
		}
	}
	return false
}
