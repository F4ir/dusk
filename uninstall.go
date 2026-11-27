package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Leftover is een map of snelkoppeling die een app achterlaat.
type Leftover struct {
	Path string
	Kind string // install, appdata, programdata, documents, shortcut
	Size int64
}

type arpEntry struct {
	DisplayName, Publisher, InstallLocation, UninstallString, DisplayIcon string
}

const arpScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$l=@(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*','HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*','HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*' | ? { $_.DisplayName } | % {[pscustomobject]@{DisplayName=[string]$_.DisplayName;Publisher=[string]$_.Publisher;InstallLocation=[string]$_.InstallLocation;UninstallString=[string]$_.UninstallString;DisplayIcon=[string]$_.DisplayIcon}})
ConvertTo-Json -InputObject $l -Compress`

// Mappen die nooit weg mogen, ook niet als ze toevallig zo heten.
var protectedNames = map[string]bool{
	"microsoft": true, "windows": true, "common files": true, "packages": true, "temp": true, "programs": true,
	"intel": true, "amd": true, "nvidia": true, "nvidia corporation": true, "google": true, "mozilla": true,
	"apple": true, "adobe": true, "logitech": true, "elgato": true, "corsair": true, "valve": true, "steam": true,
	"windowsapps": true, "installer": true, "package cache": true, "microsoft.net": true, "dotnet": true,
	"crashdumps": true, "d3dscache": true, "nuget": true, "npm": true, "pip": true, "start menu": true,
	// eigen spullen van de gebruiker
	"files": true, "music": true, "videos": true, "pictures": true, "downloads": true, "desktop": true, "documents": true,
	"games": true, "my games": true, "saved games": true, "projects": true, "backup": true, "backups": true,
	"code": true, "github": true, "source": true, "repos": true,
	// mappen met game-bibliotheken
	"steamlibrary": true, "epic games": true, "ubisoft": true, "ubisoft game launcher": true, "gog galaxy": true,
	"ea games": true, "electronic arts": true, "origin games": true, "riot games": true, "rockstar games": true,
}

// noUninstall: onderdelen waar andere programma's en games op draaien.
func noUninstall(c CatApp) bool {
	return c.Cat == "runtime" || c.ID == "Microsoft.Edge" || c.ID == "Microsoft.WindowsTerminal"
}

func dirOf(s string) string {
	s = strings.Trim(strings.TrimSpace(s), `"`)
	s = regexp.MustCompile(`,\s*-?\d+$`).ReplaceAllString(s, "")
	if i := strings.Index(strings.ToLower(s), ".exe"); i >= 0 {
		s = s[:i+4]
	}
	if strings.HasSuffix(strings.ToLower(s), ".exe") || strings.HasSuffix(strings.ToLower(s), ".ico") {
		s = filepath.Dir(s)
	}
	return strings.TrimRight(s, `\`)
}

// roots: plekken waar apps hun bestanden laten staan.
func leftoverRoots() map[string]string {
	home := os.Getenv("USERPROFILE")
	return map[string]string{
		os.Getenv("ProgramFiles"):      "install",
		os.Getenv("ProgramFiles(x86)"): "install",
		os.Getenv("APPDATA"):           "appdata",
		os.Getenv("LOCALAPPDATA"):      "appdata",
		home + `\AppData\LocalLow`:     "appdata",
		os.Getenv("ProgramData"):       "programdata",
		home + `\Documents`:            "documents",
	}
}

// safeToDelete: alleen mappen onder een bekende plek, minstens één niveau diep,
// en niet met een algemene naam (Microsoft, Program Files zelf, ...).
func safeToDelete(p string) (string, bool) {
	clean := filepath.Clean(p)
	low := strings.ToLower(clean) + `\`
	for _, bad := range []string{`\package cache\`, `\windowsapps\`, `\installer\`, `\windows\`, `\system32\`} {
		if strings.Contains(low, bad) {
			return "", false // installatie-cache en Windows zelf nooit
		}
	}
	for root, kind := range leftoverRoots() {
		if root == "" || root == "." {
			continue
		}
		rel, err := filepath.Rel(filepath.Clean(root), clean)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			continue
		}
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			if protectedNames[strings.ToLower(part)] && !strings.Contains(rel, string(filepath.Separator)) {
				return "", false // bv. AppData\Roaming\Microsoft zelf
			}
		}
		// vendor-map (Google, Logitech, ...) zelf nooit; vendor\app wel
		if protectedNames[strings.ToLower(filepath.Base(clean))] {
			return "", false
		}
		return kind, true
	}
	return "", false
}

func appNames(c CatApp, arp *arpEntry) []string {
	names := append([]string{c.Name}, arpNames[c.ID]...)
	if arp != nil {
		names = append(names, arp.DisplayName)
		if d := dirOf(arp.InstallLocation); d != "" {
			names = append(names, filepath.Base(d))
		}
		if d := dirOf(arp.DisplayIcon); d != "" && !strings.EqualFold(filepath.Base(d), "bin") {
			names = append(names, filepath.Base(d))
		}
	}
	// "Mozilla Firefox (x64 en-US)" -> "Mozilla Firefox"; ook zonder spaties ("OBS Studio" -> "obs-studio")
	var out []string
	seen := map[string]bool{}
	add := func(n string) {
		n = strings.TrimSpace(n)
		if len(n) < 3 || seen[strings.ToLower(n)] || protectedNames[strings.ToLower(n)] {
			return
		}
		seen[strings.ToLower(n)] = true
		out = append(out, n)
	}
	for _, n := range names {
		n = regexp.MustCompile(`\s*\(.*?\)|\s+\d[\d.]*$`).ReplaceAllString(n, "")
		add(n)
		add(strings.ReplaceAll(n, " ", ""))
		add(strings.ReplaceAll(strings.ToLower(n), " ", "-"))
	}
	return out
}

func vendorOf(arp *arpEntry) string {
	if arp == nil {
		return ""
	}
	v := regexp.MustCompile(`(?i),?\s*(inc\.?|llc|ltd\.?|gmbh|corporation|corp\.?|co\.?|b\.v\.)$`).ReplaceAllString(strings.TrimSpace(arp.Publisher), "")
	return strings.TrimSpace(v)
}

func folderSize(p string) int64 {
	var n int64
	filepath.WalkDir(p, func(_ string, e os.DirEntry, err error) error {
		if err == nil && !e.IsDir() {
			if fi, err := e.Info(); err == nil {
				n += fi.Size()
			}
		}
		return nil
	})
	return n
}

// findLeftovers zoekt alle mappen en snelkoppelingen die bij een app horen.
func findLeftovers(c CatApp, arp *arpEntry) []Leftover {
	var out []Leftover
	seen := map[string]bool{}
	add := func(p, kind string) {
		p = filepath.Clean(p)
		if seen[strings.ToLower(p)] {
			return
		}
		fi, err := os.Stat(p)
		if err != nil || !fi.IsDir() {
			return
		}
		k, ok := safeToDelete(p)
		if !ok {
			return
		}
		if kind == "" {
			kind = k
		}
		// map die al in een gevonden map zit niet nog een keer noemen
		for _, o := range out {
			if strings.HasPrefix(strings.ToLower(p), strings.ToLower(o.Path)+`\`) {
				return
			}
		}
		seen[strings.ToLower(p)] = true
		out = append(out, Leftover{p, kind, folderSize(p)})
	}
	if arp != nil {
		if d := dirOf(arp.InstallLocation); d != "" {
			add(d, "install")
		}
	}
	// de map van de uninstaller (bv. C:\Program Files\Foo\uninstall.exe) is vaak de programmamap
	if arp != nil && !strings.Contains(strings.ToLower(arp.UninstallString), "msiexec") {
		if d := dirOf(arp.UninstallString); d != "" {
			add(d, "install")
		}
		if d := dirOf(arp.DisplayIcon); d != "" && !strings.EqualFold(filepath.Base(d), "bin") {
			add(d, "install")
		}
	}
	names := appNames(c, arp)
	vendor := vendorOf(arp)
	// "Google Chrome" van "Google" -> ook "Google\Chrome"
	if vendor != "" {
		for _, n := range append([]string{}, names...) {
			if rest := strings.TrimSpace(strings.TrimPrefix(n, vendor)); rest != n && len(rest) >= 3 {
				names = append(names, rest)
			}
		}
	}
	for root := range leftoverRoots() {
		if root == "" {
			continue
		}
		for _, n := range names {
			add(filepath.Join(root, n), "")
			if vendor != "" && !strings.EqualFold(vendor, n) {
				add(filepath.Join(root, vendor, n), "")
			}
		}
	}
	// snelkoppelingen op het bureaublad en in Start die naar de app wijzen (op naam)
	home := os.Getenv("USERPROFILE")
	for _, dir := range []string{home + `\Desktop`, os.Getenv("PUBLIC") + `\Desktop`, os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs`, os.Getenv("ProgramData") + `\Microsoft\Windows\Start Menu\Programs`} {
		for _, n := range names {
			for _, pat := range []string{filepath.Join(dir, n+".lnk"), filepath.Join(dir, n, "*.lnk")} {
				m, _ := filepath.Glob(pat)
				for _, f := range m {
					if !seen[strings.ToLower(f)] {
						seen[strings.ToLower(f)] = true
						out = append(out, Leftover{f, "shortcut", 0})
					}
				}
			}
		}
	}
	return out
}

func readARP() []arpEntry {
	out, _ := runOut(psArgs(arpScript))
	var l []arpEntry
	json.Unmarshal([]byte(strings.TrimSpace(out)), &l)
	return l
}

// UninstallPlan: per app wat Dusk zou weghalen (voor het "weet je het zeker?"-venster).
func (a *App) UninstallPlan(ids []string) map[string][]Leftover {
	arps := readARP()
	res := map[string][]Leftover{}
	for _, id := range ids {
		for _, c := range catalog {
			if c.ID != id || noUninstall(c) {
				continue
			}
			var arp *arpEntry
			for i := range arps {
				if nameMatches(arps[i].DisplayName, c) {
					arp = &arps[i]
					break
				}
			}
			res[id] = findLeftovers(c, arp)
		}
	}
	return res
}

// removeLeftover haalt een map of snelkoppeling weg (nog een keer gecontroleerd).
// Programma's die nog uit die map draaien worden eerst gesloten (anders zijn de bestanden in gebruik).
func removeLeftover(p string) (int64, error) {
	if strings.HasSuffix(strings.ToLower(p), ".lnk") {
		err := os.Remove(p)
		// lege Start-menumap van de app ook weg
		if dir := filepath.Dir(p); strings.Contains(strings.ToLower(dir), `\start menu\programs\`) {
			os.Remove(dir)
		}
		return 0, err
	}
	stopProcessesIn(p)
	if _, ok := safeToDelete(p); !ok {
		return 0, os.ErrPermission
	}
	size := folderSize(p)
	err := os.RemoveAll(p)
	if err != nil {
		size -= folderSize(p) // wat bleef staan (in gebruik) niet meetellen
	}
	return size, err
}

// stopProcessesIn sluit processen waarvan het programma in deze map staat.
func stopProcessesIn(dir string) {
	if _, ok := safeToDelete(dir); !ok {
		return
	}
	d := strings.ReplaceAll(filepath.Clean(dir), "'", "''")
	runOut(psArgs(`Get-Process | Where-Object { $_.Path -and $_.Path.StartsWith('` + d + `\', [StringComparison]::OrdinalIgnoreCase) } | Stop-Process -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 700`))
}

var reCamel = regexp.MustCompile(`([a-z])([A-Z][a-z])`)

// prettyPkg: "WinRAR.ShellExtension" -> "WinRAR", "WindowsTerminal" -> "Windows Terminal",
// "TreeSizeFreeContextMenu" -> "Tree Size Free".
func prettyPkg(n string) string {
	if strings.HasSuffix(n, "-RSXCM") {
		return "AMD Software"
	}
	if strings.Contains(n, ".") {
		parts := strings.Split(n, ".")
		n = parts[len(parts)-1]
		generic := map[string]bool{"shellextension": true, "contextmenu": true, "shellext": true, "extension": true}
		if generic[strings.ToLower(n)] && len(parts) > 1 {
			n = parts[len(parts)-2]
		}
	}
	for _, suf := range []string{"ContextMenu", "ShellExtension", "Shell"} {
		if strings.HasSuffix(n, suf) && len(n) > len(suf) {
			n = strings.TrimSuffix(n, suf)
		}
	}
	return reCamel.ReplaceAllString(n, "$1 $2")
}
