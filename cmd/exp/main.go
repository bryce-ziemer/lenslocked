package main

import (
	"bryce-ziemer/github.com/lenslocked/models"
	"fmt"
)

func main() {
	gs := models.GalleryService{}
	fmt.Println(gs.Images(2))

}
