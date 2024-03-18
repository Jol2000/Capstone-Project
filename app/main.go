package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Functionality
// Add Item Display
func defualtDisplay(content *fyne.Container) {
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
	helpIconWidget := widget.NewButton("Help", func() {})

	searchIcon := canvas.NewImageFromFile("../images/search_icon.png")
	searchIcon.FillMode = canvas.ImageFillOriginal
	searchIconWidget := widget.NewButton("Search", func() {})

	filterIcon := canvas.NewImageFromFile("../images/filter_icon.png")
	filterIcon.FillMode = canvas.ImageFillOriginal
	filterIconWidget := widget.NewButton("Filter", func() {})

	addIcon := canvas.NewImageFromFile("../images/add_icon.png")
	addIcon.FillMode = canvas.ImageFillOriginal
	addIconWidget := widget.NewButton("Add", func() {})
	collectionToolBar := container.New(layout.NewHBoxLayout(), typeIcon, collectionName, layout.NewSpacer(), searchIconWidget, filterIconWidget, addIconWidget)
	collectionToolBarContainer := container.NewVBox(collectionToolBar)

	// Test List
	var testData = []string{"Dog", "Cat", "Bear", "Horse", "Eagle", "Snake", "Tiger", "Turtle", "Frog", "Magpie", "Fox", "Zebra", "Kangaroo", "Dingo", "Giraffe", "Elephant", "Lion"}

	// Collection List
	collectionList := widget.NewList(
		func() int {
			return len(testData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(testData[i])
		})

	// Icons
	editIcon := canvas.NewImageFromFile("../images/edit_icon.png")
	editIcon.FillMode = canvas.ImageFillOriginal

	// Item Tool bar
	itemToolBar := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), helpIconWidget, filterIcon, editIcon)
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

	// Container build
	collectionListContainer := container.NewBorder(collectionToolBarContainer, nil, nil, nil, collectionList)
	collectionViewLayout := container.NewGridWithColumns(2, collectionListContainer, itemViewLayout)
	content = container.NewBorder(TopContentContainer, nil, nil, nil, collectionViewLayout)
}

func main() {

	a := app.New()
	w := a.NewWindow("Treasure It Desktop")

	// Setting Content to window
	content := container.NewBorder(nil, nil, nil, nil)
	defualtDisplay(content)
	w.SetContent(content)
	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}
