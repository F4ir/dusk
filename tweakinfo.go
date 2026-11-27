package main

// Uitgebreide uitleg per tweak (de "i"-knop): wat doet het precies, waarom, en waar let je op.
var tweakInfo = map[string][2]string{ // en, nl
	"dark": {
		"Switches Windows and every app that follows the system theme to dark colors. Easier on your eyes in the evening and on OLED screens it uses a bit less power. Apps with their own theme setting are not changed.",
		"Zet Windows en alle apps die het systeemthema volgen op donkere kleuren. Fijner voor je ogen in de avond, en op OLED-schermen kost het iets minder stroom. Apps met een eigen thema-instelling veranderen niet.",
	},
	"transp": {
		"Windows blurs and shows the background through the taskbar, Start and some windows. That blur costs a little GPU power and makes text slightly harder to read. Turning it off gives solid colors. You won't notice a speed difference on a strong PC, but older PCs and laptops on battery do.",
		"Windows maakt de taakbalk, Start en sommige vensters wazig doorzichtig. Dat vervagen kost een beetje kracht van je videokaart en maakt tekst iets minder goed leesbaar. Uit = effen kleuren. Op een snelle pc merk je geen verschil in snelheid, op oudere pc's en laptops op de accu wel.",
	},
	"accentbar": {
		"Normally window title bars are plain white or black. With this on they get your Windows accent color, so you see at a glance which window is active.",
		"Normaal zijn titelbalken van vensters gewoon wit of zwart. Met deze tweak krijgen ze je accentkleur van Windows. Zo zie je meteen welk venster actief is.",
	},
	"fastmenu": {
		"Windows waits 400 milliseconds before a sub-menu (like 'Open with' or 'Send to') opens. This sets that wait to 0, so menus pop open the moment you hover. Purely a feel thing: everything feels snappier.",
		"Windows wacht 400 milliseconden voordat een submenu opent (zoals 'Openen met' of 'Kopiëren naar'). Deze tweak zet dat op 0. Menu's springen dan meteen open als je eroverheen gaat. Puur gevoel: alles voelt vlotter.",
	},
	"taskleft": {
		"Windows 11 puts the Start button and app icons in the middle of the taskbar. This moves them to the left, like Windows 10, so the Start button is always in the same corner.",
		"Windows 11 zet de Startknop en de app-iconen in het midden van de taakbalk. Deze tweak zet ze links, zoals in Windows 10. Dan zit de Startknop altijd in dezelfde hoek.",
	},
	"endtask": {
		"Adds 'End task' to the right-click menu of apps in the taskbar. It kills a frozen app in one click, without opening Task Manager. Unsaved work in that app is lost, just like with Task Manager.",
		"Voegt 'Taak beëindigen' toe aan het rechtsklikmenu van apps in de taakbalk. Een vastgelopen app sluit je dan met één klik, zonder Taakbeheer. Werk dat niet is opgeslagen ben je kwijt, net als bij Taakbeheer.",
	},
	"seconds": {
		"The clock in the taskbar also shows seconds. Handy for timing things. It uses a tiny bit more power because the clock redraws every second.",
		"De klok in de taakbalk laat ook seconden zien. Handig om iets te timen. Kost een heel klein beetje meer stroom, omdat de klok elke seconde opnieuw getekend wordt.",
	},
	"taskview": {
		"Removes the Task View button (the two squares) from the taskbar. You can still open Task View with Windows + Tab.",
		"Haalt de knop Taakweergave (de twee vierkantjes) uit de taakbalk. Taakweergave open je nog steeds met Windows + Tab.",
	},
	"searchicon": {
		"Replaces the wide search box in the taskbar with a small magnifying-glass icon. Saves room for your apps. Searching still works the same, also by just pressing the Windows key and typing.",
		"Maakt van de brede zoekbalk in de taakbalk een klein vergrootglas. Dat geeft meer ruimte voor je apps. Zoeken werkt hetzelfde, ook door gewoon op de Windows-toets te drukken en te typen.",
	},
	"widgets": {
		"Turns off the Widgets board (news, weather, stocks) with a Windows policy. It disappears from the taskbar and stops loading news in the background. Needs admin because it is a PC-wide setting.",
		"Zet het Widgets-bord (nieuws, weer, aandelen) uit met een Windows-beleidsregel. Het verdwijnt uit de taakbalk en laadt geen nieuws meer op de achtergrond. Admin nodig, want het geldt voor de hele pc.",
	},
	"ext": {
		"By default Windows hides the end of file names like .exe or .pdf. With extensions shown you can see what a file really is. This protects you against tricks like 'photo.jpg.exe' that look like a picture but are a program.",
		"Standaard verbergt Windows het einde van bestandsnamen, zoals .exe of .pdf. Met extensies aan zie je wat een bestand echt is. Dat beschermt je tegen trucs zoals 'foto.jpg.exe', die eruitzien als een foto maar een programma zijn.",
	},
	"hidden": {
		"Shows files and folders that Windows normally hides, like the AppData folder where many games keep saves and mods. Hidden items look a bit faded so you can tell them apart.",
		"Laat bestanden en mappen zien die Windows normaal verbergt, zoals de map AppData waar veel games hun saves en mods bewaren. Verborgen dingen zien er iets vager uit, zodat je ze herkent.",
	},
	"classicmenu": {
		"Windows 11 shows a short right-click menu and hides the rest behind 'Show more options'. This brings back the full Windows 10 menu right away, so things like 7-Zip and 'Open with' are one click away. It works by registering an empty helper that makes Explorer skip the new menu.",
		"Windows 11 laat een kort rechtsklikmenu zien en verstopt de rest achter 'Meer opties weergeven'. Deze tweak geeft je meteen het volledige menu van Windows 10, dus 7-Zip en 'Openen met' zitten één klik dichterbij. Het werkt door een lege hulp te registreren, waardoor Verkenner het nieuwe menu overslaat.",
	},
	"thispc": {
		"Explorer normally opens on 'Home' with recent files. With this on it opens on 'This PC', showing your drives right away.",
		"Verkenner opent normaal op 'Start' met recente bestanden. Met deze tweak opent hij op 'Deze pc', met meteen je schijven.",
	},
	"compact": {
		"Makes the rows in Explorer less tall, so you see more files at once without scrolling. The touch-friendly spacing of Windows 11 is not needed with a mouse.",
		"Maakt de rijen in Verkenner minder hoog, zodat je meer bestanden tegelijk ziet zonder te scrollen. De ruime afstand van Windows 11 is bedoeld voor touchscreens, met een muis heb je die niet nodig.",
	},
	"fullpath": {
		"Shows the full folder path (like C:\\Users\\you\\Downloads) in the title of Explorer windows and in Alt+Tab, so you always know where you are.",
		"Laat het volledige pad van de map zien (zoals C:\\Users\\jij\\Downloads) in de titel van Verkenner en bij Alt+Tab. Zo weet je altijd waar je bent.",
	},
	"norecent": {
		"Explorer 'Home' lists files and folders you opened recently. With this off, that list stays empty. Good for privacy, for example when you share your screen.",
		"'Start' in Verkenner laat zien welke bestanden en mappen je kort geleden opende. Met deze tweak blijft die lijst leeg. Fijn voor je privacy, bijvoorbeeld als je je scherm deelt.",
	},
	"shortcut": {
		"When you create a shortcut, Windows adds '- Shortcut' to its name. This tweak stops that, so the shortcut just gets the name of the file.",
		"Als je een snelkoppeling maakt, zet Windows er '- Snelkoppeling' achter. Deze tweak stopt dat, dan krijgt de snelkoppeling gewoon de naam van het bestand.",
	},
	"mouse": {
		"Windows 'Enhance pointer precision' makes the cursor move further when you move the mouse faster. That feels smooth on a desktop, but in shooters it breaks muscle memory: the same hand movement gives a different distance. This turns it off, so the same movement always gives the same distance.",
		"'Verbeterde aanwijzerprecisie' van Windows laat de muis verder gaan als je sneller beweegt. Op het bureaublad voelt dat soepel, maar in shooters klopt je spiergeheugen dan niet: dezelfde beweging geeft een andere afstand. Deze tweak zet het uit, zodat dezelfde beweging altijd even ver gaat.",
	},
	"sticky": {
		"Pressing Shift five times, or holding Shift for 8 seconds, opens pop-ups for Sticky Keys and Filter Keys. In games that kicks you to the desktop. This turns off those keyboard shortcuts (not the features themselves if you need them).",
		"Vijf keer op Shift drukken, of Shift 8 seconden vasthouden, opent popups voor Plaktoetsen en Filtertoetsen. In een game kom je dan ineens op je bureaublad. Deze tweak zet die sneltoetsen uit (de functies zelf kun je nog gewoon gebruiken als je ze nodig hebt).",
	},
	"gamemode": {
		"When a game runs, Game Mode makes Windows pause Windows Update installs and driver installs, and gives the game more priority. It is on by default in Windows 11, but some tools turn it off.",
		"Als er een game draait, zorgt Game Mode dat Windows geen updates of drivers installeert en de game meer voorrang geeft. In Windows 11 staat het standaard aan, maar sommige programma's zetten het uit.",
	},
	"gamedvr": {
		"Xbox Game Bar can keep recording the last minutes of your game in the background, so you can save a clip afterwards. That recording costs FPS and GPU power all the time. If you use OBS, ShadowPlay or AMD ReLive instead, turn this off.",
		"De Xbox Game Bar kan op de achtergrond steeds de laatste minuten van je game opnemen, zodat je achteraf een clip kunt opslaan. Dat opnemen kost de hele tijd FPS en kracht van je videokaart. Gebruik je OBS, ShadowPlay of AMD ReLive, zet dit dan uit.",
	},
	"windowed": {
		"Games in windowed or borderless mode used to have more input lag than real fullscreen. This Windows 11 setting uses a newer way to show those games, which lowers the lag and allows Auto HDR and VRR.",
		"Games in een venster of randloos hadden vroeger meer vertraging dan echt volledig scherm. Deze instelling van Windows 11 laat die games op een nieuwere manier zien. Dat geeft minder vertraging en maakt Auto HDR en VRR mogelijk.",
	},
	"hags": {
		"Hardware-accelerated GPU scheduling lets your graphics card manage its own work queue instead of the CPU doing it. This can lower latency a bit and is needed for DLSS Frame Generation. On most modern cards it is good; if a game stutters after turning it on, turn it off again. A restart is needed.",
		"Met hardware-versnelde GPU-planning regelt je videokaart zelf de volgorde van zijn werk, in plaats van de processor. Dat kan de vertraging iets verlagen, en het is nodig voor DLSS Frame Generation. Op de meeste moderne kaarten is het goed. Hapert een game daarna, zet het dan weer uit. Herstart nodig.",
	},
	"netthrottle": {
		"While music or video plays, Windows limits how fast network packets are handled, so the media doesn't stutter. On modern PCs that limit is not needed and can add ping spikes in online games. This removes the limit and gives games a bit more CPU time. A restart is needed.",
		"Als er muziek of video speelt, remt Windows af hoe snel netwerkpakketjes verwerkt worden, zodat de media niet hapert. Op moderne pc's is dat niet nodig, en het kan pieken in je ping geven in online games. Deze tweak haalt die rem weg en geeft games iets meer processortijd. Herstart nodig.",
	},
	"startup": {
		"After you sign in, Windows waits about 10 seconds before it starts the apps in your startup list (Discord, Steam and so on). It does that to make the desktop appear faster. This tweak removes the wait, so your startup apps open right away. Your PC is ready to use sooner, but the first seconds may feel a bit busier.",
		"Na het inloggen wacht Windows ongeveer 10 seconden voordat hij de apps uit je opstartlijst start (Discord, Steam enzovoort). Dat doet hij zodat het bureaublad sneller verschijnt. Deze tweak haalt die wachttijd weg, dan openen je opstart-apps meteen. Je pc is eerder klaar voor gebruik, maar de eerste seconden kan het iets drukker voelen.",
	},
	"bgapps": {
		"Apps from the Microsoft Store (like Weather, Photos or Your Phone) can keep running in the background to update tiles and send notifications. That uses memory, CPU and battery. This stops them in the background. Normal programs like Discord and Steam are not affected.",
		"Apps uit de Microsoft Store (zoals Weer, Foto's of Telefoonkoppeling) mogen op de achtergrond blijven draaien om tegels te verversen en meldingen te sturen. Dat kost geheugen, processor en accu. Deze tweak stopt ze op de achtergrond. Gewone programma's zoals Discord en Steam merken hier niks van.",
	},
	"autoend": {
		"When you shut down or restart and an app hangs, Windows shows 'This app is preventing shutdown' and waits for you. With this on, Windows closes those apps by itself, so shutting down never gets stuck. Unsaved work in a frozen app is lost.",
		"Als je afsluit of herstart en een app hangt, laat Windows 'Deze app voorkomt het afsluiten' zien en wacht op jou. Met deze tweak sluit Windows die apps zelf, zodat afsluiten nooit blijft hangen. Werk dat niet is opgeslagen in een vastgelopen app ben je kwijt.",
	},
	"ads": {
		"Windows shows 'suggestions' in Start, Settings and on the lock screen. Those are ads for apps and Microsoft services, and Windows can even install promoted apps by itself. This turns all of that off.",
		"Windows laat 'suggesties' zien in Start, Instellingen en op het vergrendelscherm. Dat is reclame voor apps en diensten van Microsoft. Windows kan zelfs vanzelf aangeprezen apps installeren. Deze tweak zet dat allemaal uit.",
	},
	"lockads": {
		"The lock screen shows 'fun facts', tips and ads on top of the picture. This keeps just your picture or Spotlight image, without the text and ads.",
		"Het vergrendelscherm laat weetjes, tips en reclame zien over de foto heen. Met deze tweak houd je alleen je foto of Spotlight-plaatje, zonder tekst en reclame.",
	},
	"startrec": {
		"The Start menu has a 'Recommended' part with recently opened files and suggested apps. This empties it, so Start only shows your pinned apps. Nice for privacy when someone looks over your shoulder.",
		"Het Startmenu heeft een stuk 'Aanbevolen' met bestanden die je kort geleden opende en voorgestelde apps. Deze tweak maakt dat leeg, dan laat Start alleen je vastgezette apps zien. Fijn voor je privacy als iemand meekijkt.",
	},
	"adid": {
		"Windows gives you an advertising ID so apps can show ads that match what you do. Turning it off means apps get no ID, so ads are not personal. You don't see fewer ads, just less tracking.",
		"Windows geeft je een reclame-ID, zodat apps reclame kunnen laten zien die past bij wat jij doet. Uit = apps krijgen geen ID, dus de reclame is niet persoonlijk. Je ziet niet minder reclame, je wordt alleen minder gevolgd.",
	},
	"feedback": {
		"Every now and then Windows pops up asking for your opinion (Feedback Hub surveys). This sets the number of questions to zero.",
		"Af en toe vraagt Windows in een popup om je mening (enquêtes van de Feedback-hub). Deze tweak zet het aantal vragen op nul.",
	},
	"bing": {
		"When you type in Start, Windows also sends your search to Bing and shows web results between your apps and files. This makes Start search only your PC: faster, cleaner, and nothing you type goes to Microsoft. Needs admin because it uses a Windows policy.",
		"Als je typt in Start, stuurt Windows je zoekopdracht ook naar Bing en laat hij webresultaten zien tussen je apps en bestanden. Deze tweak laat Start alleen op je pc zoeken: sneller, overzichtelijker, en wat je typt gaat niet naar Microsoft. Admin nodig, want het is een beleidsregel van Windows.",
	},
	"copilot": {
		"Turns off the built-in Windows Copilot assistant with a Windows policy. The Copilot button and the Windows + C shortcut stop working. This does not remove the separate Copilot app; you can remove that under Clean up.",
		"Zet de ingebouwde Windows Copilot-assistent uit met een beleidsregel. De Copilot-knop en de sneltoets Windows + C werken dan niet meer. De losse Copilot-app wordt niet verwijderd, die kun je weghalen bij Opschonen.",
	},
	"recall": {
		"Recall is a feature on Copilot+ PCs that takes a screenshot of your screen every few seconds, so you can search back in what you did. This policy turns it off and stops those screenshots. On PCs without Recall it does nothing.",
		"Recall is een functie op Copilot+ pc's die om de paar seconden een screenshot van je scherm maakt, zodat je later kunt terugzoeken wat je deed. Deze beleidsregel zet het uit en stopt die screenshots. Op pc's zonder Recall doet het niks.",
	},
	"activity": {
		"Windows keeps a history of apps you used and files you opened (Timeline / activity history) and can send it to your Microsoft account. This policy stops Windows from collecting it.",
		"Windows houdt een geschiedenis bij van apps die je gebruikte en bestanden die je opende (activiteitengeschiedenis), en kan die naar je Microsoft-account sturen. Deze beleidsregel stopt het verzamelen.",
	},
	"telemetry": {
		"Windows sends diagnostic data to Microsoft about how your PC is used and when things crash. This policy sets it to the lowest level Windows allows on your edition. On Home and Pro the lowest level is 'Required', on Enterprise it is fully off.",
		"Windows stuurt diagnosegegevens naar Microsoft over hoe je pc gebruikt wordt en wanneer er iets crasht. Deze beleidsregel zet dat op het laagste niveau dat jouw versie van Windows toestaat. Bij Home en Pro is dat 'Vereist', bij Enterprise kan het helemaal uit.",
	},
	"pcicon": {
		"Puts the classic 'This PC' icon on your desktop, so you can open your drives with one double-click. It only shows a shortcut; nothing else changes.",
		"Zet het bekende icoon 'Deze pc' op je bureaublad. Zo open je je schijven met één dubbelklik. Het is alleen een snelkoppeling, verder verandert er niks.",
	},
	"noacrylic": {
		"Windows blurs your lock screen picture behind the sign-in box. This policy turns that blur off, so your wallpaper stays sharp when you sign in.",
		"Windows maakt je achtergrond wazig achter het inlogvak. Deze beleidsregel zet dat vervagen uit, zodat je achtergrond scherp blijft als je inlogt.",
	},
	"chaticon": {
		"Some Windows 11 versions have a Teams 'Chat' button in the taskbar. Most people never use it. This hides it.",
		"Sommige versies van Windows 11 hebben een Teams-knop 'Chat' in de taakbalk. De meeste mensen gebruiken die nooit. Deze tweak verbergt hem.",
	},
	"alttab": {
		"By default Alt+Tab also shows your last Microsoft Edge tabs as separate items, which makes the list long and messy. This makes Alt+Tab show only real windows again.",
		"Standaard laat Alt+Tab ook je laatste Edge-tabbladen zien als losse vakjes. Daardoor wordt de lijst lang en rommelig. Deze tweak laat Alt+Tab weer alleen echte vensters zien.",
	},
	"noshake": {
		"If you grab a window by the title bar and shake it, Windows minimizes all other windows ('Aero Shake'). That often happens by accident while dragging. This turns it off.",
		"Als je een venster bij de titelbalk pakt en schudt, minimaliseert Windows alle andere vensters ('Aero Shake'). Dat gebeurt vaak per ongeluk tijdens het slepen. Deze tweak zet het uit.",
	},
	"nogallery": {
		"Windows 11 added 'Gallery' to the Explorer side bar, showing all your photos. If you don't use it, this removes it so the side bar is shorter.",
		"Windows 11 heeft 'Galerij' toegevoegd aan de zijbalk van Verkenner, met al je foto's. Gebruik je die niet, dan haalt deze tweak hem weg en wordt de zijbalk korter.",
	},
	"numlock": {
		"Num Lock is often off after you sign in, so the number pad moves the cursor instead of typing numbers. This makes Num Lock turn on automatically when you sign in.",
		"Num Lock staat na het inloggen vaak uit. Dan beweegt het numerieke toetsenblok de cursor in plaats van cijfers te typen. Deze tweak zet Num Lock vanzelf aan als je inlogt.",
	},
	"noreboot": {
		"Windows Update sometimes restarts your PC by itself to finish updates, even while you are gaming or have work open (outside your active hours). With this policy it waits until you restart yourself. Updates still install; you just choose when to restart.",
		"Windows Update herstart je pc soms vanzelf om updates af te maken, ook als je aan het gamen bent of werk open hebt (buiten je actieve uren). Met deze beleidsregel wacht hij tot je zelf herstart. Updates worden nog steeds geïnstalleerd, jij kiest alleen wanneer je herstart.",
	},
	"longpaths": {
		"Old Windows programs can't handle file paths longer than 260 characters. Deep folders, like modpacks, node_modules or long game folders, then give 'path too long' errors. This lets modern apps use longer paths.",
		"Oude Windows-programma's kunnen geen bestandspaden aan die langer zijn dan 260 tekens. Diepe mappen, zoals modpacks, node_modules of lange gamemappen, geven dan de fout 'pad te lang'. Met deze tweak kunnen moderne apps langere paden gebruiken.",
	},
	"fso": {
		"Windows 10 and 11 run most 'fullscreen' games as a borderless window behind the scenes ('fullscreen optimizations'). That allows quick Alt+Tab and overlays. Some older DirectX 9/11 games run smoother or with less input lag in real exclusive fullscreen. This tells Windows to respect real fullscreen again. If a game behaves worse, turn it off.",
		"Windows 10 en 11 laten de meeste 'volledig scherm'-games stiekem als randloos venster draaien ('volledig-scherm-optimalisaties'). Dat maakt snel Alt+Tab en overlays mogelijk. Sommige oudere DirectX 9/11-games lopen soepeler of met minder vertraging in echt exclusief volledig scherm. Deze tweak laat Windows echt volledig scherm weer respecteren. Doet een game het slechter, zet het dan weer uit.",
	},
	"power": {
		"The 'Balanced' power plan lets your CPU slow down and park cores when it thinks you don't need them, and waking them up costs a little time. 'Ultimate Performance' keeps everything at full speed all the time, for the smoothest frame times. It uses more power and makes the PC a bit warmer, so it is not for laptops on battery. Off = back to Balanced.",
		"Het energieplan 'Gebalanceerd' laat je processor afremmen en kernen in slaap zetten als Windows denkt dat je ze niet nodig hebt. Wakker worden kost dan een klein beetje tijd. 'Ultieme prestaties' houdt alles de hele tijd op volle snelheid, voor de soepelste frametijden. Het kost meer stroom en je pc wordt iets warmer, dus niet voor laptops op de accu. Uit = terug naar Gebalanceerd.",
	},
	"fastboot": {
		"'Fast Startup' does not really shut down Windows: it saves part of it to disk, like a half hibernation, so the next start is a few seconds faster. Because of that, drivers and updates are never fully reloaded, which causes odd problems (devices not working, updates that keep asking for a restart, issues with dual boot). With this off, shutting down is a real fresh start. On an SSD you barely notice the speed difference.",
		"'Snel opstarten' sluit Windows niet echt af: hij bewaart een deel op de schijf, een soort halve slaapstand, zodat de volgende start een paar seconden sneller is. Daardoor worden drivers en updates nooit helemaal opnieuw geladen. Dat geeft rare problemen (apparaten die niet werken, updates die steeds om een herstart vragen, gedoe met dual boot). Met deze tweak is afsluiten een echte frisse start. Op een SSD merk je het snelheidsverschil nauwelijks.",
	},
	"hibernate": {
		"Hibernation saves your memory to a hidden file (hiberfil.sys) on C:, which is often several GB (around 40% of your RAM). On a desktop you rarely use hibernation; normal sleep still works. Turning it off deletes that file and frees the space. It also turns off Fast Startup.",
		"Slaapstand (hibernate) bewaart je geheugen in een verborgen bestand (hiberfil.sys) op C:. Dat is vaak een paar GB (ongeveer 40% van je RAM). Op een desktop gebruik je hibernate bijna nooit, gewone slaapstand werkt nog steeds. Uitzetten verwijdert dat bestand en maakt de ruimte vrij. Snel opstarten gaat dan ook uit.",
	},
	"edgebg": {
		"Even if you never open it, Microsoft Edge starts with Windows ('Startup boost') and keeps running in the background after you close it. That costs memory and startup time. These policies stop both. Edge still works normally when you open it.",
		"Ook als je hem nooit opent, start Microsoft Edge mee met Windows ('Opstartversnelling') en blijft hij op de achtergrond draaien als je hem sluit. Dat kost geheugen en opstarttijd. Deze beleidsregels stoppen allebei. Edge werkt gewoon als je hem opent.",
	},
	"doupload": {
		"Windows 'Delivery Optimization' shares pieces of Windows updates with other PCs, also over the internet, using your upload speed. This makes Windows download updates only from Microsoft and stop uploading them to others. Good if your upload is slow or you game online.",
		"'Delivery Optimization' van Windows deelt stukjes van updates met andere pc's, ook via internet, en gebruikt daarvoor je uploadsnelheid. Deze tweak laat Windows updates alleen bij Microsoft halen en stopt het uploaden naar anderen. Fijn als je upload traag is of als je online gamet.",
	},
	"scoobe": {
		"After some updates Windows shows a full-screen page 'Let's finish setting up your device', pushing OneDrive, Microsoft 365 and Game Pass. This turns that page off.",
		"Na sommige updates laat Windows een schermvullende pagina zien: 'Laten we je apparaat verder instellen'. Die probeert je OneDrive, Microsoft 365 en Game Pass aan te smeren. Deze tweak zet die pagina uit.",
	},
	"nocombine": {
		"Windows 11 groups all windows of one app into one icon. With this on, every window gets its own button with its title, like older Windows. Handy when you have many Discord, Chrome or Explorer windows open.",
		"Windows 11 stopt alle vensters van één app in één icoon. Met deze tweak krijgt elk venster een eigen knop met zijn titel, zoals in oudere Windows. Handig als je veel Discord-, Chrome- of Verkenner-vensters open hebt.",
	},
	"copilotbtn": {
		"Some Windows 11 versions put a Copilot button next to the clock. This hides it. Copilot itself is not removed; use 'Copilot off' or Clean up for that.",
		"Sommige versies van Windows 11 zetten een Copilot-knop naast de klok. Deze tweak verbergt hem. Copilot zelf wordt niet weggehaald, gebruik daarvoor 'Copilot uit' of Opschonen.",
	},
	"snapflyout": {
		"When you move your mouse over the maximize button, Windows shows a small menu with window layouts. Many people find it gets in the way. This turns off that pop-up; snapping with Windows + arrow keys still works.",
		"Als je met de muis over de maximaliseerknop gaat, laat Windows een menuutje met vensterindelingen zien. Veel mensen vinden dat in de weg zitten. Deze tweak zet die popup uit. Snappen met Windows + pijltjes werkt nog gewoon.",
	},
	"sepproc": {
		"Normally all Explorer windows run in the same process as your taskbar and desktop. If one folder window crashes (for example because of a broken shell extension), the taskbar disappears too. With this on, folder windows run separately, so a crash only closes that window. It uses a little more memory.",
		"Normaal draaien alle Verkenner-vensters in hetzelfde proces als je taakbalk en bureaublad. Crasht één mapvenster (bijvoorbeeld door een kapotte extensie), dan verdwijnt je taakbalk ook. Met deze tweak draaien mapvensters apart, dus een crash sluit alleen dat venster. Kost een klein beetje meer geheugen.",
	},
	"gamebar": {
		"The Xbox Game Bar opens with Windows + G or the Xbox button on a controller. If you never use it, it often pops up by accident during games. This stops those shortcuts from opening it.",
		"De Xbox Game Bar opent met Windows + G of de Xbox-knop op een controller. Gebruik je hem nooit, dan springt hij vaak per ongeluk open tijdens het gamen. Deze tweak zorgt dat die sneltoetsen hem niet meer openen.",
	},
	"keyrepeat": {
		"When you hold a key, Windows waits a moment before repeating it. This sets that wait to the shortest. Nice for moving through text or menus quickly.",
		"Als je een toets vasthoudt, wacht Windows even voordat hij gaat herhalen. Deze tweak zet die wachttijd op het kortst. Fijn om snel door tekst of menu's te gaan.",
	},
	"nolockscreen": {
		"After starting up or locking your PC, Windows first shows a lock screen picture that you have to click or swipe away. This policy skips it, so you go straight to typing your PIN or password. Your PC is still locked just as safely.",
		"Na het opstarten of vergrendelen laat Windows eerst een vergrendelscherm met foto zien, die je weg moet klikken. Deze beleidsregel slaat dat over, zodat je meteen je pincode of wachtwoord kunt typen. Je pc is nog steeds net zo veilig vergrendeld.",
	},
	"trackprogs": {
		"Windows keeps a list of which apps you start, to show 'most used' in Start and search. With this off it stops keeping that list. Start still works the same.",
		"Windows houdt bij welke apps je opent, om 'meest gebruikt' te laten zien in Start en zoeken. Met deze tweak houdt hij die lijst niet meer bij. Start werkt verder hetzelfde.",
	},
	"inking": {
		"To 'improve' typing suggestions and handwriting, Windows can send samples of what you type and write to Microsoft. This turns that off. Spelling correction on your own PC keeps working.",
		"Om suggesties bij typen en handschrift te 'verbeteren', kan Windows voorbeelden van wat je typt en schrijft naar Microsoft sturen. Deze tweak zet dat uit. Spellingcontrole op je eigen pc blijft werken.",
	},
	"acctnotif": {
		"The Start menu shows reminders near your profile picture, like 'back up your files' or 'sign in to your Microsoft account'. This hides them.",
		"Het Startmenu laat herinneringen zien bij je profielfoto, zoals 'maak een back-up van je bestanden' of 'log in met je Microsoft-account'. Deze tweak verbergt ze.",
	},
	"edgesidebar": {
		"Microsoft Edge has a sidebar with Copilot and shortcuts on the right. This policy turns it off, so pages get the full width.",
		"Microsoft Edge heeft rechts een zijbalk met Copilot en snelkoppelingen. Deze beleidsregel zet hem uit, zodat pagina's de volle breedte krijgen.",
	},
	"homenav": {
		"Explorer's side bar starts with 'Home', a page with recent and favourite files. If you go straight to your folders anyway, this hides it. The page itself still exists.",
		"De zijbalk van Verkenner begint met 'Start', een pagina met recente en favoriete bestanden. Ga je toch altijd meteen naar je mappen, dan verbergt deze tweak hem. De pagina zelf blijft bestaan.",
	},
	"onedrivenav": {
		"Hides OneDrive from the Explorer side bar. OneDrive keeps syncing and you can still open the folder from your user folder. Useful if you don't use OneDrive.",
		"Verbergt OneDrive in de zijbalk van Verkenner. OneDrive blijft gewoon synchroniseren en de map staat nog in je gebruikersmap. Handig als je OneDrive niet gebruikt.",
	},
	"driveletters": {
		"Puts the drive letter in front of the drive name in 'This PC', so you see (C:) and (D:) at a glance instead of at the end.",
		"Zet de stationsletter vóór de naam van de schijf in 'Deze pc'. Zo zie je (C:) en (D:) meteen, in plaats van aan het eind.",
	},
	"copydetails": {
		"When copying files Windows shows a small window with 'More details' folded in. This makes it always open, with transfer speed and the speed graph.",
		"Bij kopiëren laat Windows een klein venster zien met 'Meer details' ingeklapt. Met deze tweak staat het altijd open, met de snelheid en de grafiek.",
	},
	"folderdisc": {
		"Explorer looks at the files in every folder to guess if it's pictures, music or documents, and changes the columns to match. In big folders (like Downloads) that makes opening slow. This tells Explorer to treat every folder as general items, which loads faster. You can still change the view per folder.",
		"Verkenner kijkt naar de bestanden in elke map om te raden of het foto's, muziek of documenten zijn, en past de kolommen aan. In grote mappen (zoals Downloads) maakt dat openen traag. Deze tweak laat Verkenner elke map als gewone items behandelen, dat laadt sneller. Je kunt de weergave per map nog gewoon veranderen.",
	},
	"hovertime": {
		"Windows waits 400 ms before showing a window preview when your mouse rests on a taskbar button. This sets it to 10 ms, so previews appear almost right away. It also makes tooltips quicker.",
		"Windows wacht 400 ms voordat hij een voorbeeld van een venster laat zien als je muis op een taakbalkknop staat. Deze tweak zet dat op 10 ms, dan verschijnen voorbeelden bijna meteen. Ook tooltips komen sneller.",
	},
	"morepins": {
		"The Start menu is split between pinned apps and 'Recommended'. This layout gives more rows to your pinned apps and fewer to recommendations.",
		"Het Startmenu is verdeeld tussen vastgemaakte apps en 'Aanbevolen'. Deze indeling geeft meer rijen aan je vastgemaakte apps en minder aan aanbevelingen.",
	},
	"snapassist": {
		"After you snap a window to one side, Windows shows your other windows to fill the other half. If you find that annoying, this turns it off. Snapping itself keeps working.",
		"Als je een venster naar één kant snapt, laat Windows je andere vensters zien om de andere helft te vullen. Vind je dat vervelend, dan zet deze tweak het uit. Snappen zelf blijft werken.",
	},
	"spotlighticon": {
		"With Windows Spotlight as desktop background, Windows puts a 'Learn about this picture' icon on your desktop that you can't simply delete. This hides it.",
		"Met Windows Spotlight als bureaubladachtergrond zet Windows een icoon 'Meer informatie over deze foto' op je bureaublad dat je niet zomaar kunt verwijderen. Deze tweak verbergt het.",
	},
	"smartclip": {
		"When you copy a phone number or a date, Windows 11 can show a small pop-up with 'suggested actions' (call, add to calendar). This turns that off.",
		"Als je een telefoonnummer of datum kopieert, kan Windows 11 een popupje laten zien met 'voorgestelde acties' (bellen, aan agenda toevoegen). Deze tweak zet dat uit.",
	},
	"tipsnotif": {
		"Windows sends notifications like 'Get even more out of Windows' or tips about features. This turns off that kind of notification; other notifications keep working.",
		"Windows stuurt meldingen zoals 'Haal nog meer uit Windows' of tips over functies. Deze tweak zet dat soort meldingen uit; andere meldingen blijven werken.",
	},
	"werr": {
		"When a program crashes, Windows Error Reporting collects data and sends it to Microsoft. This policy turns that off. Crash dumps are no longer kept for this either, which also saves some disk space.",
		"Als een programma crasht, verzamelt Windows Foutrapportage gegevens en stuurt die naar Microsoft. Deze beleidsregel zet dat uit. Daarvoor worden ook geen crashbestanden meer bewaard, wat wat ruimte scheelt.",
	},
	"wudrivers": {
		"Windows Update also installs drivers, and sometimes replaces a newer driver from AMD, NVIDIA or your laptop maker with an older one. With this policy Windows Update leaves drivers alone; you update them yourself (for example in Dusk's Hardware tab). Windows updates themselves keep coming.",
		"Windows Update installeert ook drivers, en vervangt soms een nieuwere driver van AMD, NVIDIA of je laptopmerk door een oudere. Met deze beleidsregel laat Windows Update drivers met rust; die werk je zelf bij (bijvoorbeeld in het Hardware-tabblad van Dusk). Gewone Windows-updates blijven gewoon komen.",
	},
	"edgeshortcut": {
		"Every time Microsoft Edge updates, it can put its shortcut back on your desktop, even if you deleted it. This policy stops that.",
		"Elke keer als Microsoft Edge bijwerkt, kan hij zijn snelkoppeling terugzetten op je bureaublad, ook als je die had weggehaald. Deze beleidsregel stopt dat.",
	},
	"verbose": {
		"Instead of only 'Please wait', Windows shows what it is doing while starting, signing in and shutting down (for example 'Stopping services'). Handy to see where it gets stuck.",
		"In plaats van alleen 'Even geduld' laat Windows zien wat hij doet tijdens opstarten, inloggen en afsluiten (bijvoorbeeld 'Services stoppen'). Handig om te zien waar het blijft hangen.",
	},
}
