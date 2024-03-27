package models

import "fmt"

type Item struct {
	Collection  string
	Name        string
	Description string
}

func NewItem(collection string, name string) Item {
	return Item{collection, name, ""}
}

func (t Item) String() string {
	return fmt.Sprintf("%s - $t", t.Collection, t.Name)
}
