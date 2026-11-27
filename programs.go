package main

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Program is een programma uit "Apps en onderdelen" van Windows, ook als het niet in Dusk staat.
type Program struct {
	Key, Name, Publisher, Version, InstallDate, Location string
	Uninstall, Quiet, Icon, Hive, KeyName, CatalogID     string
	Logo                                                 string // icoon uit Dusk als het programma in de lijst staat
	SizeMB                                               int
	Protected                                            bool // drivers, runtimes, Windows: niet via Dusk weghalen
}

const progScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$src=@(@('HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall','Machine\X64'),@('HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall','Machine\X86'),@('HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall','User\X64'))
$l=@(foreach($s in $src){ Get-ChildItem $s[0] | % { $p=Get-ItemProperty $_.PSPath
  if(-not $p.DisplayName -or $p.SystemComponent -eq 1 -or $p.ParentKeyName -or $p.ReleaseType -match 'Update|Hotfix' -or -not $p.UninstallString){ return }
  [pscustomobject]@{Key=$_.Name;KeyName=$_.PSChildName;Hive=$s[1];Name=[string]$p.DisplayName;Publisher=[string]$p.Publisher;Version=[string]$p.DisplayVersion;InstallDate=[string]$p.InstallDate;Location=[string]$p.InstallLocation;Uninstall=[string]$p.UninstallString;Quiet=[string]$p.QuietUninstallString;Icon=[string]$p.DisplayIcon;SizeMB=[int]([int64]$p.EstimatedSize/1024)} } })
ConvertTo-Json -InputObject $l -Compress`

var reProtected = regexp.MustCompile(`(?i)^microsoft (visual c\+\+|\.net|windows (desktop runtime|sdk|app runtime)|update health|gameinput|xna|edge( webview2)?|asp\.net|office click-to-run)|redistributable|\bruntime\b|\bdriver(s)?\b|chipset|realtek|amd software|amd ryzen master sdk|nvidia (graphics|physx|hd audio|frameview|usbc)|intel\(?r?\)? .*(management engine|serial io|rapid storage|chipset|wireless bluetooth|software installer)|directx|vcredist|windows (subsystem|pc health)|application compatibility`)

// tools mét "driver" in de naam zijn gewone apps
var reNotProtected = regexp.MustCompile(`(?i)support assistant|driver uninstaller|driver booster|driver easy|install manager`)

func (a *App) Programs() []Program {
	out, _ := runOut(psArgs(progScript))
	var list []Program
	json.Unmarshal([]byte(strings.TrimSpace(out)), &list)
	seen := map[string]bool{}
	var res []Program
	for _, p := range list {
		k := strings.ToLower(p.Name + "|" + p.Version)
		if seen[k] {
			continue
		}
		seen[k] = true
		p.Protected = reProtected.MatchString(p.Name) && !reNotProtected.MatchString(p.Name)
		for _, c := range catalog {
			if nameMatches(p.Name, c) {
				p.CatalogID, p.Logo = c.ID, c.Icon
				break
			}
		}
		res = append(res, p)
	}
	return res
}

func (a *App) findProgram(key string) *Program {
	for _, p := range a.Programs() {
		if p.Key == key {
			p := p
			return &p
		}
	}
	return nil
}

// ProgramPlan: restjes per programma (voor het "weet je het zeker?"-venster).
func (a *App) ProgramPlan(keys []string) map[string][]Leftover {
	progs := a.Programs()
	res := map[string][]Leftover{}
	for _, k := range keys {
		for _, p := range progs {
			if p.Key != k || p.Protected {
				continue
			}
			c := CatApp{Name: p.Name, ID: p.CatalogID}
			if p.CatalogID != "" {
				for _, x := range catalog {
					if x.ID == p.CatalogID {
						c = x
					}
				}
			}
			res[k] = findLeftovers(c, &arpEntry{p.Name, p.Publisher, p.Location, p.Uninstall, p.Icon})
		}
	}
	return res
}

// uninstallArgs: eerst via winget (kent stille opties), anders de eigen uninstaller van het programma.
func programUninstallArgs(p Program) ([]string, []string) {
	wg := []string{"winget", "uninstall", "--id", `ARP\` + p.Hive + `\` + p.KeyName, "--exact", "--silent", "--accept-source-agreements", "--disable-interactivity"}
	var own []string
	cmd := p.Quiet
	if cmd == "" {
		cmd = p.Uninstall
	}
	if m := regexp.MustCompile(`(?i)msiexec(\.exe)?\s+/[ix]\s*(\{[0-9a-f-]+\})`).FindStringSubmatch(cmd); m != nil {
		own = []string{"msiexec.exe", "/x", m[2], "/qn", "/norestart"}
	} else if cmd != "" {
		own = []string{"cmd.exe", "/c", cmd}
	}
	return wg, own
}

// UninstallPrograms verwijdert programma's (ook die niet in Dusk staan) en daarna de gekozen restjes.
func (a *App) UninstallPrograms(keys []string, leftovers []string) string {
	progs := a.Programs()
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		ok, fail := 0, 0
		for _, k := range keys {
			var p *Program
			for i := range progs {
				if progs[i].Key == k {
					p = &progs[i]
				}
			}
			if p == nil || p.Protected {
				continue
			}
			a.log("» " + a.tr("uninstall") + ": " + p.Name)
			a.emit("prog", Progress{k, "bezig", -1, "removing", ""})
			wg, own := programUninstallArgs(*p)
			code, err := streamCmd(wg, func(l string) { a.log("  " + l) })
			if (err != nil || !wingetOK(code)) && own != nil {
				// winget kon het niet: de eigen uninstaller (kan een venster laten zien)
				if isAdmin() {
					code, err = streamCmd(own, func(l string) { a.log("  " + l) })
				} else {
					code, err = runElevated(own, func(l string) { a.log("  " + l) })
				}
				if err == nil && (code == 0 || code == 3010 || code == 1605) {
					code = 0
				}
			}
			if err != nil || (!wingetOK(code) && code != 0) {
				fail++
				a.log("  ✗ " + a.tr("failed"))
				a.emit("prog", Progress{k, "fout", -1, "failed", ""})
				continue
			}
			ok++
			a.emit("prog", Progress{k, "klaar", 100, "removed", ""})
		}
		var freed int64
		var removed []string
		removeAll(a, leftovers, &freed, &removed)
		infos, _ := a.ScanApps()
		a.emit("apps", infos)
		a.emit("done", map[string]any{"kind": "uninstall", "ok": ok, "fail": fail, "freed": freed, "removed": removed})
	}()
	return ""
}
