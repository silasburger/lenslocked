package main

import (
	"html/template"
	"os"
)

type User struct {
	Name      string
	Age       int
	Inventory map[string]string
}

type Users = []User

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	inventory := map[string]string{"potion": "strength", "shield": "bronze"}

	// user := User{
	// 	Name:      "John Smith",
	// 	Age:       37,
	// 	Inventory: inventory,
	// }

	users := Users{
		{
			Name:      "John Smith",
			Age:       37,
			Inventory: map[string]string{"potion": "red", "shield": "wood"},
		},
		{
			Name:      "Smith",
			Age:       3,
			Inventory: inventory,
		},
	}

	err = t.Execute(os.Stdout, users)
	if err != nil {
		panic(err)
	}
}
