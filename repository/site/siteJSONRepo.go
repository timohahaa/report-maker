package site

import (
	"encoding/json"
	"os"
)

type SiteJSONRepo struct {
}

func (rep *SiteJSONRepo) GetAllSites() ([]Site, error) {
	//just read from the json file
	data, err := os.ReadFile("rawJson/Netbox/sites.json")
	if err != nil {
		return nil, err
	}
	sitePage := &SitePage{}
	err = json.Unmarshal(data, sitePage)
	if err != nil {
		return nil, err
	}
	return sitePage.Results, nil
}
