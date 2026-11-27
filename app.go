package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const version = "1.9.1"

// App: alles wat de knoppen in het venster kunnen aanroepen.
type App struct {
	ctx  context.Context
	mu   sync.Mutex
	busy bool
	last map[string]AppState
	lang string
	hw   *HardwareInfo
}

type SysInfo struct {
	Windows, Build, CPU, PC, User, Version string
	RAM                                    int
	Admin, Winget                          bool
}

type AppInfo struct {
	ID, Name, Cat, Icon, Color, DescEN, DescNL string
	Installed                                  bool
	Version, Available                         string
	Pop                                        int // 1 = populairst, 0 = niet in de top
	Note                                       string
}

type TweakInfo struct {
	ID, Group                                      string
	NameEN, DescEN, NoteEN, NameNL, DescNL, NoteNL string
	InfoEN, InfoNL                                 string
	Changes                                        []string // wat er in het register verandert
	Admin, Explorer, On                            bool
}

type BloatInfo struct {
	NameEN, NameNL, Pattern, DescEN, DescNL string
}

// Progress: Step is een sleutel die de app vertaalt
// (queued, starting, downloading, installing, admin, done, failed, removing, removed).
type Progress struct {
	ID     string `json:"id"`
	Phase  string `json:"phase"` // wacht, bezig, klaar, fout
	Pct    int    `json:"pct"`
	Step   string `json:"step"`
	Detail string `json:"detail"`
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// SetLang: taal voor de regels in het activiteitenlog ("en" of "nl").
func (a *App) SetLang(l string) { a.lang = l }

var logText = map[string][2]string{ // en, nl
	"install":        {"Installing", "Installeren"},
	"update":         {"Updating", "Updaten"},
	"remove":         {"Removing", "Verwijderen"},
	"done":           {"done", "klaar"},
	"removed":        {"removed", "weg"},
	"failed":         {"failed", "mislukt"},
	"rescan":         {"Scanning apps again...", "Apps opnieuw scannen..."},
	"tweakon":        {"Tweak on", "Tweak aan"},
	"tweakoff":       {"Tweak off", "Tweak uit"},
	"explorer":       {"Explorer restarted", "Verkenner herstart"},
	"needadm":        {"Needs admin rights, asking Windows for permission...", "Heeft admin-rechten nodig, Windows vraagt om toestemming..."},
	"declined":       {"Permission was not given", "Geen toestemming gekregen"},
	"drivers":        {"Installing drivers", "Drivers installeren"},
	"closing":        {"Closing", "Sluiten"},
	"wingetinst":     {"Installing winget (App Installer)", "winget (App Installer) installeren"},
	"repair":         {"Repairing", "Repareren"},
	"uninstall":      {"Uninstalling", "Verwijderen"},
	"inuse":          {"in use, skipped", "in gebruik, overgeslagen"},
	"cleanjunk":      {"Cleaning junk files", "Rommelbestanden opruimen"},
	"restore":        {"Creating a restore point", "Herstelpunt maken"},
	"restore_recent": {"Windows already made one in the last 24 hours", "Windows heeft er de afgelopen 24 uur al een gemaakt"},
}

func (a *App) tr(k string) string {
	t := logText[k]
	if a.lang == "nl" {
		return t[1]
	}
	return t[0]
}

func (a *App) System() SysInfo {
	s := sysInfo()
	s.Admin, s.Winget, s.Version = isAdmin(), wingetAvailable(), version
	return s
}

func (a *App) appInfos(st map[string]AppState) []AppInfo {
	out := make([]AppInfo, 0, len(catalog))
	for _, c := range catalog {
		s := st[c.ID]
		out = append(out, AppInfo{c.ID, c.Name, c.Cat, c.Icon, c.Color, c.DescEN, c.DescNL, s.Installed, s.Version, s.Available, popRank(c.ID), s.Note})
	}
	return out
}

// ScanApps vraagt winget wat er op de pc staat (duurt een paar seconden).
func (a *App) ScanApps() ([]AppInfo, error) {
	st, err := scanApps()
	a.mu.Lock()
	a.last = st
	a.mu.Unlock()
	return a.appInfos(st), err
}

func (a *App) emit(name string, data any) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, name, data)
	}
}

func (a *App) log(s string) {
	if a.ctx == nil { // zonder venster (test via de opdrachtregel): gewoon printen
		fmt.Println(s)
		return
	}
	a.emit("log", s)
}

func (a *App) claim() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.busy {
		return false
	}
	a.busy = true
	return true
}

func (a *App) release() {
	a.mu.Lock()
	a.busy = false
	a.mu.Unlock()
}

// InstallApps installeert of updatet apps één voor één, met voortgang.
// Geeft "busy" terug als er al iets loopt.
func (a *App) InstallApps(ids []string) string {
	if !a.claim() {
		return "busy"
	}
	a.mu.Lock()
	st := a.last
	a.mu.Unlock()
	var list []CatApp
	for _, id := range ids {
		for _, c := range catalog {
			if c.ID == id {
				list = append(list, c)
			}
		}
	}
	for _, c := range list {
		a.emit("app", Progress{c.ID, "wacht", -1, "queued", ""})
	}
	go func() {
		defer a.release()
		ok, fail := 0, 0
		for _, c := range list {
			code, err := a.installOne(c, st[c.ID])
			if err == nil && wingetOK(code) {
				ok++
				a.log("  ✓ " + a.tr("done"))
				a.emit("app", Progress{c.ID, "klaar", 100, "done", ""})
				continue
			}
			fail++
			detail := hex(code)
			if err != nil {
				detail = err.Error()
			}
			a.log("  ✗ " + a.tr("failed") + " (" + detail + ")")
			a.emit("app", Progress{c.ID, "fout", -1, "failed", detail})
		}
		if ok > 0 {
			for _, n := range fixPinnedIcons() {
				a.log("  ✓ " + a.tr("iconfix") + ": " + n)
			}
		}
		a.log(a.tr("rescan"))
		infos, _ := a.ScanApps()
		a.emit("apps", infos)
		a.emit("done", map[string]any{"kind": "apps", "ok": ok, "fail": fail})
	}()
	return ""
}

// installOne draait winget voor één app. Lukt het niet omdat het
// installatieprogramma admin wil (0x8A150006), dan nog een keer als admin.
func (a *App) installOne(c CatApp, s AppState) (int, error) {
	what := a.tr("install")
	if s.Installed {
		what = a.tr("update")
	}
	a.log("» " + what + ": " + c.Name)
	a.emit("app", Progress{c.ID, "bezig", -1, "starting", ""})
	downloaded := false
	onLine := func(l string) {
		if p := progress(l); p >= 0 && (strings.ContainsAny(l, "█▒") || reSize.MatchString(l)) {
			if p >= 100 {
				downloaded = true
				a.emit("app", Progress{c.ID, "bezig", -1, "installing", ""})
			} else if !downloaded {
				a.emit("app", Progress{c.ID, "bezig", p, "downloading", ""})
			}
			return
		}
		a.log("  " + l)
		low := strings.ToLower(l)
		if strings.Contains(low, "install") && !strings.Contains(low, "installer hash") && !strings.Contains(low, "found") {
			downloaded = true
			a.emit("app", Progress{c.ID, "bezig", -1, "installing", ""})
		}
	}
	args := installArgs(c, s)
	if tag, ok := strings.CutPrefix(s.Note, "pyrt:"); ok && s.Installed {
		// Python via de Install Manager: de runtime zelf bijwerken
		args = []string{"py", "install", "-y", "--update", tag}
	}
	code, err := streamCmd(args, onLine)
	if err == nil && uint32(code) == 0x8A150006 && !isAdmin() {
		a.log("  " + a.tr("needadm"))
		a.emit("app", Progress{c.ID, "bezig", -1, "admin", ""})
		downloaded = false
		code, err = runElevated(args, onLine)
		if err == errCancelled {
			a.log("  " + a.tr("declined"))
		}
	}
	return code, err
}

func hex(code int) string {
	const d = "0123456789ABCDEF"
	u := uint32(code)
	b := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		b[i] = d[u&15]
		u >>= 4
	}
	return "0x" + string(b)
}

func makeTweakInfo(t Tweak) TweakInfo {
	var ch []string
	for _, v := range t.Vals {
		on := fmtVal(v.On)
		off := "(removed)"
		if v.Off != nil {
			off = fmtVal(v.Off)
		}
		if v.Kind == "key" {
			ch = append(ch, v.Root+`\`+v.Path+"  (key exists = on)")
			continue
		}
		ch = append(ch, v.Root+`\`+v.Path+`\`+v.Name+" = "+on+"   (off: "+off+")")
	}
	switch t.Fn {
	case "power":
		ch = []string{"powercfg /duplicatescheme " + "e9a42b02-d5df-448d-aa00-03f14749eb61 (Ultimate Performance)", "powercfg /setactive e9a42b02-...   (off: setactive Balanced 381b4222-...)"}
	case "hibernate":
		ch = []string{"powercfg /hibernate off   (off: powercfg /hibernate on)"}
	}
	inf := tweakInfo[t.ID]
	return TweakInfo{t.ID, t.Group, t.NameEN, t.DescEN, t.NoteEN, t.NameNL, t.DescNL, t.NoteNL, inf[0], inf[1], ch, t.Admin, t.Explorer, tweakState(t)}
}

func fmtVal(v any) string {
	switch x := v.(type) {
	case int:
		if x < 0 || x > 9999 {
			return "0x" + strings.TrimLeft(hex(x), "0x")
		}
		return strconv.Itoa(x)
	case string:
		return `"` + x + `"`
	}
	return ""
}

func (a *App) Tweaks() []TweakInfo {
	out := make([]TweakInfo, 0, len(tweaks))
	for _, t := range tweaks {
		out = append(out, makeTweakInfo(t))
	}
	return out
}

// SetTweak zet een tweak aan of uit en geeft de nieuwe stand terug.
func (a *App) SetTweak(id string, on bool) (TweakInfo, error) {
	for _, t := range tweaks {
		if t.ID != id {
			continue
		}
		err := setTweak(t, on)
		if err == nil {
			k, name := "tweakon", t.NameEN
			if !on {
				k = "tweakoff"
			}
			if a.lang == "nl" {
				name = t.NameNL
			}
			a.log("» " + a.tr(k) + ": " + name)
		}
		return makeTweakInfo(t), err
	}
	return TweakInfo{}, nil
}

func (a *App) RestartExplorer() error {
	a.log("» " + a.tr("explorer"))
	return restartExplorer()
}

func (a *App) ScanBloat() ([]BloatInfo, error) {
	names, err := listAppx()
	out := []BloatInfo{}
	for _, b := range matchBloat(names) {
		out = append(out, BloatInfo{b.NameEN, b.NameNL, b.Pattern, b.DescEN, b.DescNL})
	}
	return out, err
}

// RemoveBloat haalt de gekozen ingebouwde apps weg.
func (a *App) RemoveBloat(patterns []string) string {
	if !a.claim() {
		return "busy"
	}
	for _, p := range patterns {
		a.emit("bloat", Progress{p, "wacht", -1, "queued", ""})
	}
	go func() {
		defer a.release()
		ok, fail := 0, 0
		for _, p := range patterns {
			a.log("» " + a.tr("remove") + ": " + p)
			a.emit("bloat", Progress{p, "bezig", -1, "removing", ""})
			code, err := streamCmd(removeAppxArgs(p), func(l string) { a.log("  " + l) })
			if err == nil && code == 0 {
				ok++
				a.log("  ✓ " + a.tr("removed"))
				a.emit("bloat", Progress{p, "klaar", 100, "removed", ""})
			} else {
				fail++
				a.log("  ✗ " + a.tr("failed"))
				a.emit("bloat", Progress{p, "fout", -1, "failed", ""})
			}
		}
		a.emit("done", map[string]any{"kind": "bloat", "ok": ok, "fail": fail})
	}()
	return ""
}

// ---------- hardware en drivers ----------

func (a *App) Hardware() (HardwareInfo, error) {
	h, err := hardwareInfo()
	if err == nil {
		a.mu.Lock()
		a.hw = &h
		a.mu.Unlock()
	}
	return h, err
}

// InstallOEM downloadt pakketten van de maker (HP SoftPaqs) en installeert ze stil.
func (a *App) InstallOEM(ups []OEMUpdate) string {
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		ok, fail := 0, 0
		for _, u := range ups {
			a.log("» " + a.tr("install") + ": " + u.Name + " " + u.Version)
			a.emit("oem", Progress{u.ID, "bezig", 0, "downloading", ""})
			file, err := download(u.URL, func(p int) { a.emit("oem", Progress{u.ID, "bezig", p, "downloading", ""}) })
			if err == nil {
				a.emit("oem", Progress{u.ID, "bezig", -1, "installing", ""})
				args := []string{file, "/s"}
				var code int
				if isAdmin() {
					code, err = streamCmd(args, func(l string) { a.log("  " + l) })
				} else {
					code, err = runElevated(args, func(l string) { a.log("  " + l) })
				}
				os.Remove(file)
				if err == nil && code != 0 && code != 3010 && code != 1641 {
					err = errExit(code)
				}
			}
			if err != nil {
				fail++
				a.log("  ✗ " + err.Error())
				a.emit("oem", Progress{u.ID, "fout", -1, "failed", err.Error()})
			} else {
				ok++
				a.log("  ✓ " + a.tr("done"))
				a.emit("oem", Progress{u.ID, "klaar", 100, "done", ""})
			}
		}
		a.emit("done", map[string]any{"kind": "oem", "ok": ok, "fail": fail})
	}()
	return ""
}

// ScanDrivers zoekt via Windows Update naar drivers die nog missen (±30 s).
func (a *App) ScanDrivers() ([]DriverUpdate, error) { return scanDrivers() }

// InstallDrivers installeert de gekozen driver-updates (vraagt om admin als dat nodig is).
func (a *App) InstallDrivers(ids []string) string {
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		a.log("» " + a.tr("drivers"))
		for _, id := range ids {
			a.emit("driver", Progress{id, "bezig", -1, "installing", ""})
		}
		results := map[string]bool{}
		onLine := func(l string) {
			a.log("  " + l)
			if strings.HasPrefix(l, "RESULT ") {
				p := strings.SplitN(strings.TrimPrefix(l, "RESULT "), "|", 2)
				if len(p) == 2 {
					results[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1]) == "2" || strings.TrimSpace(p[1]) == "3"
				}
			}
		}
		args := driverInstallArgs(ids)
		var err error
		if isAdmin() {
			_, err = streamCmd(args, onLine)
		} else {
			a.log("  " + a.tr("needadm"))
			_, err = runElevated(args, onLine)
		}
		ok, fail := 0, 0
		for _, id := range ids {
			if err == nil && results[id] {
				ok++
				a.emit("driver", Progress{id, "klaar", 100, "done", ""})
			} else {
				fail++
				a.emit("driver", Progress{id, "fout", -1, "failed", ""})
			}
		}
		a.emit("done", map[string]any{"kind": "drivers", "ok": ok, "fail": fail})
	}()
	return ""
}

func (a *App) OpenDeviceManager() error { return openDeviceManager() }

// RunAsAdmin start Dusk opnieuw als administrator en sluit deze.
func (a *App) RunAsAdmin() error {
	if err := relaunchAdmin(); err != nil {
		return err
	}
	runtime.Quit(a.ctx)
	return nil
}

func hasArg(s string) bool {
	for _, x := range os.Args[1:] {
		if x == s {
			return true
		}
	}
	return false
}

// NoIntro: na "Start als admin" de intro niet nog een keer laten zien.
func (a *App) NoIntro() bool { return hasArg("--no-intro") }

// ---------- zelftest (verborgen venster, voor testen op afstand) ----------

func (a *App) SelfTest() bool { return hasArg("--selftest") }

// SelfTestDone schrijft het verslag van de webpagina weg en sluit af.
func (a *App) SelfTestDone(report string) {
	os.WriteFile(os.TempDir()+`\dusk-selftest.txt`, []byte(report), 0o644)
	runtime.Quit(a.ctx)
}

// ---------- opstart-apps, rommel, herstelpunt ----------

func (a *App) Startup() []StartupItem { return ListStartup() }

func (a *App) SetStartup(key string, on bool) error { return SetStartup(key, on) }

func (a *App) ScanJunk() []JunkCat { return ScanJunk() }

// CleanJunk ruimt de gekozen soorten op en geeft terug hoeveel bytes vrij kwamen.
func (a *App) CleanJunk(keys []string) int64 {
	a.log("» " + a.tr("cleanjunk"))
	freed := CleanJunk(keys)
	a.log("  ✓ " + strconv.FormatInt(freed/1024/1024, 10) + " MB")
	return freed
}

// CreateRestorePoint maakt een Windows-herstelpunt (vraagt om admin als dat nodig is).
// Antwoord: "ok", "recent" (vandaag al gemaakt), "cancelled" of een foutmelding.
func (a *App) CreateRestorePoint() string {
	a.log("» " + a.tr("restore"))
	script := `$ErrorActionPreference='Stop'
try { Enable-ComputerRestore -Drive "$env:SystemDrive\" } catch {}
try { Checkpoint-Computer -Description 'Dusk' -RestorePointType MODIFY_SETTINGS -WarningVariable w 3>$null; if($w){ 'RECENT' } else { 'OK' } }
catch { 'ERR ' + $_.Exception.Message }`
	var lines []string
	args := psArgs(script)
	var err error
	if isAdmin() {
		_, err = streamCmd(args, func(l string) { lines = append(lines, l) })
	} else {
		_, err = runElevated(args, func(l string) { lines = append(lines, l) })
	}
	if err == errCancelled {
		return "cancelled"
	}
	out := strings.Join(lines, " ")
	switch {
	case strings.Contains(out, "RECENT"):
		a.log("  ✓ " + a.tr("restore_recent"))
		return "recent"
	case strings.Contains(out, "OK"):
		a.log("  ✓ " + a.tr("done"))
		return "ok"
	}
	if err != nil {
		out = err.Error()
	}
	a.log("  ✗ " + out)
	return "ERR " + strings.TrimPrefix(out, "ERR ")
}

// InstallGPUDriver downloadt de nieuwste NVIDIA-driver en installeert hem stil.
func (a *App) InstallGPUDriver(name, url string) string {
	if !strings.HasPrefix(url, "https://") || !strings.Contains(url, ".nvidia.com/") {
		return "ERR bad url"
	}
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		a.log("» " + a.tr("install") + ": " + name)
		a.emit("gpu", Progress{name, "bezig", 0, "downloading", ""})
		file, err := download(url, func(p int) { a.emit("gpu", Progress{name, "bezig", p, "downloading", ""}) })
		code := 0
		if err == nil {
			a.emit("gpu", Progress{name, "bezig", -1, "installing", ""})
			args := []string{file, "-s", "-noreboot", "-noeula"}
			if isAdmin() {
				code, err = streamCmd(args, func(l string) { a.log("  " + l) })
			} else {
				code, err = runElevated(args, func(l string) { a.log("  " + l) })
			}
			os.RemoveAll(filepath.Dir(file))
			if err == nil && code != 0 && code != 1 && code != 3010 { // 1 = klaar, herstart aangeraden
				err = errExit(code)
			}
		}
		ok, fail := 1, 0
		if err != nil {
			ok, fail = 0, 1
			a.log("  ✗ " + err.Error())
			a.emit("gpu", Progress{name, "fout", -1, "failed", err.Error()})
		} else {
			a.log("  ✓ " + a.tr("done"))
			a.emit("gpu", Progress{name, "klaar", 100, "done", ""})
		}
		a.emit("done", map[string]any{"kind": "drivers", "ok": ok, "fail": fail})
	}()
	return ""
}

// RunOEMTool laat de officiële tool van het merk (Dell Command Update, Lenovo System
// Update) op de achtergrond alle drivers en BIOS bijwerken. Antwoord "notool" als de tool er niet is.
func (a *App) RunOEMTool(brand string) string {
	args := oemToolArgs(brand)
	if args == nil {
		return "notool"
	}
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		a.log("» " + a.tr("oemrun") + " (" + brand + ")")
		a.emit("oemtool", Progress{brand, "bezig", -1, "installing", ""})
		var code int
		var err error
		if isAdmin() {
			code, err = streamCmd(args, func(l string) { a.log("  " + l) })
		} else {
			code, err = runElevated(args, func(l string) { a.log("  " + l) })
		}
		// Dell: 0 = klaar, 1/5 = herstart nodig, 500 = niks te doen
		okCode := code == 0 || code == 1 || code == 5 || code == 500 || code == 3010
		if err == nil && !okCode {
			err = errExit(code)
		}
		ok, fail := 1, 0
		if err != nil {
			ok, fail = 0, 1
			a.log("  ✗ " + err.Error())
			a.emit("oemtool", Progress{brand, "fout", -1, "failed", err.Error()})
		} else {
			a.log("  ✓ " + a.tr("done"))
			a.emit("oemtool", Progress{brand, "klaar", 100, "done", ""})
		}
		a.emit("done", map[string]any{"kind": "drivers", "ok": ok, "fail": fail})
	}()
	return ""
}

// OpenVendorApp opent de driver-app van de maker (AMD Software, Intel DSA) als die er is.
func (a *App) OpenVendorApp(vendor string) bool { return openVendorApp(vendor) }

// UninstallApps verwijdert apps met winget en ruimt daarna de gekozen restjes op.
func (a *App) UninstallApps(ids []string, leftovers []string) string {
	if !a.claim() {
		return "busy"
	}
	a.mu.Lock()
	st := a.last
	a.mu.Unlock()
	go func() {
		defer a.release()
		ok, fail := 0, 0
		var freed int64
		var removed []string
		for _, id := range ids {
			var c CatApp
			for _, x := range catalog {
				if x.ID == id {
					c = x
				}
			}
			if c.ID == "" || noUninstall(c) {
				continue
			}
			s := st[id]
			a.log("» " + a.tr("uninstall") + ": " + c.Name)
			a.emit("app", Progress{id, "bezig", -1, "removing", ""})
			wid := s.RealID
			if wid == "" {
				wid = id
			}
			args := []string{"winget", "uninstall", "--id", wid, "--exact", "--silent", "--accept-source-agreements", "--disable-interactivity"}
			code, err := streamCmd(args, func(l string) { a.log("  " + l) })
			if err == nil && !wingetOK(code) && !isAdmin() {
				a.log("  " + a.tr("needadm"))
				code, err = runElevated(args, func(l string) { a.log("  " + l) })
			}
			if err != nil || !wingetOK(code) {
				fail++
				d := hex(code)
				if err != nil {
					d = err.Error()
				}
				a.log("  ✗ " + a.tr("failed") + " (" + d + ")")
				a.emit("app", Progress{id, "fout", -1, "failed", d})
				continue
			}
			ok++
			a.emit("app", Progress{id, "klaar", 100, "removed", ""})
		}
		// restjes: alleen wat de gebruiker heeft aangevinkt, en alleen als het nog bestaat.
		// Eerst alles afsluiten wat nog uit die mappen draait, dan pas weghalen.
		removeAll(a, leftovers, &freed, &removed)
		a.log(a.tr("rescan"))
		infos, _ := a.ScanApps()
		a.emit("apps", infos)
		a.emit("done", map[string]any{"kind": "uninstall", "ok": ok, "fail": fail, "freed": freed, "removed": removed})
	}()
	return ""
}

// removeAll sluit eerst alle processen in de restjes-mappen, en haalt ze dan weg (met een tweede poging).
func removeAll(a *App, leftovers []string, freed *int64, removed *[]string) {
	for _, p := range leftovers {
		if !strings.HasSuffix(strings.ToLower(p), ".lnk") {
			stopProcessesIn(p)
		}
	}
	for _, p := range leftovers {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		n, err := removeLeftover(p)
		if err != nil {
			time.Sleep(1500 * time.Millisecond)
			var n2 int64
			n2, err = removeLeftover(p)
			n += n2
		}
		*freed += n
		if err == nil {
			*removed = append(*removed, p)
			a.log("  ✓ " + a.tr("removed") + ": " + p)
		} else {
			a.log("  ✗ " + p + " (" + a.tr("inuse") + ")")
		}
	}
}
