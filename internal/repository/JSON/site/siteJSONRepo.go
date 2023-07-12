package site

import (
	"encoding/json"
	"os"

	"github.com/timohahaa/report-maker/internal/models"
)

type SiteJSONRepo struct {
}

func (rep *SiteJSONRepo) GetAllSites() ([]models.Site, error) {
	//just read from the json file
	data, err := os.ReadFile("rawJson/Netbox/sites.json")
	if err != nil {
		return nil, err
	}
	sitePage := &models.SitePage{}
	err = json.Unmarshal(data, sitePage)
	if err != nil {
		return nil, err
	}
	return sitePage.Results, nil
}
