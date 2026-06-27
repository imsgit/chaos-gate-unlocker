//go:build !nowebview

package main

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static void paint_dark(void *win) {
	if (!win) {
		return;
	}
	GtkWidget *child = gtk_bin_get_child(GTK_BIN(win));
	if (child && WEBKIT_IS_WEB_VIEW(child)) {
		const GdkRGBA bg = {0.082, 0.082, 0.082, 1.0};
		webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(child), &bg);
	}
}
*/
import "C"

import webview "github.com/webview/webview_go"

func openWindow(title, html string) {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	C.paint_dark(w.Window())
	w.SetHtml(html)
	w.Run()
}
