package main

import (
	"fmt"

	"github.com/timohahaa/report-maker/internal/models"
	"github.com/timohahaa/report-maker/internal/repository/JSON/site"
)

type SiteRepository interface {
	GetAllSites() ([]models.Site, error)
}

func main() {
	jsonSiteRepo := site.SiteJSONRepo{}
	sites, err := jsonSiteRepo.GetAllSites()
	fmt.Println(err)
	fmt.Println(len(sites))
	for _, site := range sites {
		fmt.Printf("%+v\n\n\n", site)
	}

}
