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
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// type Items struct {
// 	Results []Item `json:"Movies"`
// }

// type Item struct {
// 	Title string `json:"Title"`
// 	Plot  string `json:"Plot"`
// 	Genre string `json:"Genre"`
// }

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

var itemsData models.Items
var collectionData = binding.NewUntypedList()
var currentItemID int
var itemFileData = binding.NewUntypedList()
var collectionsFilter []string
var editing = false
var itemImagePlaceholder = canvas.NewImageFromFile("data/images/defualtImageIcon.jpg")

func main() {
	itemImagePlaceholder.FillMode = canvas.ImageFillContain
	dataTest, _ := DecodeMovieData()
	EncodeMovieData(dataTest)
	itemsData = dataTest

	labelSelectedID := -1
	fileSelectedID := -1

	a := app.New()
	w := a.NewWindow("Treasure It Desktop")

	// Collection List View
	// Tool Bar
	typeIcon := canvas.NewImageFromFile("../images/collection_icon.png")
	typeIcon.FillMode = canvas.ImageFillOriginal
	collectionName := canvas.NewText("Animals", color.Black)
	collectionName.TextSize = 30
	helpIcon := canvas.NewImageFromFile("../images/help_icon.png")
	helpIcon.FillMode = canvas.ImageFillOriginal
	searchIcon := canvas.NewImageFromFile("../images/search_icon.png")
	searchIcon.FillMode = canvas.ImageFillOriginal
	filterIcon := canvas.NewImageFromFile("../images/filter_icon.png")
	filterIcon.FillMode = canvas.ImageFillOriginal
	addIcon := canvas.NewImageFromFile("../images/add_icon.png")
	addIcon.FillMode = canvas.ImageFillOriginal

	createbtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		CreateItemForm(w)
	})

	filterbtn := widget.NewButtonWithIcon("", theme.ContentRedoIcon(), func() {
		FilterCollectionsForm(w)
	})

	uploadImgBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		ImageUploadForm(w)
	})

	//Settings
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
					a.Settings().SetTheme(theme.DarkTheme())
				} else {
					a.Settings().SetTheme(theme.LightTheme())
				}
			}
		}, w)

		form.Show() // Show the form dialog
	})

	//Search bar
	collectionSearchBar := widget.NewEntry()
	searchBarHelpBtn := widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
		dialog.ShowInformation("Search Bar Help", "The search bar filters the currently selected collection(s), use a comma (,) between multiple search criterea options.", w)
	})

	// Item list binding
	for _, t := range itemsData.Items {
		collectionData.Append(t)
	}

	// Collection List
	collectionList := widget.NewListWithData(
		collectionData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			diu, _ := di.(binding.Untyped).Get()
			item := diu.(models.Item)
			o.(*widget.Label).SetText(item.Name)
		})

	// Label List
	// Item Labels
	inputLabelData := []string{}

	itemLabelData := binding.NewStringList()

	for _, t := range inputLabelData {
		itemLabelData.Append(t)
	}

	itemLabelsList := widget.NewListWithData(
		itemLabelData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(di.(binding.String))
		},
	)

	// Tag List
	// Item Tags
	inputTagData := []string{}

	itemTagData := binding.NewStringList()

	for _, t := range inputTagData {
		itemTagData.Append(t)
	}

	itemTagsList := widget.NewListWithData(
		itemTagData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(di.(binding.String))
		},
	)

	//File List

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

	// Set onItemSelected callback for list items
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

	// Create a container for the list
	listContainer := container.NewBorder(
		widget.NewLabel("Files:"),
		fileRemoveButton,
		nil,
		nil,
		fileList,
	)
	w.SetOnDropped(dropHandler)

	//------------------------------------------------------------------

	// Icons
	editIcon := canvas.NewImageFromFile("../images/edit_icon.png")
	editIcon.FillMode = canvas.ImageFillOriginal

	itemData := widget.NewLabel("Select a movie")
	itemData.Wrapping = fyne.TextWrapWord
	// Container build
	//collectionListContainer := container.NewBorder(collectionToolBarContainer, nil, nil, nil, collectionList)
	//collectionViewLayout := container.NewGridWithColumns(2, collectionListContainer, itemViewLayout)

	// Content
	// Item Name and Description
	itemName := widget.NewLabel("Name")
	itemDescription := widget.NewLabel("Description")
	itemDescription.Wrapping = fyne.TextWrapWord
	itemNameDescriptionContainer := container.NewBorder(itemName, nil, nil, nil, itemDescription)
	// Item image
	//itemImage := canvas.NewImageFromFile()

	// Edit Item

	// Label Add
	labelEntry := widget.NewEntry()
	labelAddButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		rawData, _ := collectionData.GetValue(currentItemID)
		if data, ok := rawData.(models.Item); ok {
			data.AddLabel(labelEntry.Text)
			itemLabelData.Append(labelEntry.Text)
			collectionData.SetValue(currentItemID, data)
			itemsData.UpdateItem(data)
			labelEntry.Text = ""
			//EncodeMovieData(itemsData)
			labelEntry.Refresh()
		}
	})

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
	labelEntry.OnChanged = func(input string) {
		labelAddButton.Disable()

		if len(input) >= 3 {
			labelAddButton.Enable()
		}
	}

	itemLabelsList.OnSelected = func(id widget.ListItemID) {
		fileList.UnselectAll()
		labelSelectedID = id
		labelRemoveButton.Enable()
	}

	labelEntry.Hide()
	labelAddButton.Hide()
	labelAddButton.Disable()
	labelRemoveButton.Hide()
	labelRemoveButton.Disable()

	// Tag Add
	tagEntry := widget.NewEntry()
	tagAddButton := widget.NewButton("  +  ", func() {
		rawData, _ := collectionData.GetValue(currentItemID)
		if data, ok := rawData.(models.Item); ok {
			data.AddLabel(tagEntry.Text)
			itemTagData.Append(tagEntry.Text)
			collectionData.SetValue(currentItemID, data)
			itemsData.UpdateItem(data)
			tagEntry.Text = ""
			EncodeMovieData(itemsData)
			tagEntry.Refresh()
		}
	})
	tagAddButton.Disable()
	tagEntry.OnChanged = func(input string) {
		tagAddButton.Disable()

		if len(input) >= 3 {
			tagAddButton.Enable()
		}
	}

	// Edit button
	var editItemButton *widget.Button
	editItemButton = widget.NewButton("Edit", func() {
		if !editing {
			editing = true
			labelEntry.Show()
			labelAddButton.Show()
			labelRemoveButton.Show()
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
			collectionList.FocusGained()
			saveNameDescription(itemNameDescriptionContainer, currentItemID)
			EncodeMovieData(itemsData)
			editItemButton.SetText("Edit")
		}
	})

	exitbtn := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
		a.Quit()
	})

	menubtn := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		// Toggle visibility of home, createbtn, searchbtn
		if createbtn.Visible() {
			//home.Hide()
			collectionSearchBar.Hide()
			editItemButton.Hide()
			createbtn.Hide()
			filterbtn.Hide()
			settingbtn.Hide()
			exitbtn.Hide()
		} else {
			//home.Hide()
			collectionSearchBar.Show()
			editItemButton.Show()
			createbtn.Show()
			filterbtn.Show()
			settingbtn.Show()
			exitbtn.Show()
		}
	})

	// Top Content Bar
	tiTitle := canvas.NewText("Treasure It", theme.ForegroundColor())
	tiTitle.TextSize = 24
	tiTitle.TextStyle = fyne.TextStyle{Bold: true}
	burgerMenu := container.NewHBox(editItemButton, createbtn, filterbtn, settingbtn, uploadImgBtn, menubtn)
	TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), burgerMenu)
	// Formatting
	nameDescriptionImageContainer := container.NewHSplit(itemNameDescriptionContainer, itemImagePlaceholder)
	labelAddRemoveButtonContainer := container.NewHBox(labelAddButton, labelRemoveButton)
	itemLabelListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, labelAddRemoveButtonContainer, labelEntry), nil, nil, itemLabelsList)
	itemTagListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, tagAddButton, tagEntry), nil, nil, itemTagsList)
	_ = itemTagListWithEntry
	labelTagListContainer := container.NewHSplit(itemLabelListWithEntry, listContainer)
	itemDataContainer := container.NewVSplit(nameDescriptionImageContainer, labelTagListContainer)
	dataDisplayContainer := container.NewHSplit(container.NewBorder(container.NewBorder(nil, nil, nil, searchBarHelpBtn, collectionSearchBar), nil, nil, nil, collectionList), itemDataContainer)
	dataDisplayContainer.Offset = 0.3
	//TopContentContainer := container.NewVBox()
	content := container.NewBorder(TopContent, nil, nil, nil, dataDisplayContainer)

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

	collectionSearchBar.OnCursorChanged = func() {
		collectionList.UnselectAll()
	}

	collectionSearchBar.OnChanged = func(searchInput string) {
		if searchInput == "" {
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

		searchInputs := strings.Split(searchInput, ",")

		var addedItems []string
		//addedItems = append(addedItems, "")

		for _, item := range itemsData.Items {
			if !slices.Contains(collectionsFilter, item.Collection) {
				fmt.Println("Filtered")
				continue
			}
			for _, searchSplit := range searchInputs {
				searchSplit = strings.Trim(searchSplit, " ")
				if searchSplit == "" {
					continue
				}
				// Name search
				if strings.Contains(item.Name, searchSplit) {
					used := false
					for _, itemUsed := range addedItems {
						if strings.Contains(itemUsed, item.Name) {
							used = true
						}
					}
					if !used {
						collectionData.Append(item)
						addedItems = append(addedItems, item.Name)
					}
				}
				// Collection search
				if strings.Contains(item.Collection, searchSplit) {
					used := false
					for _, itemUsed := range addedItems {
						if strings.Contains(itemUsed, item.Name) {
							used = true
						}
					}
					if !used {
						collectionData.Append(item)
						addedItems = append(addedItems, item.Name)
					}
				}
				// Label search
				for _, label := range item.Labels {
					if strings.Contains(label, searchSplit) {
						used := false
						for _, itemUsed := range addedItems {
							if strings.Contains(itemUsed, item.Name) {
								used = true
							}
						}
						if !used {
							collectionData.Append(item)
							addedItems = append(addedItems, item.Name)
						}
					}
				}
			}
		}
	}

	collectionList.Select(0)
	// Setting Content to window
	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}

func SaveData(itemsData *models.Items) {
	fmt.Println("Save")
	var savedData models.Items
	data, _ := collectionData.Get()

	for _, data := range data {
		if item, ok := data.(models.Item); ok {
			fmt.Println("Desc: ", item.Description)
			savedData.AddItem(models.NewItem(
				item.Collection,
				item.Name,
				item.Description,
				item.Labels,
				item.Tags,
				item.Files,
				item.Image))
		}
	}

	itemsData = &savedData
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
		fmt.Println(data.Files)
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
		SaveData(&itemsData)
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

func updateCollectionDataDescription(s string) {
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.Description = s
		collectionData.SetValue(currentItemID, data)
		itemsData.UpdateItem(data)
	}
}

func updateCollectionDataName(s string) {
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.Name = s
		collectionData.SetValue(currentItemID, data)
		itemsData.UpdateItem(data)
	}
}

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

func UpdateData() {
	resetData, _ := collectionData.Get()
	resetData = resetData[:0]
	collectionData.Set(resetData)
	for _, t := range itemsData.Items {
		collectionData.Append(t)
	}
}

// CreateItemForm creates a form dialog for creating a new item
func CreateItemForm(window fyne.Window) {
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
				UpdateData()
				EncodeMovieData(itemsData)
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
			} else {
				dialog.ShowError(errors.New("Name, Date, and Description are required."), window)
			}
		}
	}, window)

	form.Resize(fyne.NewSize(400, 300)) // Adjust the size of the form dialog
	form.Show()
}

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

// ImageUploadForm creates a form to upload an image for an item
func ImageUploadForm(window fyne.Window) {
	form := dialog.NewFileOpen(
		func(file fyne.URIReadCloser, err error) {
			handleImageDrop(file.URI().Path())
		}, window)

	form.Resize(fyne.NewSize(500, 500)) // Adjust the size of the form dialog
	form.Show()
}
