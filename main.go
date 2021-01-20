package main

import (
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

	// go service.CrawlCompany()
	// go service.CrawlAirports()

	// timer := time.Tick(60 * 1e9)

	// for {
	// 	select {
	// 	case <-timer:
	// 		fmt.Println("crawl dynamic flight info.")
	// 	}
	// }
	service.CrawlHndDynFlight("https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfarv.json")
	service.CrawlHndDynFlight("https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfdep.json")
	service.CrawlHndDynFlight("https://tokyo-haneda.com/app_resource/flight/data/int/hdacfarv.json")
	service.CrawlHndDynFlight("https://tokyo-haneda.com/app_resource/flight/data/int/hdacfdep.json")
}
