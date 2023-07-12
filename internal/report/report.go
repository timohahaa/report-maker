package report

import (
	"fmt"
	"strings"

	"github.com/timohahaa/report-maker/internal/models"
	"github.com/timohahaa/report-maker/internal/repository/JSON/asset"
	"github.com/timohahaa/report-maker/internal/repository/JSON/device"
	"github.com/xuri/excelize/v2"
)

/*
// reportRow represents a row in an report (Excel spreadsheet)
// DeviceModel AND Site uniquely identify a row in a sphreadsheet
// each site has multiple devices, and each device is present in multiple sites
// so their combo is unique
type reportRow struct {
	DeviceModel    string
	Site           string
	ActiveCount    int     //how many devices are active
	ZIPcount       int     //how many spare parts and devices of that type are available at the site
	ZIPcoverage    float64 // ZIPcount/ActiveCount
	Specifications string
}
*/

type reportRow []any

//excelize works with interface arrays for writing data to rows
//here is how a single row of a report will look like:
//
// DeviceModel |  Site  | ActiveCount | ZIPcount | ZIPcoverage | Specifications
//   string    | string |     int     |    int   |   float64   |     string
//
//AciveCount - how many devices are active (in NetBox)
//ZIPcount - how many spare parts and devices of that type are available at the site (SnipeIT)
//ZIPcoverage - ZIPcount/ActiveCount

// creates a map of device models present on the site to their count
func createDeviceMap(site *models.Site, deviceRepo models.DeviceRepository) (map[string]int, error) {
	devMap := make(map[string]int)
	devices, err := deviceRepo.GetAllDevices()
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if device.Site.Id == site.Id {
			devMap[device.DeviceType.Model] += 1
		}
	}
	return devMap, nil
}

// counts how many devices of a certain model are present in a ZIP of a certain site
func countZIPdevices(site *models.Site, deviceModel string, assets []models.Asset) int {
	ZIPcount := 0
	for _, asset := range assets {
		if asset.Location == site.Facility && strings.ToLower(asset.Model) == deviceModel {
			ZIPcount++
		}
	}
	return ZIPcount
}

// creates an array of ready to be inserted report rows for a particular site
func createReportRows(site *models.Site) ([]reportRow, error) {
	deviceRepo := &device.DeviceJSONRepo{}
	devMap, err := createDeviceMap(site, deviceRepo)
	if err != nil {
		return nil, err
	}
	rows := []reportRow{}
	assetRepo := &asset.AssetJSONRepo{}
	assets, err := assetRepo.GetAllAssets()
	if err != nil {
		return nil, err
	}
	for devModel, activeCount := range devMap {
		ZIPcount := countZIPdevices(site, devModel, assets)
		if err != nil {
			continue
		}
		row := []interface{}{
			devModel,
			site.Facility + " " + site.Name,
			activeCount, ZIPcount,
			float64(ZIPcount) / float64(activeCount),
			"потом тут обязательно будут спецификации",
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func CreateReport(siteRepo models.SiteRepository) error {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	err := file.SetColWidth("Sheet1", "A", "H", 20)
	if err != nil {
		return err
	}

	sites, err := siteRepo.GetAllSites()
	if err != nil {
		return err
	}

	rows := []reportRow{}
	for _, site := range sites {
		siteRows, err := createReportRows(&site)
		if err != nil {
			return err
		}
		rows = append(rows, siteRows...)
	}

	for idx, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		file.SetSheetRow("Sheet1", cell, &row)
	}
	if err := file.SaveAs("Report.xlsx"); err != nil {
		return err
	}
	return nil
}
