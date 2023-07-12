package report

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/timohahaa/report-maker/internal/models"
	"github.com/timohahaa/report-maker/internal/repository/JSON/asset"
	"github.com/timohahaa/report-maker/internal/repository/JSON/device"
	"github.com/timohahaa/report-maker/internal/repository/JSON/site"
	"github.com/xuri/excelize/v2"
)

const (
	ZIPstatus string = "ЗИП"
)

type reportRow []any

//excelize works with interface arrays for writing data to rows
//here is how a single row of a report will look like:
//
// DeviceModel |  Site  | ActiveCount | ZIPcount | ZIPcoverage | Specifications
//   string    | string |     int     |    int   |   float64   |     string
//   A col     | B col  |    C col    |  D col   |    E col    |     F col
//
//AciveCount - how many devices are active (in NetBox)
//ZIPcount - how many spare parts and devices of that type are available at the site (SnipeIT)
//ZIPcoverage - ZIPcount/ActiveCount
//of cource, this is not a concrete structure, it might be changed to any way you like

// creates a map of device models present on the site, mapped to their count
func createDeviceMap(site *models.Site, devices []models.Device) map[string]int {
	devMap := make(map[string]int)
	for _, device := range devices {
		if device.Site.Id == site.Id {
			devMap[device.DeviceType.Model] += 1
		}
	}
	return devMap
}

// counts how many devices of a certain model are present in a ZIP of a certain site
func countZIPdevices(site *models.Site, deviceModel string, assets []models.Asset) int {
	ZIPcount := 0
	for _, asset := range assets {
		if asset.Location == site.Facility && strings.ToLower(asset.Model) == deviceModel && asset.Status == ZIPstatus {
			ZIPcount++
		}
	}
	return ZIPcount
}

// creates an array of ready to be inserted report rows for a particular site
func createReportRows(siteRepo models.SiteRepository, deviceRepo models.DeviceRepository, assetRepo models.AssetRepository) ([]reportRow, error) {
	sites, err := siteRepo.GetAllSites()
	if err != nil {
		return nil, err
	}
	devices, err := deviceRepo.GetAllDevices()
	if err != nil {
		return nil, err
	}
	assets, err := assetRepo.GetAllAssets()
	if err != nil {
		return nil, err
	}

	rows := []reportRow{{"Модель", "Сайт", "Актив, шт", "ЗИП, шт", "% ЗИП", "Описание"}}
	for _, site := range sites {
		//make a map of device models -> their count for a particular site
		deviceMap := createDeviceMap(&site, devices)
		//for each device model find its count in a ZIP for a particular site
		//then append a row to rows slice
		for model, modelCount := range deviceMap {
			ZIPcount := countZIPdevices(&site, model, assets)
			ZIPcoverage := strconv.FormatFloat(float64(ZIPcount)/float64(modelCount), 'f', -1, 64) + "%"
			row := []any{model, site.Facility, modelCount, ZIPcount, ZIPcoverage, "specs"}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func CreateReport() error {
	siteRepo := &site.SiteJSONRepo{}
	deviceRepo := &device.DeviceJSONRepo{}
	assetRepo := &asset.AssetJSONRepo{}

	rows, err := createReportRows(siteRepo, deviceRepo, assetRepo)
	if err != nil {
		return err
	}
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	//this is what i call "ЧИСТАЯ АРХИТЕКТУРА"
	err = file.SetColWidth("Sheet1", "A", "A", 20)
	err = file.SetColWidth("Sheet1", "B", "B", 15)
	err = file.SetColWidth("Sheet1", "C", "C", 10)
	err = file.SetColWidth("Sheet1", "D", "D", 10)
	err = file.SetColWidth("Sheet1", "E", "E", 10)
	err = file.SetColWidth("Sheet1", "F", "F", 20)
	if err != nil {
		return err
	}
	for idx, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		file.SetSheetRow("Sheet1", cell, &row)
	}
	if err := file.AddChart("Sheet1", "G1", &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{
			{
				Name:       "placeholder",
				Categories: "Sheet1!$A$2:$A$10",
				Values:     "Sheet1!$C$2:$C$10",
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Title: excelize.ChartTitle{
			Name: "Количество девайсов на сайте",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowPercent: false,
		},
	}); err != nil {
		return err
	}
	if err := file.SaveAs("Report.xlsx"); err != nil {
		return err
	}
	return nil
}
