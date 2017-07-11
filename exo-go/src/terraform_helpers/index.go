package terraformHelpers

import (
	"fmt"
	"log"

	"github.com/hoisie/mustache"
)

// RenderTemplates renders templates
func RenderTemplates() {
	data := mustache.Render("hello {{c}}", map[string]string{"c": "world"})
	fmt.Println(data)
	readFile()
}

func readFile() {
	data, err := Asset("templates/main.tf")
	if err != nil {
		log.Fatalf("Failed tm read template: %s", err)
	}
	fmt.Println(string(data))
}
