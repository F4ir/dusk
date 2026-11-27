package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	if len(os.Args) == 3 && (os.Args[1] == "--test-install" || os.Args[1] == "--test-uninstall") {
		// test via de opdrachtregel: echte installatie of verwijdering, met uitvoer in de console
		a := &App{}
		st, _ := scanApps()
		a.last = st
		id := os.Args[2]
		fmt.Printf("voor: %+v\n", st[id])
		if os.Args[1] == "--test-install" {
			a.InstallApps([]string{id})
		} else {
			var paths []string
			for _, l := range a.UninstallPlan([]string{id})[id] {
				fmt.Printf("  restje: %-11s %s\n", l.Kind, l.Path)
				paths = append(paths, l.Path)
			}
			a.UninstallApps([]string{id}, paths)
		}
		time.Sleep(time.Second)
		for a.busy {
			time.Sleep(500 * time.Millisecond)
		}
		st, _ = scanApps()
		fmt.Printf("na: %+v\n", st[id])
		return
	}
	if len(os.Args) == 3 && os.Args[1] == "--test-repair" {
		a := &App{}
		fmt.Println("antwoord:", a.Repair(os.Args[2]))
		time.Sleep(time.Second)
		for a.busy {
			time.Sleep(300 * time.Millisecond)
		}
		return
	}
	if len(os.Args) == 4 && os.Args[1] == "--ctxset" { // testen: --ctxset <naam> on|off
		a := &App{}
		for _, c := range a.ContextMenu() {
			if c.Name == os.Args[2] {
				fmt.Println("voor:", c.Enabled, "fout:", a.SetContextItem(c.ID, os.Args[3] == "on"))
			}
		}
		for _, c := range a.ContextMenu() {
			if c.Name == os.Args[2] {
				fmt.Println("na:", c.Enabled)
			}
		}
		return
	}
	if len(os.Args) == 4 && os.Args[1] == "--tweak" { // testen: --tweak <id> on|off
		a := &App{}
		t, err := a.SetTweak(os.Args[2], os.Args[3] == "on")
		fmt.Println(t.ID, "aan:", t.On, "fout:", err)
		return
	}
	for _, a := range os.Args[1:] {
		if a == "--check" {
			check()
			return
		}
	}
	cleanupOldExe()
	front, _ := fs.Sub(assets, "frontend")
	app := &App{}
	err := wails.Run(&options.App{
		Title:            "Dusk",
		Width:            1240,
		Height:           800,
		MinWidth:         1000,
		MinHeight:        660,
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 9, G: 7, B: 18, A: 255},
		AssetServer:      &assetserver.Options{Assets: front},
		OnStartup:        app.startup,
		StartHidden:      hasArg("--selftest"),
		Bind:             []any{app},
		Windows: &windows.Options{
			Messages:            wv2Messages(),
			Theme:               windows.Dark,
			DisableWindowIcon:   false,
			WebviewUserDataPath: userDataPath(),
		},
	})
	if err != nil {
		fmt.Println("Fout:", err)
	}
}

// userDataPath: instellingen in %APPDATA%\Dusk, ook als de exe ergens anders staat.
func userDataPath() string {
	if d, err := os.UserConfigDir(); err == nil {
		return d + `\Dusk`
	}
	return ""
}

// check print wat Dusk ziet, zonder iets te veranderen (om te testen).
func check() {
	a := &App{}
	fmt.Printf("%+v\n", a.System())
	apps, err := a.ScanApps()
	fmt.Println("--- apps (fout:", err, ")")
	for _, x := range apps {
		if x.Installed {
			fmt.Printf("  %-18s %-14s %s\n", x.Name, x.Version, x.Available)
		}
	}
	fmt.Println("--- tweaks")
	for _, t := range a.Tweaks() {
		fmt.Printf("  %-40s %v\n", t.NameEN, t.On)
	}
	if hasArg("--programs") {
		for _, p := range a.Programs() {
			fmt.Printf("  [%-5v] %-45.45s %-28.28s %6d MB  %s\n", p.Protected, p.Name, p.Publisher, p.SizeMB, p.CatalogID)
		}
	}
	if hasArg("--ctx") {
		for _, c := range a.ContextMenu() {
			fmt.Printf("  [%v] %-5s %-11s %-40.40s %s\n", c.Enabled, c.Kind, c.Where, c.Name, c.App)
		}
		for _, c := range a.Privacy() {
			fmt.Printf("  privacy %s aan=%v\n", c.ID, c.On)
			for _, x := range c.Apps {
				fmt.Printf("     [%v] inuse=%v %-30.30s %s\n", x.Allowed, x.InUse, x.Name, x.LastUsed)
			}
		}
	}
	if hasArg("--safetest") {
		pf, ad, pd := os.Getenv("ProgramFiles"), os.Getenv("APPDATA"), os.Getenv("ProgramData")
		for p, want := range map[string]bool{
			pd + `\Package Cache\{abc}`: false, pf: false, pf + `\Common Files`: false, pf + `\Google`: false,
			pf + `\Google\Chrome`: true, ad + `\Microsoft`: false, ad + `\discord`: true, `C:\Windows\System32\foo`: false,
			os.Getenv("USERPROFILE") + `\Documents`: false, os.Getenv("USERPROFILE") + `\Documents\Obsidian`: true,
		} {
			_, ok := safeToDelete(p)
			fmt.Printf("  safe %-5v wil %-5v %s\n", ok, want, p)
		}
	}
	if hasArg("--pplan") {
		var keys []string
		names := map[string]string{}
		for _, p := range a.Programs() {
			if !p.Protected {
				keys = append(keys, p.Key)
				names[p.Key] = p.Name
			}
		}
		for k, l := range a.ProgramPlan(keys) {
			fmt.Println("  # " + names[k])
			for _, x := range l {
				fmt.Printf("     %-11s %6d MB  %s\n", x.Kind, x.Size/1024/1024, x.Path)
			}
		}
	}
	if hasArg("--plan") {
		plan := a.UninstallPlan([]string{"Discord.Discord", "OBSProject.OBSStudio", "Spotify.Spotify", "PuTTY.PuTTY", "Valve.Steam", "Modrinth.ModrinthApp", "Bitwarden.Bitwarden", "Logitech.GHUB", "Git.Git"})
		for id, l := range plan {
			fmt.Println("  " + id)
			for _, x := range l {
				fmt.Printf("     %-11s %6d MB  %s\n", x.Kind, x.Size/1024/1024, x.Path)
			}
		}
	}
	if hasArg("--startup") {
		for _, x := range a.Startup() {
			fmt.Printf("  [%v] verborgen=%-5v %-28s %-30s %-20s %s\n", x.Enabled, x.Essential, x.Name, x.Desc, x.Maker, x.Source)
		}
	}
	if hasArg("--junk") {
		for _, j := range a.ScanJunk() {
			fmt.Printf("  %-8s %6d MB  %d bestanden\n", j.Key, j.Size/1024/1024, j.Files)
		}
	}
	if hasArg("--devices") {
		h, _ := a.Hardware()
		for _, d := range h.Devices {
			fmt.Printf("  %-9s %-50.50s %-22.22s %-16s %s p=%d g=%v\n", d.Cat, d.Name, d.Maker, d.Ver, d.Date, d.Problem, d.Generic)
		}
		fmt.Printf("  monitors %+v\n  batteries %+v\n  software %+v\n", h.Monitors, h.Batteries, h.Software)
		var q []GPUQuery
		for _, g := range h.GPUs {
			if !g.Virtual {
				sw := ""
				for _, s := range h.Software {
					if s.Name == "AMD Software" {
						sw = s.Ver
					}
				}
				q = append(q, GPUQuery{g.Name, g.Vendor, g.Ver, sw})
			}
		}
		st := a.CheckDrivers(q)
		fmt.Printf("  drivers: %+v\n", st)
	}
	if hasArg("--hw") {
		h, err := a.Hardware()
		fmt.Printf("--- hardware (fout: %v)\n%+v\n", err, h)
		d, err := a.ScanDrivers()
		fmt.Printf("--- drivers (fout: %v)\n%+v\n", err, d)
	}
	b, err := a.ScanBloat()
	fmt.Println("--- rommel (fout:", err, ")")
	for _, x := range b {
		fmt.Println("  ", x.NameEN)
	}
}

// wv2Messages: als WebView2 ontbreekt, vraagt Dusk om het te installeren (in de taal van Windows).
func wv2Messages() *windows.Messages {
	if !windowsIsDutch() {
		return nil
	}
	return &windows.Messages{
		InstallationRequired: "Dusk heeft Microsoft WebView2 nodig. Nu downloaden en installeren?",
		UpdateRequired:       "Dusk heeft een nieuwere versie van Microsoft WebView2 nodig. Nu bijwerken?",
		MissingRequirements:  "Onderdeel ontbreekt",
		Webview2NotInstalled: "WebView2 is niet geïnstalleerd, Dusk kan niet starten.",
		Error:                "Fout",
		FailedToInstall:      "WebView2 installeren is mislukt.",
		DownloadPage:         "Downloadpagina openen?",
		PressOKToInstall:     "Druk op OK om te installeren.",
		ContactAdmin:         "Vraag de beheerder van deze pc om WebView2 te installeren.",
		InvalidFixedWebview2: "WebView2 is ongeldig.",
		WebView2ProcessCrash: "WebView2 is gecrasht.",
	}
}
