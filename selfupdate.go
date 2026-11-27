package main

import (
	"crypto/sha256"
	hexenc "encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// De nieuwste versie staat als latest.json bij de laatste GitHub-release:
// {"version":"1.8","url":"https://.../DuskApp.exe","sha256":"...","notesEN":"...","notesNL":"..."}
var updateURL = "https://github.com/F4ir/dusk/releases/latest/download/latest.json"

type UpdateInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA256  string `json:"sha256"`
	NotesEN string `json:"notesEN"`
	NotesNL string `json:"notesNL"`
	Newer   bool   `json:"newer"`
	Err     string `json:"err"`
}

// CheckUpdate kijkt of er een nieuwere Dusk is.
func (a *App) CheckUpdate() UpdateInfo {
	var u UpdateInfo
	b, code, err := get(updateURL)
	if err != nil || code != 200 {
		u.Err = "offline"
		return u
	}
	if err := json.Unmarshal(b, &u); err != nil || u.Version == "" {
		u.Err = "bad"
		return u
	}
	u.Newer = newer(version, u.Version) && strings.HasPrefix(u.URL, "https://") && len(u.SHA256) == 64
	return u
}

// ApplyUpdate downloadt de nieuwe versie, controleert de SHA-256, zet hem op de plek van deze
// exe (de oude wordt .old) en start hem. Daarna sluit deze Dusk.
func (a *App) ApplyUpdate() string {
	u := a.CheckUpdate()
	if !u.Newer {
		return "none"
	}
	file, err := download(u.URL, func(p int) { a.emit("selfupdate", p) })
	if err != nil {
		return "ERR " + err.Error()
	}
	defer os.Remove(file)
	f, err := os.Open(file)
	if err != nil {
		return "ERR " + err.Error()
	}
	h := sha256.New()
	io.Copy(h, f)
	f.Close()
	if !strings.EqualFold(hexenc.EncodeToString(h.Sum(nil)), u.SHA256) {
		return "ERR checksum" // niet het goede bestand: niks vervangen
	}
	exe, err := os.Executable()
	if err != nil {
		return "ERR " + err.Error()
	}
	old := exe + ".old"
	os.Remove(old)
	// een draaiende exe mag je wel hernoemen, niet overschrijven
	if err := os.Rename(exe, old); err != nil {
		return "ERR " + err.Error()
	}
	if err := copyFile(file, exe); err != nil {
		os.Rename(old, exe) // terugzetten
		return "ERR " + err.Error()
	}
	if err := exec.Command(exe, "--no-intro", "--updated").Start(); err != nil {
		return "ERR " + err.Error()
	}
	runtime.Quit(a.ctx)
	return ""
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// cleanupOldExe: na een update staat de vorige versie nog als .old naast de exe.
func cleanupOldExe() {
	if exe, err := os.Executable(); err == nil {
		os.Remove(exe + ".old")
	}
}

// JustUpdated: is Dusk net bijgewerkt? (dan laat de app een melding zien)
func (a *App) JustUpdated() bool { return hasArg("--updated") }
