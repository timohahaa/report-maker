package asset

import (
	"encoding/json"
	"os"

	"github.com/timohahaa/report-maker/internal/models"
)

type AssetJSONRepo struct {
}

func (rep *AssetJSONRepo) GetAllAssets() ([]models.Asset, error) {
	//just read from the json file
	data, err := os.ReadFile("rawJson/SnipeIT/assets.json")
	if err != nil {
		return nil, err
	}
	AssetList := &models.AssetList{}
	err = json.Unmarshal(data, AssetList)
	if err != nil {
		return nil, err
	}
	return AssetList.Assets, nil
}
