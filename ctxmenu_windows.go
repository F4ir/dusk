//go:build windows

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Rechtsklikmenu: items die apps toevoegen. Uitzetten zonder admin, per gebruiker:
//   - opdrachten (shell\<verb>): een lege override in HKCU\Software\Classes verbergt hem
//   - extensies (shellex\ContextMenuHandlers): CLSID in HKCU ...\Shell Extensions\Blocked
type CtxItem struct {
	ID, Name, App, Kind, Where string
	Enabled                    bool
}

var ctxPlaces = []struct{ path, where string }{
	{`*`, "files"}, {`AllFilesystemObjects`, "files"}, {`Directory`, "folders"}, {`Folder`, "folders"},
	{`Directory\Background`, "background"}, {`DesktopBackground`, "desktop"}, {`Drive`, "drives"},
	{`exefile`, "programs"}, {`lnkfile`, "shortcuts"}, {`CompressedFolder`, "files"},
}

// extraPlaces: alle bestandssoorten onder SystemFileAssociations (afbeeldingen, video, .zip, ...).
func extraPlaces() []struct{ path, where string } {
	var out []struct{ path, where string }
	k, err := registry.OpenKey(registry.CLASSES_ROOT, "SystemFileAssociations", registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil
	}
	defer k.Close()
	subs, _ := k.ReadSubKeyNames(0)
	for _, sub := range subs {
		out = append(out, struct{ path, where string }{`SystemFileAssociations\` + sub, "files"})
	}
	return out
}

// standaard-opdrachten van Windows zelf
var builtinVerbs = map[string]bool{
	"open": true, "opennewwindow": true, "opennewtab": true, "openas": true, "edit": true, "print": true, "printto": true,
	"runas": true, "runasuser": true, "find": true, "explore": true, "cmd": true, "powershell": true, "pintohome": true,
	"pintostartscreen": true, "properties": true, "uwp": true, "windows.share": true, "updateencryptionsettings": true,
	"updateencryptionsettingswork": true, "encrypt-bde": true, "encrypt-bde-elev": true, "manage-bde": true,
	"unlock-bde": true, "resume-bde": true, "resume-bde-elev": true, "change-passphrase": true, "change-pin": true,
	"opennewprocess": true, "anycode": true, "windows.recyclebin.empty": true, "importfrom": true,
}

const blockedKey = `Software\Microsoft\Windows\CurrentVersion\Shell Extensions\Blocked`

func sysDir(p string) bool {
	p = strings.ToLower(os.ExpandEnv(strings.ReplaceAll(strings.Trim(p, `" `), "%SystemRoot%", os.Getenv("SystemRoot"))))
	win := strings.ToLower(os.Getenv("SystemRoot"))
	return p == "" || strings.HasPrefix(p, win) || strings.Contains(p, `\windowsapps\`)
}

func isBlocked(clsid string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, blockedKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue(clsid)
	return err == nil
}

func (a *App) ContextMenu() []CtxItem {
	var out []CtxItem
	seen := map[string]bool{}
	for _, pl := range append(append([]struct{ path, where string }{}, ctxPlaces...), extraPlaces()...) {
		// opdrachten
		if k, err := registry.OpenKey(registry.CLASSES_ROOT, pl.path+`\shell`, registry.READ); err == nil {
			verbs, _ := k.ReadSubKeyNames(0)
			k.Close()
			for _, v := range verbs {
				if builtinVerbs[strings.ToLower(v)] {
					continue
				}
				id := "verb|" + pl.path + `\shell\` + v
				enabled := true
				if o, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\`+pl.path+`\shell\`+v, registry.QUERY_VALUE); err == nil {
					if _, _, err := o.GetStringValue("DuskDisabled"); err == nil {
						enabled = false
					}
					o.Close()
				}
				vk, err := registry.OpenKey(registry.CLASSES_ROOT, pl.path+`\shell\`+v, registry.READ)
				if err != nil && enabled {
					continue
				}
				name := v
				cmd := ""
				if err == nil {
					if s, _, e := vk.GetStringValue("MUIVerb"); e == nil && s != "" && !strings.HasPrefix(s, "@") {
						name = s
					} else if s, _, e := vk.GetStringValue(""); e == nil && s != "" && !strings.HasPrefix(s, "@") {
						name = s
					}
					if ck, e := registry.OpenKey(vk, "command", registry.QUERY_VALUE); e == nil {
						cmd, _, _ = ck.GetStringValue("")
						ck.Close()
					}
					vk.Close()
				}
				exe := exeFromCommand(cmd)
				if enabled && (cmd == "" || sysDir(exe)) {
					continue // van Windows zelf, of geen echte opdracht
				}
				desc, maker := fileInfo(exe)
				app := cleanMaker(desc, maker, exe)
				if seen[id] {
					continue
				}
				seen[id] = true
				out = append(out, CtxItem{id, strings.ReplaceAll(name, "&", ""), app, "verb", pl.where, enabled})
			}
		}
		// extensies
		if k, err := registry.OpenKey(registry.CLASSES_ROOT, pl.path+`\shellex\ContextMenuHandlers`, registry.READ); err == nil {
			hs, _ := k.ReadSubKeyNames(0)
			k.Close()
			for _, h := range hs {
				clsid := h
				if hk, err := registry.OpenKey(registry.CLASSES_ROOT, pl.path+`\shellex\ContextMenuHandlers\`+h, registry.QUERY_VALUE); err == nil {
					if s, _, e := hk.GetStringValue(""); e == nil && strings.HasPrefix(s, "{") {
						clsid = s
					}
					hk.Close()
				}
				if !strings.HasPrefix(clsid, "{") || seen["ext|"+strings.ToLower(clsid)] {
					continue
				}
				dll, name := "", h
				if ck, err := registry.OpenKey(registry.CLASSES_ROOT, `CLSID\`+clsid, registry.READ); err == nil {
					if s, _, e := ck.GetStringValue(""); e == nil && s != "" {
						name = s
					}
					if ik, e := registry.OpenKey(ck, "InprocServer32", registry.QUERY_VALUE); e == nil {
						dll, _, _ = ik.GetStringValue("")
						ik.Close()
					}
					ck.Close()
				}
				if dll == "" || sysDir(expandWin(dll)) {
					continue
				}
				dll = expandWin(dll)
				desc, maker := fileInfo(dll)
				app := cleanMaker(maker, desc, dll)
				if strings.HasPrefix(strings.ToLower(maker), "microsoft") && !strings.Contains(strings.ToLower(dll), "onedrive") {
					continue
				}
				seen["ext|"+strings.ToLower(clsid)] = true
				// "Rightmenu class", "{...}" of leeg: dan de naam van het programma
				if strings.HasPrefix(name, "{") || name == "" || strings.HasSuffix(strings.ToLower(name), " class") {
					name = app
				}
				out = append(out, CtxItem{"ext|" + clsid, name, app, "ext", pl.where, !isBlocked(clsid)})
			}
		}
	}
	out = append(out, packagedMenus()...)
	out = mergeCtx(out)
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].App+out[i].Name) < strings.ToLower(out[j].App+out[j].Name)
	})
	return out
}

// SetContextItem zet een item aan of uit (alleen voor deze gebruiker, zonder admin).
func (a *App) SetContextItem(id string, on bool) error {
	for _, one := range strings.Split(id, ";") {
		if err := setCtxOne(one, on); err != nil {
			return err
		}
	}
	return nil
}

func setCtxOne(id string, on bool) error {
	kind, rest, _ := strings.Cut(id, "|")
	switch kind {
	case "verb":
		path := `Software\Classes\` + rest
		if on {
			k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ENUMERATE_SUB_KEYS)
			if err != nil {
				return nil
			}
			defer k.Close()
			if _, _, e := k.GetStringValue("DuskDisabled"); e != nil {
				return nil // niet door Dusk uitgezet
			}
			k.DeleteValue("LegacyDisable")
			k.DeleteValue("DuskDisabled")
			// was het onze eigen lege sleutel (geen opdracht en geen andere waarden)? dan weg
			subs, _ := k.ReadSubKeyNames(0)
			vals, _ := k.ReadValueNames(0)
			if len(subs) == 0 && len(vals) == 0 {
				k.Close()
				return registry.DeleteKey(registry.CURRENT_USER, path)
			}
			return nil
		}
		// bestaat de sleutel al (app per gebruiker) dan komen de waarden erbij; anders een lege override
		k, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.SET_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		k.SetStringValue("DuskDisabled", "1")
		return k.SetStringValue("LegacyDisable", "")
	case "ext":
		k, _, err := registry.CreateKey(registry.CURRENT_USER, blockedKey, registry.SET_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		if on {
			return k.DeleteValue(rest)
		}
		return k.SetStringValue(rest, "Dusk")
	}
	return nil
}

// cleanMaker kiest een nette naam: geen "TODO: <Company name>" of leeg, anders de bestandsnaam.
func cleanMaker(first, second, file string) string {
	for _, v := range []string{first, second} {
		v = strings.TrimSpace(v)
		if v != "" && !strings.Contains(strings.ToLower(v), "todo") && !strings.Contains(v, "<") {
			return v
		}
	}
	// geen nette naam in het bestand: de programmamap gebruiken
	// (C:\Program Files\EaseUS\EaseUS Partition Master\bin\epmright.dll -> EaseUS Partition Master)
	f := strings.Trim(expandWin(file), `" `)
	for d := filepath.Dir(f); len(d) > 3; d = filepath.Dir(d) {
		b := filepath.Base(d)
		lb := strings.ToLower(b)
		if lb == "bin" || lb == "x64" || lb == "x86" || lb == "lib" || lb == "shell" || strings.Contains(lb, "menu") || strings.Contains(lb, "shellex") ||
			lb == "common" || lb == "shared" || lb == "plugins" || lb == "extensions" || regexp.MustCompile(`^[\d.]+$`).MatchString(b) {
			continue
		}
		if p := filepath.Dir(d); strings.EqualFold(p, os.Getenv("ProgramFiles")) || strings.EqualFold(p, os.Getenv("ProgramFiles(x86)")) || strings.EqualFold(filepath.Base(p), "Programs") {
			return b
		}
		if pp := filepath.Dir(filepath.Dir(d)); strings.EqualFold(pp, os.Getenv("ProgramFiles")) || strings.EqualFold(pp, os.Getenv("ProgramFiles(x86)")) {
			return b
		}
	}
	base := filepath.Base(f)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// expandWin vervangt %ProgramFiles% en andere %VAR% door hun waarde.
func expandWin(s string) string {
	return reEnv.ReplaceAllStringFunc(s, func(m string) string {
		if v := os.Getenv(strings.Trim(m, "%")); v != "" {
			return v
		}
		return m
	})
}

// mergeCtx: hetzelfde item op meerdere plekken (map + achtergrond) wordt één regel.
func mergeCtx(list []CtxItem) []CtxItem {
	idx := map[string]int{}
	var out []CtxItem
	for _, c := range list {
		k := strings.ToLower(c.Kind + "|" + c.App + "|" + c.Name)
		if c.Kind == "ext" {
			k = strings.ToLower(c.ID)
		}
		if i, ok := idx[k]; ok {
			o := &out[i]
			if !strings.Contains(";"+o.ID+";", ";"+c.ID+";") {
				o.ID += ";" + c.ID
			}
			if !strings.Contains(","+o.Where+",", ","+c.Where+",") {
				o.Where += "," + c.Where
			}
			o.Enabled = o.Enabled && c.Enabled
			continue
		}
		idx[k] = len(out)
		out = append(out, c)
	}
	return out
}

// packagedMenus: items in het nieuwe Windows 11-menu van apps uit de Store of met MSIX
// (Windows Terminal, Notepad++, VS Code, WinRAR 7, ...). Staan in de AppxManifest van de app.
const pkgMenuScript = `$ErrorActionPreference='SilentlyContinue'
[Console]::OutputEncoding=[Text.Encoding]::UTF8
$out=@()
foreach($p in Get-AppxPackage){
  $m=Join-Path $p.InstallLocation 'AppxManifest.xml'
  if(-not (Test-Path $m)){ continue }
  $t=[IO.File]::ReadAllText($m)
  if($t -notmatch 'fileExplorerContextMenus'){ continue }
  [xml]$x=$t
  $name=[string]$x.Package.Properties.DisplayName
  if(-not $name -or $name -like 'ms-resource*'){ $name=$p.Name }
  foreach($it in $x.SelectNodes("//*[local-name()='ItemType']")){
    foreach($v in $it.SelectNodes("*[local-name()='Verb']")){
      $out+=[pscustomobject]@{Pkg=$name;Pub=[string]$x.Package.Properties.PublisherDisplayName;Type=[string]$it.Type;Clsid=[string]$v.Clsid}
    }
  }
}
ConvertTo-Json -InputObject $out -Compress`

func packagedMenus() []CtxItem {
	out, _ := runOut(psArgs(pkgMenuScript))
	var list []struct{ Pkg, Pub, Type, Clsid string }
	json.Unmarshal([]byte(strings.TrimSpace(out)), &list)
	var res []CtxItem
	for _, l := range list {
		l.Clsid = strings.TrimSpace(l.Clsid)
		if len(l.Clsid) == 36 {
			l.Clsid = "{" + l.Clsid + "}" // in AppxManifest staan ze vaak zonder haakjes
		}
		if !strings.HasPrefix(l.Clsid, "{") {
			continue
		}
		l.Pkg = prettyPkg(l.Pkg)
		where := "files"
		switch strings.ToLower(l.Type) {
		case "directory":
			where = "folders"
		case `directory\background`:
			where = "background"
		case "drive":
			where = "drives"
		}
		res = append(res, CtxItem{"ext|" + l.Clsid, l.Pkg, cleanMaker(l.Pub, l.Pkg, l.Pkg), "ext", where, !isBlocked(l.Clsid)})
	}
	return res
}
