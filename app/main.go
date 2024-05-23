package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello/models"
	"image/color"
	"time"

	//"image/color"
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
var itemFileData = binding.NewUntypedList()
var collectionsFilter []string
var editing = false
var itemImagePlaceholder = canvas.NewImageFromFile("data/images/defualtImageIcon.jpg")

// Define homePageView function
var homePageView func()

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

	showCollectionsView := func(w fyne.Window) {
		homebtn := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
			if homePageView != nil {
				homePageView()
			} else {
				fmt.Println("Error: homePageView is nil")
			}
		})

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
		burgerMenu := container.NewHBox(homebtn, helpbtn, editItemButton, createbtn, filterbtn, settingbtn, uploadImgBtn, menubtn)
		TopContent := container.New(layout.NewHBoxLayout(), tiTitle, layout.NewSpacer(), burgerMenu)
		// Formatting
		nameDescriptionImageContainer := container.NewHSplit(itemNameDescriptionContainer, itemImagePlaceholder)
		labelAddRemoveButtonContainer := container.NewHBox(labelAddButton, labelRemoveButton)
		itemLabelListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, labelAddRemoveButtonContainer, labelEntry), nil, nil, itemLabelsList)
		itemTagListWithEntry := container.NewBorder(nil, container.NewBorder(nil, nil, nil, tagAddButton, tagEntry), nil, nil, itemTagsList)
		_ = itemTagListWithEntry
		labelTagListContainer := container.NewHSplit(itemLabelListWithEntry, listContainer)
		itemDataContainer := container.NewVSplit(nameDescriptionImageContainer, labelTagListContainer)
		dataDisplayContainer := container.NewHSplit(container.NewBorder(collectionSearchBar, nil, nil, nil, collectionList), itemDataContainer)
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

		// Setting Content to window
		w.SetContent(content)
		w.Resize(fyne.NewSize(1200, 800))
	}

	homePageView := func() {
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
	homePageView()
	w.ShowAndRun()
}

func SaveData() {
	fmt.Println("Save")
	var savedData models.Items
	data, _ := collectionData.Get()

	for _, data := range data {
		if item, ok := data.(models.Item); ok {
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
		SaveData()
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
		SaveData()
	}
	itemImagePlaceholder.File = dstPath
	itemImagePlaceholder.Refresh()
	fmt.Println("Image saved:", dstPath) // Print the file directory
	SaveData()
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
	SaveData()
	EncodeMovieData(itemsData)
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
	itemNameDescriptionContainer.Objects = newContainer.Objects
	itemNameDescriptionContainer.Layout = newContainer.Layout

	// Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}

func saveNameDescription(itemNameDescriptionContainer *fyne.Container, currentItemID int) {
	nameEntry := itemNameDescriptionContainer.Objects[1].(*widget.Entry)
	descriptionEntry := itemNameDescriptionContainer.Objects[0].(*widget.Entry)

	// Update the data in the collection
	rawData, _ := collectionData.GetValue(currentItemID)
	if data, ok := rawData.(models.Item); ok {
		data.Name = nameEntry.Text
		data.Description = descriptionEntry.Text
		collectionData.SetValue(currentItemID, data)
		SaveData()
	}

	// Create new label widgets with the updated values
	nameLabel := widget.NewLabel(nameEntry.Text)
	descriptionLabel := widget.NewLabel(descriptionEntry.Text)
	// // Replace the entry fields with the new labels in the container
	newItemNameDescriptionContainer := container.NewBorder(nameLabel, nil, nil, nil, descriptionLabel)
	itemNameDescriptionContainer.Objects = newItemNameDescriptionContainer.Objects
	itemNameDescriptionContainer.Layout = newItemNameDescriptionContainer.Layout

	// // Refresh the container to reflect the changes
	itemNameDescriptionContainer.Refresh()
}

func SetNameDescription(itemNameDescriptionContainer *fyne.Container, name string, description string, editing bool) {
	if editing {
		nameEntry := widget.NewEntry()
		nameEntry.SetText(name)
		descriptionEntry := widget.NewEntry()
		descriptionEntry.SetText(description)

		newContainer := container.NewBorder(nameEntry, nil, nil, nil, descriptionEntry)
		itemNameDescriptionContainer.Objects = newContainer.Objects
		itemNameDescriptionContainer.Layout = newContainer.Layout
	} else {
		nameLabel := widget.NewLabel(name)
		descriptionLabel := widget.NewLabel(description)

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
