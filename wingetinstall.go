package main

// InstallWinget installeert winget (App Installer van Microsoft) op pc's die het nog niet hebben,
// zoals oudere Windows 10. Alles komt van de officiële GitHub-pagina van Microsoft.
const wingetInstallScript = `$ErrorActionPreference='Stop'
$ProgressPreference='SilentlyContinue'
$base='https://github.com/microsoft/winget-cli/releases/latest/download'
$d=Join-Path $env:TEMP 'dusk-winget'
New-Item -ItemType Directory -Force $d | Out-Null
try {
  "STEP w_deps"
  Invoke-WebRequest "$base/DesktopAppInstaller_Dependencies.zip" -OutFile "$d\deps.zip" -UseBasicParsing
  Expand-Archive "$d\deps.zip" "$d\deps" -Force
  $arch='x64'
  if($env:PROCESSOR_ARCHITECTURE -eq 'ARM64'){ $arch='arm64' } elseif(-not [Environment]::Is64BitOperatingSystem){ $arch='x86' }
  Get-ChildItem "$d\deps\$arch" -Filter *.appx | % { try { Add-AppxPackage -Path $_.FullName } catch { "  (" + $_.Exception.Message + ")" } }
  "STEP w_app"
  Invoke-WebRequest "$base/Microsoft.DesktopAppInstaller_8wekyb3d8bbwe.msixbundle" -OutFile "$d\winget.msixbundle" -UseBasicParsing
  Add-AppxPackage -Path "$d\winget.msixbundle"
  "STEP w_check"
  Start-Sleep 2
  $w=Join-Path $env:LOCALAPPDATA 'Microsoft\WindowsApps\winget.exe'
  if(-not (Test-Path $w)){ throw 'winget.exe not found after install' }
  & $w --version
  "STEP done"
} finally { Remove-Item -Recurse -Force $d -ErrorAction SilentlyContinue }`

func (a *App) InstallWinget() string {
	if wingetAvailable() {
		return "ok"
	}
	if !a.claim() {
		return "busy"
	}
	go func() {
		defer a.release()
		a.log("» " + a.tr("wingetinst"))
		a.emit("winget", Progress{"winget", "bezig", -1, "w_deps", ""})
		code, err := streamCmd(psArgs(wingetInstallScript), func(l string) {
			if step, ok := cutStep(l); ok {
				a.emit("winget", Progress{"winget", "bezig", -1, step, ""})
				return
			}
			a.log("  " + l)
		})
		if err == nil && code != 0 {
			err = errExit(code)
		}
		if err != nil || !wingetAvailable() {
			msg := "winget"
			if err != nil {
				msg = err.Error()
			}
			a.log("  ✗ " + msg)
			a.emit("winget", Progress{"winget", "fout", -1, "failed", msg})
			return
		}
		a.log("  ✓ " + a.tr("done"))
		a.emit("winget", Progress{"winget", "klaar", 100, "done", ""})
	}()
	return ""
}

func cutStep(l string) (string, bool) {
	const p = "STEP "
	if len(l) > len(p) && l[:len(p)] == p {
		return l[len(p):], true
	}
	return "", false
}
