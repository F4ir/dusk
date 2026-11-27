//go:build !windows

package main

// Alleen zodat de tests ook op Linux draaien.

import "os/exec"

func hideWindow(cmd *exec.Cmd)                                  {}
func isAdmin() bool                                             { return false }
func relaunchAdmin() error                                      { return nil }
func restartExplorer() error                                    { return nil }
func wingetAvailable() bool                                     { return false }
func wingetList() (string, error)                               { return "", nil }
func wingetShow(id string) string                               { return "" }
func listAppx() ([]string, error)                               { return nil, nil }
func removeAppxArgs(p string) []string                          { return nil }
func tweakState(t Tweak) bool                                   { return false }
func setTweak(t Tweak, on bool) error                           { return nil }
func sysInfo() SysInfo                                          { return SysInfo{} }
func runOut(args []string) (string, error)                      { return "", nil }
func openDeviceManager() error                                  { return nil }
func runElevated(args []string, line func(string)) (int, error) { return 0, nil }

type StartupItem struct {
	Key, Name, Desc, Maker, Command, Source string
	Enabled, Admin, Essential               bool
}
type JunkCat struct {
	Key   string
	Size  int64
	Files int
	Admin bool
}

func ListStartup() []StartupItem           { return nil }
func SetStartup(key string, on bool) error { return nil }
func ScanJunk() []JunkCat                  { return nil }
func CleanJunk(keys []string) int64        { return 0 }
func fixPinnedIcons() []string             { return nil }
func oemToolArgs(brand string) []string    { return nil }
func openVendorApp(vendor string) bool     { return false }

type PrivApp struct {
	Key, Name, Path string
	Allowed, InUse  bool
	LastUsed        string
	Packaged        bool
}
type PrivCap struct {
	ID   string
	On   bool
	Apps []PrivApp
}

func (a *App) Privacy() []PrivCap                             { return nil }
func (a *App) SetPrivacy(capID, key string, allow bool) error { return nil }

type CtxItem struct {
	ID, Name, App, Kind, Where string
	Enabled                    bool
}

func (a *App) ContextMenu() []CtxItem                  { return nil }
func (a *App) SetContextItem(id string, on bool) error { return nil }
func windowsIsDutch() bool                             { return false }
