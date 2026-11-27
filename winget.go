package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// AppState is wat winget over een app weet.
type AppState struct {
	Installed bool
	Version   string
	Available string // leeg = geen update
	RealID    string // ID waaronder winget hem kent (kan een Store-ID zijn)
	Source    string // winget, msstore of leeg
	ARP       bool   // alleen op naam gevonden: winget kent geen pakket
	Note      string // extra uitleg, bv. "bundled:21, 25" (Modrinth heeft al eigen Java)
}

type wgRow struct{ Name, ID, Version, Available, Source string }

// parseTable leest de tabel van "winget list" met de kolomposities uit de
// kop. Zo werken ook ID's met spaties (ARP\...) en elke taal van Windows.
func parseTable(out string) []wgRow {
	lines := cleanLines(out)
	var rows []wgRow
	for i := 0; i+1 < len(lines); i++ {
		d := strings.TrimSpace(lines[i+1])
		if len(d) < 10 || strings.Trim(d, "-") != "" {
			continue
		}
		head := []rune(lines[i])
		var starts []int
		for j, r := range head {
			if r != ' ' && (j == 0 || head[j-1] == ' ') {
				starts = append(starts, j)
			}
		}
		if len(starts) < 3 {
			continue
		}
		col := func(r []rune, k int) string {
			if k >= len(starts) || starts[k] >= len(r) {
				return ""
			}
			e := len(r)
			if k+1 < len(starts) && starts[k+1] < e {
				e = starts[k+1]
			}
			return strings.TrimSpace(string(r[starts[k]:e]))
		}
		for _, l := range lines[i+2:] {
			r := []rune(l)
			row := wgRow{Name: col(r, 0), ID: col(r, 1), Version: col(r, 2)}
			if len(starts) == 4 {
				if v := col(r, 3); isSource(v) {
					row.Source = v
				} else {
					row.Available = v
				}
			} else {
				row.Available, row.Source = col(r, 3), col(r, 4)
			}
			if row.ID != "" {
				rows = append(rows, row)
			}
		}
		break
	}
	return rows
}

func isSource(s string) bool {
	s = strings.ToLower(s)
	return s == "winget" || s == "msstore"
}

// Namen zoals ze in "Programma's" staan, als die anders zijn dan App.Name.
var arpNames = map[string][]string{
	"Mozilla.Firefox":                      {"Mozilla Firefox"},
	"EpicGames.EpicGamesLauncher":          {"Epic Games Launcher"},
	"Microsoft.VisualStudioCode":           {"Microsoft Visual Studio Code"},
	"Telegram.TelegramDesktop":             {"Telegram Desktop"},
	"VideoLAN.VLC":                         {"VLC media player"},
	"Mozilla.Thunderbird":                  {"Mozilla Thunderbird"},
	"Microsoft.VisualStudio.Community":     {"Visual Studio Community 2026"},
	"Elgato.StreamDeck":                    {"Elgato Stream Deck"},
	"Elgato.WaveLink":                      {"Elgato Wave Link"},
	"Ookla.Speedtest.Desktop":              {"Speedtest by Ookla"},
	"Microsoft.VCRedist.2015+.x64":         {"Microsoft Visual C++ 2015-2022 Redistributable (x64)", "Microsoft Visual C++ v14 Redistributable (x64)"},
	"EclipseAdoptium.Temurin.21.JRE":       {"Eclipse Temurin JRE with Hotspot 21"},
	"Intel.IntelDriverAndSupportAssistant": {"Intel® Driver & Support Assistant", "Intel(R) Driver & Support Assistant"},
	"Corsair.iCUE.5":                       {"iCUE"},
	"Guru3D.Afterburner":                   {"MSI Afterburner"},
	"CPUID.CPU-Z":                          {"CPUID CPU-Z"},
	"TechPowerUp.GPU-Z":                    {"TechPowerUp GPU-Z"},
	"Microsoft.DotNet.DesktopRuntime.8":    {"Microsoft Windows Desktop Runtime - 8"},
	"Microsoft.Teams":                      {"Microsoft Teams"},
	"Microsoft.OneDrive":                   {"Microsoft OneDrive"},
	"Adobe.Acrobat.Reader.64-bit":          {"Adobe Acrobat"},
	"CodecGuide.K-LiteCodecPack.Standard":  {"K-Lite Codec Pack"},
	"Wagnardsoft.DisplayDriverUninstaller": {"Display Driver Uninstaller"},
	"Nvidia.GeForceNow":                    {"NVIDIA GeForce NOW"},
	"VB-Audio.Voicemeeter":                 {"VB:Voicemeeter", "Voicemeeter"},
	"OpenAI.ChatGPT":                       {"ChatGPT"},
	"Ollama.Ollama":                        {"Ollama version"},
	"Logitech.OptionsPlus":                 {"Logi Options+", "Logitech Options+"},
}

// "PuTTY release 0.85 (64-bit)" hoort bij "PuTTY", maar "Steam Link" niet bij "Steam".
var reNameRest = regexp.MustCompile(`^\s*($|\(|-|\d|v\d|version\b|release\b|x64|x86|64-bit|32-bit|\[)`)

func nameMatches(name string, a CatApp) bool {
	name = strings.ToLower(strings.TrimRight(name, "…. "))
	for _, n := range append([]string{a.Name}, arpNames[a.ID]...) {
		n = strings.ToLower(n)
		if strings.HasPrefix(name, n) && reNameRest.MatchString(name[len(n):]) {
			return true
		}
	}
	return false
}

// matchApps koppelt de rijen van winget aan onze lijst: eerst op ID, dan op naam.
func matchApps(rows []wgRow) map[string]AppState {
	res := map[string]AppState{}
	for _, a := range catalog {
		for _, r := range rows {
			if strings.EqualFold(r.ID, a.ID) {
				res[a.ID] = AppState{Installed: true, Version: r.Version, Available: r.Available, RealID: r.ID, Source: r.Source}
				break
			}
		}
		if _, ok := res[a.ID]; ok {
			continue
		}
		for _, r := range rows {
			if nameMatches(r.Name, a) {
				arp := r.Source == "" || strings.HasPrefix(r.ID, `ARP\`) || strings.HasPrefix(r.ID, "{") || strings.HasPrefix(r.ID, `MSIX\`)
				res[a.ID] = AppState{Installed: true, Version: r.Version, Available: r.Available, RealID: r.ID, Source: r.Source, ARP: arp}
				break
			}
		}
	}
	return res
}

func cleanLines(out string) []string {
	var res []string
	for _, l := range strings.Split(out, "\n") {
		if i := strings.LastIndex(strings.TrimRight(l, "\r"), "\r"); i >= 0 {
			l = l[i+1:]
		}
		l = strings.TrimRight(l, "\r ")
		// lege regels en het draaiende streepje overslaan, maar de ------ onder de kop houden
		if tl := strings.TrimSpace(l); tl == "" || (len(tl) <= 2 && strings.Trim(tl, `-\|/`) == "") {
			continue
		}
		res = append(res, l)
	}
	return res
}

func parseShowVersion(out string) string {
	for _, l := range cleanLines(out) {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "Version:") || strings.HasPrefix(l, "Versie:") {
			return strings.TrimSpace(l[strings.Index(l, ":")+1:])
		}
	}
	return ""
}

// newer zegt of versie b nieuwer is dan a.
func newer(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x = leadNum(pa[i])
		}
		if i < len(pb) {
			y = leadNum(pb[i])
		}
		if x != y {
			return y > x
		}
	}
	return false
}

func leadNum(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func scanApps() (map[string]AppState, error) {
	out, err := wingetList()
	rows := parseTable(out)
	st := matchApps(rows)
	fixSquirrel(st)
	detectExtras(st, rows)
	// zonder bekende versie (Roblox e.d. zeggen "Unknown") kunnen we niet vergelijken; zulke apps
	// werken zichzelf bij, dus geen eindeloze "update" laten zien
	for id, s := range st {
		if v := strings.ToLower(strings.TrimSpace(s.Version)); v == "" || v == "unknown" || v == "onbekend" {
			s.Available = ""
			st[id] = s
		}
	}
	for _, a := range catalog {
		s, ok := st[a.ID]
		if !ok || !s.ARP {
			continue
		}
		// winget kent geen pakket voor deze installatie: vraag de nieuwste versie van ons ID
		if latest := parseShowVersion(wingetShow(a.ID)); latest != "" && newer(s.Version, latest) {
			s.Available = latest
			st[a.ID] = s
		}
	}
	return st, err
}

func matchBloat(names []string) []Bloat {
	var found []Bloat
	for _, b := range bloatList {
		for _, n := range names {
			if ok, _ := path.Match(strings.ToLower(b.Pattern), strings.ToLower(n)); ok {
				found = append(found, b)
				break
			}
		}
	}
	return found
}

// ---------- uitvoer van winget volgen ----------

var (
	reSize = regexp.MustCompile(`([\d.,]+)\s*(KB|MB|GB)\s*/\s*([\d.,]+)\s*(KB|MB|GB)`)
	rePct  = regexp.MustCompile(`(\d{1,3})\s*%`)
)

func toBytes(n, unit string) float64 {
	f, _ := strconv.ParseFloat(strings.ReplaceAll(n, ",", "."), 64)
	switch unit {
	case "KB":
		return f * 1e3
	case "MB":
		return f * 1e6
	}
	return f * 1e9
}

// progress haalt een percentage uit een voortgangsregel (-1 = geen).
func progress(l string) int {
	if m := reSize.FindStringSubmatch(l); m != nil {
		tot := toBytes(m[3], m[4])
		if tot > 0 {
			return int(toBytes(m[1], m[2]) / tot * 100)
		}
	}
	if m := rePct.FindStringSubmatch(l); m != nil {
		p, _ := strconv.Atoi(m[1])
		return p
	}
	return -1
}

// streamCmd draait een programma en geeft elke regel (ook \r-regels) door.
func streamCmd(args []string, line func(string)) (int, error) {
	cmd := exec.Command(args[0], args[1:]...)
	hideWindow(cmd)
	pr, pw := io.Pipe()
	cmd.Stdout, cmd.Stderr = pw, pw
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	go func() {
		cmd.Wait()
		pw.Close()
	}()
	sc := bufio.NewScanner(pr)
	sc.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
			return i + 1, data[:i], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	last := ""
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if strings.Trim(l, `-\|/ `) == "" || l == last {
			continue
		}
		last = l
		line(l)
	}
	cmd.Wait()
	if cmd.ProcessState == nil {
		return -1, nil
	}
	return cmd.ProcessState.ExitCode(), nil
}

var wingetFlags = []string{"--exact", "--silent", "--accept-package-agreements", "--accept-source-agreements", "--disable-interactivity"}

func installArgs(a CatApp, st AppState) []string {
	if st.Installed && !st.ARP {
		// updaten onder het ID waaronder winget hem kent (bv. de Store-versie van PuTTY)
		id := st.RealID
		if id == "" {
			id = a.ID
		}
		args := []string{"winget", "upgrade", "--id", id}
		if st.Source != "" {
			args = append(args, "--source", st.Source)
		}
		return append(args, wingetFlags...)
	}
	args := []string{"winget", "install", "--id", a.ID}
	if reStoreID.MatchString(a.ID) { // Microsoft Store-app (bv. 9N7R5S6B0ZZH)
		args = append(args, "--source", "msstore")
	}
	return append(args, wingetFlags...)
}

var reStoreID = regexp.MustCompile(`^[0-9X][0-9A-Z]{11}$`)

// wingetOK: exitcodes die eigenlijk goed zijn (al nieuwste versie, al geïnstalleerd).
func wingetOK(code int) bool {
	switch uint32(code) {
	case 0, 0x8A15002B, 0x8A150061:
		return true
	}
	return false
}

// Apps die zichzelf bijwerken met Squirrel (Discord, Slack) zetten hun nieuwe versie niet in
// "Programma's" van Windows. Daar blijft de oude versie staan, waardoor winget steeds een update
// ziet. De echte versie staat in de naam van de nieuwste map app-x.y.z.
var squirrelDirs = map[string]string{
	"Discord.Discord":         `Discord`,
	"SlackTechnologies.Slack": `slack`,
}

func fixSquirrel(st map[string]AppState) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return
	}
	for id, dir := range squirrelDirs {
		s, ok := st[id]
		if !ok {
			continue
		}
		entries, _ := os.ReadDir(filepath.Join(base, dir))
		best := ""
		for _, e := range entries {
			if v, ok := strings.CutPrefix(e.Name(), "app-"); ok && e.IsDir() && (best == "" || newer(best, v)) {
				best = v
			}
		}
		if best == "" || newer(best, s.Version) {
			continue
		}
		s.Version = best
		if s.Available != "" && !newer(best, s.Available) {
			s.Available = "" // de map is al even nieuw als de "update"
		}
		st[id] = s
	}
}

var (
	rePythonRow = regexp.MustCompile(`(?i)^python (3\.\d+(\.\d+)?)`)
	reJavaRow   = regexp.MustCompile(`(?i)\b(java|jdk|jre|openjdk|temurin|zulu|corretto)\b`)
	reJavaVer   = regexp.MustCompile(`(\d+)`)
)

// detectExtras: Python en Java kun je op veel manieren hebben. Telt ook mee als het een ander
// pakket is (losse Python 3.x, Oracle Java, ...). Zonder update-melding, want dat is een ander product.
func detectExtras(st map[string]AppState, rows []wgRow) {
	if pythonRuntime(st) {
		// via de Python Install Manager (py): versie en update van de Python-runtime zelf
	} else if s := st["Python.PythonInstallManager"]; !s.Installed {
		for _, r := range rows {
			if m := rePythonRow.FindStringSubmatch(r.Name); m != nil {
				st["Python.PythonInstallManager"] = AppState{Installed: true, Version: m[1], RealID: r.ID, Source: r.Source, ARP: true}
				break
			}
		}
	}
	const jid = "EclipseAdoptium.Temurin.21.JRE"
	if s := st[jid]; !s.Installed {
		for _, r := range rows {
			if reJavaRow.MatchString(r.Name) && !strings.Contains(strings.ToLower(r.Name), "javascript") {
				v := r.Version
				if m := reJavaVer.FindString(r.Name); m != "" {
					v = "Java " + m
				}
				st[jid] = AppState{Installed: true, Version: v, RealID: r.ID, Source: r.Source, ARP: true}
				return
			}
		}
		// geen Java voor Windows, maar misschien wel een eigen Java in een Minecraft-launcher
		var found []string
		for _, dir := range []string{os.Getenv("APPDATA") + `\ModrinthApp\meta\java_versions`, os.Getenv("APPDATA") + `\.minecraft\runtime`,
			os.Getenv("APPDATA") + `\PrismLauncher\java`, os.Getenv("USERPROFILE") + `\curseforge\minecraft\Install\runtime`} {
			entries, _ := os.ReadDir(dir)
			for _, e := range entries {
				if m := regexp.MustCompile(`(?i)(?:jre|jdk|java-runtime-[a-z]*|zulu)[-_]?(\d{1,2})`).FindStringSubmatch(e.Name()); m != nil {
					found = append(found, m[1])
				}
			}
		}
		if len(found) > 0 {
			seen := map[string]bool{}
			var uniq []string
			for _, f := range found {
				if !seen[f] {
					seen[f] = true
					uniq = append(uniq, f)
				}
			}
			s := st[jid]
			s.Note = "bundled:" + strings.Join(uniq, ", ")
			st[jid] = s
		}
	}
}

var reUpdatingTo = regexp.MustCompile(`(?i)updating to python (\d+\.\d+\.\d+)`)

// pythonRuntime: is er een Python via de Install Manager? Dan de runtime-versie tonen
// (bv. 3.14.7) met de update die "py install --update" zou doen (bv. 3.14.10).
func pythonRuntime(st map[string]AppState) bool {
	out, err := runOut([]string{"py", "list", "--format=json"})
	if err != nil {
		return false
	}
	var l struct {
		Versions []struct {
			Tag, Company string
			Sort         string `json:"sort-version"`
			Default      bool   `json:"default"`
		} `json:"versions"`
	}
	if i := strings.Index(out, "{"); i < 0 || json.Unmarshal([]byte(out[i:]), &l) != nil || len(l.Versions) == 0 {
		return false
	}
	v := l.Versions[0]
	for _, x := range l.Versions {
		if x.Default {
			v = x
		}
	}
	s := AppState{Installed: true, Version: v.Sort, Note: "pyrt:" + v.Tag}
	c := exec.Command("py", "install", "--update", "--dry-run", v.Tag)
	hideWindow(c)
	if b, _ := c.CombinedOutput(); len(b) > 0 { // py schrijft dit naar stderr
		upd := string(b)
		if m := reUpdatingTo.FindStringSubmatch(upd); m != nil && newer(v.Sort, m[1]) {
			s.Available = m[1]
		}
	}
	st["Python.PythonInstallManager"] = s
	return true
}
