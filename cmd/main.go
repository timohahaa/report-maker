package main

import (
	"fmt"

	"github.com/timohahaa/report-maker/repository/site"
)

func main() {
	jsonSiteRepo := site.SiteJSONRepo{}
	sites, err := jsonSiteRepo.GetAllSites()
	fmt.Println(err)
	fmt.Println(len(sites))
	for _, site := range sites {
		fmt.Printf("%+v\n\n\n", site)
	}
}
