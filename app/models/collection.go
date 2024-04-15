package models

import "fmt"

type Item struct {
	Collection  string
	Name        string
	Description string
	Tags        []string
	Labels      []string
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
