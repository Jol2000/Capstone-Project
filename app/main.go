package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Movies struct {
	Results []Movie `json:"Movies"`
}

type Movie struct {
	Title string `json:"Title"`
	Plot  string `json:"Plot"`
}

func LoadMovieData() (Movies, error) {
	data, err := ioutil.ReadFile("./testDataMovies.json")
	if err != nil {
		return Movies{}, err
	}
	var moviesResult Movies
	err = json.Unmarshal(data, &moviesResult)
	if err != nil {
		return Movies{}, err
	}
	return moviesResult, nil
}

func main() {

	moviesData, err := LoadMovieData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Movies: %s/n", moviesData.Results)

	a := app.New()
	w := a.NewWindow("Treasure It Desktop")

	// Top Content Bar
	tiTitle := canvas.NewText("Treasure It", color.Black)
	tiTitle.TextSize = 40
	searchEntry := canvas.NewText("Search", color.Black)
	searchEntry.TextSize = 30
	homeButton := canvas.NewText("Home", color.Black)
	homeButton.TextSize = 30
	collectionsButton := canvas.NewText("Collections", color.Black)
	collectionsButton.TextSize = 30
	TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), searchEntry, homeButton, collectionsButton)
	TopContentContainer := container.NewVBox(TopContent)

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
	collectionToolBar := container.New(layout.NewHBoxLayout(), typeIcon, collectionName, layout.NewSpacer(), helpIcon, searchIcon, filterIcon, addIcon)
	collectionToolBarContainer := container.NewVBox(collectionToolBar)

	// Test List
	//var testData = []string{"Dog", "Cat", "Bear", "Horse", "Eagle", "Snake", "Tiger", "Turtle", "Frog", "Magpie", "Fox", "Zebra", "Kangaroo", "Dingo", "Giraffe", "Elephant", "Lion"}

	// Collection List
	collectionList := widget.NewList(
		func() int {
			return len(moviesData.Results)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(moviesData.Results[i].Title)
		})

	// Icons
	editIcon := canvas.NewImageFromFile("../images/edit_icon.png")
	editIcon.FillMode = canvas.ImageFillOriginal

	// Item Tool bar
	itemToolBar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), helpIcon, filterIcon, editIcon)
	itemToolBarContainer := container.NewVBox(itemToolBar)

	// Item View
	itemTitleText := canvas.NewText("Horse", color.Black)
	itemDescriptionText := canvas.NewText("Decription of item", color.Black)
	horseImage := canvas.NewImageFromFile("../images/horse_image.png")

	var testData2 = []string{"Name: Horsey", "Height: 213cm", "Diet: Carrot, Apple, Hay", "Birth: 12/06/2005", "Legs: 4", "Sex: Male", "Weight: 267kg", "Age: 19", "Colour: White", "Breed: Mule", "Area: 2B", "Personality: Timid", "Health: Fair", "For Sale: No"}

	// Item Data List
	itemDataList := widget.NewList(
		func() int {
			return len(testData2)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(testData2[i])
		})

	var testData3 = []string{"Herbivore", "Fur", "Favourite", "Harmless", "Old", "Fast", "Heavy"}

	// Item Data List
	itemTagList := widget.NewList(
		func() int {
			return len(testData3)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(testData3[i])
		})

	itemMainInfoLayout := container.NewGridWithColumns(2,
		(container.NewBorder(itemTitleText, nil, nil, nil, itemDescriptionText)), horseImage)

	itemMainDataLayout := container.NewGridWithColumns(2, itemDataList, itemTagList)

	itemViewLayout := container.NewBorder(itemToolBarContainer, nil, nil, nil, container.NewGridWithRows(2, itemMainInfoLayout, itemMainDataLayout))

	itemData := widget.NewLabel("Select a movie")

	// Container build
	collectionListContainer := container.NewBorder(collectionToolBarContainer, nil, nil, nil, collectionList)
	collectionViewLayout := container.NewGridWithColumns(2, collectionListContainer, itemViewLayout)

	// Content
	content := container.NewBorder(TopContentContainer, nil, nil, nil, collectionViewLayout)
	content2 := container.NewHSplit(collectionList, itemData)

	collectionList.OnSelected = func(id widget.ListItemID) {
		itemData.SetText(moviesData.Results[id].Plot)
	}

	// Setting Content to window
	w.SetContent(content)
	w.SetContent(content2)
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}
