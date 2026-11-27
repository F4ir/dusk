// Nep-Go voor testen in een gewone browser.
(() => {
  const L = {};
  const emit = (n, d) => (L[n] || []).forEach(f => f(d));
  const POP = ['Discord.Discord','Valve.Steam','Google.Chrome','Logitech.GHUB','Microsoft.VCRedist.2015+.x64','Microsoft.VisualStudioCode'];
  const A = (Name, ID, Cat, Icon, Color, DescEN, inst, Version = '', Available = '') => ({ Name, ID, Cat, Icon, Color, DescEN, DescNL: DescEN + ' (nl)', Installed: !!inst, Version, Available, Pop: POP.indexOf(ID) + 1 });
  const apps = [
    A('Google Chrome','Google.Chrome','browser','googlechrome','#4285F4',"Google's web browser.",1,'157.0.8103.80','158.0.8127.37'),
    A('Firefox','Mozilla.Firefox','browser','firefoxbrowser','#FF7139','Fast browser that protects your privacy.',0),
    A('ChatGPT','OpenAI.ChatGPT','ai','openai','#10A37F','OpenAI\'s ChatGPT app for Windows.',0),
    A('Ollama','Ollama.Ollama','ai','ollama','#2B2B2B','Run AI models on your own PC.',0),
    A('Discord','Discord.Discord','chat','discord','#5865F2','Chat and call with friends.',1,'1.0.9269'),
    A('Claude','Anthropic.Claude','ai','claude','#D97757','The Claude app for your PC.',1,'6.28899.6.0'),
    A('Steam','Valve.Steam','games','steam','#1B2838','The biggest PC game store.',1,'2.10.91.91'),
    A('EA app','ElectronicArts.EADesktop','games','ea','#FF4747','Games from EA (FC, Battlefield, Apex).',0),
    A('Ubisoft Connect','Ubisoft.Connect','games','ubisoft','#0070FF','Games from Ubisoft.',0),
    A('Stream Deck','Elgato.StreamDeck','gear','elgato','#1A1A2E','Software for the Elgato Stream Deck.',1,'7.1.0'),
    A('Wave Link','Elgato.WaveLink','gear','elgato','#2A2350','Elgato audio mixer for your mic.',0),
    A('OBS Studio','OBSProject.OBSStudio','media','obsstudio','#302E31','Record and stream.',1,'32.2.2','33.0.3'),
    A('Bitwarden','Bitwarden.Bitwarden','tools','bitwarden','#175DDC','Password manager.',1,'2026.11.0'),
    A('HWiNFO','REALiX.HWiNFO','system','','#2F6DB5','Temperatures and all hardware info.',0),
    A('Logitech G HUB','Logitech.GHUB','gear','logitechg','#00B8FC','Settings for Logitech mice and keyboards.',1,'2026.7.1011216'),
    A('PuTTY','PuTTY.PuTTY','network','','#3A5A9B','SSH to your servers.',1,'0.85.0.0'),
    A('Advanced IP Scanner','Famatech.AdvancedIPScanner','network','','#1F8FD6','See every device on your network.',0),
    A('VS Code','Microsoft.VisualStudioCode','dev','visualstudiocode','#007ACC','Write code.',1,'1.142.0','1.146.0'),
    A('Java 21','EclipseAdoptium.Temurin.21.JRE','runtime','eclipseadoptium','#E21D66','Needed for Minecraft servers and Java apps.',0),
    A('Revo Uninstaller','RevoUninstaller.RevoUninstaller','tools','','#D9452B','Removes programs including leftovers.',0),
    A('Visual C++','Microsoft.VCRedist.2015+.x64','runtime','','#6A3BB5','Many games crash without this.',1,'14.44'),
  ];
  apps.find(a => a.ID === 'EclipseAdoptium.Temurin.21.JRE').Note = 'bundled:21, 25';
  const T = (ID, Group, NameEN, DescEN, Admin, Explorer, On, NoteEN = '') => ({ ID, Group, NameEN, DescEN, NoteEN, NameNL: NameEN + ' (nl)', DescNL: DescEN + ' (nl)', NoteNL: NoteEN, Admin, Explorer, On,
    InfoEN: 'After you sign in, Windows waits about 10 seconds before it starts the apps in your startup list. This tweak removes the wait.', InfoNL: 'Na het inloggen wacht Windows ongeveer 10 seconden voordat hij de opstart-apps start. Deze tweak haalt die wachttijd weg.',
    Changes: ['HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Serialize\\StartupDelayInMSec = 0   (off: (removed))'] });
  const tw = [T('dark','looks','Dark mode','Windows and apps in dark mode.',0,0,1), T('transp','looks','No transparency','Turns off see-through effects.',0,0,0),
    T('taskleft','taskbar','Start button on the left','Taskbar icons on the left, like Windows 10.',0,0,1), T('seconds','taskbar','Seconds in the clock','The taskbar clock also shows seconds.',0,1,0),
    T('ext','explorer','Show file extensions','See .exe, .png and .zip after file names.',0,1,1), T('compact','explorer','Compact view','Less space between files.',0,1,0),
    T('mouse','gaming','Mouse acceleration off','Your mouse moves 1-to-1.',0,0,0,'Takes effect after signing in again.'), T('hags','gaming','GPU scheduling (HAGS)','Your graphics card manages its own work.',1,0,0),
    T('startup','speed','Faster startup apps','No waiting time for startup apps.',0,0,0), T('ads','privacy','No ads and tips','No suggested apps, tips and ads.',0,0,0), T('recall','privacy','Recall off','No screenshots of what you do.',1,0,0)];
  let bloat = [['Weather','Weer','Microsoft.BingWeather'],['Clipchamp','Clipchamp','Clipchamp.Clipchamp'],['To Do','To Do','Microsoft.Todos']].map(b => ({ NameEN:b[0], NameNL:b[1], Pattern:b[2], DescEN:'Pre-installed', DescNL:'Vooraf geïnstalleerd' }));
  let drv = [{ ID:'702e', Title:'Microsoft Corporation AudioProcessingObject Driver Update (1.0.4.55794)', Class:'OtherHardware', Date:'2026-10-16', Maker:'Microsoft Corporation' }];
  const sleep = ms => new Promise(r => setTimeout(r, ms));
  window.runtime = { EventsOn:(n,f)=>(L[n]=L[n]||[]).push(f), WindowMinimise(){}, WindowToggleMaximise(){}, Quit(){}, BrowserOpenURL(u){ window.LASTURL=u; } };
  window.go = { main: { App: {
    SetLang: async () => {}, NoIntro: async () => !!window.MOCK_NOINTRO, SelfTest: async () => false,
    System: async () => ({ Windows:'Windows 11 Pro 25H2', Build:'26200', CPU:'AMD Ryzen 5 7600 6-Core Processor            ', PC:'DESKTOP-7Q2M4K', User:'Alex', Version:'1.9', RAM:16, Admin:false, Winget:true }),
    ScanApps: async () => { await sleep(window.MOCK_FAST ? 0 : 1500); return JSON.parse(JSON.stringify(apps)); },
    Tweaks: async () => JSON.parse(JSON.stringify(tw)),
    SetTweak: async (id, on) => { await sleep(200); const x = tw.find(y => y.ID === id); x.On = on; return { ...x }; },
    RestartExplorer: async () => {}, OpenDeviceManager: async () => {},
    InUse: async ids => { await sleep(800); return ids.includes('OBSProject.OBSStudio') && !window.CLOSED ? { 'OBSProject.OBSStudio': [{Name:'claude',Desc:'Claude',PIDs:[1]},{Name:'Discord',Desc:'Discord',PIDs:[2]},{Name:'msedge',Desc:'Microsoft Edge',PIDs:[3,4]}] } : {}; },
    CloseApps: async names => { window.CLOSED = names; await sleep(500); },
    Hardware: async () => { await sleep(300); if (window.MOCK_LAPTOP) return {
      CPU:'Intel(R) Core(TM) Ultra 7 155H', Board:'HP 8B41', BoardMaker:'HP', BoardProduct:'8B41', BIOS:'W70 Ver. 01.05.00', BIOSDate:'2025-03-10',
      SysMaker:'HP', SysModel:'HP OMEN Transcend 14 inch Gaming Laptop PC 14-fb0xxx', Serial:'5CD1234XYZ', Chassis:[10], OSBuild:26200, OSRelease:'25H2',
      OEM:{Brand:'hp',Maker:'HP',Model:'HP OMEN Transcend 14 inch Gaming Laptop PC 14-fb0xxx',Serial:'5CD1234XYZ',SystemID:'8B41',Laptop:true,SupportURL:'https://support.hp.com/drivers',Tool:'HP.SupportAssistant',ToolName:'HP Support Assistant',Deep:true},
      GPUs:[{Name:'Intel(R) Arc(TM) Graphics',Ver:'31.0.101.5186',Date:'2024-01-12',Vendor:'intel',Virtual:false},{Name:'NVIDIA GeForce RTX 4070 Laptop GPU',Ver:'32.0.15.6094',Date:'2024-08-14',Vendor:'nvidia',Virtual:false}],
      Disks:[{Name:'KIOXIA KXG80ZNV1T02',Type:'SSD',Health:'Healthy',Size:953}],
      RAM:[{Maker:'Micron',Part:'MT62F2G32D4DS-023',GB:32,Speed:7467,Rated:0}],
      Monitors:[{Name:'',Maker:'SDC'}], Batteries:[{Name:'Primary',Design:71000,Full:61200,Charge:84}], Software:[],
      Devices:[
        {Cat:'gpu',Name:'Intel(R) Arc(TM) Graphics',Maker:'Intel Corporation',Ver:'31.0.101.5186',Date:'2024-01-12',ID:'PCI\\VEN_8086&DEV_7D55',HWIDs:['PCI\\VEN_8086&DEV_7D55'],Problem:0},
        {Cat:'gpu',Name:'NVIDIA GeForce RTX 4070 Laptop GPU',Maker:'NVIDIA',Ver:'32.0.15.6094',Date:'2024-08-14',ID:'PCI\\VEN_10DE',HWIDs:['PCI\\VEN_10DE&DEV_2860'],Problem:0},
        {Cat:'chipset',Name:'Intel(R) Management Engine Interface #1',Maker:'Intel',Ver:'2406.5.5.0',Date:'2024-03-01',ID:'PCI\\VEN_8086&DEV_7E70',HWIDs:['PCI\\VEN_8086&DEV_7E70'],Problem:0},
        {Cat:'audio',Name:'Realtek(R) Audio',Maker:'Realtek Semiconductor Corp.',Ver:'6.0.9600.1',Date:'2023-11-02',ID:'HDAUDIO\\X',HWIDs:['HDAUDIO\\FUNC_01&VEN_10EC&DEV_0245'],Problem:0},
        {Cat:'wifi',Name:'Intel(R) Wi-Fi 7 BE200 320MHz',Maker:'Intel',Ver:'23.40.0.4',Date:'2024-05-01',ID:'PCI\\VEN_8086&DEV_272B',HWIDs:['PCI\\VEN_8086&DEV_272B'],Problem:0},
        {Cat:'bluetooth',Name:'Intel(R) Wireless Bluetooth(R)',Maker:'Intel Corporation',Ver:'23.40.0.2',Date:'2024-05-01',ID:'USB\\VID_8087',HWIDs:['USB\\VID_8087&PID_0036'],Problem:0},
        {Cat:'camera',Name:'HP 5MP Camera',Maker:'Microsoft',Ver:'10.0.26100.1',Date:'2006-06-21',ID:'USB\\VID_30C9',HWIDs:['USB\\VID_30C9&PID_00AF'],Problem:0,Generic:true},
        {Cat:'other',Name:'PCI Data Acquisition and Signal Processing Controller',Maker:'',Ver:'',Date:'',ID:'PCI\\VEN_8086&DEV_7D0B',HWIDs:['PCI\\VEN_8086&DEV_7D0B'],Problem:28},
      ] };
      return { CPU:'AMD Ryzen 5 7600 6-Core Processor            ', Board:'ASUSTeK COMPUTER INC. TUF GAMING B650-PLUS',
      GPUs:[{Name:'Parsec Virtual Display Adapter',Ver:'0.45.0.0',Date:'2024-01-25',Vendor:'',Virtual:true},{Name:'AMD Radeon RX 7700 XT',Ver:'32.0.32029.2010',Date:'2026-11-18',Vendor:'amd',Virtual:false}],
      Disks:[{Name:'Crucial P3 1TB',Type:'SSD',Health:'Healthy',Size:931},{Name:'WD Black SN770 2TB',Type:'SSD',Health:'Healthy',Size:1863}],
      BoardMaker:'ASUSTeK COMPUTER INC.', BoardProduct:'TUF GAMING B650-PLUS', BIOS:'3024', BIOSDate:'2026-06-02',
      SysMaker:'ASUS', SysModel:'System Product Name', Chassis:[3], OSBuild:26200, OSRelease:'25H2',
      OEM:{Brand:'custom',Maker:'Micro-Star International Co., Ltd.',Model:'ASUSTeK COMPUTER INC. TUF GAMING B650-PLUS',Laptop:false},
      RAM:[{Maker:'Unknown',Part:'KF560C36-16',GB:16,Speed:4800,Rated:6000},{Maker:'Unknown',Part:'KF560C36-16',GB:16,Speed:4800,Rated:6000}],
      Monitors:[{Name:'VG27AQ',Maker:'AUS'},{Name:'24G2',Maker:'AOC'}], Batteries:[], Software:[{Name:'AMD Software',Ver:'26.11.2'}],
      Devices:[
        {Cat:'gpu',Name:'AMD Radeon RX 7700 XT',Maker:'Advanced Micro Devices, Inc.',Ver:'32.0.32029.2010',Date:'2026-11-18',ID:'PCI\\VEN_1002',HWIDs:['PCI\\VEN_1002&DEV_747E'],Problem:0},
        {Cat:'chipset',Name:'AMD SMBus',Maker:'Advanced Micro Devices, Inc',Ver:'5.12.0.44',Date:'2025-09-09',ID:'PCI\\VEN_1022&DEV_790B',HWIDs:['PCI\\VEN_1022&DEV_790B'],Problem:0},
        {Cat:'chipset',Name:'AMD GPIO Controller',Maker:'Advanced Micro Devices, Inc',Ver:'2.2.0.137',Date:'2026-03-22',ID:'ACPI\\AMDI0030',HWIDs:['ACPI\\AMDI0030'],Problem:0},
        {Cat:'audio',Name:'Realtek High Definition Audio',Maker:'Realtek Semiconductor Corp.',Ver:'6.0.9998.1',Date:'2026-02-06',ID:'HDAUDIO\\X',HWIDs:['HDAUDIO\\FUNC_01&VEN_10EC&DEV_1220'],Problem:0},
        {Cat:'network',Name:'Realtek Gaming 2.5GbE Family Controller',Maker:'Realtek',Ver:'10.79.50.1003',Date:'2025-03-10',ID:'PCI\\VEN_10EC&DEV_8125',HWIDs:['PCI\\VEN_10EC&DEV_8125'],Problem:0},
        {Cat:'wifi',Name:'Intel(R) Wi-Fi 6E AX210 160MHz',Maker:'Intel',Ver:'23.30.0.6',Date:'2024-01-21',ID:'PCI\\VEN_8086&DEV_2723',HWIDs:['PCI\\VEN_8086&DEV_2723'],Problem:0},
        {Cat:'bluetooth',Name:'Intel(R) Wireless Bluetooth(R)',Maker:'Intel Corporation',Ver:'22.210.0.3',Date:'2023-08-02',ID:'USB\\VID_8087',HWIDs:['USB\\VID_8087&PID_0029'],Problem:0},
        {Cat:'usb',Name:'LIGHTSPEED Receiver',Maker:'Logitech',Ver:'1.1.91.2790',Date:'2026-06-17',ID:'USB\\VID_046D',HWIDs:['USB\\VID_046D&PID_C539'],Problem:0},
        {Cat:'other',Name:'USB Game Controller ',Maker:'',Ver:'',Date:'',ID:'USB\\VID_1234',HWIDs:['USB\\VID_1234&PID_0001'],Problem:28},
      ] }; },
    CheckDrivers: async gpus => { await sleep(window.MOCK_FAST ? 1500 : 4000); if (window.MOCK_LAPTOP) return {
        Updates:[{ID:'wu1',Title:'Intel - Display - 32.0.101.6130',Class:'Display',Date:'2024-10-15',Maker:'Intel',HWID:'PCI\\VEN_8086&DEV_7D55'}],
        GPU:{'NVIDIA GeForce RTX 4070 Laptop GPU':{Version:'621.24',URL:'https://www.nvidia.com/en-us/software/nvidia-app/',Download:'https://us.download.nvidia.com/x.exe',Newer:true,Known:true}},
        OEM:[{ID:'sp153001',Name:'HP BIOS and System Firmware (W70)',Category:'BIOS - System Firmware',Version:'01.12.00',Installed:'01.05.00',Date:'2026-07-01',URL:'x',Size:32000000,Device:'BIOS'},
             {ID:'sp152002',Name:'Realtek High-Definition (HD) Audio Driver',Category:'Driver - Audio',Version:'6.0.9700.1',Installed:'6.0.9600.1',Date:'2026-05-11',URL:'x',Size:410000000,Device:'Realtek(R) Audio'}] };
      return { Updates:[{ID:'702e',Title:'Microsoft Corporation AudioProcessingObject Driver Update (1.0.4.55794)',Class:'OtherHardware',Date:'2026-10-16',Maker:'Microsoft Corporation',HWID:'SWC\\X'}],
        GPU:{'AMD Radeon RX 7700 XT':{Version:'26.11.2',URL:'https://www.amd.com/en/support/download/drivers.html',Newer:false,Known:true}}, OEM:[] }; },
    Startup: async () => [{Key:'hklm|vpnui',Name:'vpnui',Desc:'VPN Client',Maker:'Example VPN Ltd.',Command:'x',Source:'hklm',Enabled:true,Admin:true},{Key:'hkcu|Steam',Name:'Steam',Desc:'Steam',Maker:'Valve Corporation',Command:'"C:\\Steam\\steam.exe" -silent',Source:'hkcu',Enabled:false,Admin:false},
      {Key:'hkcu|Discord',Name:'Discord',Desc:'Discord',Maker:'Discord Inc.',Command:'x',Source:'hkcu',Enabled:true,Admin:false},
      {Key:'hkcu|LGHUB',Name:'LGHUB',Desc:'G HUB',Maker:'Logitech, Inc.',Command:'x',Source:'hkcu',Enabled:true,Admin:false},
      {Key:'hklm|RtkAudUService',Name:'RtkAudUService',Desc:'Realtek HD Audio Universal Service',Maker:'Realtek Semiconductor',Command:'x',Source:'hklm',Enabled:true,Admin:true,Essential:true},{Key:'hklm|SecurityHealth',Name:'SecurityHealth',Desc:'Windows Security notification icon',Maker:'Microsoft Corporation',Command:'x',Source:'hklm',Enabled:true,Admin:true,Essential:true}],
    SetStartup: async () => {},
    Repair: async id => { (async () => { const st = {windows:['r_dism','r_sfc'],network:['r_dns','r_winsock','r_ip'],update:['r_stop','r_cache','r_start'],dns:['r_dns'],store:['r_store']}[id];
      for (const x of st) { emit('repair',{id,phase:'bezig',pct:-1,step:x}); for (let p=0;p<=100;p+=20){ emit('repair',{id,phase:'bezig',pct:p,step:''}); await sleep(120);} }
      emit('repair',{id,phase:'klaar',pct:100,step:'done'}); emit('done',{kind:'repair',ok:1,fail:0}); })(); return ''; },
    Privacy: async () => [{ID:'webcam',On:true,Apps:[{Key:'a',Name:'OBS Studio',Allowed:true,InUse:false,LastUsed:'2026-10-16'},{Key:'b',Name:'Camera',Allowed:true,LastUsed:'2026-08-31'}]},
      {ID:'microphone',On:true,Apps:[{Key:'c',Name:'Discord',Allowed:true,InUse:true,LastUsed:'2026-11-24'},{Key:'d',Name:'Zoom',Allowed:true,LastUsed:'2026-10-16'},{Key:'e',Name:'Microsoft Teams',Allowed:false,LastUsed:'2026-10-11'},{Key:'g',Name:'OBS Studio',Allowed:true,LastUsed:'2026-10-10'},{Key:'h',Name:'Settings',Allowed:true,LastUsed:'2026-10-09'},{Key:'i',Name:'Microsoft Edge',Allowed:true,LastUsed:'2026-10-08'},{Key:'j',Name:'OBS Studio',Allowed:true,LastUsed:'2026-09-15'},{Key:'k',Name:'WhatsApp',Allowed:true,LastUsed:'2026-09-03'},{Key:'l',Name:'Claude',Allowed:true}]},
      {ID:'location',On:false,Apps:[{Key:'f',Name:'Microsoft Edge',Allowed:true,LastUsed:'2026-10-14'}]}],
    SetPrivacy: async () => {},
    ContextMenu: async () => [{ID:'verb|a',Name:'Open Git Bash here',App:'Git for Windows',Kind:'verb',Where:'folders,background',Enabled:true},{ID:'ext|z',Name:'WinRAR',App:'Alexander Roshal',Kind:'ext',Where:'files,folders',Enabled:true},{ID:'ext|y',Name:'AMD Software',App:'Advanced Micro Devices, Inc.',Kind:'ext',Where:'desktop',Enabled:true},{ID:'verb|b',Name:'Open Git GUI here',App:'Git for Windows',Kind:'verb',Where:'background',Enabled:false},{ID:'ext|c',Name:'Rightmenu class',App:'',Kind:'ext',Where:'files',Enabled:true}],
    SetContextItem: async () => {},
    ExportProfile: async () => 'C:\\Users\\Alex\\Documents\\Dusk-2026-11-24.dusk',
    ImportProfile: async () => ({app:'Dusk',version:'1.8',created:'2026-11-24',apps:['Mozilla.Firefox','Discord.Discord','EpicGames.EpicGamesLauncher'],tweaks:['transp','hags','seconds']}),
    CheckUpdate: async () => ({version:'1.9',url:'https://x',sha256:'a'.repeat(64),notesEN:'Faster app scan and new repair tools.',notesNL:'Sneller apps scannen en nieuwe reparatie-tools.',newer:true}),
    ApplyUpdate: async () => { for (let p=0;p<=100;p+=25){ emit('selfupdate',p); await sleep(100);} return ''; }, JustUpdated: async () => false,
    Programs: async () => { await sleep(500); return [
      {Key:'k1',Name:'FileZilla 3.68.1',Publisher:'Tim Kosse',Version:'3.68.1',InstallDate:'20260310',SizeMB:43,Icon:'',Protected:false},
      {Key:'k2',Name:'Discord',Publisher:'Discord Inc.',Version:'1.0.9270',InstallDate:'20261124',SizeMB:152,Logo:'discord',CatalogID:'Discord.Discord',Protected:false},
      {Key:'k3',Name:'Example Tool',Publisher:'Example Software',Version:'6.0',InstallDate:'20251101',SizeMB:535,Icon:'',Protected:false},
      {Key:'k4',Name:'Example Utility',Publisher:'Example Inc.',Version:'9.0',InstallDate:'20260801',SizeMB:516,Icon:'',Protected:false},
      {Key:'k5',Name:'Realtek Audio Driver',Publisher:'Realtek Semiconductor Corp.',Version:'6.0.9998.1',InstallDate:'20260206',SizeMB:323,Icon:'',Protected:true},
      {Key:'k6',Name:'Microsoft Visual C++ v14 Redistributable (x64)',Publisher:'Microsoft Corporation',Version:'14.52',InstallDate:'20261001',SizeMB:21,Icon:'',Protected:true}]; },
    ProgramPlan: async keys => { await sleep(700); return { k1: [{Path:'C:\\Program Files\\FileZilla FTP Client',Kind:'install',Size:45e6},{Path:'C:\\Users\\Alex\\AppData\\Roaming\\FileZilla',Kind:'appdata',Size:2e5},{Path:'C:\\Users\\Alex\\Desktop\\FileZilla.lnk',Kind:'shortcut',Size:0}],
      k3: [{Path:'C:\\Users\\Alex\\AppData\\Local\\Example Tool',Kind:'install',Size:530e6},{Path:'C:\\Users\\Alex\\AppData\\Roaming\\Example Tool',Kind:'appdata',Size:12e6},{Path:'C:\\Users\\Alex\\Documents\\Example Tool',Kind:'documents',Size:1e6}] }; },
    UninstallPlan: async ids => { await sleep(500); return { 'Discord.Discord': [{Path:'C:\\Users\\Alex\\AppData\\Local\\Discord',Kind:'install',Size:1006e6},{Path:'C:\\Users\\Alex\\AppData\\Roaming\\Discord',Kind:'appdata',Size:63e6}] }; },
    UninstallPrograms: async (keys, paths) => { (async () => { for (const k of keys) { emit('prog',{id:k,phase:'bezig',pct:-1,step:'removing'}); await sleep(600); emit('prog',{id:k,phase:'klaar',pct:100,step:'removed'}); } emit('done',{kind:'uninstall',ok:keys.length,fail:0,freed:587e6,removed:paths}); })(); return ''; },
    UninstallApps: async (ids, paths) => { (async () => { await sleep(600); emit('done',{kind:'uninstall',ok:ids.length,fail:0,freed:1069e6,removed:paths}); })(); return ''; },
    InstallGPUDriver: async (n) => { (async () => { for (let p=0;p<=100;p+=10){ emit('gpu',{id:n,phase:'bezig',pct:p,step:'downloading'}); await sleep(200);} emit('gpu',{id:n,phase:'bezig',pct:-1,step:'installing'}); await sleep(800); emit('gpu',{id:n,phase:'klaar',pct:100,step:'done'}); emit('done',{kind:'drivers',ok:1,fail:0}); })(); return ''; },
    RunOEMTool: async () => 'notool', OpenVendorApp: async () => true,
    ScanJunk: async () => { await sleep(600); return [{Key:'temp',Size:389e6,Files:1686},{Key:'shader',Size:687e6,Files:511},{Key:'browser',Size:707e6,Files:4121},{Key:'discord',Size:352e6,Files:1587},{Key:'crash',Size:52e6,Files:65},{Key:'recycle',Size:0,Files:0},{Key:'wintemp',Size:23e6,Files:27,Admin:true},{Key:'wucache',Size:1e5,Files:2,Admin:true}]; },
    CleanJunk: async () => { await sleep(800); return 2187000000; },
    CreateRestorePoint: async () => { await sleep(800); return 'ok'; },
    InstallOEM: async ups => { (async () => { for (const u of ups) { for (let p=0;p<=100;p+=20){ emit('oem',{id:u.ID,phase:'bezig',pct:p,step:'downloading'}); await sleep(150);} emit('oem',{id:u.ID,phase:'klaar',pct:100,step:'done'}); } emit('done',{kind:'oem',ok:ups.length,fail:0}); })(); return ''; },
    ScanDrivers: async () => { await sleep(window.MOCK_FAST ? 300 : 3000); return JSON.parse(JSON.stringify(drv)); },
    InstallDrivers: async ids => { (async () => { emit('log','» Installing drivers'); ids.forEach(id => emit('driver',{id,phase:'bezig',pct:-1,step:'installing'})); await sleep(1500);
      ids.forEach(id => emit('driver',{id,phase:'klaar',pct:100,step:'done'})); drv = drv.filter(d => !ids.includes(d.ID)); emit('done',{kind:'drivers',ok:ids.length,fail:0}); })(); return ''; },
    ScanBloat: async () => { await sleep(window.MOCK_FAST ? 0 : 800); return JSON.parse(JSON.stringify(bloat)); },
    RemoveBloat: async pats => { (async () => { for (const p of pats) { emit('log','» Removing: '+p); emit('bloat',{id:p,phase:'bezig',pct:-1,step:'removing'}); await sleep(700);
      emit('bloat',{id:p,phase:'klaar',pct:100,step:'removed'}); emit('log','  ✓ removed'); } bloat = bloat.filter(b => !pats.includes(b.Pattern)); emit('done',{kind:'bloat',ok:pats.length,fail:0}); })(); return ''; },
    RunAsAdmin: async () => { throw 'cancelled'; },
    InstallApps: async ids => { ids.forEach(id => emit('app',{id,phase:'wacht',pct:-1,step:'queued'}));
      (async () => { for (const id of ids) { const a = apps.find(x => x.ID === id); emit('log','» '+(a.Installed?'Updating':'Installing')+': '+a.Name);
        if (id === 'OBSProject.OBSStudio') { emit('app',{id,phase:'bezig',pct:-1,step:'admin'}); await sleep(window.MOCK_SLOW?3000:600); }
        for (let p = 0; p <= 100; p += 9) { emit('app',{id,phase:'bezig',pct:p,step:'downloading'}); await sleep(window.MOCK_SLOW?400:120); }
        emit('app',{id,phase:'bezig',pct:-1,step:'installing'}); await sleep(window.MOCK_SLOW?3000:600);
        a.Installed = true; a.Version = a.Available || '1.0.0'; a.Available = ''; emit('app',{id,phase:'klaar',pct:100,step:'done'}); emit('log','  ✓ done'); }
        emit('apps', JSON.parse(JSON.stringify(apps))); emit('done',{kind:'apps',ok:ids.length,fail:0}); })(); return ''; },
  } } };
})();
