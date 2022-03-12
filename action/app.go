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

func RunApp(str chan info) {
	fmt.Println("Should open app")
	a := app.New()
	w := a.NewWindow("Screen Limit")
	w.Resize(fyne.NewSize(400, 400))
	e_pass := widget.NewEntry()
	e_pass.SetPlaceHolder("Enter the password...")
	e_submit := widget.NewButton("Submit", func() {
		i := info{
			pass: e_pass.Text,
		}
		str <- i
		e_pass.Refresh()
		w.Close()
		a.Quit()
	})
	c := container.NewVBox(e_pass, e_submit)
	fmt.Println("c ", c.Layout)
	w.SetContent(container.NewHSplit(c, c))
	w.ShowAndRun()
	fmt.Println("asd")
}
