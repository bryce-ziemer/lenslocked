package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Smith",
		Bio:  `<script>alert("Haha, you have been h4x0r3d!");</script>`,
	}

	err = t.Execute(os.Stdout, user) // template already handels encoding!

	if err != nil {
		panic(err)
	}
}
