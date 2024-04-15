package main

import (
	"fmt"
	"hello/models"
	"image/color"
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

//	func LoadMovieData() (Items, error) {
//		data, err := ioutil.ReadFile("./data/testDataMovies.json")
//		if err != nil {
//			return Items{}, err
//		}
//		var moviesResult Items
//		err = json.Unmarshal(data, &moviesResult)
//		if err != nil {
//			return Items{}, err
//		}
//		return moviesResult, nil
//	}
var itemsData []models.Item
var collectionData = binding.NewUntypedList()

func main() {

	testItem1 := models.NewItem("Animals", "Dog")
	testItem1.AddLabel("Label: Test")
	testItem2 := models.NewItem("Animals", "Cat")
	testItem2.AddTag("Test Tag")

	itemsData = []models.Item{
		testItem1,
		testItem2,
		models.NewItemwithLabelTag("Animals", "Bird", "A bird", []string{"Test Label"}, []string{"Test Tag"}),
		models.NewItem("Animals", "Horse"),
		models.NewItem("Movies", "Avatar"),
		models.NewItem("Movies", "I am Legend"),
		models.NewItem("Movies", "TMNT"),
		models.NewItem("Movies", "No country for old men"),
		models.NewItem("Movies", "Inception"),
	}

	//var currentItem models.Item
	var currentItemID int

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
	//collectionFilter := widget.NewSelectEntry()

	// Item list binding
	for _, t := range itemsData {
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
	editItemButton := widget.NewButton("Edit", func() {
		fmt.Println("Edit Press")
	})
	// Label Add
	labelEntry := widget.NewEntry()
	labelAddButton := widget.NewButton("  +  ", func() {
		rawData, _ := collectionData.GetValue(currentItemID)
		if data, ok := rawData.(models.Item); ok {
			data.Labels = append(data.Labels, labelEntry.Text)
			itemLabelData.Append(labelEntry.Text)
			collectionData.SetValue(currentItemID, data)
			labelEntry.Text = ""
			SaveData()
			labelEntry.Refresh()
		}
	})
	labelAddButton.Disable()
	labelEntry.OnChanged = func(input string) {
		labelAddButton.Disable()

		if len(input) >= 3 {
			labelAddButton.Enable()
		}
	}

	// Tag Add
	tagEntry := widget.NewEntry()
	tagAddButton := widget.NewButton("  +  ", func() {
		rawData, _ := collectionData.GetValue(currentItemID)
		if data, ok := rawData.(models.Item); ok {
			data.Labels = append(data.Labels, tagEntry.Text)
			itemTagData.Append(tagEntry.Text)
			collectionData.SetValue(currentItemID, data)
			tagEntry.Text = ""
			SaveData()
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

	// Formatting
	nameDescriptionImageContainer := container.NewHSplit(itemNameDescriptionContainer, itemImagePlaceholder)
	itemLabelListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, labelAddButton, labelEntry), nil, nil, itemLabelsList)
	itemTagListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, tagAddButton, tagEntry), nil, nil, itemTagsList)
	labelTagListContainer := container.NewHSplit(itemLabelListWithEntry, itemTagListWithEntry)
	itemDataContainer := container.NewVSplit(nameDescriptionImageContainer, labelTagListContainer)
	dataDisplayContainer := container.NewHSplit(collectionList, itemDataContainer)
	dataDisplayContainer.Offset = 0.3
	TopContentContainer := container.NewVBox(TopContent, collectionSearchBar, editItemButton)
	content := container.NewBorder(TopContentContainer, nil, nil, nil, dataDisplayContainer)

	collectionList.OnSelected = func(id widget.ListItemID) {
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
			for _, t := range itemsData {
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

		for _, item := range itemsData {
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
	var savedData []models.Item
	data, _ := collectionData.Get()

	for _, data := range data {
		if item, ok := data.(models.Item); ok {
			savedData = append(savedData, models.NewItemwithLabelTag(
				item.Collection,
				item.Name,
				item.Description,
				item.Labels,
				item.Tags))
		}
	}

	itemsData = savedData
}
