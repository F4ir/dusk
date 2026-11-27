package main

// CatApp is een programma dat via winget geïnstalleerd/geüpdatet kan worden.
// Cat is een sleutel; de app vertaalt hem (browser, chat, games, ...).
type CatApp struct {
	Name, ID, Cat, Icon, Color, DescEN, DescNL string
}

var catalog = []CatApp{
	// browsers
	{"Google Chrome", "Google.Chrome", "browser", "googlechrome", "#4285F4", "Google's web browser.", "De browser van Google."},
	{"Firefox", "Mozilla.Firefox", "browser", "firefoxbrowser", "#FF7139", "Fast browser that protects your privacy.", "Snelle browser die je privacy beschermt."},
	{"Brave", "Brave.Brave", "browser", "brave", "#FB542B", "Browser with a built-in ad blocker.", "Browser met ingebouwde adblocker."},
	{"Opera GX", "Opera.OperaGX", "browser", "operagx", "#EE2950", "Browser made for gamers.", "Browser gemaakt voor gamers."},
	{"Microsoft Edge", "Microsoft.Edge", "browser", "microsoftedge", "#0C59A4", "Microsoft's browser, built into Windows.", "De browser van Microsoft, zit in Windows."},
	{"Vivaldi", "Vivaldi.Vivaldi", "browser", "vivaldi", "#EF3939", "Browser you can set up exactly how you like.", "Browser die je helemaal zelf kunt instellen."},
	{"Zen Browser", "Zen-Team.Zen-Browser", "browser", "zenbrowser", "#2B2A33", "Calm, modern browser based on Firefox.", "Rustige, moderne browser op basis van Firefox."},
	{"Arc", "TheBrowserCompany.Arc", "browser", "arc", "#FCBFBD", "Browser with tabs in a side bar.", "Browser met tabbladen in een zijbalk."},
	// chat
	{"Discord", "Discord.Discord", "chat", "discord", "#5865F2", "Chat and call with friends.", "Chatten en bellen met vrienden."},
	{"Claude", "Anthropic.Claude", "ai", "claude", "#D97757", "The Claude app for your PC.", "De Claude-app voor je pc."},
	{"Telegram", "Telegram.TelegramDesktop", "chat", "telegram", "#26A5E4", "Fast messaging app.", "Snelle berichten-app."},
	{"WhatsApp", "9NKSQGP7F2NH", "chat", "whatsapp", "#25D366", "WhatsApp on your PC.", "WhatsApp op je pc."},
	{"Signal", "OpenWhisperSystems.Signal", "chat", "signal", "#3B45FD", "Private, encrypted messaging.", "Veilige, versleutelde berichten."},
	{"Slack", "SlackTechnologies.Slack", "chat", "slack", "#4A154B", "Team chat.", "Chatten met je team."},
	{"Zoom", "Zoom.Zoom", "chat", "zoom", "#0B5CFF", "Video meetings.", "Videobellen en meetings."},
	{"Thunderbird", "Mozilla.Thunderbird", "chat", "thunderbird", "#0A84FF", "Free e-mail app.", "Gratis mail-programma."},
	{"Microsoft Teams", "Microsoft.Teams", "chat", "microsoftteams", "#5059C9", "Chat and meetings for school and work.", "Chatten en vergaderen voor school en werk."},
	{"Proton Mail", "Proton.ProtonMail", "chat", "protonmail", "#6D4AFF", "Private, encrypted e-mail.", "Privé, versleutelde e-mail."},
	// AI
	{"ChatGPT", "OpenAI.ChatGPT", "ai", "openai", "#10A37F", "OpenAI's ChatGPT app for Windows.", "De ChatGPT-app van OpenAI voor Windows."},
	{"Perplexity", "Perplexity.Perplexity", "ai", "perplexity", "#1F8A8A", "AI search with sources.", "Zoeken met AI, met bronnen erbij."},
	{"Ollama", "Ollama.Ollama", "ai", "ollama", "#2B2B2B", "Run AI models on your own PC.", "AI-modellen op je eigen pc draaien."},
	{"LM Studio", "ElementLabs.LMStudio", "ai", "lmstudio", "#4C3BCF", "Download and chat with local AI models.", "Lokale AI-modellen downloaden en ermee chatten."},
	{"Jan", "Jan.Jan", "ai", "", "#1D1D1F", "Offline ChatGPT-style app for local models.", "Offline chat-app voor lokale modellen."},
	{"Cursor", "Anysphere.Cursor", "ai", "cursor", "#1A1A1A", "Code editor with AI built in.", "Code-editor met ingebouwde AI."},
	{"Windsurf", "Codeium.Windsurf", "ai", "windsurf", "#0B8C7E", "AI code editor.", "Code-editor met AI."},
	// games
	{"Steam", "Valve.Steam", "games", "steam", "#1B2838", "The biggest PC game store.", "De grootste winkel voor pc-games."},
	{"Epic Games", "EpicGames.EpicGamesLauncher", "games", "epicgames", "#313131", "Free games every week.", "Elke week gratis games."},
	{"EA app", "ElectronicArts.EADesktop", "games", "ea", "#FF4747", "Games from EA (FC, Battlefield, Apex).", "Games van EA (FC, Battlefield, Apex)."},
	{"Ubisoft Connect", "Ubisoft.Connect", "games", "ubisoft", "#0070FF", "Games from Ubisoft.", "Games van Ubisoft."},
	{"Battle.net", "Blizzard.BattleNet", "games", "battledotnet", "#4381C3", "Blizzard games (CoD, Overwatch, WoW).", "Blizzard-games (CoD, Overwatch, WoW)."},
	{"GOG Galaxy", "GOG.Galaxy", "games", "gogdotcom", "#86328A", "DRM-free games and all your launchers in one.", "Games zonder DRM en al je launchers in één."},
	{"Playnite", "Playnite.Playnite", "games", "", "#5B3FD1", "One library for all your games.", "Eén bibliotheek voor al je games."},
	{"Modrinth App", "Modrinth.ModrinthApp", "games", "modrinth", "#00AF5C", "Minecraft with mods, made easy.", "Minecraft met mods, makkelijk."},
	{"CurseForge", "Overwolf.CurseForge", "games", "curseforge", "#F16436", "Mods for Minecraft and more.", "Mods voor Minecraft en meer."},
	{"Parsec", "Parsec.Parsec", "games", "", "#7C46E6", "Play together on one PC over the internet.", "Samen op één pc spelen via internet."},
	{"GeForce NOW", "Nvidia.GeForceNow", "games", "nvidia", "#76B900", "Play PC games in the cloud.", "Pc-games spelen in de cloud."},
	// randapparatuur
	{"Stream Deck", "Elgato.StreamDeck", "gear", "elgato", "#1A1A2E", "Software for the Elgato Stream Deck.", "Software voor de Elgato Stream Deck."},
	{"Wave Link", "Elgato.WaveLink", "gear", "elgato", "#2A2350", "Elgato audio mixer for your mic.", "Elgato-geluidsmixer voor je microfoon."},
	{"Logitech G HUB", "Logitech.GHUB", "gear", "logitechg", "#141416", "Settings for Logitech mice and keyboards.", "Instellingen voor Logitech-muizen en toetsenborden."},
	{"SteelSeries GG", "SteelSeries.GG", "gear", "steelseries", "#FF5200", "Settings for SteelSeries gear.", "Instellingen voor SteelSeries-spullen."},
	{"Corsair iCUE", "Corsair.iCUE.5", "gear", "corsair", "#2B2A28", "RGB and settings for Corsair gear.", "RGB en instellingen voor Corsair-spullen."},
	{"Logi Options+", "Logitech.OptionsPlus", "gear", "logitech", "#00B8FC", "Settings for Logitech office mice and keyboards.", "Instellingen voor Logitech-kantoormuizen en -toetsenborden."},
	{"Logi Tune", "Logitech.LogiTune", "gear", "logitech", "#814EFA", "Settings for Logitech headsets and webcams.", "Instellingen voor Logitech-headsets en -webcams."},
	{"Elgato Camera Hub", "Elgato.CameraHub", "gear", "elgato", "#1A1A2E", "Settings for Elgato webcams.", "Instellingen voor Elgato-webcams."},
	{"Elgato 4K Capture", "Elgato.4KCaptureUtility", "gear", "elgato", "#1A1A2E", "Record with an Elgato capture card.", "Opnemen met een Elgato-capturekaart."},
	// media
	{"Spotify", "Spotify.Spotify", "media", "spotify", "#1ED760", "Music and podcasts.", "Muziek en podcasts."},
	{"VLC", "VideoLAN.VLC", "media", "vlcmediaplayer", "#FF8800", "Plays every video.", "Speelt elke video af."},
	{"OBS Studio", "OBSProject.OBSStudio", "media", "obsstudio", "#302E31", "Record and stream.", "Opnemen en streamen."},
	{"ShareX", "ShareX.ShareX", "media", "sharex", "#2885F1", "Screenshots and screen recordings.", "Screenshots en schermopnames."},
	{"Jellyfin Media Player", "Jellyfin.JellyfinMediaPlayer", "media", "jellyfin", "#00A4DC", "Watch your own movies and series.", "Je eigen films en series kijken."},
	{"Audacity", "Audacity.Audacity", "media", "audacity", "#2D33CC", "Record and edit audio.", "Geluid opnemen en bewerken."},
	{"GIMP", "GIMP.GIMP.3", "media", "gimp", "#6E6455", "Free photo editor.", "Gratis fotobewerker."},
	{"Paint.NET", "dotPDN.PaintDotNet", "media", "", "#2F6FD6", "Simple, fast image editor.", "Simpele, snelle fotobewerker."},
	{"Blender", "BlenderFoundation.Blender", "media", "blender", "#E87D0D", "3D modelling and animation.", "3D-modellen en animaties maken."},
	{"HandBrake", "HandBrake.HandBrake", "media", "", "#C4372B", "Convert and shrink videos.", "Video's omzetten en kleiner maken."},
	{"Plex", "Plex.Plex", "media", "plex", "#E5A00D", "Watch your own movies and shows.", "Je eigen films en series kijken."},
	{"Plexamp", "Plex.Plexamp", "media", "plex", "#E5A00D", "Music player for Plex.", "Muziekspeler voor Plex."},
	{"iTunes", "Apple.iTunes", "media", "itunes", "#FB5BC5", "Apple's music app and iPhone backups.", "Muziek-app van Apple en iPhone-backups."},
	{"foobar2000", "PeterPawlowski.foobar2000", "media", "foobar2000", "#2C2C2C", "Light, powerful music player.", "Lichte, krachtige muziekspeler."},
	{"MPC-HC", "clsid2.mpc-hc", "media", "", "#3D6CB7", "Light video player.", "Lichte videospeler."},
	{"K-Lite Codec Pack", "CodecGuide.K-LiteCodecPack.Standard", "media", "", "#3A73C4", "Plays every video and audio format.", "Speelt elk video- en audioformaat af."},
	{"Krita", "KDE.Krita", "media", "krita", "#3BABFF", "Free drawing and painting app.", "Gratis teken- en schilderprogramma."},
	{"Inkscape", "Inkscape.Inkscape", "media", "inkscape", "#2C2C2C", "Free vector drawing (like Illustrator).", "Gratis vectortekenen (zoals Illustrator)."},
	{"Kdenlive", "KDE.Kdenlive", "media", "kdenlive", "#527EB2", "Free video editor.", "Gratis video-editor."},
	{"Shotcut", "Meltytech.Shotcut", "media", "shotcut", "#115C77", "Simple free video editor.", "Simpele gratis video-editor."},
	{"Voicemeeter", "VB-Audio.Voicemeeter", "media", "", "#2B5C8A", "Mix all your audio, great for streaming.", "Al je geluid mixen, handig voor streamen."},
	// tools
	{"7-Zip", "7zip.7zip", "tools", "7zip", "#2B2B2B", "Pack and unpack zip files.", "Zip-bestanden in- en uitpakken."},
	{"WinRAR", "RARLab.WinRAR", "tools", "", "#7B4BC4", "Unpack rar and zip files. Keeps working, but asks you to buy a license.", "Rar- en zip-bestanden uitpakken. Blijft werken, maar vraagt om een licentie."},
	{"PowerToys", "Microsoft.PowerToys", "tools", "", "#3B6FE0", "Handy extras from Microsoft.", "Handige extra's van Microsoft."},
	{"Everything", "voidtools.Everything", "tools", "", "#E0752D", "Finds any file in a second.", "Vindt elk bestand in een tel."},
	{"Notepad++", "Notepad++.Notepad++", "tools", "notepadplusplus", "#5BAA61", "Notepad, but better.", "Kladblok, maar dan beter."},
	{"Bitwarden", "Bitwarden.Bitwarden", "tools", "bitwarden", "#175DDC", "Password manager.", "Wachtwoordkluis."},
	{"KeePassXC", "KeePassXCTeam.KeePassXC", "tools", "keepassxc", "#6CAC4D", "Offline password manager.", "Wachtwoordkluis zonder internet."},
	{"Notion", "Notion.Notion", "tools", "notion", "#2F2F2F", "Notes, tasks and wikis.", "Notities, taken en wiki's."},
	{"Obsidian", "Obsidian.Obsidian", "tools", "obsidian", "#7C3AED", "Notes that link together.", "Notities die aan elkaar linken."},
	{"LibreOffice", "TheDocumentFoundation.LibreOffice", "tools", "libreoffice", "#18A303", "Free Word and Excel.", "Gratis Word en Excel."},
	{"qBittorrent", "qBittorrent.qBittorrent", "tools", "qbittorrent", "#2F67BA", "Torrent client without ads.", "Torrent-programma zonder reclame."},
	{"Malwarebytes", "Malwarebytes.Malwarebytes", "tools", "malwarebytes", "#0D3ECC", "Scans for viruses and malware.", "Scant op virussen en malware."},
	{"Revo Uninstaller", "RevoUninstaller.RevoUninstaller", "tools", "", "#D9452B", "Removes programs including leftovers.", "Verwijdert programma's met restjes."},
	{"Rufus", "Rufus.Rufus", "tools", "", "#4F7D9A", "Make a bootable USB stick.", "Een opstartbare usb-stick maken."},
	{"Ventoy", "Ventoy.Ventoy", "tools", "", "#2A8E4B", "Many ISO files on one USB stick.", "Veel ISO's op één usb-stick."},
	{"NanaZip", "M2Team.NanaZip", "tools", "", "#3D8AD8", "Modern 7-Zip with the new Windows 11 menu.", "Moderne 7-Zip met het nieuwe Windows 11-menu."},
	{"Adobe Acrobat Reader", "Adobe.Acrobat.Reader.64-bit", "tools", "adobeacrobatreader", "#EC1C24", "Open and sign PDF files.", "Pdf's openen en ondertekenen."},
	{"SumatraPDF", "SumatraPDF.SumatraPDF", "tools", "", "#E3A21A", "Very fast, light PDF reader.", "Heel snelle, lichte pdf-lezer."},
	{"Google Drive", "Google.GoogleDrive", "tools", "googledrive", "#4285F4", "Your Google Drive files in Explorer.", "Je Google Drive-bestanden in Verkenner."},
	{"Dropbox", "Dropbox.Dropbox", "tools", "dropbox", "#0061FF", "Cloud storage and file sync.", "Cloudopslag en bestanden synchroniseren."},
	{"OneDrive", "Microsoft.OneDrive", "tools", "microsoftonedrive", "#0078D4", "Microsoft cloud storage.", "Cloudopslag van Microsoft."},
	{"1Password", "AgileBits.1Password", "tools", "1password", "#0094F5", "Popular password manager.", "Bekende wachtwoordkluis."},
	{"Proton Pass", "Proton.ProtonPass", "tools", "proton", "#6D4AFF", "Private password manager.", "Privé wachtwoordkluis."},
	{"Greenshot", "Greenshot.Greenshot", "tools", "", "#4AAA4A", "Quick screenshots with arrows and text.", "Snel screenshots met pijlen en tekst."},
	{"Flow Launcher", "Flow-Launcher.Flow-Launcher", "tools", "", "#2C8CFF", "Find apps and files super fast with a shortcut.", "Supersnel apps en bestanden vinden met een sneltoets."},
	{"AutoHotkey", "AutoHotkey.AutoHotkey", "tools", "autohotkey", "#334455", "Make your own shortcuts and macros.", "Eigen sneltoetsen en macro's maken."},
	{"Ditto", "Ditto.Ditto", "tools", "", "#3C7FC2", "Clipboard history for everything you copy.", "Klembordgeschiedenis van alles wat je kopieert."},
	{"Files", "9NGHP3DX8HDX", "tools", "", "#4F6BED", "Modern file manager with tabs.", "Moderne bestandsbeheerder met tabbladen."},
	{"EarTrumpet", "File-New-Project.EarTrumpet", "tools", "", "#4A90D9", "Volume per app, right from the taskbar.", "Volume per app, direct vanuit de taakbalk."},
	{"TranslucentTB", "CharlesMilette.TranslucentTB", "tools", "", "#7B5BD6", "Make your taskbar see-through.", "Maakt je taakbalk doorzichtig."},
	{"Rainmeter", "Rainmeter.Rainmeter", "tools", "rainmeter", "#19519B", "Widgets and stats on your desktop.", "Widgets en statistieken op je bureaublad."},
	{"Lively Wallpaper", "rocksdanister.LivelyWallpaper", "tools", "", "#4E2A84", "Animated wallpapers.", "Bewegende achtergronden."},
	// systeem
	{"HWiNFO", "REALiX.HWiNFO", "system", "", "#2F6DB5", "Temperatures and all hardware info.", "Temperaturen en alle hardware-info."},
	{"CPU-Z", "CPUID.CPU-Z", "system", "", "#4361A8", "Details about your processor and RAM.", "Alles over je processor en geheugen."},
	{"GPU-Z", "TechPowerUp.GPU-Z", "system", "", "#B5242A", "Details about your graphics card.", "Alles over je videokaart."},
	{"CrystalDiskInfo", "CrystalDewWorld.CrystalDiskInfo", "system", "", "#3B8ED8", "Checks if your drives are healthy.", "Kijkt of je schijven gezond zijn."},
	{"MSI Afterburner", "Guru3D.Afterburner", "system", "msi", "#D11F2B", "Overclocking and FPS overlay.", "Overklokken en FPS in beeld."},
	{"WinDirStat", "WinDirStat.WinDirStat", "system", "", "#3E8E41", "See what takes up your disk space.", "Zie wat je schijf vol maakt."},
	{"Intel Driver Assistant", "Intel.IntelDriverAndSupportAssistant", "system", "intel", "#0071C5", "Finds the newest Intel drivers.", "Zoekt de nieuwste Intel-drivers."},
	{"TreeSize Free", "JAMSoftware.TreeSize.Free", "system", "", "#2F8F4E", "See which folders take up the most space.", "Zie welke mappen de meeste ruimte innemen."},
	{"CrystalDiskMark", "CrystalDewWorld.CrystalDiskMark", "system", "", "#3B8ED8", "Test how fast your drives are.", "Test hoe snel je schijven zijn."},
	{"Process Explorer", "Microsoft.Sysinternals.ProcessExplorer", "system", "", "#2F5D9B", "Task Manager, but much more detailed.", "Taakbeheer, maar veel uitgebreider."},
	{"Autoruns", "Microsoft.Sysinternals.Autoruns", "system", "", "#2F5D9B", "Everything that starts with Windows.", "Alles wat met Windows mee opstart."},
	{"DDU", "Wagnardsoft.DisplayDriverUninstaller", "system", "", "#C0392B", "Fully removes old graphics drivers.", "Haalt oude videodrivers helemaal weg."},
	{"NVCleanstall", "TechPowerUp.NVCleanstall", "system", "nvidia", "#76B900", "Install NVIDIA drivers without extras.", "NVIDIA-drivers installeren zonder extra's."},
	// netwerk
	{"Advanced IP Scanner", "Famatech.AdvancedIPScanner", "network", "", "#1F8FD6", "See every device on your network.", "Zie elk apparaat in je netwerk."},
	{"WinSCP", "WinSCP.WinSCP", "network", "", "#2C6EB5", "FTP and SFTP file transfers.", "Bestanden via FTP en SFTP."},
	{"PuTTY", "PuTTY.PuTTY", "network", "", "#3A5A9B", "SSH to your servers.", "SSH naar je servers."},
	{"Wireshark", "WiresharkFoundation.Wireshark", "network", "wireshark", "#1679A7", "See all network traffic.", "Zie al het netwerkverkeer."},
	{"Tailscale", "Tailscale.Tailscale", "network", "tailscale", "#3A3A3A", "Your own private network, anywhere.", "Je eigen privénetwerk, overal."},
	{"OpenVPN Connect", "OpenVPNTechnologies.OpenVPNConnect", "network", "openvpn", "#EA7E20", "Connect to an OpenVPN server.", "Verbinden met een OpenVPN-server."},
	{"Proton VPN", "Proton.ProtonVPN", "network", "protonvpn", "#6D4AFF", "Free, private VPN.", "Gratis, privé vpn."},
	{"NordVPN", "NordSecurity.NordVPN", "network", "nordvpn", "#4687FF", "Popular paid VPN.", "Bekende betaalde vpn."},
	{"TeamViewer", "TeamViewer.TeamViewer", "network", "teamviewer", "#0E8EE9", "Control another PC remotely.", "Een andere pc op afstand besturen."},
	{"AnyDesk", "AnyDesk.AnyDesk", "network", "anydesk", "#EF443B", "Fast remote desktop.", "Snel op afstand meekijken."},
	{"Speedtest", "Ookla.Speedtest.Desktop", "network", "speedtest", "#141526", "Test your internet speed.", "Test je internetsnelheid."},
	// programmeren
	{"Visual Studio 2026", "Microsoft.VisualStudio.Community", "dev", "", "#865FC5", "Microsoft's full IDE for C#, C++ and .NET.", "De volledige IDE van Microsoft voor C#, C++ en .NET."},
	{"VS Code", "Microsoft.VisualStudioCode", "dev", "visualstudiocode", "#007ACC", "Write code.", "Code schrijven."},
	{"Git", "Git.Git", "dev", "git", "#F03C2E", "Version control for code.", "Versiebeheer voor code."},
	{"Windows Terminal", "Microsoft.WindowsTerminal", "dev", "windowsterminal", "#4D4D4D", "The modern command prompt.", "De moderne opdrachtprompt."},
	{"PowerShell 7", "Microsoft.PowerShell", "dev", "powershell", "#2671BE", "The newest PowerShell.", "De nieuwste PowerShell."},
	{"Python", "Python.PythonInstallManager", "dev", "python", "#3776AB", "Python with the official install manager (py command).", "Python met de officiële installatiebeheerder (py-commando)."},
	{"Node.js", "OpenJS.NodeJS.LTS", "dev", "nodedotjs", "#5FA04E", "JavaScript outside the browser.", "JavaScript buiten de browser."},
	// officiële updatetools van merken (alleen zichtbaar bij dat merk)
	{"HP Support Assistant", "HP.SupportAssistant", "oem", "hp", "#0096D6", "HP's own tool for drivers and BIOS for your exact model.", "HP's eigen tool voor drivers en BIOS van precies jouw model."},
	{"HP Image Assistant", "HP.ImageAssistant", "oem", "hp", "#0072C6", "HP's advanced driver tool (for HP business PCs).", "HP's uitgebreide drivertool (voor zakelijke HP-pc's)."},
	{"Dell Command Update", "Dell.CommandUpdate.Universal", "oem", "dell", "#007DB8", "Dell's own tool for drivers and BIOS for your exact model.", "Dell's eigen tool voor drivers en BIOS van precies jouw model."},
	{"Lenovo Vantage", "9WZDNCRFJ4MV", "oem", "lenovo", "#E2231A", "Lenovo's own app for drivers, BIOS and settings.", "Lenovo's eigen app voor drivers, BIOS en instellingen."},
	{"Lenovo System Update", "Lenovo.SystemUpdate", "oem", "lenovo", "#C4161C", "Driver and BIOS updates for ThinkPad and ThinkCentre.", "Driver- en BIOS-updates voor ThinkPad en ThinkCentre."},
	{"MyASUS", "9N7R5S6B0ZZH", "oem", "asus", "#00539B", "ASUS's own app for drivers and BIOS.", "De eigen app van ASUS voor drivers en BIOS."},
	{"Armoury Crate", "Asus.ArmouryCrate", "oem", "asus", "#B0001E", "Drivers, RGB and performance for ROG and TUF.", "Drivers, RGB en prestaties voor ROG en TUF."},
	{"Acer Care Center S", "9P8BB54NQNQ4", "oem", "acer", "#83B81A", "Acer's own app for drivers and updates.", "De eigen app van Acer voor drivers en updates."},
	// onderdelen die games nodig hebben
	{"Java 21", "EclipseAdoptium.Temurin.21.JRE", "runtime", "eclipseadoptium", "#E21D66", "Needed for Minecraft servers and Java apps.", "Nodig voor Minecraft-servers en Java-apps."},
	{"Visual C++", "Microsoft.VCRedist.2015+.x64", "runtime", "", "#6A3BB5", "Many games crash without this.", "Veel games crashen zonder dit."},
	{"DirectX", "Microsoft.DirectX", "runtime", "", "#2C8A3E", "Older DirectX parts for games.", "Oudere DirectX-onderdelen voor games."},
	{".NET 8 Desktop", "Microsoft.DotNet.DesktopRuntime.8", "runtime", "dotnet", "#512BD4", "Needed for many Windows apps.", "Nodig voor veel Windows-apps."},
}

// popular: de meest gebruikte apps, in deze volgorde bovenaan.
var popular = []string{
	"Discord.Discord", "Valve.Steam", "Google.Chrome", "Spotify.Spotify", "Logitech.GHUB", "Microsoft.VCRedist.2015+.x64",
	"Microsoft.VisualStudioCode", "Microsoft.VisualStudio.Community", "EpicGames.EpicGamesLauncher", "OBSProject.OBSStudio",
	"7zip.7zip", "RARLab.WinRAR", "9NKSQGP7F2NH", "Mozilla.Firefox", "Brave.Brave", "Blizzard.BattleNet",
	"Ubisoft.Connect", "Elgato.StreamDeck", "Elgato.WaveLink", "Opera.OperaGX", "VideoLAN.VLC", "Anthropic.Claude", "OpenAI.ChatGPT",
	"Bitwarden.Bitwarden", "Microsoft.DirectX", "EclipseAdoptium.Temurin.21.JRE", "Microsoft.PowerToys", "Malwarebytes.Malwarebytes",
}

func popRank(id string) int {
	for i, p := range popular {
		if p == id {
			return i + 1
		}
	}
	return 0
}

// RegVal is één registerwaarde van een tweak.
// Kind: "dword", "sz" of "key" (key = de sleutel bestaat = aan).
type RegVal struct {
	Root, Path, Name, Kind string
	On                     any
	Off                    any // nil = waarde (of sleutel) weghalen
}

// Tweak: Group is een sleutel (looks, explorer, taskbar, gaming, speed, privacy).
type Tweak struct {
	ID, Group              string
	NameEN, DescEN, NoteEN string
	NameNL, DescNL, NoteNL string
	Admin                  bool   // moet als administrator
	Explorer               bool   // Verkenner herstarten voor effect
	Fn                     string // speciale tweak zonder register (power, hibernate)
	Vals                   []RegVal
}

const (
	advKey  = `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
	cdmKey  = `Software\Microsoft\Windows\CurrentVersion\ContentDeliveryManager`
	themKey = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	accKey  = `Control Panel\Accessibility`
	relogEN = "Takes effect after signing in again."
	relogNL = "Werkt na opnieuw inloggen."
	rebEN   = "Takes effect after a restart."
	rebNL   = "Werkt na een herstart van de pc."
)

var tweaks = []Tweak{
	// ---- uiterlijk
	{ID: "dark", Group: "looks", NameEN: "Dark mode", DescEN: "Windows and apps in dark mode.", NameNL: "Donkere modus", DescNL: "Windows en apps in het donker.", Vals: []RegVal{
		{"HKCU", themKey, "AppsUseLightTheme", "dword", 0, 1},
		{"HKCU", themKey, "SystemUsesLightTheme", "dword", 0, 1},
	}},
	{ID: "transp", Group: "looks", NameEN: "No transparency", DescEN: "Turns off see-through effects. A bit faster on older PCs.", NameNL: "Geen doorzichtige effecten", DescNL: "Zet doorzichtige effecten uit. Iets sneller op oudere pc's.", Vals: []RegVal{
		{"HKCU", themKey, "EnableTransparency", "dword", 0, 1},
	}},
	{ID: "accentbar", Group: "looks", NameEN: "Accent color on title bars", DescEN: "Window title bars get your accent color.", NameNL: "Accentkleur op titelbalken", DescNL: "Titelbalken van vensters krijgen je accentkleur.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\DWM`, "ColorPrevalence", "dword", 1, 0},
	}},
	{ID: "fastmenu", Group: "looks", NameEN: "Faster menus", DescEN: "Menus open right away instead of after a short wait.", NameNL: "Snellere menu's", DescNL: "Menu's openen meteen in plaats van na een korte pauze.", NoteEN: relogEN, NoteNL: relogNL, Vals: []RegVal{
		{"HKCU", `Control Panel\Desktop`, "MenuShowDelay", "sz", "0", "400"},
	}},
	// ---- taakbalk
	{ID: "taskleft", Group: "taskbar", NameEN: "Start button on the left", DescEN: "Taskbar icons on the left, like Windows 10.", NameNL: "Startknop links", DescNL: "Taakbalk-iconen links, zoals in Windows 10.", Vals: []RegVal{
		{"HKCU", advKey, "TaskbarAl", "dword", 0, 1},
	}},
	{ID: "endtask", Group: "taskbar", NameEN: "'End task' in the taskbar", DescEN: "Right-click an app in the taskbar to close it right away.", NameNL: "'Taak beëindigen' in de taakbalk", DescNL: "Rechtsklik op een app in de taakbalk om hem meteen te sluiten.", Vals: []RegVal{
		{"HKCU", advKey + `\TaskbarDeveloperSettings`, "TaskbarEndTask", "dword", 1, 0},
	}},
	{ID: "seconds", Group: "taskbar", NameEN: "Seconds in the clock", DescEN: "The taskbar clock also shows seconds.", NameNL: "Seconden in de klok", DescNL: "De klok in de taakbalk laat ook seconden zien.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "ShowSecondsInSystemClock", "dword", 1, 0},
	}},
	{ID: "taskview", Group: "taskbar", NameEN: "Hide Task View button", DescEN: "Removes the Task View button from the taskbar.", NameNL: "Taakweergave-knop weg", DescNL: "Haalt de knop Taakweergave uit de taakbalk.", Vals: []RegVal{
		{"HKCU", advKey, "ShowTaskViewButton", "dword", 0, 1},
	}},
	{ID: "searchicon", Group: "taskbar", NameEN: "Search as icon only", DescEN: "A small search icon instead of the big search box.", NameNL: "Zoeken als klein icoon", DescNL: "Een klein zoek-icoon in plaats van de grote zoekbalk.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Search`, "SearchboxTaskbarMode", "dword", 1, 2},
	}},
	{ID: "widgets", Group: "taskbar", NameEN: "Widgets off", DescEN: "No news and weather widgets in the taskbar.", NameNL: "Widgets uit", DescNL: "Geen nieuws- en weer-widgets in de taakbalk.", Admin: true, Explorer: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Dsh`, "AllowNewsAndInterests", "dword", 0, nil},
	}},
	// ---- verkenner
	{ID: "ext", Group: "explorer", NameEN: "Show file extensions", DescEN: "See .exe, .png and .zip after file names.", NameNL: "Bestandsextensies tonen", DescNL: "Zie .exe, .png en .zip achter bestandsnamen.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "HideFileExt", "dword", 0, 1},
	}},
	{ID: "hidden", Group: "explorer", NameEN: "Show hidden files", DescEN: "Shows hidden folders and files.", NameNL: "Verborgen bestanden tonen", DescNL: "Laat verborgen mappen en bestanden zien.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "Hidden", "dword", 1, 2},
	}},
	{ID: "classicmenu", Group: "explorer", NameEN: "Classic right-click menu", DescEN: "The full menu right away, without 'Show more options'.", NameNL: "Oud rechtsklikmenu", DescNL: "Het volledige menu meteen, zonder 'Meer opties weergeven'.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "", "key", "", nil},
	}},
	{ID: "thispc", Group: "explorer", NameEN: "Explorer opens 'This PC'", DescEN: "Explorer starts at your drives instead of Home.", NameNL: "Verkenner opent 'Deze pc'", DescNL: "Verkenner start bij je schijven in plaats van bij Start.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "LaunchTo", "dword", 1, nil},
	}},
	{ID: "compact", Group: "explorer", NameEN: "Compact view", DescEN: "Less space between files, so you see more at once.", NameNL: "Compacte weergave", DescNL: "Minder ruimte tussen bestanden, je ziet er meer tegelijk.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "UseCompactMode", "dword", 1, 0},
	}},
	{ID: "fullpath", Group: "explorer", NameEN: "Full path in the title", DescEN: "Shows the whole folder path at the top of Explorer.", NameNL: "Volledig pad in de titel", DescNL: "Laat het hele pad van de map zien bovenin Verkenner.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\CabinetState`, "FullPath", "dword", 1, 0},
	}},
	{ID: "norecent", Group: "explorer", NameEN: "No recent files in Explorer", DescEN: "Home no longer shows recent and frequent files.", NameNL: "Geen recente bestanden", DescNL: "Start in Verkenner toont geen recente en vaak gebruikte bestanden meer.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer`, "ShowRecent", "dword", 0, 1},
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer`, "ShowFrequent", "dword", 0, 1},
	}},
	{ID: "shortcut", Group: "explorer", NameEN: "No '- Shortcut' in names", DescEN: "New shortcuts don't get '- Shortcut' added to the name.", NameNL: "Geen '- Snelkoppeling' in namen", DescNL: "Nieuwe snelkoppelingen krijgen geen '- Snelkoppeling' achter de naam.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplates`, "ShortcutNameTemplate", "sz", "%s.lnk", nil},
	}},
	// ---- gaming
	{ID: "mouse", Group: "gaming", NameEN: "Mouse acceleration off", DescEN: "Your mouse moves 1-to-1. Better aim in games.", NameNL: "Muisversnelling uit", DescNL: "Je muis beweegt 1-op-1. Beter richten in games.", NoteEN: relogEN, NoteNL: relogNL, Vals: []RegVal{
		{"HKCU", `Control Panel\Mouse`, "MouseSpeed", "sz", "0", "1"},
		{"HKCU", `Control Panel\Mouse`, "MouseThreshold1", "sz", "0", "6"},
		{"HKCU", `Control Panel\Mouse`, "MouseThreshold2", "sz", "0", "10"},
	}},
	{ID: "sticky", Group: "gaming", NameEN: "No Sticky Keys pop-up", DescEN: "No pop-ups when you press Shift five times or hold keys.", NameNL: "Geen plaktoetsen-popup", DescNL: "Geen popups als je vijf keer op Shift drukt of toetsen vasthoudt.", Vals: []RegVal{
		{"HKCU", accKey + `\StickyKeys`, "Flags", "sz", "506", "510"},
		{"HKCU", accKey + `\Keyboard Response`, "Flags", "sz", "122", "126"},
		{"HKCU", accKey + `\ToggleKeys`, "Flags", "sz", "58", "62"},
	}},
	{ID: "gamemode", Group: "gaming", NameEN: "Game Mode", DescEN: "Windows gives games priority.", NameNL: "Game Mode", DescNL: "Windows geeft games voorrang.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\GameBar`, "AutoGameModeEnabled", "dword", 1, 0},
	}},
	{ID: "gamedvr", Group: "gaming", NameEN: "No background recording", DescEN: "Turns off Xbox Game DVR, which records in the background and costs FPS.", NameNL: "Geen opnemen op de achtergrond", DescNL: "Zet Xbox Game DVR uit. Die neemt stiekem op en kost FPS.", Vals: []RegVal{
		{"HKCU", `System\GameConfigStore`, "GameDVR_Enabled", "dword", 0, 1},
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\GameDVR`, "AppCaptureEnabled", "dword", 0, 1},
	}},
	{ID: "windowed", Group: "gaming", NameEN: "Optimize windowed games", DescEN: "Lower input lag in games that run windowed or borderless.", NameNL: "Games in venster sneller", DescNL: "Minder vertraging in games die in een venster of randloos draaien.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\DirectX\UserGpuPreferences`, "DirectXUserGlobalSettings", "sz", "SwapEffectUpgradeEnable=1;", nil},
	}},
	{ID: "hags", Group: "gaming", NameEN: "GPU scheduling (HAGS)", DescEN: "Your graphics card manages its own work. Often a bit smoother.", NameNL: "GPU-planning (HAGS)", DescNL: "Je videokaart regelt zelf zijn werk. Vaak iets soepeler.", NoteEN: rebEN, NoteNL: rebNL, Admin: true, Vals: []RegVal{
		{"HKLM", `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, "HwSchMode", "dword", 2, 1},
	}},
	{ID: "netthrottle", Group: "gaming", NameEN: "No network throttling", DescEN: "Windows stops slowing down your network while media plays.", NameNL: "Netwerk niet afremmen", DescNL: "Windows remt je netwerk niet meer af als er media speelt.", NoteEN: rebEN, NoteNL: rebNL, Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile`, "NetworkThrottlingIndex", "dword", 0xffffffff, 10},
		{"HKLM", `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile`, "SystemResponsiveness", "dword", 10, 20},
	}},
	// ---- snelheid
	{ID: "startup", Group: "speed", NameEN: "Faster startup apps", DescEN: "No waiting time for startup apps after signing in.", NameNL: "Sneller opstarten van apps", DescNL: "Geen wachttijd voor opstart-apps na het inloggen.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\Serialize`, "StartupDelayInMSec", "dword", 0, nil},
	}},
	{ID: "bgapps", Group: "speed", NameEN: "No background apps", DescEN: "Store apps stop running in the background.", NameNL: "Geen apps op de achtergrond", DescNL: "Store-apps draaien niet meer stiekem op de achtergrond.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\BackgroundAccessApplications`, "GlobalUserDisabled", "dword", 1, 0},
	}},
	{ID: "autoend", Group: "speed", NameEN: "Close frozen apps on shutdown", DescEN: "Shutting down doesn't wait for apps that don't respond.", NameNL: "Vastgelopen apps sluiten bij afsluiten", DescNL: "Afsluiten wacht niet op apps die niet meer reageren.", Vals: []RegVal{
		{"HKCU", `Control Panel\Desktop`, "AutoEndTasks", "sz", "1", nil},
	}},
	// ---- privacy
	{ID: "ads", Group: "privacy", NameEN: "No ads and tips", DescEN: "No suggested apps, tips and ads in Start and Settings.", NameNL: "Geen reclame en tips", DescNL: "Geen voorgestelde apps, tips en reclame in Start en Instellingen.", Vals: []RegVal{
		{"HKCU", cdmKey, "SubscribedContent-338388Enabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SubscribedContent-338389Enabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SubscribedContent-353694Enabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SubscribedContent-353696Enabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SystemPaneSuggestionsEnabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SilentInstalledAppsEnabled", "dword", 0, 1},
	}},
	{ID: "lockads", Group: "privacy", NameEN: "No tips on the lock screen", DescEN: "No fun facts and ads on your lock screen.", NameNL: "Geen tips op het vergrendelscherm", DescNL: "Geen weetjes en reclame op je vergrendelscherm.", Vals: []RegVal{
		{"HKCU", cdmKey, "RotatingLockScreenOverlayEnabled", "dword", 0, 1},
		{"HKCU", cdmKey, "SubscribedContent-338387Enabled", "dword", 0, 1},
	}},
	{ID: "startrec", Group: "privacy", NameEN: "No recommendations in Start", DescEN: "Start no longer shows recent files and suggestions.", NameNL: "Geen aanbevelingen in Start", DescNL: "Start laat geen recente bestanden en suggesties meer zien.", Vals: []RegVal{
		{"HKCU", advKey, "Start_TrackDocs", "dword", 0, 1},
		{"HKCU", advKey, "Start_IrisRecommendations", "dword", 0, 1},
	}},
	{ID: "adid", Group: "privacy", NameEN: "Advertising ID off", DescEN: "Apps can't use your ID for personal ads.", NameNL: "Reclame-ID uit", DescNL: "Apps kunnen je ID niet gebruiken voor persoonlijke reclame.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\AdvertisingInfo`, "Enabled", "dword", 0, 1},
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Privacy`, "TailoredExperiencesWithDiagnosticDataEnabled", "dword", 0, 1},
	}},
	{ID: "feedback", Group: "privacy", NameEN: "No feedback questions", DescEN: "Windows stops asking for your opinion.", NameNL: "Geen feedback-vragen", DescNL: "Windows vraagt niet meer om je mening.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Siuf\Rules`, "NumberOfSIUFInPeriod", "dword", 0, nil},
	}},
	{ID: "bing", Group: "privacy", NameEN: "No Bing in Start search", DescEN: "Start only searches your PC, not the internet.", NameNL: "Bing uit in Start-zoeken", DescNL: "Start zoekt alleen op je pc, niet op internet.", Admin: true, Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Policies\Microsoft\Windows\Explorer`, "DisableSearchBoxSuggestions", "dword", 1, nil},
	}},
	{ID: "copilot", Group: "privacy", NameEN: "Copilot off", DescEN: "Turns off Windows Copilot.", NameNL: "Copilot uit", DescNL: "Zet Windows Copilot uit.", Admin: true, Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Policies\Microsoft\Windows\WindowsCopilot`, "TurnOffWindowsCopilot", "dword", 1, nil},
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\WindowsCopilot`, "TurnOffWindowsCopilot", "dword", 1, nil},
	}},
	{ID: "recall", Group: "privacy", NameEN: "Recall off", DescEN: "Windows doesn't take screenshots of what you do.", NameNL: "Recall uit", DescNL: "Windows maakt geen screenshots van wat je doet.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\WindowsAI`, "DisableAIDataAnalysis", "dword", 1, nil},
	}},
	{ID: "activity", Group: "privacy", NameEN: "Activity history off", DescEN: "Windows doesn't keep a list of what you did.", NameNL: "Activiteitengeschiedenis uit", DescNL: "Windows houdt geen lijst bij van wat je deed.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\System`, "PublishUserActivities", "dword", 0, nil},
	}},
	{ID: "telemetry", Group: "privacy", NameEN: "Minimal telemetry", DescEN: "Windows sends as little data to Microsoft as possible.", NameNL: "Telemetrie minimaal", DescNL: "Windows stuurt zo min mogelijk gegevens naar Microsoft.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\DataCollection`, "AllowTelemetry", "dword", 0, nil},
	}},
	// ---- extra (v1.3)
	{ID: "pcicon", Group: "looks", NameEN: "'This PC' on the desktop", DescEN: "Puts the This PC icon on your desktop.", NameNL: "'Deze pc' op het bureaublad", DescNL: "Zet het icoon Deze pc op je bureaublad.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\HideDesktopIcons\NewStartPanel`, "{20D04FE0-3AEA-1069-A2D8-08002B30309D}", "dword", 0, 1},
	}},
	{ID: "noacrylic", Group: "looks", NameEN: "Sharp sign-in screen", DescEN: "No blur on the sign-in screen background.", NameNL: "Scherp inlogscherm", DescNL: "Geen wazige achtergrond op het inlogscherm.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\System`, "DisableAcrylicBackgroundOnLogon", "dword", 1, nil},
	}},
	{ID: "chaticon", Group: "taskbar", NameEN: "Hide Chat button", DescEN: "Removes the Teams Chat button from the taskbar.", NameNL: "Chat-knop weg", DescNL: "Haalt de Teams Chat-knop uit de taakbalk.", Vals: []RegVal{
		{"HKCU", advKey, "TaskbarMn", "dword", 0, 1},
	}},
	{ID: "alttab", Group: "taskbar", NameEN: "Alt+Tab without browser tabs", DescEN: "Alt+Tab only shows windows, not every Edge tab.", NameNL: "Alt+Tab zonder browsertabbladen", DescNL: "Alt+Tab laat alleen vensters zien, niet elk Edge-tabblad.", Vals: []RegVal{
		{"HKCU", advKey, "MultiTaskingAltTabFilter", "dword", 3, nil},
	}},
	{ID: "noshake", Group: "taskbar", NameEN: "No 'shake to minimize'", DescEN: "Shaking a window no longer minimizes all other windows.", NameNL: "Geen 'schudden = minimaliseren'", DescNL: "Een venster schudden minimaliseert niet meer alle andere vensters.", Vals: []RegVal{
		{"HKCU", advKey, "DisallowShaking", "dword", 1, 0},
	}},
	{ID: "nogallery", Group: "explorer", NameEN: "Hide Gallery", DescEN: "Removes Gallery from the Explorer side bar.", NameNL: "Galerij verbergen", DescNL: "Haalt Galerij uit de zijbalk van Verkenner.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree", "dword", 0, nil},
	}},
	{ID: "numlock", Group: "system", NameEN: "Num Lock on at start", DescEN: "The number pad works right away after signing in.", NameNL: "Num Lock aan bij opstarten", DescNL: "Het numerieke toetsenblok werkt meteen na het inloggen.", NoteEN: relogEN, NoteNL: relogNL, Vals: []RegVal{
		{"HKCU", `Control Panel\Keyboard`, "InitialKeyboardIndicators", "sz", "2", "0"},
	}},
	{ID: "noreboot", Group: "system", NameEN: "No forced restarts for updates", DescEN: "Windows Update won't restart your PC while you're signed in.", NameNL: "Geen gedwongen herstart voor updates", DescNL: "Windows Update herstart je pc niet zolang je bent ingelogd.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU`, "NoAutoRebootWithLoggedOnUsers", "dword", 1, nil},
	}},
	{ID: "longpaths", Group: "system", NameEN: "Long file paths", DescEN: "Allows file paths longer than 260 characters.", NameNL: "Lange bestandspaden", DescNL: "Staat bestandspaden langer dan 260 tekens toe.", Admin: true, Vals: []RegVal{
		{"HKLM", `SYSTEM\CurrentControlSet\Control\FileSystem`, "LongPathsEnabled", "dword", 1, 0},
	}},
	{ID: "fso", Group: "gaming", NameEN: "Fullscreen optimizations off", DescEN: "Real exclusive fullscreen in older games.", NameNL: "Volledig-scherm-optimalisaties uit", DescNL: "Echt exclusief volledig scherm in oudere games.", Vals: []RegVal{
		{"HKCU", `System\GameConfigStore`, "GameDVR_FSEBehaviorMode", "dword", 2, 0},
		{"HKCU", `System\GameConfigStore`, "GameDVR_HonorUserFSEBehaviorMode", "dword", 1, 0},
		{"HKCU", `System\GameConfigStore`, "GameDVR_DXGIHonorFSEWindowsCompatible", "dword", 1, 0},
	}},
	{ID: "power", Group: "speed", NameEN: "Ultimate Performance power plan", DescEN: "Your CPU never slows down to save power.", NameNL: "Energieplan Ultieme prestaties", DescNL: "Je processor remt nooit af om stroom te besparen.", Fn: "power"},
	{ID: "fastboot", Group: "speed", NameEN: "Fast Startup off", DescEN: "Shutting down really shuts down. Fixes many driver and update problems.", NameNL: "Snel opstarten uit", DescNL: "Afsluiten sluit echt af. Lost veel driver- en updateproblemen op.", Admin: true, Vals: []RegVal{
		{"HKLM", `SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled", "dword", 0, 1},
	}},
	{ID: "hibernate", Group: "speed", NameEN: "Hibernation off", DescEN: "Frees several GB on your C: drive.", NameNL: "Slaapstand (hibernate) uit", DescNL: "Maakt een paar GB vrij op je C-schijf.", Admin: true, Fn: "hibernate"},
	{ID: "edgebg", Group: "speed", NameEN: "Edge doesn't run in the background", DescEN: "Edge stops starting with Windows and stays closed when you close it.", NameNL: "Edge niet op de achtergrond", DescNL: "Edge start niet meer met Windows mee en blijft dicht als je hem sluit.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Edge`, "StartupBoostEnabled", "dword", 0, nil},
		{"HKLM", `SOFTWARE\Policies\Microsoft\Edge`, "BackgroundModeEnabled", "dword", 0, nil},
	}},
	{ID: "doupload", Group: "speed", NameEN: "Don't share updates with other PCs", DescEN: "Windows stops uploading updates to strangers on the internet.", NameNL: "Updates niet delen met andere pc's", DescNL: "Windows uploadt geen updates meer naar vreemden op internet.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\DeliveryOptimization`, "DODownloadMode", "dword", 0, nil},
	}},
	{ID: "scoobe", Group: "privacy", NameEN: "No 'Finish setting up' nag", DescEN: "No more full-screen 'Let's finish setting up your PC' after updates.", NameNL: "Geen 'Instellen afronden'-scherm", DescNL: "Geen schermvullend 'Laten we je pc verder instellen' meer na updates.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\UserProfileEngagement`, "ScoobeSystemSettingEnabled", "dword", 0, 1},
	}},
	// ---- extra (v1.5)
	{ID: "nocombine", Group: "taskbar", NameEN: "Never combine taskbar buttons", DescEN: "Every window gets its own button with its name.", NameNL: "Taakbalkknoppen nooit samenvoegen", DescNL: "Elk venster krijgt een eigen knop met zijn naam.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "TaskbarGlomLevel", "dword", 2, 0},
	}},
	{ID: "copilotbtn", Group: "taskbar", NameEN: "Hide Copilot button", DescEN: "Removes the Copilot button from the taskbar.", NameNL: "Copilot-knop weg", DescNL: "Haalt de Copilot-knop uit de taakbalk.", Vals: []RegVal{
		{"HKCU", advKey, "ShowCopilotButton", "dword", 0, 1},
	}},
	{ID: "snapflyout", Group: "looks", NameEN: "No snap layouts pop-up", DescEN: "No layout menu when you hover over the maximize button.", NameNL: "Geen snap-indelingen-popup", DescNL: "Geen indelingenmenu als je over de maximaliseerknop gaat.", Vals: []RegVal{
		{"HKCU", advKey, "EnableSnapAssistFlyout", "dword", 0, 1},
	}},
	{ID: "sepproc", Group: "explorer", NameEN: "Folders in a separate process", DescEN: "If an Explorer window crashes, your taskbar and desktop keep working.", NameNL: "Mappen in een apart proces", DescNL: "Crasht een Verkenner-venster, dan blijven je taakbalk en bureaublad werken.", Explorer: true, Vals: []RegVal{
		{"HKCU", advKey, "SeparateProcess", "dword", 1, 0},
	}},
	{ID: "gamebar", Group: "gaming", NameEN: "Xbox Game Bar off", DescEN: "Windows + G and the controller button no longer open the Game Bar.", NameNL: "Xbox Game Bar uit", DescNL: "Windows + G en de controllerknop openen de Game Bar niet meer.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\GameBar`, "UseNexusForGameBarEnabled", "dword", 0, 1},
		{"HKCU", `Software\Microsoft\GameBar`, "ShowStartupPanel", "dword", 0, 1},
	}},
	{ID: "keyrepeat", Group: "gaming", NameEN: "Faster key repeat", DescEN: "Holding a key repeats it sooner.", NameNL: "Sneller toetsen herhalen", DescNL: "Een toets vasthouden herhaalt sneller.", NoteEN: relogEN, NoteNL: relogNL, Vals: []RegVal{
		{"HKCU", `Control Panel\Keyboard`, "KeyboardDelay", "sz", "0", "1"},
		{"HKCU", `Control Panel\Keyboard`, "KeyboardSpeed", "sz", "31", "31"},
	}},
	{ID: "nolockscreen", Group: "system", NameEN: "Skip the lock screen", DescEN: "Go straight to the sign-in screen, without the picture first.", NameNL: "Vergrendelscherm overslaan", DescNL: "Meteen naar het inlogscherm, zonder eerst de foto.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\Personalization`, "NoLockScreen", "dword", 1, nil},
	}},
	{ID: "trackprogs", Group: "privacy", NameEN: "Don't track app launches", DescEN: "Windows stops keeping a list of which apps you open most.", NameNL: "App-gebruik niet bijhouden", DescNL: "Windows houdt niet meer bij welke apps je het vaakst opent.", Vals: []RegVal{
		{"HKCU", advKey, "Start_TrackProgs", "dword", 0, 1},
	}},
	{ID: "inking", Group: "privacy", NameEN: "Typing data stays on your PC", DescEN: "Windows stops sending what you type and write to Microsoft.", NameNL: "Typgegevens blijven op je pc", DescNL: "Windows stuurt niet meer naar Microsoft wat je typt en schrijft.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Input\TIPC`, "Enabled", "dword", 0, 1},
		{"HKCU", `Software\Microsoft\InputPersonalization`, "RestrictImplicitTextCollection", "dword", 1, 0},
		{"HKCU", `Software\Microsoft\InputPersonalization`, "RestrictImplicitInkCollection", "dword", 1, 0},
	}},
	{ID: "acctnotif", Group: "privacy", NameEN: "No account nags in Start", DescEN: "No 'back up your PC' or 'sign in' reminders in the Start menu.", NameNL: "Geen account-meldingen in Start", DescNL: "Geen 'maak een back-up' of 'log in'-herinneringen in het Startmenu.", Vals: []RegVal{
		{"HKCU", advKey, "Start_AccountNotifications", "dword", 0, 1},
	}},
	{ID: "edgesidebar", Group: "privacy", NameEN: "No sidebar in Edge", DescEN: "Removes the Copilot sidebar from Microsoft Edge.", NameNL: "Geen zijbalk in Edge", DescNL: "Haalt de Copilot-zijbalk uit Microsoft Edge.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Edge`, "HubsSidebarEnabled", "dword", 0, nil},
	}},
	// ---- extra (v1.9)
	{ID: "homenav", Group: "explorer", NameEN: "Hide 'Home' in Explorer", DescEN: "Removes Home from the Explorer side bar.", NameNL: "'Start' verbergen in Verkenner", DescNL: "Haalt Start uit de zijbalk van Verkenner.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree", "dword", 0, nil},
	}},
	{ID: "onedrivenav", Group: "explorer", NameEN: "Hide OneDrive in Explorer", DescEN: "Removes OneDrive from the Explorer side bar.", NameNL: "OneDrive verbergen in Verkenner", DescNL: "Haalt OneDrive uit de zijbalk van Verkenner.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Classes\CLSID\{018D5C66-4533-4307-9B53-224DE2ED1FE6}`, "System.IsPinnedToNameSpaceTree", "dword", 0, nil},
	}},
	{ID: "driveletters", Group: "explorer", NameEN: "Drive letters first", DescEN: "Shows (C:) Windows instead of Windows (C:).", NameNL: "Stationsletter eerst", DescNL: "Laat (C:) Windows zien in plaats van Windows (C:).", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer`, "ShowDriveLettersFirst", "dword", 4, nil},
	}},
	{ID: "copydetails", Group: "explorer", NameEN: "Always show copy details", DescEN: "The copy window always shows speed and a graph.", NameNL: "Altijd details bij kopiëren", DescNL: "Het kopieervenster laat altijd snelheid en een grafiek zien.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\OperationStatusManager`, "EnthusiastMode", "dword", 1, nil},
	}},
	{ID: "folderdisc", Group: "explorer", NameEN: "Faster folder loading", DescEN: "Explorer stops guessing what kind of files are in a folder.", NameNL: "Mappen sneller laden", DescNL: "Verkenner raadt niet meer wat voor bestanden in een map staan.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags\AllFolders\Shell`, "FolderType", "sz", "NotSpecified", nil},
	}},
	{ID: "hovertime", Group: "taskbar", NameEN: "Faster taskbar previews", DescEN: "Window previews show up almost right away when you hover.", NameNL: "Snellere taakbalk-voorbeelden", DescNL: "Voorbeelden van vensters verschijnen bijna meteen als je eroverheen gaat.", NoteEN: relogEN, NoteNL: relogNL, Vals: []RegVal{
		{"HKCU", `Control Panel\Mouse`, "MouseHoverTime", "sz", "10", "400"},
	}},
	{ID: "morepins", Group: "taskbar", NameEN: "More pins in Start", DescEN: "More room for pinned apps, less for recommendations.", NameNL: "Meer vastgemaakte apps in Start", DescNL: "Meer ruimte voor vastgemaakte apps, minder voor aanbevelingen.", Vals: []RegVal{
		{"HKCU", advKey, "Start_Layout", "dword", 1, 0},
	}},
	{ID: "snapassist", Group: "looks", NameEN: "No snap suggestions", DescEN: "No list of other windows after you snap one to the side.", NameNL: "Geen snap-voorstellen", DescNL: "Geen lijst met andere vensters als je er een naar de zijkant snapt.", Vals: []RegVal{
		{"HKCU", advKey, "SnapAssist", "dword", 0, 1},
	}},
	{ID: "spotlighticon", Group: "looks", NameEN: "No 'Learn about this picture' icon", DescEN: "Removes the Spotlight icon from the desktop.", NameNL: "Geen 'Meer over deze foto'-icoon", DescNL: "Haalt het Spotlight-icoon van het bureaublad.", Explorer: true, Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\HideDesktopIcons\NewStartPanel`, "{2cc5ca98-6485-489a-920e-b3e88a6ccce3}", "dword", 1, 0},
	}},
	{ID: "smartclip", Group: "privacy", NameEN: "No suggested actions when copying", DescEN: "No pop-up when you copy a date or phone number.", NameNL: "Geen voorstellen bij kopiëren", DescNL: "Geen popup als je een datum of telefoonnummer kopieert.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\SmartActionPlatform\SmartClipboard`, "Disabled", "dword", 1, nil},
	}},
	{ID: "tipsnotif", Group: "privacy", NameEN: "No tips notifications", DescEN: "Windows stops sending 'tips and suggestions' notifications.", NameNL: "Geen tips-meldingen", DescNL: "Windows stuurt geen meldingen met tips en suggesties meer.", Vals: []RegVal{
		{"HKCU", `Software\Microsoft\Windows\CurrentVersion\Notifications\Settings\Windows.SystemToast.Suggested`, "Enabled", "dword", 0, nil},
	}},
	{ID: "werr", Group: "privacy", NameEN: "Error reports off", DescEN: "Windows stops sending crash reports to Microsoft.", NameNL: "Foutrapporten uit", DescNL: "Windows stuurt geen crashrapporten meer naar Microsoft.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\Windows Error Reporting`, "Disabled", "dword", 1, nil},
	}},
	{ID: "wudrivers", Group: "system", NameEN: "No drivers through Windows Update", DescEN: "Windows Update won't replace drivers you installed yourself.", NameNL: "Geen drivers via Windows Update", DescNL: "Windows Update vervangt drivers die je zelf installeerde niet meer.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`, "ExcludeWUDriversInQualityUpdate", "dword", 1, nil},
	}},
	{ID: "edgeshortcut", Group: "system", NameEN: "No new Edge shortcut on the desktop", DescEN: "Edge updates stop putting an icon back on your desktop.", NameNL: "Geen nieuwe Edge-snelkoppeling", DescNL: "Edge-updates zetten geen icoon meer terug op je bureaublad.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Policies\Microsoft\EdgeUpdate`, "CreateDesktopShortcutDefault", "dword", 0, nil},
	}},
	{ID: "verbose", Group: "system", NameEN: "Detailed startup messages", DescEN: "Shows what Windows is doing while starting and shutting down.", NameNL: "Uitgebreide opstartmeldingen", DescNL: "Laat zien wat Windows doet tijdens opstarten en afsluiten.", Admin: true, Vals: []RegVal{
		{"HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "VerboseStatus", "dword", 1, nil},
	}},
}

// Bloat is een ingebouwde app die weg mag. Pattern mag * bevatten.
type Bloat struct {
	NameEN, NameNL, Pattern, DescEN, DescNL string
}

var bloatList = []Bloat{
	{"News", "Nieuws", "Microsoft.BingNews", "News from Bing", "Nieuws van Bing"},
	{"Weather", "Weer", "Microsoft.BingWeather", "Weather app from Bing", "Weer-app van Bing"},
	{"Bing Search", "Bing Zoeken", "Microsoft.BingSearch", "Search through Bing", "Zoeken via Bing"},
	{"Tips", "Tips", "Microsoft.Getstarted", "Tips for Windows", "Tips voor Windows"},
	{"Get Help", "Hulp vragen", "Microsoft.GetHelp", "Microsoft's help app", "Hulp-app van Microsoft"},
	{"Solitaire", "Solitaire", "Microsoft.MicrosoftSolitaireCollection", "Card games with ads", "Kaartspellen met reclame"},
	{"Microsoft 365", "Microsoft 365", "Microsoft.MicrosoftOfficeHub", "Ad app for Office", "Reclame-app voor Office"},
	{"People", "Contacten", "Microsoft.People", "Old contacts app", "Oude contacten-app"},
	{"Power Automate", "Power Automate", "Microsoft.PowerAutomateDesktop", "Automate tasks", "Taken automatiseren"},
	{"To Do", "To Do", "Microsoft.Todos", "To-do list", "Takenlijstje"},
	{"Feedback Hub", "Feedback-hub", "Microsoft.WindowsFeedbackHub", "Send feedback to Microsoft", "Feedback sturen naar Microsoft"},
	{"Maps", "Kaarten", "Microsoft.WindowsMaps", "Maps app", "Kaarten-app"},
	{"Movies & TV", "Films en tv", "Microsoft.ZuneVideo", "Old video app", "Oude video-app"},
	{"Clipchamp", "Clipchamp", "Clipchamp.Clipchamp", "Video editor", "Video-editor"},
	{"Cortana", "Cortana", "Microsoft.549981C3F5F10", "Old voice assistant", "Oude spraakassistent"},
	{"Teams (free)", "Teams (gratis)", "MSTeams", "Teams for home", "Teams voor thuis"},
	{"Teams (old)", "Teams (oud)", "MicrosoftTeams", "Old Teams chat", "Oude Teams-chat"},
	{"New Outlook", "Nieuwe Outlook", "Microsoft.OutlookForWindows", "Mail app", "Mail-app"},
	{"Dev Home", "Dev Home", "Microsoft.Windows.DevHome", "Dashboard for developers", "Dashboard voor programmeurs"},
	{"Copilot", "Copilot", "Microsoft.Copilot", "Copilot app", "Copilot-app"},
	{"Mixed Reality", "Mixed Reality", "Microsoft.MixedReality.Portal", "For VR headsets", "Voor VR-brillen"},
	{"Skype", "Skype", "Microsoft.SkypeApp", "Calling app (discontinued)", "Bel-app (gestopt)"},
	{"3D Viewer", "3D Viewer", "Microsoft.Microsoft3DViewer", "View 3D models", "3D-modellen bekijken"},
	{"Paint 3D", "Paint 3D", "Microsoft.MSPaint", "Old 3D paint app", "Oude 3D-tekenapp"},
	{"Candy Crush etc.", "Candy Crush e.d.", "king.com.*", "Games with ads", "Spelletjes met reclame"},
	{"Disney+", "Disney+", "Disney.*", "Pre-installed", "Vooraf geïnstalleerd"},
	{"TikTok", "TikTok", "BytedancePte.Ltd.TikTok", "Pre-installed", "Vooraf geïnstalleerd"},
	{"Facebook / Instagram", "Facebook / Instagram", "Facebook.*", "Pre-installed", "Vooraf geïnstalleerd"},
	{"Prime Video", "Prime Video", "AmazonVideo.PrimeVideo", "Pre-installed", "Vooraf geïnstalleerd"},
}
