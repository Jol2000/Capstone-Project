package main

import (
	"encoding/json"
	"fmt"
	"hello/models"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
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
				//fmt.Println(items)
				resultData.Items = append(resultData.Items, items.Items...)
			}
		}
	}
	fmt.Println("Load data success")
	return resultData, nil
}

func EncodeMovieData(data models.Items) {
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
}

var itemsData models.Items
var collectionData = binding.NewUntypedList()
var currentItemID int

func main() {

	dataTest, _ := DecodeMovieData()
	fmt.Println(dataTest.CollectionNames())
	EncodeMovieData(dataTest)
	itemsData = dataTest

	// testItem1 := models.NewItem("Animals", "Dog")
	// testItem2 := models.NewItem("Animals", "Cat")
	// testItem2.AddTag("Test Tag")

	//collectionData.Append(testItem1)
	// itemsData.AddItem(testItem1)
	// itemsData.AddItem(testItem2)
	// itemsData.AddItem(models.NewItemwithLabelTag("Animals", "Bird", "A bird", []string{"Test Label"}, []string{"Test Tag"}))
	// itemsData.AddItem(models.NewItem("Animals", "Horse"))
	// itemsData.AddItem(models.NewItem("Movies", "Avatar"))
	// itemsData.AddItem(models.NewItem("Movies", "I am Legend"))
	// itemsData.AddItem(models.NewItem("Movies", "TMNT"))
	// itemsData.AddItem(models.NewItem("Movies", "No country for old men"))
	// itemsData.AddItem(models.NewItem("Movies", "Inception"))
	// fmt.Println("Length: ", len(itemsData.Items))

	//var currentItem models.Item
	labelSelectedID := -1

	a := app.New()
	w := a.NewWindow("Treasure It Desktop")

	// Top Content Bar
	tiTitle := canvas.NewText("Treasure It", color.Black)
	tiTitle.TextSize = 40
	searchEntry := canvas.NewText("Home", color.Black)
	searchEntry.TextSize = 30
	homeButton := canvas.NewText("Collections", color.Black)
	homeButton.TextSize = 30
	collectionsButton := canvas.NewText("Search", color.Black)
	collectionsButton.TextSize = 30
	addItemButton := widget.NewLabel("+")
	TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), addItemButton, searchEntry, homeButton, collectionsButton)

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

	//Search bar
	collectionSearchBar := widget.NewEntry()
	// Collection filter
	collectionsList := []string{"1", "2", "3"}
	collectionFilter := widget.NewSelectEntry(collectionsList)

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
	inputLabelData := []string{"L1", "L2", "L3"}

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
		})

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
		})

	//File List
	fileList := widget.NewList(
		func() int {
			// Number of items in the list
			return len(getFileNames())
		},
		func() fyne.CanvasObject {
			// Create an empty canvas object for each item
			return widget.NewLabel("")
		},
		func(index int, item fyne.CanvasObject) {
			// Update the label with file name
			item.(*widget.Label).SetText(getFileNames()[index])
		},
	)

	// Set onItemSelected callback for list items
	fileList.OnSelected = func(index int) {
		fileName := getFileNames()[index]
		openFile(fileName)
	}

	// Create a container for the list
	listContainer := container.NewVBox(
		widget.NewLabel("Files:"),
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
	itemNameDescriptionContainer := container.NewBorder(itemName, nil, nil, nil, itemDescription)
	// Item image
	//itemImage := canvas.NewImageFromFile()
	itemImagePlaceholder := widget.NewLabel("Image Placeholder")
	// Edit Item

	// Label Add
	labelEntry := widget.NewEntry()
	labelAddButton := widget.NewButton("  +  ", func() {
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

	labelRemoveButton := widget.NewButton("  -  ", func() {
		if labelSelectedID != -1 {
			rawData, _ := collectionData.GetValue(currentItemID)
			if data, ok := rawData.(models.Item); ok {
				data.RemoveLabelID(labelSelectedID)
				itemLabelData.Set(data.Labels)
				collectionData.SetValue(currentItemID, data)
				itemsData.UpdateItem(data)
				//EncodeMovieData(itemsData)
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
	editing := false
	var editItemButton *widget.Button
	editItemButton = widget.NewButton("Edit", func() {
		if !editing {
			editing = true
			labelEntry.Show()
			labelAddButton.Show()
			labelRemoveButton.Show()
			editNameDescription(itemNameDescriptionContainer)
			editItemButton.SetText("Save")
		} else {
			editing = false
			labelEntry.Hide()
			labelAddButton.Hide()
			labelRemoveButton.Hide()
			EncodeMovieData(itemsData)
			saveNameDescription(itemNameDescriptionContainer, currentItemID)
			editItemButton.SetText("Edit")
		}
	})

	// Formatting
	nameDescriptionImageContainer := container.NewHSplit(itemNameDescriptionContainer, itemImagePlaceholder)
	labelAddRemoveButtonContainer := container.NewHBox(labelAddButton, labelRemoveButton)
	itemLabelListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, labelAddRemoveButtonContainer, labelEntry), nil, nil, itemLabelsList)
	itemTagListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, tagAddButton, tagEntry), nil, nil, itemTagsList)
	_ = itemTagListWithEntry
	labelTagListContainer := container.NewHSplit(itemLabelListWithEntry, listContainer)
	itemDataContainer := container.NewVSplit(nameDescriptionImageContainer, labelTagListContainer)
	dataDisplayContainer := container.NewHSplit(collectionList, itemDataContainer)
	dataDisplayContainer.Offset = 0.3
	TopContentContainer := container.NewVBox(TopContent, collectionSearchBar, editItemButton, collectionFilter)
	content := container.NewBorder(TopContentContainer, nil, nil, nil, dataDisplayContainer)

	collectionList.OnSelected = func(id widget.ListItemID) {
		labelSelectedID = -1
		rawData, _ := collectionData.GetValue(id)
		currentItemID = id

		if data, ok := rawData.(models.Item); ok {
			// Set current Item
			//currentItem = data

			// Label Data
			labelData, _ := itemLabelData.Get()
			labelData = labelData[:0]
			itemLabelData.Set(labelData)
			for _, label := range data.Labels {
				itemLabelData.Append(label)
			}
			// Tag Data
			tagData, _ := itemTagData.Get()
			tagData = tagData[:0]
			itemTagData.Set(tagData)
			for _, tag := range data.Tags {
				itemTagData.Append(tag)
			}
			itemName.SetText(data.Name)
			itemDescription.SetText(data.Description)
		} else {
			fmt.Println("Data not found")
		}
	}

	collectionSearchBar.OnCursorChanged = func() {
		collectionList.UnselectAll()
	}

	collectionSearchBar.OnChanged = func(searchInput string) {
		if searchInput == "" {
			resetData, _ := collectionData.Get()
			resetData = resetData[:0]
			collectionData.Set(resetData)
			for _, t := range itemsData.Items {
				collectionData.Append(t)
			}
			return
		}
		searchData, _ := collectionData.Get()

		searchData = searchData[:0]
		collectionData.Set(searchData)

		searchInputs := strings.Split(searchInput, ",")

		var addedItems []string
		//addedItems = append(addedItems, "")

		for _, item := range itemsData.Items {
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
				// Tag search
				for _, tag := range item.Tags {
					if strings.Contains(tag, searchSplit) {
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

	// Setting Content to window
	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}

func SaveData() {
	fmt.Println("Save")
	var savedData models.Items
	data, _ := collectionData.Get()

	for _, data := range data {
		if item, ok := data.(models.Item); ok {
			savedData.AddItem(models.NewItemwithLabelTag(
				item.Collection,
				item.Name,
				item.Description,
				item.Labels,
				item.Tags))
		}
	}

	itemsData = savedData
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
func handleDrop(uri string) {
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
}

// Function to open a file
func openFile(fileName string) {
	filePath := filepath.Join(".", "data", "files", fileName)
	cmd := fmt.Sprintf("xdg-open %s", filePath) // On Linux, use xdg-open to open the file
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
}

// dropHandler handles dropped files onto the window
// dropHandler handles dropped files onto the window
func dropHandler(pos fyne.Position, uris []fyne.URI) {
	for _, uri := range uris {
		handleDrop(uri.Path())
	}
	// Refresh the list after handling drops
	//fileList.Refresh()
}

func editNameDescription(itemNameDescriptionContainer *fyne.Container) {
	// Assuming itemNameDescriptionContainer contains two labels (Name and Description)
	nameLabel := itemNameDescriptionContainer.Objects[1].(*widget.Label)
	descriptionLabel := itemNameDescriptionContainer.Objects[0].(*widget.Label)

	// Create entry fields to replace the labels
	nameEntry := widget.NewEntry()
	nameEntry.SetText(nameLabel.Text)

	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetText(descriptionLabel.Text)

	// Create a new container to hold both entry fields
	newContainer := container.NewBorder(nameEntry, nil, nil, nil, descriptionEntry)

	// Replace the content of the existing container with the new container
	itemNameDescriptionContainer.Objects = []fyne.CanvasObject{newContainer}

	// Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}

func saveNameDescription(itemNameDescriptionContainer *fyne.Container, currentItemID int) {
	nameEntry := itemNameDescriptionContainer.Objects[1].(*widget.Entry)
	descriptionEntry := itemNameDescriptionContainer.Objects[3].(*widget.Entry)

	// Update the data in the collection
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.Name = nameEntry.Text
		data.Description = descriptionEntry.Text
		SaveData()
	}

	// Create new label widgets with the updated values
	nameLabel := widget.NewLabel(nameEntry.Text)
	descriptionLabel := widget.NewLabel(descriptionEntry.Text)

	// Replace the entry fields with the new labels in the container
	newItemNameDescriptionContainer := container.NewBorder(nameLabel, nil, nil, nil, descriptionLabel)
	itemNameDescriptionContainer.Objects = []fyne.CanvasObject{newItemNameDescriptionContainer}

	// Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}
