package main

import (
	"fmt"

	"github.com/timohahaa/report-maker/internal/models"
	"github.com/timohahaa/report-maker/internal/report"
	"github.com/timohahaa/report-maker/internal/repository/JSON/site"
)

func printA(repo models.AssetRepository) {
	assets, err := repo.GetAllAssets()
	fmt.Println(err)
	fmt.Println(len(assets))
	for _, asset := range assets {
		fmt.Printf("%+v\n\n\n", asset)
	}
}

func printS(repo models.SiteRepository) {
	sites, err := repo.GetAllSites()
	fmt.Println(err)
	fmt.Println(len(sites))
	for _, site := range sites {
		fmt.Printf("%+v\n\n\n", site)
	}
}

func printD(repo models.DeviceRepository) {
	devs, err := repo.GetAllDevices()
	fmt.Println(err)
	fmt.Println(len(devs))
	for _, dev := range devs {
		fmt.Printf("%+v\n\n\n", dev)
	}
}

func printDT(repo models.DeviceTypeRepository) {
	dTypes, err := repo.GetAllDeviceTypes()
	fmt.Println(err)
	fmt.Println(len(dTypes))
	for _, dt := range dTypes {
		fmt.Printf("%+v\n\n\n", dt)
	}
}

func main() {
	siteJSONRepo := &site.SiteJSONRepo{}
	err := report.CreateReport(siteJSONRepo)
	fmt.Println(err)
}
