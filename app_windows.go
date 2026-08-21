//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func trimProcessMemory() {
	// Minimal WorkingSet trim like Gotion on Windows
	const PROCESS_SET_QUOTA = 0x0100
	h, _ := windows.GetCurrentProcess()
	_ = windows.SetProcessWorkingSetSize(h, ^uintptr(0), ^uintptr(0))
	_ = unsafe.Pointer(nil)
}
