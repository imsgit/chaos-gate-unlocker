package tooltip

import "fyne.io/fyne/v2"

func OverlayShown(obj fyne.CanvasObject) bool {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	if c == nil {
		return false
	}
	return c.Overlays().Top() != nil
}
