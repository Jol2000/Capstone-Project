package main

import (
	"encoding/json"
	"fmt"
	"hello/models"
	"image/color"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Items struct {
	Results []Item `json:"Movies"`
}

type Item struct {
	Title string `json:"Title"`
	Plot  string `json:"Plot"`
	Genre string `json:"Genre"`
}

func LoadMovieData() (Items, error) {
	data, err := ioutil.ReadFile("./data/testDataMovies.json")
	if err != nil {
		return Items{}, err
	}
	var moviesResult Items
	err = json.Unmarshal(data, &moviesResult)
	if err != nil {
		return Items{}, err
	}
	return moviesResult, nil
}

func main() {

	testItem1 := models.NewItem("Animals", "Dog")
	testItem1.AddLabel("Label: Test")
	testItem2 := models.NewItem("Animals", "Cat")
	testItem2.AddTag("Test Tag")

	itemsData := []models.Item{
		testItem1,
		testItem2,
		models.NewItemwithLabelTag("Animals", "Bird", []string{"Test Label"}, []string{"Test Tag"}),
		models.NewItem("Animals", "Horse"),
		models.NewItem("Movies", "Avatar"),
		models.NewItem("Movies", "I am Legend"),
		models.NewItem("Movies", "TMNT"),
		models.NewItem("Movies", "No country for old men"),
		models.NewItem("Movies", "Inception"),
	}

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
	TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), searchEntry, homeButton, collectionsButton)

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

	// Item list binding
	collectionData := binding.NewUntypedList()
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

	itemLabelData := binding.NewUntypedList()

	for _, t := range inputLabelData {
		itemLabelData.Append(t)
	}

	itemLabelsList := widget.NewListWithData(
		itemLabelData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			diu, _ := di.(binding.Untyped).Get()
			label := diu.(string)
			o.(*widget.Label).SetText(label)
		})

	// Tag List
	// Item Tags
	inputTagData := []string{"T1", "T2", "T3"}

	itemTagData := binding.NewUntypedList()

	for _, t := range inputTagData {
		itemTagData.Append(t)
	}

	itemTagsList := widget.NewListWithData(
		itemTagData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, o fyne.CanvasObject) {
			diu, _ := di.(binding.Untyped).Get()
			tag := diu.(string)
			o.(*widget.Label).SetText(tag)
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
	labelTagListContainer := container.NewHSplit(itemLabelsList, itemTagsList)
	dataDisplayContainer := container.NewHSplit(collectionList, labelTagListContainer)
	dataDisplayContainer.Offset = 0.3
	TopContentContainer := container.NewVBox(TopContent, collectionSearchBar)
	content := container.NewBorder(TopContentContainer, nil, nil, nil, dataDisplayContainer)

	collectionList.OnSelected = func(id widget.ListItemID) {
		rawData, _ := collectionData.GetValue(id)

		if data, ok := rawData.(models.Item); ok {
			// Label Data
			labelData, _ := itemLabelData.Get()
			labelData = labelData[:0]
			itemLabelData.Set(labelData)
			fmt.Println("Labels: ", data.Labels)
			for _, label := range data.Labels {
				itemLabelData.Append(label)
			}
			fieldValue := data.Collection
			itemData.SetText(fieldValue)
		} else {
			fmt.Println("Data not found")
		}
	}

	collectionSearchBar.OnCursorChanged = func() {
		collectionList.UnselectAll()
	}

	collectionSearchBar.OnChanged = func(searchInput string) {
		if searchInput == "" {
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
		addedItems = append(addedItems, "")

		for _, item := range itemsData {
			for _, searchSplit := range searchInputs {
				searchSplit = strings.Trim(searchSplit, " ")
				if searchSplit == "" {
					continue
				}
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
			}
		}
	}

	// Setting Content to window
	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}
