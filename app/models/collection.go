package models

import "fmt"

type Item struct {
	Collection  string
	Name        string
	Description string
	Tags        []string
}

func NewItem(collection string, name string) Item {
	return Item{
		Collection:  collection,
		Name:        name,
		Description: "",
		Tags:        []string{},
	}
}

func (t Item) String() string {
	return fmt.Sprintf("%s - $t", t.Collection, t.Name)
}
