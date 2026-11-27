package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OEMInfo: wie heeft deze pc gemaakt, en waar haal je de drivers voor precies dit model.
type OEMInfo struct {
	Brand      string // hp, dell, lenovo, asus, acer, msi, microsoft, custom, other
	Maker      string
	Model      string
	Serial     string
	SystemID   string // HP: platform-ID (BaseBoard Product)
	Laptop     bool
	SupportURL string
	Tool       string // winget-ID van de officiële updatetool
	ToolName   string
	Deep       bool // Dusk kan zelf de updates van dit merk vergelijken (HP)
}

func detectOEM(maker, model, serial, boardMaker, boardProduct string, chassis []int, hasBattery bool) OEMInfo {
	o := OEMInfo{Maker: strings.TrimSpace(maker), Model: strings.TrimSpace(model), Serial: strings.TrimSpace(serial), SystemID: strings.TrimSpace(boardProduct)}
	for _, c := range chassis { // 8-14, 30-32 = draagbaar
		if (c >= 8 && c <= 14) || (c >= 30 && c <= 32) {
			o.Laptop = true
		}
	}
	if hasBattery {
		o.Laptop = true
	}
	m := strings.ToLower(o.Maker)
	bad := o.Serial == "" || strings.Contains(strings.ToLower(o.Serial), "default") || strings.Contains(strings.ToLower(o.Serial), "o.e.m")
	switch {
	case strings.Contains(m, "hewlett") || m == "hp" || strings.HasPrefix(m, "hp "):
		o.Brand, o.Tool, o.ToolName, o.Deep = "hp", "HP.SupportAssistant", "HP Support Assistant", true
		o.SupportURL = "https://support.hp.com/drivers"
	case strings.Contains(m, "dell"):
		o.Brand, o.Tool, o.ToolName = "dell", "Dell.CommandUpdate.Universal", "Dell Command | Update"
		o.SupportURL = "https://www.dell.com/support/home/drivers"
		if !bad {
			o.SupportURL = "https://www.dell.com/support/home/product-support/servicetag/" + o.Serial + "/drivers"
		}
	case strings.Contains(m, "lenovo"):
		o.Brand, o.Tool, o.ToolName = "lenovo", "9WZDNCRFJ4MV", "Lenovo Vantage"
		o.SupportURL = "https://pcsupport.lenovo.com/"
		if !bad {
			o.SupportURL = "https://pcsupport.lenovo.com/search?query=" + o.Serial
		}
	case strings.Contains(m, "asus"):
		o.Brand, o.Tool, o.ToolName = "asus", "9N7R5S6B0ZZH", "MyASUS"
		if regexp.MustCompile(`(?i)rog|tuf`).MatchString(o.Model) {
			o.Tool, o.ToolName = "Asus.ArmouryCrate", "Armoury Crate"
		}
		o.SupportURL = "https://www.asus.com/support/"
	case strings.Contains(m, "acer"):
		o.Brand, o.Tool, o.ToolName = "acer", "9P8BB54NQNQ4", "Acer Care Center S"
		o.SupportURL = "https://www.acer.com/support/drivers-and-manuals"
		if !bad {
			o.SupportURL += "?search=" + o.Serial
		}
	case strings.Contains(m, "microsoft") && strings.Contains(strings.ToLower(o.Model), "surface"):
		o.Brand = "microsoft"
		o.SupportURL = "https://support.microsoft.com/surface/download-drivers-and-firmware-for-surface"
	case strings.Contains(m, "micro-star") || strings.Contains(m, "msi"):
		if o.Laptop {
			o.Brand, o.SupportURL = "msi", "https://www.msi.com/support"
		} else {
			o.Brand = "custom"
		}
	default:
		// zelfbouw-pc's zeggen vaak "To Be Filled By O.E.M." of de naam van het moederbord
		if strings.Contains(m, "o.e.m") || strings.Contains(m, "system manufacturer") || m == "" ||
			strings.Contains(strings.ToLower(boardMaker), strings.Fields(m + " x")[0]) {
			o.Brand = "custom"
		} else {
			o.Brand = "other"
		}
	}
	if o.Brand == "custom" {
		o.Model = strings.TrimSpace(boardMaker + " " + boardProduct)
	}
	return o
}

// ---------- HP: dezelfde lijst die HP Image Assistant gebruikt ----------

// OEMUpdate is een pakket van de maker (driver, BIOS of firmware) dat nieuwer is.
type OEMUpdate struct {
	ID, Name, Category, Version, Installed, Date, URL string
	Size                                              int64
	Device                                            string
}

type hpRef struct {
	Devices []struct {
		DeviceID  string `xml:"DeviceId"`
		HWID      string `xml:"HardwareId"`
		Name      string `xml:"FriendlyName"`
		DriverVer string `xml:"DriverVersion"`
		Sol       []struct {
			IDRef string `xml:"IdRef,attr"`
		} `xml:"Solutions>UpdateInfo"`
	} `xml:"Devices>Device"`
	BIOS []struct {
		IDRef string      `xml:"IdRef,attr"`
		Fam   []hpBiosFam `xml:"IsApplicable>BIOS"`
	} `xml:"SystemInfo>System>Solutions>UpdateInfo"`
	Sols []struct {
		ID       string `xml:"Id"`
		Name     string `xml:"Name"`
		Category string `xml:"Category"`
		Version  string `xml:"Version"`
		URL      string `xml:"Url"`
		Date     string `xml:"DateReleased"`
		Size     int64  `xml:"Size"`
	} `xml:"Solutions>UpdateInfo"`
}

type hpBiosFam struct {
	Family string `xml:"Family,attr"`
}

// hpRefURL: https://hpia.hpcloud.hp.com/ref/896d/896d_64_11.0.26h2.cab
func hpRefURL(sysID string, win11 bool, release string) string {
	id := strings.ToLower(sysID)
	v := "10.0"
	if win11 {
		v = "11.0"
	}
	return fmt.Sprintf("https://hpia.hpcloud.hp.com/ref/%s/%s_64_%s.%s.cab", id, id, v, strings.ToLower(release))
}

// hpUpdates haalt de referentielijst voor dit HP-model op en vergelijkt hem met de hardware.
func hpUpdates(o OEMInfo, devs []Device, bios string, win11 bool, release string) ([]OEMUpdate, error) {
	if !regexp.MustCompile(`^[0-9A-Fa-f]{4}$`).MatchString(o.SystemID) {
		return nil, fmt.Errorf("unknown HP platform ID %q", o.SystemID)
	}
	// probeer de huidige Windows-versie, en anders oudere (HP loopt soms achter)
	rels := []string{release, "26h2", "25h2", "24h2", "23h2", "22h2"}
	var data []byte
	for _, r := range rels {
		if r == "" {
			continue
		}
		b, code, err := get(hpRefURL(o.SystemID, win11, r))
		if err == nil && code == 200 {
			data = b
			break
		}
	}
	if data == nil {
		return nil, fmt.Errorf("no HP reference file for platform %s", o.SystemID)
	}
	xmlData, err := extractCab(data)
	if err != nil {
		return nil, err
	}
	return hpMatch(xmlData, devs, bios)
}

// hpMatch vergelijkt de HP-referentielijst met de apparaten en de BIOS van deze pc.
func hpMatch(xmlData []byte, devs []Device, bios string) ([]OEMUpdate, error) {
	var ref hpRef
	if err := xml.Unmarshal(xmlData, &ref); err != nil {
		return nil, err
	}
	sols := map[string]int{}
	for i, s := range ref.Sols {
		sols[s.ID] = i
	}
	var out []OEMUpdate
	added := map[string]bool{}
	add := func(id, installed, dev string) {
		i, ok := sols[id]
		if !ok || added[id] {
			return
		}
		s := ref.Sols[i]
		added[id] = true
		url := s.URL
		if !strings.HasPrefix(url, "http") {
			url = "https://" + url
		}
		out = append(out, OEMUpdate{s.ID, s.Name, s.Category, s.Version, installed, s.Date, url, s.Size, dev})
	}
	// drivers: apparaat van HP-lijst dat op deze pc zit, met een oudere driver
	for _, rd := range ref.Devices {
		pat := strings.ToUpper(rd.DeviceID)
		if pat == "" || len(rd.Sol) == 0 {
			continue
		}
		for _, d := range devs {
			hit := strings.HasPrefix(strings.ToUpper(d.ID), pat)
			for _, h := range d.HWIDs {
				if strings.EqualFold(h, pat) || strings.HasPrefix(strings.ToUpper(h), pat+"&") {
					hit = true
				}
			}
			if hit && d.Ver != "" && newer(d.Ver, rd.DriverVer) {
				add(rd.Sol[0].IDRef, d.Ver, d.Name)
			}
		}
	}
	// BIOS: "U70 Ver. 01.20.00" -> familie U70, versie 01.20.00
	if m := regexp.MustCompile(`^(\w+)\s+Ver\.\s*([\d.]+)`).FindStringSubmatch(strings.TrimSpace(bios)); m != nil {
		for _, b := range ref.BIOS {
			for _, f := range b.Fam {
				if strings.EqualFold(f.Family, m[1]) {
					if i, ok := sols[b.IDRef]; ok && newer(m[2], ref.Sols[i].Version) {
						add(b.IDRef, m[2], "BIOS")
					}
				}
			}
		}
	}
	return out, nil
}

// extractCab pakt een .cab met één bestand uit (Windows heeft tar.exe/expand.exe).
func extractCab(data []byte) ([]byte, error) {
	dir, err := os.MkdirTemp("", "dusk-cab-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	cab := filepath.Join(dir, "x.cab")
	if err := os.WriteFile(cab, data, 0o644); err != nil {
		return nil, err
	}
	out := filepath.Join(dir, "out")
	os.Mkdir(out, 0o755)
	if _, err := runOut([]string{"tar.exe", "-xf", cab, "-C", out}); err != nil {
		runOut([]string{"expand.exe", cab, "-F:*", out + `\`})
	}
	files, _ := os.ReadDir(out)
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".xml") {
			return os.ReadFile(filepath.Join(out, f.Name()))
		}
	}
	return nil, fmt.Errorf("cab is empty")
}
