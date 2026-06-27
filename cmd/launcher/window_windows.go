package main

/*
#cgo windows LDFLAGS: -luser32
#include <windows.h>

extern void *webview_get_native_handle(void *w, int kind);

#define CG_NATIVE_HANDLE_BROWSER_CONTROLLER 2

typedef struct { BYTE A, R, G, B; } cg_color;

typedef struct cg_controller cg_controller;
typedef struct {
	HRESULT(STDMETHODCALLTYPE *QueryInterface)(cg_controller *, const GUID *, void **);
	ULONG(STDMETHODCALLTYPE *AddRef)(cg_controller *);
	ULONG(STDMETHODCALLTYPE *Release)(cg_controller *);
	void *pad[24];
	HRESULT(STDMETHODCALLTYPE *put_DefaultBackgroundColor)(cg_controller *, cg_color);
} cg_controller_vtbl;
struct cg_controller {
	cg_controller_vtbl *lpVtbl;
};

static void cg_paint_dark(void *wv) {
	if (!wv) {
		return;
	}
	cg_controller *ctrl = (cg_controller *)webview_get_native_handle(wv, CG_NATIVE_HANDLE_BROWSER_CONTROLLER);
	if (!ctrl) {
		return;
	}
	GUID iid = {0xc979903e, 0xd4ca, 0x4228, {0x92, 0xeb, 0x47, 0xee, 0x3f, 0xa9, 0x6e, 0xab}};
	cg_controller *c2 = NULL;
	if (ctrl->lpVtbl->QueryInterface(ctrl, &iid, (void **)&c2) == S_OK && c2) {
		cg_color bg = {255, 21, 21, 21};
		c2->lpVtbl->put_DefaultBackgroundColor(c2, bg);
		c2->lpVtbl->Release(c2);
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
	"reflect"
	"sync"
	"unsafe"

	webview "github.com/webview/webview_go"
)

func openWindow(title, html string) {
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)

	handle := *(*unsafe.Pointer)(unsafe.Pointer(reflect.ValueOf(w).Pointer()))
	C.cg_paint_dark(handle)
	C.cg_set_app_icon(w.Window())
	C.cg_center(w.Window())

	var once sync.Once
	_ = w.Bind("__cgReady", func() {
		once.Do(func() { C.cg_show(w.Window()) })
	})
	w.Init(`(function(){function r(){window.__cgReady&&window.__cgReady()}` +
		`if(document.readyState==='loading')document.addEventListener('DOMContentLoaded',r);else r();` +
		`setTimeout(r,3000)})()`)

	w.SetHtml(html)
	w.Run()
}
