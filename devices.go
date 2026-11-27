package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Device is één stuk hardware met zijn driver.
type Device struct {
	Cat, Name, Class, Maker, Ver, Date, ID string
	HWIDs                                  []string
	Problem                                int  // 0 = werkt, 28 = geen driver, ...
	Generic                                bool // standaard driver van Windows (werkt, maar niet van de maker)
	Basic                                  bool // videokaart zonder echte driver (Microsoft Basic Display Adapter)
}

type Monitor struct{ Name, Maker string }

type Battery struct {
	Name         string
	Design, Full int // mWh
	Charge       int // %
}

const devScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$ent=@{}; $prob=@()
foreach($e in Get-CimInstance Win32_PnPEntity){ $ent[$e.DeviceID]=$e; if($e.ConfigManagerErrorCode -and $e.ConfigManagerErrorCode -notin 22,45 -and $e.DeviceID -notmatch '^(ROOT|SWD)\\'){ $prob+=$e } }
$out=@()
foreach($d in Get-CimInstance Win32_PnPSignedDriver){
  if(-not $d.DeviceName -or -not $d.DeviceID -or $d.DeviceID -match '^(ROOT|SWD|HTREE|UMB|STORAGE\\VOLUME)\\'){ continue }
  $e=$ent[$d.DeviceID]
  $out+=[pscustomobject]@{Name=$d.DeviceName.Trim();Class=[string]$d.DeviceClass;Maker=[string]$d.DriverProviderName;Ver=[string]$d.DriverVersion;Date=$(if($d.DriverDate){$d.DriverDate.ToString('yyyy-MM-dd')}else{''});ID=$d.DeviceID;HWIDs=@($e.HardwareID);Problem=[int]$e.ConfigManagerErrorCode}
}
foreach($e in $prob){ if(-not ($out | ? { $_.ID -eq $e.DeviceID })){ $out+=[pscustomobject]@{Name=$(if($e.Name){$e.Name.Trim()}else{$e.DeviceID});Class=[string]$e.PNPClass;Maker='';Ver='';Date='';ID=$e.DeviceID;HWIDs=@($e.HardwareID);Problem=[int]$e.ConfigManagerErrorCode} } }
$mon=@(Get-CimInstance -Namespace root\wmi WmiMonitorID | % { [pscustomobject]@{Name=(($_.UserFriendlyName | ? {$_ -gt 0} | % {[char]$_}) -join '');Maker=(($_.ManufacturerName | ? {$_ -gt 0} | % {[char]$_}) -join '')} })
$bat=@()
$full=@(Get-CimInstance -Namespace root\wmi BatteryFullChargedCapacity); $des=@(Get-CimInstance -Namespace root\wmi BatteryStaticData)
$i=0; foreach($b in Get-CimInstance Win32_Battery){ $bat+=[pscustomobject]@{Name=[string]$b.Name;Charge=[int]$b.EstimatedChargeRemaining;Full=[int]$full[$i].FullChargedCapacity;Design=[int]$des[$i].DesignedCapacity}; $i++ }
$sw=@(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*','HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*' | ? { $_.DisplayName -match '^(AMD Software|AMD Chipset Software|NVIDIA Graphics Driver|NVIDIA App|Intel.*(Graphics|Chipset))' } | % {[pscustomobject]@{Name=$_.DisplayName;Ver=[string]$_.DisplayVersion}})
[pscustomobject]@{Devices=$out;Monitors=$mon;Batteries=$bat;Software=$sw} | ConvertTo-Json -Depth 5 -Compress`

type devScan struct {
	Devices   []Device
	Monitors  []Monitor
	Batteries []Battery
	Software  []struct{ Name, Ver string }
}

// welke apparaat-klassen we laten zien, en onder welke groep
var classCat = map[string]string{
	"DISPLAY": "gpu", "MEDIA": "audio", "NET": "network", "BLUETOOTH": "bluetooth", "SYSTEM": "chipset",
	"HDC": "storage", "SCSIADAPTER": "storage", "USB": "usb", "HIDCLASS": "input", "KEYBOARD": "input",
	"MOUSE": "input", "CAMERA": "camera", "IMAGE": "camera", "BIOMETRIC": "other", "SENSOR": "other",
}

var reSkipDev = regexp.MustCompile(`(?i)virtual|wan miniport|microsoft kernel debug|teredo|6to4|isatap|vpn|tap-windows|hyper-v|remote desktop|root hub|generic usb hub|composite device|hid-compliant (consumer|system|vendor)|usb input device|hid keyboard device|hid-compliant mouse|standard ps/2`)

func scanDevices() (devScan, error) {
	var sc devScan
	out, err := runOut(psArgs(devScript))
	if strings.TrimSpace(out) == "" {
		return sc, err
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &sc); err != nil {
		return sc, err
	}
	seen := map[string]bool{}
	var list []Device
	for _, d := range sc.Devices {
		cat, ok := classCat[strings.ToUpper(d.Class)]
		ms := strings.HasPrefix(strings.ToLower(d.Maker), "microsoft")
		if d.Problem != 0 {
			if !ok {
				cat = "other"
			}
		} else {
			if !ok || reSkipDev.MatchString(d.Name) {
				continue
			}
			// standaard Windows-drivers alleen tonen bij echte hardware (PCI-kaarten,
			// USB-headsets en webcams), niet de vele Bluetooth-protocollen e.d.
			if ms {
				id := strings.ToUpper(d.ID)
				switch {
				case cat == "gpu" && strings.Contains(strings.ToLower(d.Name), "basic"):
					d.Basic = true
				// bruggen, root-poorten, USB-, SATA- en NVMe-controllers werken altijd met de
				// Windows-driver; alleen netwerk, geluid en camera met een Windows-driver tonen
				case (strings.HasPrefix(id, `PCI\`) && (cat == "network" || cat == "audio")),
					strings.HasPrefix(id, `USB\`) && (cat == "audio" || cat == "camera"):
					d.Generic = true
				default:
					continue
				}
			}
		}
		if cat == "network" && regexp.MustCompile(`(?i)wi-?fi|wireless|wlan|802\.11`).MatchString(d.Name) {
			cat = "wifi"
		}
		key := d.Name + "|" + d.Ver
		if seen[key] {
			continue
		}
		seen[key] = true
		d.Cat = cat
		list = append(list, d)
	}
	sc.Devices = list
	return sc, nil
}

// ---------- nieuwste drivers opzoeken ----------

type DriverStatus struct {
	Updates []DriverUpdate    // via Windows Update (met hardware-ID's)
	GPU     map[string]Latest // per videokaart-naam
	OEM     []OEMUpdate       // pakketten van de maker (HP)
	OEMErr  string
	Err     string
}

type Latest struct {
	Version, URL string
	Download     string // directe download (NVIDIA), leeg = niet in Dusk te installeren
	Newer        bool   // er is een nieuwere versie
	Known        bool   // we weten de nieuwste versie echt
}

var httpc = &http.Client{Timeout: 15 * time.Second}

func get(url string) ([]byte, int, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Dusk")
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	return b, resp.StatusCode, nil
}

// nvidiaVersion: Windows-driverversie 32.0.15.6094 -> NVIDIA-versie 560.94
func nvidiaVersion(win string) string {
	p := strings.Split(win, ".")
	if len(p) < 4 {
		return ""
	}
	d := p[2] + p[3]
	if len(d) < 5 {
		return ""
	}
	d = d[len(d)-5:]
	return strings.TrimLeft(d[:3], "0") + "." + d[3:]
}

func nvidiaLatest(gpuName string) (Latest, error) {
	b, _, err := get("https://www.nvidia.com/Download/API/lookupValueSearch.aspx?TypeID=3")
	if err != nil {
		return Latest{}, err
	}
	var lv struct {
		Values []struct {
			Parent string `xml:"ParentID,attr"`
			Name   string `xml:"Name"`
			Value  string `xml:"Value"`
		} `xml:"LookupValues>LookupValue"`
	}
	if err := xml.Unmarshal(b, &lv); err != nil {
		return Latest{}, err
	}
	want := strings.ToLower(strings.Join(strings.Fields(gpuName), " "))
	psid, pfid := "", ""
	for _, v := range lv.Values {
		if strings.ToLower(strings.TrimSpace(v.Name)) == want {
			psid, pfid = v.Parent, v.Value
			break
		}
	}
	if pfid == "" {
		return Latest{}, fmt.Errorf("gpu not in NVIDIA list")
	}
	b, _, err = get("https://gfwsl.geforce.com/services_toolkit/services/com/nvidia/services/AjaxDriverService.php?func=DriverManualLookup&psid=" + psid + "&pfid=" + pfid + "&osID=57&languageCode=1033&isWHQL=1&dch=1&sort1=0&numberOfResults=1")
	if err != nil {
		return Latest{}, err
	}
	var r struct {
		IDS []struct {
			DownloadInfo struct {
				Version, DetailsURL, DownloadURL string
			} `json:"downloadInfo"`
		}
	}
	if err := json.Unmarshal(b, &r); err != nil || len(r.IDS) == 0 {
		return Latest{}, fmt.Errorf("no NVIDIA driver info")
	}
	di := r.IDS[0].DownloadInfo
	dl := di.DownloadURL
	if !strings.HasPrefix(dl, "https://") || !strings.Contains(dl, ".nvidia.com/") {
		dl = ""
	}
	return Latest{Version: di.Version, URL: "https://www.nvidia.com/en-us/software/nvidia-app/", Download: dl, Known: true}, nil
}

// amdLatest: AMD zet bij elke Adrenalin-versie release notes online
// (RN-RAD-WIN-26-11-2). We kijken of er een nieuwere bestaat dan de geïnstalleerde.
func amdLatest(installed string) Latest {
	l := Latest{URL: "https://www.amd.com/en/support/download/drivers.html"}
	p := strings.Split(installed, ".")
	if len(p) < 3 {
		return l
	}
	y, _ := strconv.Atoi(p[0])
	m, _ := strconv.Atoi(p[1])
	n, _ := strconv.Atoi(p[2])
	if y == 0 || m == 0 {
		return l
	}
	exists := func(y, m, n int) bool {
		_, code, err := get(fmt.Sprintf("https://www.amd.com/en/resources/support-articles/release-notes/RN-RAD-WIN-%d-%d-%d.html", y, m, n))
		return err == nil && code == 200
	}
	if !exists(y, m, n) { // we kunnen AMD niet bereiken of de naam klopt niet: dan weten we het niet
		return l
	}
	l.Known, l.Version = true, installed
	for i := 0; i < 24; i++ {
		switch {
		case exists(y, m, n+1):
			n++
		case m < 12 && exists(y, m+1, 1):
			m, n = m+1, 1
		case m < 11 && exists(y, m+2, 1):
			m, n = m+2, 1
		case m == 12 && exists(y+1, 1, 1):
			y, m, n = y+1, 1, 1
		default:
			i = 99
			continue
		}
		l.Version, l.Newer = fmt.Sprintf("%d.%d.%d", y, m, n), true
	}
	return l
}

const driverScanScript2 = `[Console]::OutputEncoding=[Text.Encoding]::UTF8
$s=New-Object -ComObject Microsoft.Update.Session
$r=$s.CreateUpdateSearcher().Search("IsInstalled=0 and Type='Driver' and IsHidden=0")
$l=@(foreach($u in $r.Updates){[pscustomobject]@{ID=$u.Identity.UpdateID;Title=$u.Title;Class=[string]$u.DriverClass;Date=$(if($u.DriverVerDate){$u.DriverVerDate.ToString('yyyy-MM-dd')}else{''});Maker=[string]$u.DriverManufacturer;HWID=[string]$u.DriverHardwareID;Model=[string]$u.DriverModel}})
ConvertTo-Json -InputObject $l -Compress`

// CheckDrivers zoekt voor alle hardware de nieuwste drivers (Windows Update,
// NVIDIA en AMD). Duurt ongeveer 30 seconden.
func (a *App) CheckDrivers(gpus []GPUQuery) DriverStatus {
	st := DriverStatus{GPU: map[string]Latest{}}
	var wg sync.WaitGroup
	var mu sync.Mutex
	a.mu.Lock()
	h := a.hw
	a.mu.Unlock()
	if h != nil && h.OEM.Deep {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ups, err := hpUpdates(h.OEM, h.Devices, h.BIOS, h.OSBuild >= 22000, h.OSRelease)
			mu.Lock()
			st.OEM = ups
			if err != nil {
				st.OEMErr = err.Error()
			}
			mu.Unlock()
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		out, err := runOut(psArgs(driverScanScript2))
		list := []DriverUpdate{}
		if e := json.Unmarshal([]byte(strings.TrimSpace(out)), &list); e != nil && err == nil {
			err = e
		}
		mu.Lock()
		st.Updates = list
		if err != nil && len(list) == 0 {
			st.Err = err.Error()
		}
		mu.Unlock()
	}()
	for _, g := range gpus {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			var l Latest
			switch g.Vendor {
			case "nvidia":
				if x, err := nvidiaLatest(g.Name); err == nil {
					l = x
					l.Newer = newer(nvidiaVersion(g.Ver), l.Version)
				}
			case "amd":
				if g.Software != "" {
					l = amdLatest(g.Software)
				} else {
					l = Latest{URL: "https://www.amd.com/en/support/download/drivers.html"}
				}
			case "intel":
				l = Latest{URL: "https://www.intel.com/content/www/us/en/support/intel-driver-support-assistant.html"}
			}
			mu.Lock()
			st.GPU[g.Name] = l
			mu.Unlock()
		}()
	}
	wg.Wait()
	return st
}

type GPUQuery struct{ Name, Vendor, Ver, Software string }
