//go:build windows

package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// Privacy: welke apps mogen je camera, microfoon en locatie gebruiken (Instellingen > Privacy).
const consentKey = `Software\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\`

type PrivApp struct {
	Key, Name, Path string
	Allowed, InUse  bool
	LastUsed        string // yyyy-mm-dd, leeg = nooit
	Packaged        bool
}

type PrivCap struct {
	ID   string
	On   bool // aan voor deze gebruiker
	Apps []PrivApp
}

var privCaps = []string{"webcam", "microphone", "location"}

func readConsent(k registry.Key, name string) (allowed, inUse bool, last string) {
	v, _, _ := k.GetStringValue("Value")
	allowed = !strings.EqualFold(v, "Deny")
	start, _, e1 := k.GetIntegerValue("LastUsedTimeStart")
	stop, _, e2 := k.GetIntegerValue("LastUsedTimeStop")
	if e1 == nil && start > 0 {
		t := time.Unix(0, (int64(start)-116444736000000000)*100)
		last = t.Format("2006-01-02")
		inUse = e2 == nil && stop == 0
	}
	return
}

func (a *App) Privacy() []PrivCap {
	var out []PrivCap
	for _, c := range privCaps {
		pc := PrivCap{ID: c, On: true}
		root, err := registry.OpenKey(registry.CURRENT_USER, consentKey+c, registry.READ)
		if err != nil {
			out = append(out, pc)
			continue
		}
		if v, _, err := root.GetStringValue("Value"); err == nil && strings.EqualFold(v, "Deny") {
			pc.On = false
		}
		subs, _ := root.ReadSubKeyNames(0)
		for _, s := range subs {
			if s == "NonPackaged" {
				continue
			}
			k, err := registry.OpenKey(root, s, registry.READ)
			if err != nil {
				continue
			}
			al, use, last := readConsent(k, s)
			k.Close()
			name := s
			if i := strings.Index(name, "_"); i > 0 { // Microsoft.WindowsCamera_8wekyb3d8bbwe -> WindowsCamera
				name = name[:i]
			}
			if i := strings.LastIndex(name, "."); i >= 0 {
				name = name[i+1:]
			}
			pc.Apps = append(pc.Apps, PrivApp{Key: s, Name: name, Allowed: al, InUse: use, LastUsed: last, Packaged: true})
		}
		if np, err := registry.OpenKey(root, "NonPackaged", registry.READ); err == nil {
			names, _ := np.ReadSubKeyNames(0)
			for _, s := range names {
				k, err := registry.OpenKey(np, s, registry.READ)
				if err != nil {
					continue
				}
				al, use, last := readConsent(k, s)
				k.Close()
				path := strings.ReplaceAll(s, "#", `\`)
				desc, _ := fileInfo(path)
				name := desc
				if name == "" {
					name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
				}
				if _, err := os.Stat(path); err != nil && last == "" {
					continue // oud programma dat er niet meer is en nooit gebruikt
				}
				pc.Apps = append(pc.Apps, PrivApp{Key: "NonPackaged\\" + s, Name: name, Path: path, Allowed: al, InUse: use, LastUsed: last})
			}
			np.Close()
		}
		root.Close()
		pc.Apps = mergePrivApps(pc.Apps)
		sort.Slice(pc.Apps, func(i, j int) bool {
			if pc.Apps[i].InUse != pc.Apps[j].InUse {
				return pc.Apps[i].InUse
			}
			return pc.Apps[i].LastUsed > pc.Apps[j].LastUsed
		})
		out = append(out, pc)
	}
	return out
}

// SetPrivacy zet toegang voor één app (key) of voor alles (key leeg) aan of uit.
func (a *App) SetPrivacy(capID, key string, allow bool) error {
	ok := false
	for _, c := range privCaps {
		if c == capID {
			ok = true
		}
	}
	if !ok || strings.Contains(key, "..") {
		return os.ErrInvalid
	}
	v := "Allow"
	if !allow {
		v = "Deny"
	}
	keys := strings.Split(key, "|")
	if key == "" {
		keys = []string{""}
	}
	for _, one := range keys {
		path := consentKey + capID
		if one != "" {
			path += `\` + one
		}
		k, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.SET_VALUE)
		if err != nil {
			return err
		}
		err = k.SetStringValue("Value", v)
		k.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

var nicePackaged = map[string]string{
	"windowscamera": "Camera", "immersivecontrolpanel": "Settings", "screensketch": "Snipping Tool",
	"windowssoundrecorder": "Sound Recorder", "whatsappdesktop": "WhatsApp", "windowsmaps": "Maps",
	"people": "People", "microsoftedge": "Microsoft Edge", "photos": "Photos", "clipchamp": "Clipchamp",
}

// mergePrivApps: elke Discord-versie staat apart in Windows; samenvoegen tot één regel per naam.
// De Key wordt dan "key1|key2|...", zodat aan/uit voor allemaal tegelijk geldt.
func mergePrivApps(list []PrivApp) []PrivApp {
	idx := map[string]int{}
	var out []PrivApp
	for _, p := range list {
		if n, ok := nicePackaged[strings.ToLower(p.Name)]; ok {
			p.Name = n
		}
		k := strings.ToLower(p.Name)
		if i, ok := idx[k]; ok {
			o := &out[i]
			o.Key += "|" + p.Key
			o.Allowed = o.Allowed || p.Allowed
			o.InUse = o.InUse || p.InUse
			if p.LastUsed > o.LastUsed {
				o.LastUsed = p.LastUsed
			}
			continue
		}
		idx[k] = len(out)
		out = append(out, p)
	}
	return out
}
