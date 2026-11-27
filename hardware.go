package main

import (
	"crypto/sha256"
	hexenc "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GPU struct {
	Name, Ver, Date, Pnp, Vendor string
	Virtual                      bool
}

type Disk struct {
	Name, Type, Health string
	Size               int
}

type Problem struct {
	Name, Class string
	Code        int
}

type Drv struct {
	Class, Name, Maker, Ver, Date string
}

type RAMStick struct {
	Maker, Part string
	GB, Speed   int // Speed = snelheid waarop hij nu draait (MT/s)
	Rated       int // snelheid uit het typenummer (0 = onbekend)
}

type HardwareInfo struct {
	CPU, Board, BoardMaker, BoardProduct string
	BIOS, BIOSDate                       string
	Chipset, ChipsetVer                  string // bv. "AMD Chipset Software", "8.09.28.1124"
	GPUs                                 []GPU
	Disks                                []Disk
	Problems                             []Problem
	Drivers                              []Drv
	RAM                                  []RAMStick
	Devices                              []Device
	OEM                                  OEMInfo
	SysMaker, SysModel, Serial           string
	Chassis                              []int
	OSBuild                              int
	OSRelease                            string
	Monitors                             []Monitor
	Batteries                            []Battery
	Software                             []struct{ Name, Ver string }
}

type DriverUpdate struct {
	ID, Title, Class, Date, Maker, HWID, Model string
}

var errCancelled = errors.New("cancelled")

// psArgs zet het script in een gewoon .ps1-bestand in de tijdelijke map en geeft de opdracht
// om dat te draaien. (Eerder -EncodedCommand, maar dat ziet Windows Defender als verdacht.)
func psArgs(script string) []string {
	sum := sha256.Sum256([]byte(script))
	path := filepath.Join(os.TempDir(), "dusk-"+hexenc.EncodeToString(sum[:6])+".ps1")
	if _, err := os.Stat(path); err != nil {
		// UTF-8 met BOM, zodat PowerShell 5 tekens als é goed leest
		os.WriteFile(path, append([]byte{0xEF, 0xBB, 0xBF}, []byte(script)...), 0o644)
	}
	return []string{"powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", path}
}

const hwScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$g=@(Get-CimInstance Win32_VideoController | % {[pscustomobject]@{Name=$_.Name;Ver=$_.DriverVersion;Date=$(if($_.DriverDate){$_.DriverDate.ToString('yyyy-MM-dd')}else{''});Pnp=$_.PNPDeviceID}})
$b=Get-CimInstance Win32_BaseBoard
$p=@(Get-CimInstance Win32_PnPEntity | ? {$_.ConfigManagerErrorCode -ne 0 -and $_.ConfigManagerErrorCode -ne 22 -and $_.ConfigManagerErrorCode -ne 45} | % {[pscustomobject]@{Name=$(if($_.Name){$_.Name}else{$_.DeviceID});Code=[int]$_.ConfigManagerErrorCode;Class=[string]$_.PNPClass}})
$cpu=(Get-CimInstance Win32_Processor | select -First 1).Name
$d=@(Get-PhysicalDisk | % {[pscustomobject]@{Name=$_.FriendlyName;Size=[int]($_.Size/1GB);Type=[string]$_.MediaType;Health=[string]$_.HealthStatus}})
$ram=@(Get-CimInstance Win32_PhysicalMemory | % {[pscustomobject]@{Maker=[string]$_.Manufacturer;Part=([string]$_.PartNumber).Trim();GB=[int]($_.Capacity/1GB);Speed=[int]$(if($_.ConfiguredClockSpeed){$_.ConfiguredClockSpeed}else{$_.Speed})}})
$bios=Get-CimInstance Win32_BIOS
$cs=Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*','HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*' | ? { $_.DisplayName -match '^(AMD Chipset Software|Intel\(R\) Chipset Device Software|Intel.*Chipset)' } | select -First 1
$cs2=Get-CimInstance Win32_ComputerSystem; $enc=Get-CimInstance Win32_SystemEnclosure; $osr=Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion'
[pscustomobject]@{SysMaker=[string]$cs2.Manufacturer;SysModel=[string]$cs2.Model;Serial=[string]$bios.SerialNumber;Chassis=@($enc.ChassisTypes | % {[int]$_});OSBuild=[int]$osr.CurrentBuild;OSRelease=[string]$osr.DisplayVersion;CPU=$cpu;Board=(($b.Manufacturer+' '+$b.Product).Trim());BoardMaker=[string]$b.Manufacturer;BoardProduct=[string]$b.Product;BIOS=[string]$bios.SMBIOSBIOSVersion;BIOSDate=$(if($bios.ReleaseDate){$bios.ReleaseDate.ToString('yyyy-MM-dd')}else{''});Chipset=[string]$cs.DisplayName;ChipsetVer=[string]$cs.DisplayVersion;GPUs=$g;Disks=$d;Problems=$p;RAM=$ram} | ConvertTo-Json -Depth 4 -Compress`

func hardwareInfo() (HardwareInfo, error) {
	var h HardwareInfo
	out, err := runOut(psArgs(hwScript))
	if err != nil && out == "" {
		return h, err
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &h); err != nil {
		return h, err
	}
	if sc, err := scanDevices(); err == nil {
		h.Devices, h.Monitors, h.Batteries, h.Software = sc.Devices, sc.Monitors, sc.Batteries, sc.Software
	}
	for i := range h.RAM {
		h.RAM[i].Rated = ratedSpeed(h.RAM[i].Part)
	}
	h.OEM = detectOEM(h.SysMaker, h.SysModel, h.Serial, h.BoardMaker, h.BoardProduct, h.Chassis, len(h.Batteries) > 0)
	for i := range h.GPUs {
		g := &h.GPUs[i]
		p := strings.ToUpper(g.Pnp)
		switch {
		case strings.Contains(p, "VEN_10DE"):
			g.Vendor = "nvidia"
		case strings.Contains(p, "VEN_1002"):
			g.Vendor = "amd"
		case strings.Contains(p, "VEN_8086"):
			g.Vendor = "intel"
		}
		g.Virtual = !strings.HasPrefix(p, "PCI\\")
	}
	return h, nil
}

const driverScanScript = `[Console]::OutputEncoding=[Text.Encoding]::UTF8
$s=New-Object -ComObject Microsoft.Update.Session
$r=$s.CreateUpdateSearcher().Search("IsInstalled=0 and Type='Driver' and IsHidden=0")
$l=@(foreach($u in $r.Updates){[pscustomobject]@{ID=$u.Identity.UpdateID;Title=$u.Title;Class=[string]$u.DriverClass;Date=$(if($u.DriverVerDate){$u.DriverVerDate.ToString('yyyy-MM-dd')}else{''});Maker=[string]$u.DriverManufacturer}})
ConvertTo-Json -InputObject $l -Compress`

func scanDrivers() ([]DriverUpdate, error) {
	out, err := runOut(psArgs(driverScanScript))
	if err != nil && strings.TrimSpace(out) == "" {
		return nil, err
	}
	list := []DriverUpdate{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &list); err != nil {
		return nil, err
	}
	return list, nil
}

// driverInstallArgs: download en installeer de drivers met deze UpdateID's.
// Schrijft "RESULT <id> | <code>" (2 = gelukt, 3 = gelukt met fouten).
func driverInstallArgs(ids []string) []string {
	q := make([]string, len(ids))
	for i, id := range ids {
		q[i] = "'" + strings.ReplaceAll(id, "'", "") + "'"
	}
	script := `$ids=@(` + strings.Join(q, ",") + `)
$s=New-Object -ComObject Microsoft.Update.Session
"Searching..."
$r=$s.CreateUpdateSearcher().Search("IsInstalled=0 and Type='Driver' and IsHidden=0")
$col=New-Object -ComObject Microsoft.Update.UpdateColl
foreach($u in $r.Updates){ if($ids -contains $u.Identity.UpdateID){ $u.AcceptEula(); [void]$col.Add($u); "Found: " + $u.Title } }
if($col.Count -eq 0){ "Nothing found"; exit 2 }
"Downloading..."
$d=$s.CreateUpdateDownloader(); $d.Updates=$col; [void]$d.Download()
"Installing..."
$i=$s.CreateUpdateInstaller(); $i.Updates=$col; $ir=$i.Install()
for($k=0;$k -lt $col.Count;$k++){ "RESULT " + $col.Item($k).Identity.UpdateID + " | " + $ir.GetUpdateResult($k).ResultCode }
if($ir.RebootRequired){ "Restart needed" }`
	return psArgs(script)
}

var (
	reRated    = regexp.MustCompile(`(?:^|[^0-9])(2133|2400|2666|2933|3000|3200|3333|3466|3600|3733|3800|4000|4133|4266|4400|4600|4800|5200|5600|6000|6200|6400|6600|6800|7200|7600|8000)(?:[^0-9]|$)`)
	reKingston = regexp.MustCompile(`^KF[45](\d{2})C`)
)

// ratedSpeed haalt de snelheid waarvoor het geheugen gemaakt is uit het typenummer
// (Corsair CMK16GX4M2B3600C18 = 3600, Kingston KF436C17 = 3600). 0 = onbekend.
func ratedSpeed(part string) int {
	part = strings.ToUpper(strings.TrimSpace(part))
	if m := reKingston.FindStringSubmatch(part); m != nil {
		n, _ := strconv.Atoi(m[1])
		return n * 100
	}
	if m := reRated.FindStringSubmatch(part); m != nil {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}

type errExit int

func (e errExit) Error() string { return "exit code " + strconv.Itoa(int(e)) }

// download haalt een bestand op naar een tijdelijke map en meldt de voortgang in procenten.
func download(url string, progress func(int)) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Dusk")
	cl := &http.Client{Timeout: 60 * time.Minute}
	resp, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}
	dir, _ := os.MkdirTemp("", "dusk-dl-")
	name := filepath.Base(resp.Request.URL.Path)
	if name == "" || name == "/" {
		name = "setup.exe"
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	total, done, last := resp.ContentLength, int64(0), -1
	buf := make([]byte, 256<<10)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := f.Write(buf[:n]); err != nil {
				return "", err
			}
			done += int64(n)
			if total > 0 {
				if p := int(done * 100 / total); p != last {
					last = p
					progress(p)
				}
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return "", rerr
		}
	}
	return path, nil
}
