package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Blocker is een programma dat bestanden van een app open heeft, waardoor
// de installer weigert (bijv. Discord met de virtuele camera van OBS).
type Blocker struct {
	Name string // procesnaam, bv. "Discord"
	Desc string // nette naam, bv. "Discord"
	PIDs []int
}

// processen die we nooit voorstellen om te sluiten
var neverClose = map[string]bool{
	"explorer": true, "dwm": true, "sihost": true, "svchost": true, "ctfmon": true, "runtimebroker": true,
	"searchhost": true, "startmenuexperiencehost": true, "shellexperiencehost": true, "textinputhost": true,
	"msedgewebview2": true, "dusk": true, "applicationframehost": true, "lockapp": true, "widgets": true,
}

const inUseScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$apps=ConvertFrom-Json '__APPS__'
$roots='HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall','HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall','HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
$entries=@(foreach($r in $roots){ Get-ChildItem $r | % { Get-ItemProperty $_.PSPath } })
$dirs=@{}
foreach($a in $apps){
  $list=@()
  foreach($e in $entries){ foreach($n in $a.Names){ if($e.DisplayName -and $e.DisplayName -like "$n*"){
    foreach($v in $e.InstallLocation,$e.UninstallString,$e.DisplayIcon){ if($v){ $v=($v -replace '"','') -replace ',\s*-?\d+$',''; if($v -match '\.(exe|ico|dll)(\s|$)'){ $v=Split-Path ($v -replace '\.(exe|ico|dll)\s.*$','.$1') }; if($v -and $v.Length -gt 3 -and $v -notmatch '^[A-Za-z]:\\(Windows|Program Files( \(x86\))?|ProgramData)\\?$'){ $list+=$v.TrimEnd('\') } } } } } }
  $dirs[$a.ID]=@($list | Sort-Object -Unique)
}
$me=(Get-Process -Id $PID).SessionId
$out=@()
foreach($pr in Get-Process | ? { $_.SessionId -eq $me }){
  $mods=$null
  try { $mods=@($pr.Modules | % { $_.FileName }) } catch {}
  if(-not $mods){ continue }
  foreach($a in $apps){ foreach($d in $dirs[$a.ID]){
    if(@($mods | ? { $_ -like "$d\*" }).Count){ $out+=[pscustomobject]@{App=$a.ID;Name=$pr.ProcessName;Id=$pr.Id;Desc=[string]$pr.MainModule.FileVersionInfo.FileDescription}; break }
  } }
}
ConvertTo-Json -InputObject @($out) -Compress`

// InUse zoekt per app welke programma's zijn bestanden open hebben.
func (a *App) InUse(ids []string) (map[string][]Blocker, error) {
	type appQ struct {
		ID    string
		Names []string
	}
	var q []appQ
	for _, id := range ids {
		for _, c := range catalog {
			if c.ID == id {
				q = append(q, appQ{c.ID, append([]string{c.Name}, arpNames[c.ID]...)})
			}
		}
	}
	res := map[string][]Blocker{}
	if len(q) == 0 {
		return res, nil
	}
	js, _ := json.Marshal(q)
	out, err := runOut(psArgs(strings.Replace(inUseScript, "__APPS__", strings.ReplaceAll(string(js), "'", "''"), 1)))
	if strings.TrimSpace(out) == "" {
		return res, err
	}
	var hits []struct {
		App, Name, Desc string
		Id              int
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &hits); err != nil {
		return res, err
	}
	self := os.Getpid()
	for _, h := range hits {
		if h.Id == self || neverClose[strings.ToLower(h.Name)] {
			continue
		}
		list := res[h.App]
		found := false
		for i := range list {
			if strings.EqualFold(list[i].Name, h.Name) {
				list[i].PIDs = append(list[i].PIDs, h.Id)
				found = true
			}
		}
		if !found {
			desc := strings.TrimSpace(h.Desc)
			if desc == "" {
				desc = h.Name
			}
			list = append(list, Blocker{h.Name, desc, []int{h.Id}})
		}
		res[h.App] = list
	}
	return res, nil
}

// CloseApps sluit programma's netjes, en na een paar seconden hard als ze blijven hangen.
func (a *App) CloseApps(names []string) error {
	for _, n := range names {
		if neverClose[strings.ToLower(n)] {
			continue
		}
		a.log("» " + a.tr("closing") + ": " + n)
		c := exec.Command("taskkill", "/IM", n+".exe", "/T")
		hideWindow(c)
		c.Run()
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !anyRunning(names) {
			return nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	for _, n := range names {
		if neverClose[strings.ToLower(n)] {
			continue
		}
		c := exec.Command("taskkill", "/F", "/IM", n+".exe", "/T")
		hideWindow(c)
		c.Run()
	}
	time.Sleep(800 * time.Millisecond)
	return nil
}

func anyRunning(names []string) bool {
	for _, n := range names {
		out, _ := runOut([]string{"tasklist", "/FI", "IMAGENAME eq " + n + ".exe", "/NH"})
		if strings.Contains(strings.ToLower(out), strings.ToLower(n)+".exe") {
			return true
		}
	}
	return false
}
