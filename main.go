package main

import (
	"context"
	"fmt"
	"hda/service"
	"os"
)

func init() {

	err := service.LoadAirports()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = service.LoadCompanies()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {

	go service.CrawlCompany()
	go service.CrawlAirports()

	ctx, cancel := context.WithCancel(context.Background())
	go service.CrawlHndDynFlight(cancel)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("System error, quit...")
			return
		}
	}
}
