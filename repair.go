package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Reparaties: één PowerShell-script per reparatie, als admin. Regels die met "STEP " beginnen
// zijn de stappen (de app vertaalt ze), "PCT n" is de voortgang.
var repairScripts = map[string]string{
	// Windows-bestanden: eerst DISM (haalt goede bestanden bij Windows Update), dan SFC
	"windows": `$ErrorActionPreference='Continue'
"STEP dism"
DISM.exe /Online /Cleanup-Image /RestoreHealth 2>&1 | % { $_ }
"STEP sfc"
sfc.exe /scannow 2>&1 | % { ($_ -replace "\x00","") } | ? { $_.Trim() }
"STEP done"`,
	// Windows Update: diensten stoppen, de caches opzij zetten, diensten weer starten
	"update": `$ErrorActionPreference='Continue'
"STEP stop"
foreach($s in 'wuauserv','bits','cryptsvc','msiserver','usosvc'){ Stop-Service $s -Force -ErrorAction SilentlyContinue }
"STEP cache"
$d=Get-Date -Format yyyyMMdd-HHmmss
if(Test-Path "$env:SystemRoot\SoftwareDistribution"){ Rename-Item "$env:SystemRoot\SoftwareDistribution" "SoftwareDistribution.dusk-$d" -ErrorAction SilentlyContinue }
if(Test-Path "$env:SystemRoot\System32\catroot2"){ Rename-Item "$env:SystemRoot\System32\catroot2" "catroot2.dusk-$d" -ErrorAction SilentlyContinue }
"STEP start"
foreach($s in 'cryptsvc','bits','wuauserv','msiserver'){ Start-Service $s -ErrorAction SilentlyContinue }
"STEP done"`,
	// Internet: DNS-cache, Winsock en TCP/IP terugzetten, nieuw IP-adres vragen
	"network": `$ErrorActionPreference='Continue'
"STEP dns"
ipconfig /flushdns
"STEP winsock"
netsh winsock reset
netsh int ip reset
"STEP ip"
ipconfig /release | Out-Null
ipconfig /renew | Out-Null
"STEP done"`,
	// Microsoft Store: cache leegmaken (wsreset opent daarna de Store)
	"store": `"STEP store"
Start-Process wsreset.exe -Wait
"STEP done"`,
	// Geluid: de audiodiensten opnieuw starten (geen geluid, crackle, apparaat weg)
	"audio": `"STEP audio"
Restart-Service AudioEndpointBuilder -Force
Restart-Service Audiosrv -Force
"STEP done"`,
	// Klok: tijd opnieuw ophalen bij de tijdserver
	"time": `"STEP time"
Set-Service w32time -StartupType Manual
Start-Service w32time -ErrorAction SilentlyContinue
w32tm /resync /force
"STEP done"`,
	// Printer: wachtrij leegmaken (vastgelopen printopdrachten)
	"print": `"STEP print"
Stop-Service Spooler -Force
Remove-Item "$env:SystemRoot\System32\spool\PRINTERS\*" -Force -ErrorAction SilentlyContinue
Start-Service Spooler
"STEP done"`,
	// Zoeken: de zoekindex van Windows opnieuw opbouwen
	"search": `"STEP search"
Stop-Service WSearch -Force
$d="$env:ProgramData\Microsoft\Search\Data\Applications\Windows"
Remove-Item "$d\Windows.edb","$d\Windows.db","$d\Windows-gather.db","$d\Windows-usn.db" -Force -ErrorAction SilentlyContinue
Start-Service WSearch
"STEP done"`,
	// Bluetooth: de bluetooth-dienst opnieuw starten
	"bluetooth": `"STEP bluetooth"
Restart-Service bthserv -Force
"STEP done"`,
	// Schijf: C: controleren op fouten (alleen lezen, kan tijdens gebruik)
	"disk": `"STEP disk"
chkdsk.exe $env:SystemDrive /scan 2>&1 | % { $_ }
"STEP done"`,
	// Iconen en miniaturen: cache weg en Verkenner opnieuw (de taakbalk is even weg)
	"icons": `"STEP icons"
Stop-Process -Name explorer -Force
Start-Sleep -Milliseconds 800
Remove-Item "$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*.db","$env:LOCALAPPDATA\Microsoft\Windows\Explorer\thumbcache_*.db","$env:LOCALAPPDATA\IconCache.db" -Force -ErrorAction SilentlyContinue
Start-Process explorer.exe
"STEP done"`,
	// Verkenner opnieuw starten (taakbalk of bureaublad hangt)
	"explorer": `"STEP explorer"
Stop-Process -Name explorer -Force
Start-Sleep -Milliseconds 800
Start-Process explorer.exe
"STEP done"`,
	// Alleen de DNS-cache (snel, geen herstart nodig)
	"dns": `"STEP dns"
ipconfig /flushdns
"STEP done"`,
}

// reparaties zonder admin
var repairUser = map[string]bool{"dns": true, "icons": true, "explorer": true}

var reRepairPct = regexp.MustCompile(`(\d{1,3}(?:[.,]\d)?)\s?%`)

// Repair voert een reparatie uit. Antwoord "busy" als er al iets loopt.
func (a *App) Repair(id string) string {
	script, ok := repairScripts[id]
	if !ok {
		return "ERR unknown"
	}
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		a.log("» " + a.tr("repair") + ": " + id)
		a.emit("repair", Progress{id, "bezig", -1, "starting", ""})
		lastPct := -1
		onLine := func(l string) {
			l = strings.ReplaceAll(l, "\x00", "")
			if step, ok := strings.CutPrefix(strings.TrimSpace(l), "STEP "); ok {
				lastPct = -1
				a.emit("repair", Progress{id, "bezig", -1, "r_" + step, ""})
				return
			}
			if m := reRepairPct.FindStringSubmatch(l); m != nil {
				f, _ := strconv.ParseFloat(strings.ReplaceAll(m[1], ",", "."), 64)
				if p := int(f); p != lastPct && p <= 100 {
					lastPct = p
					a.emit("repair", Progress{id, "bezig", p, "", ""})
				}
				return
			}
			if strings.TrimSpace(l) != "" {
				a.log("  " + strings.TrimSpace(l))
			}
		}
		args := psArgs(script)
		var code int
		var err error
		if repairUser[id] || isAdmin() {
			code, err = streamCmd(args, onLine)
		} else {
			code, err = runElevated(args, onLine)
		}
		if err == nil && code != 0 && id != "windows" && id != "disk" && id != "time" { // sfc/dism geven ook codes bij "gerepareerd"
			err = errExit(code)
		}
		if err != nil {
			msg := err.Error()
			if err == errCancelled {
				msg = "cancelled"
			}
			a.log("  ✗ " + msg)
			a.emit("repair", Progress{id, "fout", -1, "failed", msg})
			a.emit("done", map[string]any{"kind": "repair", "ok": 0, "fail": 1})
			return
		}
		a.log("  ✓ " + a.tr("done"))
		a.emit("repair", Progress{id, "klaar", 100, "done", ""})
		a.emit("done", map[string]any{"kind": "repair", "ok": 1, "fail": 0})
	}()
	return ""
}
