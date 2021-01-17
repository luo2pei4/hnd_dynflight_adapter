package main

import (
	"fmt"
	"hda/service"
)

func main() {

	err := service.CrawlCompany()

	if err != nil {
		fmt.Println(err.Error())
	}
}
