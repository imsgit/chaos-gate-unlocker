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

static void cg_fit_to_native(void *hwnd, int baseW, int baseH) {
	if (!hwnd) {
		return;
	}
	HWND h = (HWND)hwnd;
	HMODULE user32 = GetModuleHandleW(L"user32.dll");

	UINT dpi = 96;
	typedef UINT(WINAPI * GetDpiForWindowFn)(HWND);
	GetDpiForWindowFn getDpi = user32 ? (GetDpiForWindowFn)GetProcAddress(user32, "GetDpiForWindow") : NULL;
	if (getDpi) {
		UINT d = getDpi(h);
		if (d) {
			dpi = d;
		}
	}

	int scaleTenths = (int)((dpi * 10 + 48) / 96);
	if (scaleTenths < 10) {
		scaleTenths = 10;
	}

	RECT r;
	r.left = 0;
	r.top = 0;
	r.right = MulDiv(baseW, scaleTenths, 10);
	r.bottom = MulDiv(baseH, scaleTenths, 10);

	DWORD style = (DWORD)GetWindowLongPtrW(h, GWL_STYLE);
	DWORD exStyle = (DWORD)GetWindowLongPtrW(h, GWL_EXSTYLE);
	typedef BOOL(WINAPI * AdjForDpiFn)(LPRECT, DWORD, BOOL, DWORD, UINT);
	AdjForDpiFn adj = user32 ? (AdjForDpiFn)GetProcAddress(user32, "AdjustWindowRectExForDpi") : NULL;
	if (adj) {
		adj(&r, style, FALSE, exStyle, dpi);
	} else {
		AdjustWindowRectEx(&r, style, FALSE, exStyle);
	}

	int ww = r.right - r.left, wh = r.bottom - r.top;
	int sw = GetSystemMetrics(SM_CXSCREEN), sh = GetSystemMetrics(SM_CYSCREEN);
	SetWindowPos(h, NULL, (sw - ww) / 2, (sh - wh) / 2, ww, wh,
		SWP_NOZORDER | SWP_NOACTIVATE);
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

	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)

	C.cg_set_app_icon(w.Window())

	var once sync.Once
	_ = w.Bind("__cgReady", func() {
		once.Do(func() {
			C.cg_fit_to_native(w.Window(), 800, 600)
			C.cg_show(w.Window())
		})
	})

	w.Init(`(function(){function r(){window.__cgReady&&window.__cgReady()}` +
		`requestAnimationFrame(function(){requestAnimationFrame(r)});setTimeout(r,3000)})()`)

	w.Navigate(url)
	w.Run()
}
