package main

/*
#cgo windows LDFLAGS: -luser32
#include <windows.h>

static void cg_hide(void *hwnd) {
	if (hwnd) {
		ShowWindow((HWND)hwnd, SW_HIDE);
	}
}

static void cg_show(void *hwnd) {
	if (!hwnd) {
		return;
	}
	HWND h = (HWND)hwnd;
	ShowWindow(h, SW_SHOW);

	HWND fg = GetForegroundWindow();
	DWORD fgThread = GetWindowThreadProcessId(fg, NULL);
	DWORD curThread = GetCurrentThreadId();
	BOOL attached = (fg && fgThread != curThread && AttachThreadInput(curThread, fgThread, TRUE));
	SetForegroundWindow(h);
	if (attached) {
		AttachThreadInput(curThread, fgThread, FALSE);
	}
}

static void cg_set_app_icon(void *hwnd) {
	if (!hwnd) {
		return;
	}
	HWND h = (HWND)hwnd;
	HINSTANCE inst = GetModuleHandleW(NULL);
	HICON big = (HICON)LoadImageW(inst, L"APP", IMAGE_ICON,
		GetSystemMetrics(SM_CXICON), GetSystemMetrics(SM_CYICON), LR_DEFAULTCOLOR);
	HICON small = (HICON)LoadImageW(inst, L"APP", IMAGE_ICON,
		GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON), LR_DEFAULTCOLOR);
	if (big) {
		SendMessageW(h, WM_SETICON, ICON_BIG, (LPARAM)big);
	}
	if (small) {
		SendMessageW(h, WM_SETICON, ICON_SMALL, (LPARAM)small);
	}
}

static void cg_center(void *hwnd) {
	if (!hwnd) {
		return;
	}
	HWND h = (HWND)hwnd;
	RECT rc;
	if (GetWindowRect(h, &rc)) {
		int ww = rc.right - rc.left, wh = rc.bottom - rc.top;
		int sw = GetSystemMetrics(SM_CXSCREEN), sh = GetSystemMetrics(SM_CYSCREEN);
		SetWindowPos(h, NULL, (sw - ww) / 2, (sh - wh) / 2, 0, 0, SWP_NOSIZE | SWP_NOZORDER);
	}
}
*/
import "C"

import (
	"sync"

	webview "github.com/webview/webview_go"
)

func openWindow(title, html string) {
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)

	C.cg_hide(w.Window())
	C.cg_set_app_icon(w.Window())
	C.cg_center(w.Window())

	var once sync.Once
	_ = w.Bind("__cgReady", func() {
		once.Do(func() { C.cg_show(w.Window()) })
	})

	w.Init(`(function(){function r(){window.__cgReady&&window.__cgReady()}` +
		`requestAnimationFrame(function(){requestAnimationFrame(r)});setTimeout(r,3000)})()`)

	w.SetHtml(html)
	w.Run()
}
