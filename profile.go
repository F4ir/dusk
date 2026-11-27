package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Profile: je apps en tweaks in één bestand, om een nieuwe pc (of die van een vriend) in te stellen.
type Profile struct {
	App     string   `json:"app"`
	Version string   `json:"version"`
	Created string   `json:"created"`
	Apps    []string `json:"apps"`   // winget-ID's uit de lijst van Dusk
	Tweaks  []string `json:"tweaks"` // tweaks die aan staan
}

// ExportProfile slaat je geïnstalleerde apps en aan-staande tweaks op. Antwoord: pad, of leeg bij annuleren.
func (a *App) ExportProfile() (string, error) {
	p := Profile{App: "Dusk", Version: version, Created: time.Now().Format("2006-01-02")}
	a.mu.Lock()
	st := a.last
	a.mu.Unlock()
	if st == nil {
		st, _ = scanApps()
	}
	for _, c := range catalog {
		if st[c.ID].Installed && c.Cat != "oem" {
			p.Apps = append(p.Apps, c.ID)
		}
	}
	for _, t := range tweaks {
		if tweakState(t) {
			p.Tweaks = append(p.Tweaks, t.ID)
		}
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "Dusk-" + time.Now().Format("2006-01-02") + ".dusk",
		Filters:         []runtime.FileFilter{{DisplayName: "Dusk (*.dusk)", Pattern: "*.dusk"}},
	})
	if err != nil || path == "" {
		return "", err
	}
	b, _ := json.MarshalIndent(p, "", "  ")
	return path, os.WriteFile(path, b, 0o644)
}

// ImportProfile laat je een .dusk-bestand kiezen en geeft de inhoud terug (de app laat eerst zien wat er gebeurt).
func (a *App) ImportProfile() (*Profile, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{{DisplayName: "Dusk (*.dusk)", Pattern: "*.dusk"}},
	})
	if err != nil || path == "" {
		return nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Profile
	if err := json.Unmarshal(b, &p); err != nil || p.App != "Dusk" {
		return nil, os.ErrInvalid
	}
	return &p, nil
}
