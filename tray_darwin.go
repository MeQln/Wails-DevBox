//go:build darwin

package main

/*
#include <stdlib.h>

// start_tray 由 tray_darwin_helper.go 的 ObjC 实现提供
void start_tray(const unsigned char* data, int len, const char* tooltip);
*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//export handleTrayAction
func handleTrayAction(tag C.int) {
	// 非阻塞发送：Cocoa 主线程不得阻塞在 channel 上。
	select {
	case actionCh <- trayAction(tag):
	default:
		// channel 满时丢弃事件，不影响主线程响应
	}
}

func startTray(iconData []byte, tooltip string) {
	cTooltip := C.CString(tooltip)
	C.start_tray((*C.uchar)(unsafe.Pointer(&iconData[0])), C.int(len(iconData)), cTooltip)
	// cTooltip 不释放：ObjC 的 dispatch_async 块异步读取该指针，
	// defer free 会在块执行前释放。一次性分配（~20 字节），应用生命周期内有效。
}

// startTrayActionHandler 在 darwin 上启动菜单事件分发 goroutine。
func startTrayActionHandler(ctx context.Context) {
	go func() {
		for act := range actionCh {
			switch act {
			case actionShow:
				runtime.WindowShow(ctx)
			case actionHide:
				runtime.WindowHide(ctx)
			case actionQuit:
				runtime.Quit(ctx)
			}
		}
	}()
}