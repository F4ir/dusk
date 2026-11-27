package main

import (
	"os"
	"testing"
)

func TestParseTable(t *testing.T) {
	out := "   - \r   \\ \r" + `Name                                    Id                                       Version          Available        Source
-----------------------------------------------------------------------------------------------------------------------------
Logitech G HUB                          Logitech.GHUB                            2026.7.1011216                    winget
PuTTY release 0.85 (64-bit)             XPFNZKSKLBP7RJ                           0.85.0.0                          msstore
Discord                                 ARP\User\X64\Discord                     1.0.9269
Google Chrome                           ARP\Machine\X86\Google Chrome            157.0.8103.80
OBS Studio                              OBSProject.OBSStudio                     32.2.2           33.0.3           winget
Steam Link                              Valve.SteamLink                          1.3                               winget
Steam                                   Valve.Steam                              2.10.91.91                        winget
Mozilla Firefox (x64 en-US)             Mozilla.Firefox                          < 150.0          151.0            winget
`
	st := matchApps(parseTable(out))
	chk := func(id string, inst bool, ver, avail, real, src string, arp bool) {
		s := st[id]
		if s.Installed != inst || s.Version != ver || s.Available != avail || s.RealID != real || s.Source != src || s.ARP != arp {
			t.Errorf("%s: %+v", id, s)
		}
	}
	chk("Logitech.GHUB", true, "2026.7.1011216", "", "Logitech.GHUB", "winget", false)
	chk("PuTTY.PuTTY", true, "0.85.0.0", "", "XPFNZKSKLBP7RJ", "msstore", false)
	chk("Discord.Discord", true, "1.0.9269", "", `ARP\User\X64\Discord`, "", true)
	chk("Google.Chrome", true, "157.0.8103.80", "", `ARP\Machine\X86\Google Chrome`, "", true)
	chk("OBSProject.OBSStudio", true, "32.2.2", "33.0.3", "OBSProject.OBSStudio", "winget", false)
	chk("Valve.Steam", true, "2.10.91.91", "", "Valve.Steam", "winget", false)
	chk("Mozilla.Firefox", true, "< 150.0", "151.0", "Mozilla.Firefox", "winget", false)
	if a := installArgs(CatApp{ID: "PuTTY.PuTTY"}, st["PuTTY.PuTTY"]); a[3] != "XPFNZKSKLBP7RJ" || a[5] != "msstore" {
		t.Errorf("putty args %v", a)
	}
}

func TestRated(t *testing.T) {
	for part, want := range map[string]int{"CMK16GX4M2B3600C18": 3600, "KF436C17BBK2/16": 3600, "F4-3200C16-8GVKB": 3200, "M378A1K43CB2-CTD": 0, "KF560C40-16": 6000} {
		if got := ratedSpeed(part); got != want {
			t.Errorf("%s: %d, wil %d", part, got, want)
		}
	}
}

func TestNvidiaVersion(t *testing.T) {
	if v := nvidiaVersion("32.0.15.6094"); v != "560.94" {
		t.Error(v)
	}
	if v := nvidiaVersion("31.0.15.3623"); v != "536.23" {
		t.Error(v)
	}
}

func TestNvidiaLatestLive(t *testing.T) {
	l, err := nvidiaLatest("NVIDIA GeForce RTX 4070")
	t.Logf("4070: %+v %v", l, err)
	a := amdLatest("26.11.2")
	t.Logf("amd 26.11.2: %+v", a)
	b := amdLatest("26.9.1")
	t.Logf("amd 26.9.1: %+v", b)
}

func TestHPMatch(t *testing.T) {
	data, err := os.ReadFile("/tmp/hpref/896d_64_11.0.26H2.xml")
	if err != nil {
		t.Skip("geen HP-testbestand")
	}
	devs := []Device{
		{Name: "AMD Audio CoProcessor", ID: `ACP\DEVTYPE_0004&DEVREV_0000&VEN_1022&DEV_15E2\3&1`, HWIDs: []string{`ACP\DEVTYPE_0004&DEVREV_0000&VEN_1022&DEV_15E2`}, Ver: "6.0.0.40"},
		{Name: "AMD Audio CoProcessor (nieuw)", ID: `ACP\X`, HWIDs: []string{`ACP\DEVTYPE_0004&DEVREV_0000&VEN_1022&DEV_15E2`}, Ver: "9.0.0.0"},
	}
	ups, err := hpMatch(data, devs, "U70 Ver. 01.01.00")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range ups {
		t.Logf("%+v", u)
	}
	if len(ups) < 2 {
		t.Errorf("verwacht driver + BIOS, kreeg %d", len(ups))
	}
}

func TestFixSquirrel(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir+"/Discord/app-1.0.9269", 0o755)
	os.MkdirAll(dir+"/Discord/app-1.0.9270", 0o755)
	t.Setenv("LOCALAPPDATA", dir)
	st := map[string]AppState{"Discord.Discord": {Installed: true, Version: "1.0.9269", Available: "1.0.9270"}}
	fixSquirrel(st)
	if s := st["Discord.Discord"]; s.Version != "1.0.9270" || s.Available != "" {
		t.Errorf("%+v", s)
	}
	st = map[string]AppState{"Discord.Discord": {Installed: true, Version: "1.0.9269", Available: "1.0.9275"}}
	fixSquirrel(st)
	if s := st["Discord.Discord"]; s.Available != "1.0.9275" {
		t.Errorf("echte update weg: %+v", s)
	}
}

func TestExtras(t *testing.T) {
	st := map[string]AppState{}
	detectExtras(st, []wgRow{{Name: "Python 3.14.7", ID: `ARP\User\X64\pymanager`, Version: "3.14-64"}, {Name: "Java 8 Update 401", ID: "{x}", Version: "8.0.4010"}})
	if s := st["Python.PythonInstallManager"]; !s.Installed || s.Version != "3.14.7" || s.Available != "" {
		t.Errorf("python %+v", s)
	}
	if s := st["EclipseAdoptium.Temurin.21.JRE"]; !s.Installed || s.Version != "Java 8" {
		t.Errorf("java %+v", s)
	}
}

func TestPrettyPkg(t *testing.T) {
	for in, want := range map[string]string{"WinRAR.ShellExtension": "WinRAR", "Microsoft.WindowsTerminal": "Windows Terminal",
		"JAMSoftware.TreeSizeFreeContextMenu": "Tree Size Free", "AdvancedMicroDevicesInc-RSXCM": "AMD Software", "Microsoft.WindowsNotepad": "Windows Notepad", "Paint": "Paint"} {
		if got := prettyPkg(in); got != want {
			t.Errorf("%s: %q", in, got)
		}
	}
}
