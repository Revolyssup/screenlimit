package action

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type info struct {
	pass string
}

var a = app.NewWithID("")
var w = a.NewWindow("Screen Limit")

func getWindowSingleton(str chan info) fyne.Window {
	w.Resize(fyne.NewSize(400, 400))
	e_pass := widget.NewEntry()
	e_pass.SetPlaceHolder("Enter the password...")
	e_submit := widget.NewButton("Submit", func() {
		i := info{
			pass: e_pass.Text,
		}
		str <- i
		e_pass.Refresh()
		w.Hide()
		fmt.Println("closing window after submit")
	})
	c := container.NewVBox(e_pass, e_submit)
	w.SetContent(container.NewHSplit(c, c))
	return w
}
