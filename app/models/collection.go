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
}

func NewItem(collection string, name string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: "",
		Tags:        []string{},
		Labels:      []string{},
	}
}

func NewItemwithLabelTag(collection string, name string, description string, labels []string, tags []string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: description,
		Tags:        tags,
		Labels:      labels,
	}
}

func (item Item) AddTag(tag string) {
	item.Tags = append(item.Tags, tag)
}

func (item Item) AddLabel(label string) {
	item.Labels = append(item.Labels, label)
	fmt.Println("Added label: ", label)
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
		fmt.Println("Length1: ", len(t.Items))
		t.Items = append(t.Items, itemInput)
		fmt.Println("Length2: ", len(t.Items))
		fmt.Println("Added Item: ", itemInput.Name)
	}
}

func (t *Items) AddItems(itemInputs []Item) {
	for _, item := range itemInputs {
		t.AddItem(item)
	}
}
