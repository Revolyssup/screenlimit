package action

import (
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	file, _ := os.Open("./download.png")
	fileopen, _ := ioutil.ReadAll(file)
	image := canvas.NewImageFromResource(&fyne.StaticResource{
		StaticName:    "ashish",
		StaticContent: fileopen,
	})
	// image.Resize(fyne.Size{
	// 	Height: 1000,
	// 	Width:  1000,
	// })
	c := container.NewVBox(e_pass, e_submit)
	con := container.NewVSplit(image, c)
	w.SetContent(con)

	return w
}
