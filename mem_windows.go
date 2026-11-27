//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type memStatus struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

var procMem = windows.NewLazySystemDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")

func uintptrOf(m *memStatus) uintptr { return uintptr(unsafe.Pointer(m)) }

// shellExecInfo = SHELLEXECUTEINFOW (zit niet in x/sys).
type shellExecInfo struct {
	Size       uint32
	Mask       uint32
	Hwnd       windows.Handle
	Verb       *uint16
	File       *uint16
	Parameters *uint16
	Directory  *uint16
	Show       int32
	InstApp    windows.Handle
	IDList     uintptr
	Class      *uint16
	KeyClass   windows.Handle
	HotKey     uint32
	IconOrMon  windows.Handle
	Process    windows.Handle
}

var procShellExecEx = windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW")

func shellExecuteEx(info *shellExecInfo) error {
	info.Size = uint32(unsafe.Sizeof(*info))
	r, _, err := procShellExecEx.Call(uintptr(unsafe.Pointer(info)))
	if r == 0 {
		return err
	}
	return nil
}
