package models

import "fmt"

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Collection  string   `json:"collection"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Labels      []string `json:"labels"`
	Files       []File   `json:"files"`
	Image       string   `json:"image"`
}

type File struct {
	FileName     string `json:"fileName"`
	FileLocation string `json:"fileLocation"`
}

func NewFile(name string, location string) File {
	if name == "" {
		return File{
			FileName:     location,
			FileLocation: location,
		}
	} else {
		return File{
			FileName:     name,
			FileLocation: location,
		}
	}
}

func NewItem(collection string, name string, description string, labels []string, tags []string, files []File, image string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: "",
		Files:       files,
		Tags:        tags,
		Labels:      labels,
		Image:       image,
	}
}

func NewBasicItem(collection string, name string, description string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: description,
		Files:       []File{},
		Tags:        []string{},
		Labels:      []string{},
	}
}

func NewItemwithLabelTag(collection string, name string, description string, labels []string, tags []string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: description,
		Files:       []File{},
		Tags:        tags,
		Labels:      labels,
	}
}

func (item *Item) AddTag(tag string) {
	item.Tags = append(item.Tags, tag)
}

func (item *Item) AddFile(file File) {
	item.Files = append(item.Files, file)
}

func (item *Item) AddImagePath(imagePath string) {
	item.Image = imagePath
}

func (item *Item) RemoveFileID(fileID int) {
	// Check if labelID is within the range of Labels slice
	if fileID >= 0 && fileID < len(item.Files) {
		// Remove the label at labelID position
		item.Files = append(item.Files[:fileID], item.Files[fileID+1:]...)
	}
}

func (item *Item) AddLabel(label string) {
	item.Labels = append(item.Labels, label)
	fmt.Println("Updated Labels: ", item.Labels)
}

func (item *Item) RemoveLabelID(labelID int) {
	// Check if labelID is within the range of Labels slice
	if labelID >= 0 && labelID < len(item.Labels) {
		// Remove the label at labelID position
		item.Labels = append(item.Labels[:labelID], item.Labels[labelID+1:]...)
	}
}

func (t Item) String() string {
	return fmt.Sprintf("%s - $t", t.Collection, t.Name, t.Labels, t.Tags)
}

func (t *Items) AddItem(itemInput Item) {
	contains := false
	for _, item := range t.Items {
		if item.Name == itemInput.Name {
			contains = true
		}
	}
	if !contains {
		t.Items = append(t.Items, itemInput)
	}
}

func (t *Items) AddItems(itemInputs []Item) {
	for _, item := range itemInputs {
		t.AddItem(item)
	}
}

func (t *Item) UpdateItem(itemInput Item) {
	*t = itemInput
}

func (t *Items) UpdateItem(itemInput Item) {
	for i, item := range t.Items {
		if item.Name == itemInput.Name {
			t.Items[i].UpdateItem(itemInput)
			fmt.Println("Updated Item: ", item.Name)
		}
	}
}

func (t Items) PrintData() {
	for _, item := range t.Items {
		fmt.Println("Coll: ", item.Collection)
		fmt.Println("Name: ", item.Name)
		fmt.Println("Desc: ", item.Description)

		fmt.Println("Labels: ")
		for _, label := range item.Labels {
			fmt.Println(label)
		}

		fmt.Println("Tags: ")
		for _, tag := range item.Tags {
			fmt.Println(tag)
		}
	}
}

func (t Items) PrintItemData(itemName string) {
	for _, item := range t.Items {
		if item.Name == itemName {
			fmt.Println("Coll: ", item.Collection)
			fmt.Println("Name: ", item.Name)
			fmt.Println("Desc: ", item.Description)

			fmt.Println("Labels: ")
			for _, label := range item.Labels {
				fmt.Println(label)
			}

			fmt.Println("Tags: ")
			for _, tag := range item.Tags {
				fmt.Println(tag)
			}
		}
	}
}

func (t Item) PrintItemData() {
	fmt.Println("Coll: ", t.Collection)
	fmt.Println("Name: ", t.Name)
	fmt.Println("Desc: ", t.Description)

	fmt.Println("Labels: ")
	for _, label := range t.Labels {
		fmt.Println(label)
	}

	fmt.Println("Tags: ")
	for _, tag := range t.Tags {
		fmt.Println(tag)
	}
}

func (items Items) CollectionCount() (numberOfCollections int) {
	count := 0
	var countedCollections []string
	for _, item := range items.Items {
		used := false
		for _, collection := range countedCollections {
			if item.Collection == collection {
				used = true
			}
		}
		if !used {
			countedCollections = append(countedCollections, item.Collection)
			count++
		}
	}
	return count
}

func (items Items) CollectionNames() (namesOfCollections []string) {
	var countedCollections []string
	for _, item := range items.Items {
		used := false
		for _, collection := range countedCollections {
			if item.Collection == collection {
				used = true
			}
		}
		if !used {
			countedCollections = append(countedCollections, item.Collection)
		}
	}
	return countedCollections
}

func (items Items) FilterCollection(collectionFilter string) (output Items) {
	var filteredCollection Items
	for _, item := range items.Items {
		if item.Collection == collectionFilter {
			filteredCollection.AddItem(item)
		}
	}
	return filteredCollection
}
