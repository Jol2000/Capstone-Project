package main

import (
	"fmt"
	"path/filepath"

	//"fyne.io/fyne/layout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Item represents the structure of each item data
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Image string  `json:"image"`
}

// ImageFilter implements the fyne.FileFilter interface for image files
type ImageFilter struct{}

// Matches checks if the file extension is one of the supported image types
func (f ImageFilter) Matches(uri fyne.URI) bool {
	ext := filepath.Ext(uri.Path())
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// Description returns the description for the image filter
func (f ImageFilter) Description() string {
	return "Image Files (.jpg, .jpeg, .png, .gif)"
}

// UploadImage opens a file dialog for the user to upload an image file
func UploadImage(window fyne.Window) {
	openFile := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			// Upload the file
			// Here, you can save the image to a specific location
			fmt.Println("Image uploaded successfully:", reader.URI().Path())
			DisplayImage(window, reader.URI().Path())
		} else {
			fmt.Println("Error uploading image:", err)
		}
	}, window)

	openFile.SetFilter(ImageFilter{}) // Set filter for image file types
	openFile.Show()
}

func DisplayImage(window fyne.Window, imagePath string) {

	myApp := app.New()
	myWindow := myApp.NewWindow("Image Viewer")

	img := canvas.NewImageFromFile(imagePath)
	img.FillMode = canvas.ImageFillOriginal // Ensure the image is displayed in its original size
	img.SetMinSize(fyne.NewSize(1, 1))      // Ensure the image size is not limited by container constraints

	// Adjust the size of the image to fit the window dimensions
	img.Resize(fyne.NewSize(window.Canvas().Size().Width, window.Canvas().Size().Height))

	// Add the image to the window content
	container := container.NewCenter(img)

	myWindow.SetContent(container)

	myWindow.Resize(fyne.NewSize(img.MinSize().Width, img.MinSize().Height)) // Set window size to match image size
	myWindow.Show()
}

func makeBanner(home *widget.Button, createbtn *widget.Button, searchbtn *widget.Button, settingbtn *widget.Button, exitbtn *widget.Button, uploadBtn *widget.Button) fyne.CanvasObject {
	title := canvas.NewText("Treasure It", theme.ForegroundColor())
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}

	menubtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		// Toggle visibility of home, createbtn, searchbtn
		if home.Visible() {
			home.Hide()
			createbtn.Hide()
			searchbtn.Hide()
			settingbtn.Hide()
			exitbtn.Hide()
		} else {
			home.Show()
			createbtn.Show()
			searchbtn.Show()
			settingbtn.Show()
			exitbtn.Show()
		}
	})

	left := container.NewHBox(menubtn, home, createbtn, searchbtn, settingbtn, exitbtn, uploadBtn, title)

	return left
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Treasure It")

	home := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	createbtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
	searchbtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {})
	settingbtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		var themeOption string // Variable to store selected theme option

		// Create a radio group for selecting theme
		radio := widget.NewRadioGroup([]string{"Dark", "Light"}, func(selected string) {
			themeOption = selected // Update themeOption with selected theme
		})

		// Create a form dialog with radio buttons
		form := dialog.NewForm("Settings", "Ok", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Theme", radio), // Add radio group to the form
		}, func(bool) {
			// Function to handle submission of form
			if themeOption == "Dark" {
				myApp.Settings().SetTheme(theme.DarkTheme())
			} else {
				myApp.Settings().SetTheme(theme.LightTheme())
			}
		}, myWindow)

		form.Show() // Show the form dialog
	})
	exitbtn := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
		myApp.Quit()
	})
	uploadBtn := widget.NewButton("Upload Image", func() {
		UploadImage(myWindow)
	})

	myWindow.SetContent(container.NewVBox(makeBanner(home, createbtn, searchbtn, settingbtn, exitbtn, uploadBtn)))
	myWindow.SetPadded(false)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
