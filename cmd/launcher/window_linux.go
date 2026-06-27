//go:build !nowebview

package main

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static void paint_dark(void *win, int ww, int wh) {
	if (!win) {
		return;
	}
	GtkWindow *window = GTK_WINDOW(win);
	gtk_window_set_position(window, GTK_WIN_POS_CENTER);

	GtkWidget *child = gtk_bin_get_child(GTK_BIN(window));
	if (child && WEBKIT_IS_WEB_VIEW(child)) {
		const GdkRGBA bg = {0.082, 0.082, 0.082, 1.0};
		webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(child), &bg);
	}

	GdkScreen *screen = gtk_window_get_screen(window);
	if (!screen) {
		return;
	}
	GdkDisplay *display = gdk_screen_get_display(screen);
	GdkMonitor *mon = gdk_display_get_primary_monitor(display);
	if (!mon) {
		mon = gdk_display_get_monitor(display, 0);
	}
	if (!mon) {
		return;
	}
	GdkRectangle geo;
	gdk_monitor_get_geometry(mon, &geo);
	gtk_window_move(window, geo.x + (geo.width - ww) / 2, geo.y + (geo.height - wh) / 2);
}
*/
import "C"

import webview "github.com/webview/webview_go"

func openWindow(title, html string) {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	C.paint_dark(w.Window(), 800, 600)
	w.SetHtml(html)
	w.Run()
}
