package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Treasure It")

	//l := widget.NewLabel("Treasure It")

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MenuIcon(), func() {}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			newWin := myApp.NewWindow("Settings")

			darkbtn := widget.NewButton("Dark", func() {
				myApp.Settings().SetTheme(theme.DarkTheme())
			})

			lightbtn := widget.NewButton("Light", func() {
				myApp.Settings().SetTheme(theme.LightTheme())
			})

			backbtn := widget.NewButton("Back", func() {
				newWin.Close()
			})

			themes := container.NewVBox(darkbtn, lightbtn, backbtn)

			newWin.SetContent(themes)
			newWin.Resize(fyne.NewSize(500, 300))
			newWin.Show()
		}),
		widget.NewToolbarAction(theme.LogoutIcon(), func() {
			myApp.Quit()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	content := container.NewBorder(toolbar, nil, nil, nil)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
