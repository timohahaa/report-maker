package report

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/timohahaa/report-maker/internal/models"
	"github.com/timohahaa/report-maker/internal/repository/JSON/asset"
	"github.com/timohahaa/report-maker/internal/repository/JSON/device"
	"github.com/timohahaa/report-maker/internal/repository/JSON/site"
	"github.com/xuri/excelize/v2"
)

const (
	ZIPstatus                     string = "ЗИП"
	minimumZIPcount                      = 1
	criticalZIPCoveragePercentage        = 0.3
)

type reportRow []any

type row struct {
	sphreadsheetRow reportRow
	//indicates if ZIP count is too low and a cell in a sphreadsheet should be marked as red
	//check const criticalZIPCoveragePercentage
	hasCritialZIPcount bool
}

//excelize works with interface arrays for writing data to rows
//here is how a single row of a report will look like:
//
//    Site     |  Device Model  | Device Category| ActiveCount| ZIPcount | ZIPcoverage | Device Specifications
//   string    |     string     |     string     |    int     |    int   |   float64   |        string
//   A col     |     B col      |     C col      |    D col   |   E col  |    F col    |        G col
//
//Device Category - i.e Switch/Server/Storage
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

// counts how many devices of a certain model are present in a ZIP of a certain site, and returns their count, specification and category
func retrieveZIPdevices(site *models.Site, deviceModel string, assets []models.Asset) (int, string, string) {
	ZIPcount := 0
	var specs, category string
	for _, asset := range assets {
		if asset.Location == site.Facility && strings.ToLower(asset.Model) == deviceModel /*&& asset.Status == ZIPstatus*/ {
			ZIPcount++
			specs = asset.Specifications
			category = asset.Category
		}
	}
	return ZIPcount, specs, category
}

// creates an array of ready to be inserted report rows for a particular site
func createReportRows(siteRepo models.SiteRepository, deviceRepo models.DeviceRepository, assetRepo models.AssetRepository) ([]row, error) {
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
	headers := []any{"Сайт", "Модель", "Категория", "Актив, шт", "ЗИП, шт", "% ЗИП", "Описание"}
	rows := []row{{headers, false}} //init a first row - headers
	for _, site := range sites {
		//make a map of device models -> their count for a particular site
		deviceMap := createDeviceMap(&site, devices)
		//for each device model find its count in a ZIP for a particular site
		//then append a row to rows slice
		for model, modelCount := range deviceMap {
			ZIPcount, specs, category := retrieveZIPdevices(&site, model, assets)
			ZIPcoverage := math.Round(float64(ZIPcount) / float64(modelCount) * 100.0) // * 100 to get percentages, round to make it look nice
			sphreadsheetRow := []any{
				site.Facility,
				model,
				category,
				modelCount,
				ZIPcount,
				strconv.FormatFloat(ZIPcoverage, 'f', -1, 64) + "%",
				specs,
			}
			hasCriticalZIPcount := ZIPcount >= minimumZIPcount && ZIPcoverage <= criticalZIPCoveragePercentage
			rows = append(rows, row{sphreadsheetRow, hasCriticalZIPcount})
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
	//this is what I call "ЧИСТАЯ АРХИТЕКТУРА"
	styleAllignRight, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", WrapText: true},
	})
	/*
		styleAllignLeft, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "left", WrapText: true},
		})
	*/
	styleLowZIPcount, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", WrapText: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FF0000"},
			Pattern: 1,
		},
	})
	err = file.SetColWidth("Sheet1", "A", "A", 15) //site
	err = file.SetColWidth("Sheet1", "B", "B", 20) //model
	err = file.SetColWidth("Sheet1", "C", "C", 15) //category
	err = file.SetColWidth("Sheet1", "D", "D", 10) //active count
	err = file.SetColWidth("Sheet1", "E", "E", 10) //ZIP count
	err = file.SetColWidth("Sheet1", "F", "F", 15) //ZIP coverage
	err = file.SetColWidth("Sheet1", "G", "G", 20) //specs
	err = file.SetColStyle("Sheet1", "A:G", styleAllignRight)
	//err = file.SetColStyle("Sheet1", "D:G", styleAllignRight)
	if err != nil {
		return err
	}
	//add the rows to sphreadsheet
	for idx, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		file.SetSheetRow("Sheet1", cell, &row.sphreadsheetRow)
		if row.hasCritialZIPcount {
			// I decided to ignore the error, because its just a format, if it fails is is not THAT bad...
			_ = file.SetCellStyle("Sheet1", "A"+cell[1:], "G"+cell[1:], styleLowZIPcount)
		}

	}
	err = file.SetCellFormula("Sheet1", "D13", "=SUM(D2:D10)")
	if err != nil {
		return err
	}
	//add a chart for active devices on a site
	if err := file.AddChart("Sheet1", "I1", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Активы",
				Categories: "Sheet1!$A$2:$A$10",
				Values:     "Sheet1!$D$2:$D$10",
			},
		},
		Dimension: excelize.ChartDimension{
			Width:  600,
			Height: 350,
		},
		XAxis: excelize.ChartAxis{
			Font: excelize.Font{
				Color: "404040", //Hex RGB of a color
			},
		},
		YAxis: excelize.ChartAxis{
			Font: excelize.Font{
				Color: "404040",
			},
		},
		Legend: excelize.ChartLegend{
			Position:      "top_right",
			ShowLegendKey: false,
		},
		Title: excelize.ChartTitle{
			Name: "Общее количество активов на сайте",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     false,
			ShowSerName:     false,
			ShowVal:         true,
		},
		ShowBlanksAs: "zero",
	}); err != nil {
		return err
	}
	//add a chart for ZIP devices on a site
	if err := file.AddChart("Sheet1", "I20", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "ЗИП",
				Categories: "Sheet1!$A$2:$A$10",
				Values:     "Sheet1!$E$2:$E$10",
			},
		},
		Dimension: excelize.ChartDimension{
			Width:  600,
			Height: 350,
		},
		XAxis: excelize.ChartAxis{
			Font: excelize.Font{
				Color: "404040", //Hex RGB of a color
			},
		},
		YAxis: excelize.ChartAxis{
			Font: excelize.Font{
				Color: "404040",
			},
		},
		Legend: excelize.ChartLegend{
			Position:      "top_right",
			ShowLegendKey: false,
		},
		Title: excelize.ChartTitle{
			Name: "Общее количество ЗИП на сайте",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     false,
			ShowSerName:     false,
			ShowVal:         true,
		},
		ShowBlanksAs: "zero",
	}); err != nil {
		return err
	}
	//add chart for a ZIP to active ratio on a site
	if err := file.AddChart("Sheet1", "I40", &excelize.Chart{
		Type: excelize.ColPercentStacked,
		Series: []excelize.ChartSeries{
			{
				Name:       "Активы",
				Categories: "Sheet1!$A$2:$A$10",
				Values:     "Sheet1!$D$2:$D$10",
			},
			{
				Name:       "ЗИП",
				Categories: "Sheet1!$A$2:$A$10",
				Values:     "Sheet1!$E$2:$E$10",
			},
		},
		Dimension: excelize.ChartDimension{
			Width:  600,
			Height: 350,
		},
		XAxis: excelize.ChartAxis{
			Font: excelize.Font{
				Color: "404040", //Hex RGB of a color
			},
		},
		YAxis: excelize.ChartAxis{
			MajorGridLines: true,
			//MinorGridLines: true,
			Font: excelize.Font{
				Color: "404040",
			},
		},
		Legend: excelize.ChartLegend{
			Position:      "top_right",
			ShowLegendKey: false,
		},
		Title: excelize.ChartTitle{
			Name: "ЗИП / Активы",
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     false,
			ShowSerName:     false,
			ShowVal:         false,
		},
		ShowBlanksAs: "zero",
	}); err != nil {
		return err
	}

	if err := file.SaveAs("Report.xlsx"); err != nil {
		return err
	}
	return nil
}
