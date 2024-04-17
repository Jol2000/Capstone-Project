package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	//"fyne.io/fyne/layout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

// Define a variable to store the currently displayed items
var items []Item

// Define a variable to store the currently selected folder name
var selectedFolderName string

// Define a variable to store the currently displayed screen
var currentScreen fyne.CanvasObject

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Treasure It")

	// Variable to keep track of the selected folder index
	//var selectedFolderIndex int

	folders, err := LoadFolders()
	if err != nil {
		fmt.Println("Error loading folders:", err)
		folders = []Folder{}
	}

	// Create a list view for displaying folders
	listView := widget.NewList(
		func() int {
			return len(folders)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Placeholder label
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(folders[i].Name)
		},
	)

	home := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	createbtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		CreateFolderForm(myWindow, &folders, listView)
	})
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
		}, func(submitted bool) {
			// Function to handle submission of form
			if submitted && themeOption != "" { // Check if a theme option is selected
				if themeOption == "Dark" {
					myApp.Settings().SetTheme(theme.DarkTheme())
				} else {
					myApp.Settings().SetTheme(theme.LightTheme())
				}
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

	// Define the behavior when a folder is selected
	listView.OnSelected = func(i widget.ListItemID) {
		folderName := folders[i].Name
		fmt.Println("Opening folder:", folderName)

		// Show the folder details screen
		showFolderDetails(selectedFolderName, listView, folders, myWindow)
	}

	listView.Refresh()

	content := container.NewVBox(makeBanner(home, createbtn, searchbtn, settingbtn, exitbtn, uploadBtn), listView)

	myWindow.SetContent(content)
	myWindow.SetPadded(false)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}

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

const (
	dataFileName = "folders.json"
)

// Folder represents the structure of each folder
type Folder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateFolderForm creates a form dialog for creating a new folder
func CreateFolderForm(window fyne.Window, folders *[]Folder, listView *widget.List) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter Collection Name")

	form := dialog.NewForm("Create Collection", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Name:", nameEntry),
	}, func(submitted bool) {
		if submitted {
			name := nameEntry.Text
			if name != "" {
				newFolder := Folder{ID: len(*folders) + 1, Name: name}
				*folders = append(*folders, newFolder)
				saveFolders(*folders)
				// Update UI to reflect new folder
				listView.Refresh()
			}
		}
	}, window)

	form.Show()
}

// LoadFolders loads the saved folders from the data file
func LoadFolders() ([]Folder, error) {
	data, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		return nil, err
	}

	var folders []Folder
	err = json.Unmarshal(data, &folders)
	if err != nil {
		return nil, err
	}

	return folders, nil
}

// saveFolders saves the folders to the data file
func saveFolders(folders []Folder) error {
	data, err := json.MarshalIndent(folders, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dataFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func showFolderDetails(folderName string, listView *widget.List, folders []Folder, myWindow fyne.Window) {
	// Create a list view for displaying items
	itemListView := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Placeholder label
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[i].Name)
		},
	)

	// Define the behavior when an item is selected
	itemListView.OnSelected = func(i widget.ListItemID) {

	}

	createItemBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		createItemInFolder(folderName, listView, folders, myWindow)
	})

	// Create a button to go back to the folder screen
	backBtn := widget.NewButtonWithIcon("Back to Main Menu", theme.NavigateBackIcon(), func() {
		//mainContent := container.NewVBox(makeBanner(home, createbtn, searchbtn, settingbtn, exitbtn, uploadBtn), listView)
		//myWindow.SetContent(mainContent)
	})

	// Create a container to hold the folder details, items list, create item button, and back button
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Collection: %s", folderName)),
		createItemBtn,
		backBtn,
		itemListView,
	)

	// Replace the content of the main window
	myWindow.SetContent(content)

	// Refresh the list view to reflect the changes
	itemListView.Refresh()
}

func createItemInFolder(folderName string, listView *widget.List, folders []Folder, window fyne.Window) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter item name")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Enter item price")

	form := dialog.NewForm("Create Item", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Name:", nameEntry),
		widget.NewFormItem("Price:", priceEntry),
	}, func(submitted bool) {
		if submitted {
			name := nameEntry.Text
			priceStr := priceEntry.Text
			if name != "" && priceStr != "" {
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					dialog.ShowError(errors.New("invalid price"), window)
					return
				}

				// Create the item
				newItem := Item{
					ID:    len(items) + 1,
					Name:  name,
					Price: price,
					Image: "",
				}

				// Append the item to the items list
				items = append(items, newItem)

				// Save the items list
				saveItems(items)

				// Update the folder details view to reflect the changes
				showFolderDetails(folderName, listView, folders, window)
			}
		}
	}, window)

	form.Show()
}

const itemsFileName = "items.json"

// saveItems saves the items to the data file
func saveItems(items []Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(itemsFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
