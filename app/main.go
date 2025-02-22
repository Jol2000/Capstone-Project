package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello/models"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

// Global Variables
var itemsData models.Items
var collectionData = binding.NewUntypedList()
var currentItemID int
var itemFileData = binding.NewUntypedList()
var collectionsFilter []string
var viewsFilter []string
var editing = false
var itemImagePlaceholder = canvas.NewImageFromFile("data/images/defualtImageIcon.jpg")
var collectionList *widget.List
var collectionSearchBar *widget.Entry
var homePageView func()

func main() {
	// Decode JSON data
	dataTest, _ := DecodeMovieData()
	// Store data in itemsData storage Variable
	itemsData = dataTest
	// Adds Item data to collection list
	UpdateCollectionSearch("")

	// Index Selection variables
	labelSelectedID := -1
	fileSelectedID := -1

	// Initialise app and window
	a := app.New()
	w := a.NewWindow("Treasure It Desktop")

	// Collections view page
	showCollectionsView := func(w fyne.Window) {
		// Home Button: Navigates to Home view page
		homebtn := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
			if homePageView != nil {
				homePageView()
			} else {
				fmt.Println("Error: homePageView is nil")
			}
		})

		// Help Button: Displays help dialog prompt
		helpbtn := widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
			helpText := `
		Welcome to Treasure It Desktop App!
		
		Instructions:
		1. Navigate through the collections using the list on the left.
		2. Click on a collection to view its items.
		3. Select an item to view its details on the right.
		4. To add a new item, click on the "+" button.
		5. You can filter collections using the filter button.
		6. To upload an image for an item, click on the upload image button.
		7. To edit an item, click on the edit button.
		8. Use the search bar to search for specific items.
		9. Enjoy organizing your treasures!

		For further assistance, please refer to the user manual or contact support.
		`

			// Create a dialog to display the help text
			dialog.ShowInformation("Help", helpText, w)
		})

		// Create Button: Calls create item function
		createbtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			CreateItemForm(w, collectionSearchBar)
		})

		// Filter Button: Calls Filter collections function
		filterbtn := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
			collectionList.UnselectAll()
			FilterCollectionsForm(w)
		})

		// Upload Image Button: Calls Upload Image function
		uploadImgBtn := widget.NewButtonWithIcon("", theme.FileImageIcon(), func() {
			ImageUploadForm(w)
		})

		// Print to Text Button: Calls Print Data Function
		printToTextBtn := widget.NewButtonWithIcon("", theme.ContentPasteIcon(), func() {
			PrintDataForm(w)
		})

		//Settings Button: Opens display settings form
		settingbtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
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
						a.Settings().SetTheme(theme.DarkTheme())
					} else {
						a.Settings().SetTheme(theme.LightTheme())
					}
				}
			}, w)

			form.Show() // Show the form dialog
		})

		// Search Entry: Used to search collection(s)
		collectionSearchBar = widget.NewEntry()

		// Search Bar Help Button: Displays search bar help dialog
		searchBarHelpBtn := widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
			dialog.ShowInformation("Search Bar Help", "The search bar filters the currently selected collection(s), use a comma (,) between multiple search criterea options.", w)
		})

		// Collection List
		collectionList = widget.NewListWithData(
			collectionData,
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(di binding.DataItem, o fyne.CanvasObject) {
				diu, _ := di.(binding.Untyped).Get()
				item := diu.(models.Item)
				o.(*widget.Label).SetText(item.Name)
			})

		// Import Help Text
		helpMessage := `Excel file data is imported utilizing Column Header values.
Use 'Collection', 'Name', 'Description' for relative Item information.
Any alternative Columns will have their data added to that item's Label data 
(in the format: [Header: Data]).
Cells without a Header will be added to that item's Label data 
(in the format: [Data]).`
		// Import Button: Calls Import Excel Function
		importButton := widget.NewButtonWithIcon("Import Excel", theme.ContentPasteIcon(), func() {
			infoDialog := dialog.NewInformation("Excel Import Information", helpMessage, w)
			infoDialog.SetOnClosed(func() {
				openExcel(w)
			})
			infoDialog.Show()
		})

		// Label List variable
		itemLabelData := binding.NewStringList()
		// Label List widget
		itemLabelsList := widget.NewListWithData(
			itemLabelData,
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(di binding.DataItem, o fyne.CanvasObject) {
				o.(*widget.Label).Bind(di.(binding.String))
			},
		)

		//File List widget
		fileList := widget.NewListWithData(
			itemFileData,
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(di binding.DataItem, o fyne.CanvasObject) {
				diu, _ := di.(binding.Untyped).Get()
				file := diu.(models.File)
				o.(*widget.Label).SetText(file.FileName)
			},
		)

		// Set openFile callback for File list items
		fileList.OnSelected = func(index int) {
			if !editing {
				itemLabelsList.UnselectAll()
				fileRaw, _ := itemFileData.GetValue(index)
				if data, ok := fileRaw.(models.File); ok {
					openFile(data.FileLocation)
				}
			} else {
				fileSelectedID = index
			}
		}

		// File Remove Button: Removes selected file from Item data
		fileRemoveButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			if fileSelectedID != -1 {
				rawData, _ := collectionData.GetValue(currentItemID)
				if data, ok := rawData.(models.Item); ok {
					data.RemoveFileID(fileSelectedID)
					resetData, _ := itemFileData.Get()
					resetData = resetData[:0]
					itemFileData.Set(resetData)
					for _, t := range data.Files {
						itemFileData.Append(t)
					}
					collectionData.SetValue(currentItemID, data)
					itemsData.UpdateItem(data)
					EncodeMovieData(itemsData)
					fileList.Refresh()
				}
			} else {
				fmt.Println("Please select a file")
			}
		})
		fileRemoveButton.Hide()

		// Item Remove Button: Removes selected Item from Item Data
		itemRemoveButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			if currentItemID != -1 {
				removeName := ""
				var collectionUpdate models.Items
				rawData, _ := collectionData.GetValue(currentItemID)
				if data, ok := rawData.(models.Item); ok {
					removeName = data.Name
					fmt.Println("Length:", collectionData.Length())
					for i := 0; i < collectionData.Length(); i++ {
						rawData, _ := collectionData.GetValue(i)
						if data, ok := rawData.(models.Item); ok {
							if data.Name != removeName {
								collectionUpdate.AddItem(data)
							}
						}
					}
					resetData, _ := collectionData.Get()
					resetData = resetData[:0]
					collectionData.Set(resetData)
					for _, t := range collectionUpdate.Items {
						collectionData.Append(t)
					}
					itemsData.RemoveItem(removeName)
					EncodeMovieData(itemsData)
				}
			} else {
				fmt.Println("Please select an Item")
			}
		})
		itemRemoveButton.Hide()

		// File List Container
		listContainer := container.NewBorder(
			nil,
			fileRemoveButton,
			nil,
			nil,
			fileList,
		)
		// File Drop handler
		w.SetOnDropped(dropHandler)

		// MAIN CONTENT
		// Item Name and Description with Container
		itemName := widget.NewLabel("")
		itemDescription := widget.NewLabel("")
		itemDescription.Wrapping = fyne.TextWrapWord
		itemNameDescriptionContainer := container.NewBorder(itemName, nil, nil, nil, itemDescription)

		// Label Add Entry and Button
		labelEntry := widget.NewEntry()
		labelAddButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			rawData, _ := collectionData.GetValue(currentItemID)
			if data, ok := rawData.(models.Item); ok {
				data.AddLabel(labelEntry.Text)
				itemLabelData.Append(labelEntry.Text)
				collectionData.SetValue(currentItemID, data)
				itemsData.UpdateItem(data)
				labelEntry.Text = ""
				labelEntry.Refresh()
			}
		})

		// Label Remove Button: Removes selected Label from Item Data
		labelRemoveButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			if labelSelectedID != -1 {
				rawData, _ := collectionData.GetValue(currentItemID)
				if data, ok := rawData.(models.Item); ok {
					data.RemoveLabelID(labelSelectedID)
					itemLabelData.Set(data.Labels)
					collectionData.SetValue(currentItemID, data)
					itemsData.UpdateItem(data)
					EncodeMovieData(itemsData)
					labelEntry.Refresh()
				}
			} else {
				fmt.Println("Please select a label")
			}
		})
		// Enables Label Entry
		labelEntry.OnChanged = func(input string) {
			labelAddButton.Disable()

			if len(input) >= 3 {
				labelAddButton.Enable()
			}
		}

		// Updates Label Selected variable
		itemLabelsList.OnSelected = func(id widget.ListItemID) {
			fileList.UnselectAll()
			labelSelectedID = id
			labelRemoveButton.Enable()
		}

		// Hide edit buttons
		labelEntry.Hide()
		labelAddButton.Hide()
		labelAddButton.Disable()
		labelRemoveButton.Hide()
		labelRemoveButton.Disable()

		// Image Placeholder format
		itemImagePlaceholder.FillMode = canvas.ImageFillContain
		// Data View Formating
		nameDescriptionImageContainer := container.NewHSplit(itemNameDescriptionContainer, itemImagePlaceholder)
		labelAddRemoveButtonContainer := container.NewHBox(labelAddButton, labelRemoveButton)
		itemLabelListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, labelAddRemoveButtonContainer, labelEntry), nil, nil, itemLabelsList)
		labelTagListContainer := container.NewHSplit(itemLabelListWithEntry, listContainer)
		itemDataContainerSplit := container.NewVSplit(nameDescriptionImageContainer, labelTagListContainer)
		itemDataContainer := container.NewBorder(nil, nil, nil, nil, itemDataContainerSplit)

		// Collection View Edit Button: Calls Filter collections form function
		viewEditBtn := widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
			FilterDataViewForm(w, itemNameDescriptionContainer, itemImagePlaceholder, itemLabelListWithEntry, listContainer, itemDataContainer)
		})

		// Edit button: Toggles Views Between Edit and Save
		var editItemButton *widget.Button
		editItemButton = widget.NewButton("Edit", func() {
			if !editing {
				editing = true
				labelEntry.Show()
				labelAddButton.Show()
				labelRemoveButton.Show()
				itemRemoveButton.Show()
				collectionList.FocusLost()
				fileRemoveButton.Show()
				editNameDescription(itemNameDescriptionContainer)
				editItemButton.SetText("Save")
			} else {
				editing = false
				labelEntry.Hide()
				labelAddButton.Hide()
				labelRemoveButton.Hide()
				fileRemoveButton.Hide()
				itemRemoveButton.Hide()
				collectionList.FocusGained()
				saveNameDescription(itemNameDescriptionContainer, currentItemID)
				EncodeMovieData(itemsData)
				editItemButton.SetText("Edit")
			}
		})

		// Exit Button: Closes Application
		exitbtn := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
			a.Quit()
		})

		// Menu Button: Toggles pop out menu visibiliy
		menubtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
			// Toggle visibility of home, createbtn, searchbtn
			if createbtn.Visible() {
				editItemButton.Hide()
				createbtn.Hide()
				filterbtn.Hide()
				settingbtn.Hide()
				exitbtn.Hide()
				uploadImgBtn.Hide()
				viewEditBtn.Hide()
				printToTextBtn.Hide()
				importButton.Hide()
			} else {
				editItemButton.Show()
				createbtn.Show()
				filterbtn.Show()
				settingbtn.Show()
				exitbtn.Show()
				uploadImgBtn.Show()
				viewEditBtn.Show()
				printToTextBtn.Show()
				importButton.Show()
			}
		})

		// Top Content Bar
		tiTitle := canvas.NewText("Treasure It", theme.ForegroundColor())
		tiTitle.TextSize = 24
		tiTitle.TextStyle = fyne.TextStyle{Bold: true}
		burgerMenu := container.NewHBox(homebtn, helpbtn, editItemButton, createbtn, filterbtn, uploadImgBtn, viewEditBtn, printToTextBtn, importButton, settingbtn, menubtn)
		TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), burgerMenu)

		// Formatting
		dataDisplayContainer := container.NewHSplit(container.NewBorder(container.NewBorder(nil, nil, nil, searchBarHelpBtn, collectionSearchBar), itemRemoveButton, nil, nil, collectionList), itemDataContainer)
		dataDisplayContainer.Offset = 0.3
		content := container.NewBorder(TopContent, nil, nil, nil, dataDisplayContainer)

		// Collection List Select Update
		collectionList.OnSelected = func(id widget.ListItemID) {
			fileList.UnselectAll()
			itemLabelsList.UnselectAll()
			labelSelectedID = -1
			rawData, _ := collectionData.GetValue(id)
			currentItemID = id

			if data, ok := rawData.(models.Item); ok {
				// Label Data
				labelData, _ := itemLabelData.Get()
				labelData = labelData[:0]
				itemLabelData.Set(labelData)
				for _, label := range data.Labels {
					itemLabelData.Append(label)
				}
				// Tag Data
				fileData, _ := itemFileData.Get()
				fileData = fileData[:0]
				itemFileData.Set(fileData)
				for _, file := range data.Files {
					itemFileData.Append(file)
				}
				// Update Name and Description
				SetNameDescription(itemNameDescriptionContainer, data.Name, data.Description, editing)
				if data.Image == "" {
					itemImagePlaceholder.File = "data/images/defualtImageIcon.jpg"
				} else {
					itemImagePlaceholder.File = data.Image
				}
				itemImagePlaceholder.Refresh()
			} else {
				fmt.Println("Data not found")
			}
		}

		// Deselect Collection on cursor change
		collectionSearchBar.OnCursorChanged = func() {
			collectionList.UnselectAll()
		}

		// Search Bar on change calls Update collection function with search input
		collectionSearchBar.OnChanged = func(searchInput string) {
			UpdateCollectionSearch(searchInput)
		}

		collectionList.Select(0)
		// Setting Content to window
		w.SetContent(content)
	}

	// Home Page View
	homePageView = func() {
		// Function to show the homepage view
		totalCollections := len(itemsData.CollectionNames())
		totalItems := len(itemsData.Items)

		overview := canvas.NewText("Welcome to Treasure It!", theme.ForegroundColor())
		overview.TextSize = 40
		overview.TextStyle = fyne.TextStyle{Bold: true}

		collectionLabelText := canvas.NewText("Total Collections", theme.ForegroundColor())
		collectionLabelText.TextSize = 16
		collectionNumberText := canvas.NewText(fmt.Sprintf("%d", totalCollections), theme.ForegroundColor())
		collectionNumberText.TextSize = 22
		collectionNumberText.TextStyle = fyne.TextStyle{Bold: true}

		collectionBorder := canvas.NewRectangle(color.Transparent)
		collectionBorder.StrokeColor = theme.ForegroundColor()
		collectionBorder.StrokeWidth = 2
		collectionBorder.Resize(fyne.NewSize(250, 120))

		collectionContainer := container.NewCenter(
			container.NewVBox(
				container.NewPadded(container.NewHBox(collectionLabelText)),
				container.NewHBox(layout.NewSpacer(), collectionNumberText, layout.NewSpacer()),
			),
		)

		collectionContent := container.NewMax(
			collectionBorder,
			collectionContainer,
		)

		itemLabelText := canvas.NewText("Total Items", theme.ForegroundColor())
		itemLabelText.TextSize = 16
		itemNumberText := canvas.NewText(fmt.Sprintf("%d", totalItems), theme.ForegroundColor())
		itemNumberText.TextSize = 22
		itemNumberText.TextStyle = fyne.TextStyle{Bold: true}

		itemBorder := canvas.NewRectangle(color.Transparent)
		itemBorder.StrokeColor = theme.ForegroundColor()
		itemBorder.StrokeWidth = 2
		itemBorder.Resize(fyne.NewSize(250, 120))

		itemContainer := container.NewCenter(
			container.NewVBox(
				container.NewPadded(container.NewHBox(itemLabelText)),
				container.NewHBox(layout.NewSpacer(), itemNumberText, layout.NewSpacer()),
			),
		)

		itemContent := container.NewMax(
			itemBorder,
			itemContainer,
		)

		date := time.Now().Format("Monday, January 2, 2006")
		dateLabel := canvas.NewText(fmt.Sprintf("Date: %s", date), theme.ForegroundColor())
		dateLabel.TextSize = 16

		timeLabel := canvas.NewText("", theme.ForegroundColor())
		timeLabel.TextSize = 16

		go func() {
			for {
				currentTime := time.Now()
				timeStr := currentTime.Format("15:04:05")
				timeLabel.Text = fmt.Sprintf("Time: %s", timeStr)
				time.Sleep(1 * time.Second)
			}
		}()

		dateTimeContainer := container.NewCenter(
			container.NewVBox(
				dateLabel,
				timeLabel,
			),
		)

		viewCollectionsButton := widget.NewButton("View Collections", func() {
			showCollectionsView(w)
		})
		viewCollectionsButton.Importance = widget.HighImportance

		toolbar := widget.NewToolbar(
			widget.NewToolbarAction(theme.HomeIcon(), func() {

			}),
		)

		settingbtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
			var themeOption string

			radio := widget.NewRadioGroup([]string{"Dark", "Light"}, func(selected string) {
				themeOption = selected
			})

			form := dialog.NewForm("Settings", "Ok", "Cancel", []*widget.FormItem{
				widget.NewFormItem("Theme", radio),
			}, func(submitted bool) {
				if submitted && themeOption != "" {
					if themeOption == "Dark" {
						a.Settings().SetTheme(theme.DarkTheme())
					} else {
						a.Settings().SetTheme(theme.LightTheme())
					}
				}
			}, w)

			form.Show()
		})

		exitbtn := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
			a.Quit()
		})

		topLeftButtons := container.NewHBox(toolbar, layout.NewSpacer(), settingbtn, exitbtn)
		header := container.NewHBox(overview, layout.NewSpacer())

		overviewContainer := container.NewVBox(
			header,
			widget.NewSeparator(),
			container.NewCenter(
				container.NewHBox(
					collectionContent,
					itemContent,
				),
			),
			dateTimeContainer,
			viewCollectionsButton,
		)
		overviewContainer.Add(widget.NewSeparator())

		footerText := canvas.NewText("© 2024 Treasure It | Privacy Policy | Terms of Service | Contact Us", theme.ForegroundColor())
		footerText.TextSize = 12

		footer := container.NewCenter(footerText)
		content := container.NewBorder(topLeftButtons, footer, nil, nil, container.NewCenter(overviewContainer))

		w.SetContent(content)
		w.Resize(fyne.NewSize(800, 600))

	}

	w.Resize(fyne.NewSize(1200, 800))
	homePageView()
	w.ShowAndRun()
}

// Save data Function:
func SaveData(itemsData *models.Items) {
	data, _ := collectionData.Get()

	for _, data := range data {
		if item, ok := data.(models.Item); ok {
			itemsData.UpdateItem(item)
		}
	}
}

// Function to get the list of file names in the /data/files folder
func getFileNames() []string {
	var fileNames []string
	files, err := os.ReadDir("./data/files")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return fileNames
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

// Function to handle dropped files
func handleFileDrop(uri string) {
	srcFile, err := os.Open(uri)
	if err != nil {
		fmt.Println("Error opening dropped file:", err)
		return
	}
	defer srcFile.Close()

	dstPath := filepath.Join(".", "data", "files", filepath.Base(uri))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("Error copying file contents:", err)
		return
	}

	// Update the data in the collection
	newFile := models.NewFile("", filepath.Base(uri))
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.AddFile(newFile)
		collectionData.SetValue(currentItemID, data)
		SaveData(&itemsData)
	}
	itemFileData.Append(newFile)
	fmt.Println("File saved:", filepath.Base(uri)) // Print the file directory
}

// Function to handle image upload
func handleImageDrop(uri string) {
	srcFile, err := os.Open(uri)
	if err != nil {
		fmt.Println("Error opening dropped file:", err)
		return
	}
	defer srcFile.Close()

	dstPath := filepath.Join(".", "data", "images", filepath.Base(uri))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("Error copying file contents:", err)
		return
	}

	// Update the data in the collection
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.AddImagePath(dstPath)
		collectionData.SetValue(currentItemID, data)
		fmt.Println(data.Image)
		itemsData.UpdateItem(data)
	}
	itemImagePlaceholder.File = dstPath
	itemImagePlaceholder.Refresh()
	fmt.Println("Image saved:", dstPath) // Print the file directory
	SaveData(&itemsData)
	EncodeMovieData(itemsData)
}

// Function to open a file
func openFile(fileName string) {
	filePath := filepath.Join(".", "data", "files", fileName)

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", filePath)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening file:", err)
		}
	case "linux":
		cmd := exec.Command("xdg-open", filePath)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening file:", err)
		}
	default:
		fmt.Println("Unsupported operating system")
		return
	}
}

// dropHandler handles dropped files onto the window
func dropHandler(pos fyne.Position, uris []fyne.URI) {
	for _, uri := range uris {
		handleFileDrop(uri.Path())
	}
	// Refresh the list after handling drops
	SaveData(&itemsData)
	EncodeMovieData(itemsData)
}

func editNameDescription(itemNameDescriptionContainer *fyne.Container) {
	// Assuming itemNameDescriptionContainer contains two labels (Name and Description)
	nameLabel := itemNameDescriptionContainer.Objects[1].(*widget.Label)
	descriptionLabel := itemNameDescriptionContainer.Objects[0].(*widget.Label)

	// Create entry fields to replace the labels
	nameEntry := widget.NewEntry()
	nameEntry.Wrapping = fyne.TextWrapWord
	nameEntry.SetText(nameLabel.Text)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.Wrapping = fyne.TextWrapWord
	//descriptionEntry.Scroll = container.ScrollVerticalOnly
	descriptionEntry.SetText(descriptionLabel.Text)

	nameEntry.OnChanged = func(s string) {
		updateCollectionDataName(s)
	}

	descriptionEntry.OnChanged = func(s string) {
		updateCollectionDataDescription(s)
	}

	// Create a new container to hold both entry fields
	newContainer := container.NewBorder(nameEntry, nil, nil, nil, descriptionEntry)

	// Replace the content of the existing container with the new container
	itemNameDescriptionContainer.Objects = newContainer.Objects
	itemNameDescriptionContainer.Layout = newContainer.Layout

	// Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}

// Update collection data Description
func updateCollectionDataDescription(s string) {
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.Description = s
		collectionData.SetValue(currentItemID, data)
		itemsData.UpdateItem(data)
	}
}

// Update collection data Name
func updateCollectionDataName(s string) {
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		originalName := data.Name
		data.Name = s
		collectionData.SetValue(currentItemID, data)
		itemsData.UpdateItemName(data, originalName)
	}
}

// Save current Name and Description
func saveNameDescription(itemNameDescriptionContainer *fyne.Container, currentItemID int) {
	nameEntry := itemNameDescriptionContainer.Objects[1].(*widget.Entry)
	descriptionEntry := itemNameDescriptionContainer.Objects[0].(*widget.Entry)

	SaveData(&itemsData)
	EncodeMovieData(itemsData)

	// Create new label widgets with the updated values
	nameLabel := widget.NewLabel(nameEntry.Text)
	nameLabel.Wrapping = fyne.TextWrapWord
	descriptionLabel := widget.NewLabel(descriptionEntry.Text)
	descriptionLabel.Wrapping = fyne.TextWrapWord
	// // Replace the entry fields with the new labels in the container
	newItemNameDescriptionContainer := container.NewBorder(nameLabel, nil, nil, nil, descriptionLabel)
	itemNameDescriptionContainer.Objects = newItemNameDescriptionContainer.Objects
	itemNameDescriptionContainer.Layout = newItemNameDescriptionContainer.Layout

	// // Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}

// Set Name and Description
func SetNameDescription(itemNameDescriptionContainer *fyne.Container, name string, description string, editing bool) {
	if editing {
		nameEntry := widget.NewMultiLineEntry()
		nameEntry.Wrapping = fyne.TextWrapWord
		nameEntry.SetText(name)
		descriptionEntry := widget.NewMultiLineEntry()
		descriptionEntry.Wrapping = fyne.TextWrapWord
		descriptionEntry.SetText(description)
		descriptionEntry.OnChanged = func(s string) {
			updateCollectionDataDescription(s)
		}

		newContainer := container.NewBorder(nameEntry, nil, nil, nil, descriptionEntry)
		itemNameDescriptionContainer.Objects = newContainer.Objects
		itemNameDescriptionContainer.Layout = newContainer.Layout
	} else {
		nameLabel := widget.NewLabel(name)
		nameLabel.Wrapping = fyne.TextWrapWord
		descriptionLabel := widget.NewLabel(description)
		descriptionLabel.Wrapping = fyne.TextWrapWord

		newContainer := container.NewBorder(nameLabel, nil, nil, nil, descriptionLabel)
		itemNameDescriptionContainer.Objects = newContainer.Objects
		itemNameDescriptionContainer.Layout = newContainer.Layout
	}
}

// Updates current collection bind.UntypedList object.
func UpdateData() {
	resetData, _ := collectionData.Get()
	resetData = resetData[:0]
	collectionData.Set(resetData)
	for _, t := range itemsData.Items {
		collectionData.Append(t)
	}
}

// CreateItemForm creates a form dialog for creating a new item
func CreateItemForm(window fyne.Window, collectionSearchBar *widget.Entry) {
	collectionEntry := widget.NewEntry()
	collectionEntry.SetPlaceHolder("Collection")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Enter Description")
	descriptionEntry.Resize(fyne.NewSize(300, 100)) // Set the initial size of the description entry

	form := dialog.NewForm("Create Collection", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Collection:", collectionEntry),
		widget.NewFormItem("Name:", nameEntry),
		widget.NewFormItem("Description:", descriptionEntry),
	}, func(submitted bool) {
		if submitted {
			name := nameEntry.Text
			collection := collectionEntry.Text
			description := descriptionEntry.Text
			if name != "" && collection != "" && description != "" {
				newItem := models.NewBasicItem(collection, name, description)
				itemsData.AddItem(newItem)
				EncodeMovieData(itemsData)
				collectionSearchBarinput := collectionSearchBar.Text
				UpdateCollectionSearch(collectionSearchBarinput)
			} else {
				dialog.ShowError(errors.New("Name, Date, and Description are required."), window)
			}
		}
	}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

// EditItemForm creates a form dialog for editing an item
func EditItemForm(window fyne.Window, itemID int) {
	var itemCollection string
	var itemName string
	var itemDescription string

	rawData, _ := collectionData.GetValue(itemID)
	if data, ok := rawData.(models.Item); ok {
		itemCollection = data.Collection
		itemName = data.Name
		itemDescription = data.Description
	}
	collectionEntry := widget.NewEntry()
	collectionEntry.SetPlaceHolder(itemCollection)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(itemName)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder(itemDescription)
	descriptionEntry.Resize(fyne.NewSize(300, 100)) // Set the initial size of the description entry

	form := dialog.NewForm("Create Collection", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Collection:", collectionEntry),
		widget.NewFormItem("Name:", nameEntry),
		widget.NewFormItem("Description:", descriptionEntry),
	}, func(submitted bool) {
		if submitted {
			name := nameEntry.Text
			collection := collectionEntry.Text
			description := descriptionEntry.Text
			if name != "" && collection != "" && description != "" {
				newItem := models.NewBasicItem(collection, name, description)
				itemsData.AddItem(newItem)
				UpdateData()
				EncodeMovieData(itemsData)
				collectionSearchBarinput := collectionSearchBar.Text
				UpdateCollectionSearch(collectionSearchBarinput)
			} else {
				dialog.ShowError(errors.New("Name, Date, and Description are required."), window)
			}
		}
	}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

// Filters collection for filter settings
func FilterCollectionUpdate() {
	resetData, _ := collectionData.Get()
	resetData = resetData[:0]
	collectionData.Set(resetData)
	for _, t := range itemsData.Items {
		if slices.Contains(collectionsFilter, t.Collection) {
			collectionData.Append(t)
		}
	}
}

// Updates View settings for Item based on sellection.
func FilterViewUpdate(nameDescription *fyne.Container, image *canvas.Image, labelList *fyne.Container, fileList *fyne.Container, dataDisplayContainer *fyne.Container) {
	topDataSplit := container.NewHSplit(nameDescription, image)
	topData := container.NewBorder(nil, nil, nil, nil, topDataSplit)
	bottomDataSplit := container.NewHSplit(labelList, fileList)
	bottomData := container.NewBorder(nil, nil, nil, nil, bottomDataSplit)
	topDataShow := true
	bottomDataShow := true

	if !slices.Contains(viewsFilter, "Name/Description") && !slices.Contains(viewsFilter, "Image") {
		topDataShow = false
	} else if slices.Contains(viewsFilter, "Name/Description") && slices.Contains(viewsFilter, "Image") {
		topData.Show()
	} else if slices.Contains(viewsFilter, "Name/Description") {
		topData = nameDescription
		topData.Show()
	} else {
		topData = container.NewBorder(nil, nil, nil, nil, image)
	}

	if !slices.Contains(viewsFilter, "Labels") && !slices.Contains(viewsFilter, "Files") {
		bottomDataShow = false
	} else if slices.Contains(viewsFilter, "Labels") && slices.Contains(viewsFilter, "Files") {
		topData.Show()
	} else if slices.Contains(viewsFilter, "Labels") {
		bottomData = labelList
	} else {
		bottomData = fileList
	}

	if topDataShow && bottomDataShow {
		dataDisplaySplit := container.NewVSplit(topData, bottomData)
		dataDisplayContainerNew := container.NewBorder(nil, nil, nil, nil, dataDisplaySplit)
		dataDisplayContainer.Layout = dataDisplayContainerNew.Layout
		dataDisplayContainer.Objects = dataDisplayContainerNew.Objects
	} else if topDataShow {
		dataDisplayContainerNew := container.NewBorder(nil, nil, nil, nil, topData)
		dataDisplayContainer.Layout = dataDisplayContainerNew.Layout
		dataDisplayContainer.Objects = dataDisplayContainerNew.Objects
	} else if bottomDataShow {
		dataDisplayContainerNew := container.NewBorder(nil, nil, nil, nil, bottomData)
		dataDisplayContainer.Layout = dataDisplayContainerNew.Layout
		dataDisplayContainer.Objects = dataDisplayContainerNew.Objects
	} else {
		dataDisplayContainerNew := container.NewBorder(nil, nil, nil, nil)
		dataDisplayContainer.Layout = dataDisplayContainerNew.Layout
		dataDisplayContainer.Objects = dataDisplayContainerNew.Objects
	}
}

// FilterCollectionsForm creates a form to filter collections
func FilterCollectionsForm(window fyne.Window) {
	collections := itemsData.CollectionNames()

	var formItems []*widget.FormItem
	for _, collection := range collections {
		collectionsCheck := widget.NewCheck(collection, nil)
		if slices.Contains(collectionsFilter, collection) {
			collectionsCheck.SetChecked(true)
		}
		formItems = append(formItems, widget.NewFormItem("", collectionsCheck))
	}

	form := dialog.NewForm("Filter Collections", "Filter", "Cancel", formItems,
		func(submitted bool) {
			if submitted {
				var collectionsFiltered []string
				for index, item := range formItems {
					// Cast the widget in each form item to a *widget.Check
					checkbox, ok := item.Widget.(*widget.Check)
					if ok {
						// Check if the checkbox is checked
						if checkbox.Checked {
							fmt.Printf("%s is selected\n", checkbox.Text)
							collectionsFiltered = append(collectionsFiltered, collections[index])
						}
					}
				}
				fmt.Println(collectionsFiltered)
				collectionsFilter = collectionsFiltered
				FilterCollectionUpdate()
			}
		}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

// FilterCollectionsForm creates a form to filter collections
func FilterDataViewForm(window fyne.Window, nameDescription *fyne.Container, image *canvas.Image, labelList *fyne.Container, fileList *fyne.Container, dataDisplayContainer *fyne.Container) {
	views := []string{"Name/Description", "Image", "Labels", "Files"}

	var formItems []*widget.FormItem
	for _, view := range views {
		viewsCheck := widget.NewCheck(view, nil)
		if slices.Contains(viewsFilter, view) {
			viewsCheck.SetChecked(true)
		}
		formItems = append(formItems, widget.NewFormItem("", viewsCheck))
	}

	form := dialog.NewForm("Edit Views", "Confirm", "Cancel", formItems,
		func(submitted bool) {
			if submitted {
				var viewsFiltered []string
				for index, item := range formItems {
					// Cast the widget in each form item to a *widget.Check
					checkbox, ok := item.Widget.(*widget.Check)
					if ok {
						// Check if the checkbox is checked
						if checkbox.Checked {
							fmt.Printf("%s is selected\n", checkbox.Text)
							viewsFiltered = append(viewsFiltered, views[index])
						}
					}
				}
				fmt.Println(viewsFiltered)
				viewsFilter = viewsFiltered
				FilterViewUpdate(nameDescription, image, labelList, fileList, dataDisplayContainer)
			}
		}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

// Logic for data print.
func PrintDataForm(window fyne.Window) {
	printableData := []string{"Name", "Description", "Labels", "Files"}

	var formItems []*widget.FormItem
	for _, view := range printableData {
		viewsCheck := widget.NewCheck(view, nil)
		formItems = append(formItems, widget.NewFormItem("", viewsCheck))
	}

	form := dialog.NewForm("Print Options", "Confirm", "Cancel", formItems,
		func(submitted bool) {
			if submitted {
				var printDataFiltered []string
				for index, item := range formItems {
					// Cast the widget in each form item to a *widget.Check
					checkbox, ok := item.Widget.(*widget.Check)
					if ok {
						// Check if the checkbox is checked
						if checkbox.Checked {
							fmt.Printf("%s is selected\n", checkbox.Text)
							printDataFiltered = append(printDataFiltered, printableData[index])
						}
					}
				}
				//viewsFilter = viewsFiltered
				CreatePrintFile(collectionData, printDataFiltered, window)
			}
		}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

// ImageUploadForm creates a form to upload an image for an item
func ImageUploadForm(window fyne.Window) {
	form := dialog.NewFileOpen(
		func(file fyne.URIReadCloser, err error) {
			if err != nil {
				// Handle error, such as upload canceled
				fmt.Println("Error:", err)
				return
			}
			if file == nil {
				// Upload canceled
				fmt.Println("Upload canceled")
				return
			}
			handleImageDrop(file.URI().Path())
		}, window)

	form.Resize(fyne.NewSize(500, 500)) // Adjust the size of the form dialog
	form.Show()
}

func CreatePrintFile(currentData binding.UntypedList, printOptions []string, window fyne.Window) {
	// Retrieve the data from the binding
	dataList, err := currentData.Get()
	if err != nil {
		fmt.Println("Error getting data:", err)
		return
	}

	// Convert the data to a list of Items
	itemsList := models.Items{}
	for _, item := range dataList {
		if typedItem, ok := item.(models.Item); ok {
			itemsList.AddItem(typedItem)
		} else {
			fmt.Println("Data is not of type Item")
			return
		}
	}

	// Group items by their collection
	collectionMap := make(map[string][]models.Item)
	for _, item := range itemsList.Items {
		collectionMap[item.Collection] = append(collectionMap[item.Collection], item)
	}

	// Get the current date and format it
	currentDate := time.Now().Format("2006-01-02") // YYYY-MM-DD format
	fileName := fmt.Sprintf("prints/print_%s.txt", currentDate)

	// Create the output file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write grouped data to the file
	for collection, items := range collectionMap {
		_, _ = file.WriteString(fmt.Sprintf("Collection: %s\n", collection))

		// Sort items alphabetically by name for consistent output
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})

		for _, item := range items {
			for _, option := range printOptions {
				switch option {
				case "Name":
					_, _ = file.WriteString(fmt.Sprintf("\tName: %s\n", item.Name))
				case "Description":
					_, _ = file.WriteString(fmt.Sprintf("\tDescription: %s\n", item.Description))
					// Add more cases for other options if needed
				case "Labels":
					_, _ = file.WriteString(fmt.Sprintf("\tLabels:\n"))
					for _, label := range item.Labels {
						_, _ = file.WriteString(fmt.Sprintf("\t\t%s\n", label))
					}
				case "Files":
					_, _ = file.WriteString(fmt.Sprintf("\tFiles:\n"))
					for _, fileData := range item.Files {
						_, _ = file.WriteString(fmt.Sprintf("\t\t%s\n", fileData.FileName))
					}
				}
			}
			_, _ = file.WriteString("\n")
		}
		_, _ = file.WriteString("\n")
	}
	dialog.ShowInformation("Text File Created", "Path: "+fileName, window)
}

// Opens excel input
func openExcel(window fyne.Window) {
	dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		importedItems, err := readExcel(filePath, window)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		itemsData.AddItems(importedItems)
		EncodeMovieData(itemsData)
		collectionSearchBarinput := collectionSearchBar.Text
		UpdateCollectionSearch(collectionSearchBarinput)
		if len(importedItems) != 0 {
			dialog.ShowInformation("Import Successful", fmt.Sprintf("%d items imported.", len(importedItems)), window)
		}
	}, window).Show()
}

// Reads excel input into collection storage
func readExcel(filePath string, w fyne.Window) ([]models.Item, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from Excel file: %w", err)
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("no data found in Excel file")
	}

	headers := rows[0]
	hasCollectionHeader := false
	hasNameHeader := false
	for _, header := range headers {
		if header == "Collection" {
			fmt.Println("Collection data found")
			hasCollectionHeader = true
		}
		if header == "Name" {
			fmt.Println("Name data found")
			hasNameHeader = true
		}
	}

	if !hasCollectionHeader || !hasNameHeader {
		dialog.ShowInformation("Invalid Import Data", "Header types missing (Collection or Name).", w)
		return nil, err
	}

	var items []models.Item

	for _, row := range rows[1:] {
		labels := []string{}
		var collection, name, description string
		for i, cell := range row {
			var header string
			if i < len(headers) {
				header = headers[i]
			}

			switch header {
			case "Collection":
				collection = cell
			case "Name":
				name = cell
			case "Description":
				description = cell
			default:
				if header == "" {
					labels = append(labels, cell)
				} else {
					if cell != "" {
						labels = append(labels, fmt.Sprintf("%s: %s", header, cell))
					}
				}
			}
		}
		item := models.NewItem(collection, name, description, labels, nil, nil, "")
		items = append(items, item)
	}
	return items, nil
}

// Reads Item Data from JSON
func DecodeMovieData() (models.Items, error) {
	var resultData models.Items
	items, _ := ioutil.ReadDir("./data/collections")
	for _, item := range items {
		if item.IsDir() {
			// subitems, _ := ioutil.ReadDir(item.Name())
			// for _, subitem := range subitems {
			// 	if !subitem.IsDir() {
			// 		// handle file there
			// 		fmt.Println(item.Name() + "/" + subitem.Name())
			// 	}
			// }
		} else {
			if strings.Split((item.Name()), ".")[1] == "JSON" {
				fmt.Println("Loading: ", item.Name())
				var result models.Items

				data, err := ioutil.ReadFile("./data/collections/" + item.Name())
				if err != nil {
					fmt.Println("Read data failure", err)
					return result, err
				}
				var items models.Items
				err = json.Unmarshal(data, &items)
				if err != nil {
					fmt.Println("Load data failure: ", err)
					return result, err
				}
				resultData.Items = append(resultData.Items, items.Items...)
			}
		}
	}
	fmt.Println("Load data success")
	return resultData, nil
}

// Writes Items Data to JSON
func EncodeMovieData(data models.Items) {
	fmt.Println("Encoding Data to JSON...")
	clearFolder("data/collections")
	collections := data.CollectionNames()
	for _, collection := range collections {
		file, errs := os.Create("data/collections/" + strings.ToLower(collection) + ".JSON")
		if errs != nil {
			fmt.Println("Failed to create file:", errs)
			return
		}
		defer file.Close()

		filteredCollection := data.FilterCollection(collection)
		encodedItem, err := json.MarshalIndent(filteredCollection, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		_, errs = file.WriteString(string(encodedItem))
		if errs != nil {
			fmt.Println("Failed to write to file:", errs) //print the failed message
			return
		}
	}
	fmt.Println("Data Encoded to JSON.")
}

// clears input folder of files
func clearFolder(folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("failed to read folder: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(folder, file.Name())
		err := os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("failed to remove file %s: %w", filePath, err)
		}
	}

	return nil
}

// Search bar logic
func UpdateCollectionSearch(searchCriterea string) {
	if searchCriterea == "" {
		dataTest, _ := DecodeMovieData()
		itemsData = dataTest
		for _, t := range itemsData.Items {
			collectionData.Append(t)
		}
		FilterCollectionUpdate()
		return
	}
	searchData, _ := collectionData.Get()

	searchData = searchData[:0]
	collectionData.Set(searchData)

	searchInputs := strings.Split(searchCriterea, ",")

	var validItems []models.Item

	for _, item := range itemsData.Items {
		if !slices.Contains(collectionsFilter, item.Collection) {
			continue
		}
		for _, searchSplit := range searchInputs {
			removeInput := false
			labelSearch := false
			if strings.Contains(searchSplit, ":") {
				searchSplit = strings.TrimLeft(searchSplit, ":")
				fmt.Println(searchSplit)
				labelSearch = true
			}
			searchSplit = strings.Trim(searchSplit, " ")
			if searchSplit == "" {
				continue
			}
			if searchSplit[0] == '-' {
				removeInput = true
				searchSplit = strings.TrimLeft(searchSplit, "-")
				searchSplit = strings.Trim(searchSplit, " ")
			}
			if searchSplit == "" {
				continue
			}

			// Name search
			if !labelSearch {
				if strings.Contains(item.Name, searchSplit) {
					used := false
					for i, itemUsed := range validItems {
						if strings.Contains(itemUsed.Name, item.Name) {
							if removeInput {
								validItems = append(validItems[:i], validItems[i+1:]...)
								continue
							}
							used = true
						}
					}
					if !used && !removeInput {
						validItems = append(validItems, item)
					}
				}
			}

			// Label search
			for _, label := range item.Labels {
				if strings.Contains(label, searchSplit) {
					used := false
					for i, itemUsed := range validItems {
						if strings.Contains(itemUsed.Name, item.Name) {
							if removeInput {
								validItems = append(validItems[:i], validItems[i+1:]...)
								continue
							}
							used = true
						}
					}
					if !used && !removeInput {
						validItems = append(validItems, item)
					}
				}
			}
		}
	}

	for _, item := range validItems {
		collectionData.Append(item)
	}
}
