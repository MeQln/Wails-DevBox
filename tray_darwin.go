//go:build darwin

package main

/*
#include <stdlib.h>

// start_tray 由 tray_darwin_helper.go 的 ObjC 实现提供
void start_tray(const unsigned char* data, int len, const char* tooltip);
*/
import "C"
import "unsafe"

//export handleTrayAction
func handleTrayAction(tag C.int) {
	actionCh <- trayAction(tag)
}

func startTray(iconData []byte, tooltip string) {
	cTooltip := C.CString(tooltip)
	defer C.free(unsafe.Pointer(cTooltip))
	C.start_tray((*C.uchar)(unsafe.Pointer(&iconData[0])), C.int(len(iconData)), cTooltip)
}