package main

/*
#cgo windows LDFLAGS: -luser32
#include <windows.h>

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
		GetSystemMetrics(SM_CXICON), GetSystemMetrics(SM_CYICON), LR_DEFAULTCOLOR | LR_SHARED);
	HICON small = (HICON)LoadImageW(inst, L"APP", IMAGE_ICON,
		GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON), LR_DEFAULTCOLOR | LR_SHARED);
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

static void cg_webview2_missing(void) {
	MessageBoxW(NULL,
		L"Microsoft Edge WebView2 runtime was not found.\n\n"
		L"Install it from:\nhttps://developer.microsoft.com/microsoft-edge/webview2/\n\n"
		L"or use the native build of Chaos Gate Unlocker.",
		L"Chaos Gate Unlocker", MB_ICONERROR | MB_OK);
}
*/
import "C"

import (
	"sync"

	webview "github.com/webview/webview_go"
)

func openWindow(title, url string) {
	w := webview.New(false)
	defer w.Destroy()
	hwnd := w.Window()
	if hwnd == nil {
		C.cg_webview2_missing()
		return
	}

	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	C.cg_set_app_icon(hwnd)
	C.cg_center(hwnd)
	_ = w.Bind("__cgReady", sync.OnceFunc(func() { C.cg_show(hwnd) }))

	w.Init(`(function(){function r(){window.__cgReady&&window.__cgReady()}` +
		`requestAnimationFrame(function(){requestAnimationFrame(r)});setTimeout(r,3000)})()`)

	w.Navigate(url)
	w.Run()
}
