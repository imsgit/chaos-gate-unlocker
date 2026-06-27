//go:build !nowebview

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
	void *pad[23];
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

static void cg_strip_icon(void *hwnd) {
	if (!hwnd) {
		return;
	}
	HWND h = (HWND)hwnd;
	SendMessageW(h, WM_SETICON, ICON_SMALL, 0);
	SendMessageW(h, WM_SETICON, ICON_BIG, 0);
	SetWindowLongPtrW(h, GWL_EXSTYLE, GetWindowLongPtrW(h, GWL_EXSTYLE) | WS_EX_DLGMODALFRAME);
	SetWindowPos(h, NULL, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_FRAMECHANGED);
}
*/
import "C"

import (
	"reflect"
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
	C.cg_strip_icon(w.Window())

	w.SetHtml(html)
	w.Run()
}
