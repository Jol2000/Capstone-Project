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

	itemsData := []models.Item{
		models.NewItem("Animals", "Dog"),
		models.NewItem("Animals", "Cat"),
		models.NewItem("Animals", "Bird"),
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

	// Icons
	editIcon := canvas.NewImageFromFile("../images/edit_icon.png")
	editIcon.FillMode = canvas.ImageFillOriginal

	// Item Tool bar
	//itemToolBar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), helpIcon, filterIcon, editIcon)
	//itemToolBarContainer := container.NewVBox(itemToolBar)

	// Item View
	// itemTitleText := canvas.NewText("Horse", color.Black)
	// itemDescriptionText := canvas.NewText("Decription of item", color.Black)
	// horseImage := canvas.NewImageFromFile("../images/horse_image.png")

	//var testData2 = []string{"Name: Horsey", "Height: 213cm", "Diet: Carrot, Apple, Hay", "Birth: 12/06/2005", "Legs: 4", "Sex: Male", "Weight: 267kg", "Age: 19", "Colour: White", "Breed: Mule", "Area: 2B", "Personality: Timid", "Health: Fair", "For Sale: No"}

	// Item Data List
	// itemDataList := widget.NewList(
	// 	func() int {
	// 		return len(testData2)
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widget.NewLabel("Template")
	// 	},
	// 	func(i widget.ListItemID, o fyne.CanvasObject) {
	// 		o.(*widget.Label).SetText(testData2[i])
	// 	})

	//var testData3 = []string{"Herbivore", "Fur", "Favourite", "Harmless", "Old", "Fast", "Heavy"}

	// Item Data List
	//itemTagList := widget.NewList(
	// func() int {
	// 	return len(testData3)
	// },
	// func() fyne.CanvasObject {
	// 	return widget.NewLabel("Template")
	// },
	// func(i widget.ListItemID, o fyne.CanvasObject) {
	// 	o.(*widget.Label).SetText(testData3[i])
	// })

	//itemMainInfoLayout := container.NewGridWithColumns(2,
	//(container.NewBorder(itemTitleText, nil, nil, nil, itemDescriptionText)), horseImage)

	//itemMainDataLayout := container.NewGridWithColumns(2, itemDataList, itemTagList)

	//itemViewLayout := container.NewBorder(itemToolBarContainer, nil, nil, nil, container.NewGridWithRows(2, itemMainInfoLayout, itemMainDataLayout))

	itemData := widget.NewLabel("Select a movie")
	itemData.Wrapping = fyne.TextWrapWord
	// Container build
	//collectionListContainer := container.NewBorder(collectionToolBarContainer, nil, nil, nil, collectionList)
	//collectionViewLayout := container.NewGridWithColumns(2, collectionListContainer, itemViewLayout)

	// Content
	dataDisplayContainer := container.NewHSplit(collectionList, itemData)
	dataDisplayContainer.Offset = 0.3
	TopContentContainer := container.NewVBox(TopContent, collectionSearchBar)
	content := container.NewBorder(TopContentContainer, nil, nil, nil, dataDisplayContainer)

	collectionList.OnSelected = func(id widget.ListItemID) {
		rawData, _ := collectionData.GetValue(id)

		if data, ok := rawData.(models.Item); ok {
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
