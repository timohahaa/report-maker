package models

//asset represents the Snipe-It asset

// asset list represents the response from Snipe-It API then requesting for assets
// header is an array of an array (weird), it holds fileds, that describe the assets (just skip it)
// data is an array of assets
type AssetList struct {
	//header [][]string `json:"header"`
	Assets []Asset `json:"data"`
}

type Asset struct {
	Name           string `json:"Asset Name"`
	Tag            int    `json:"Asset Tag"`
	Model          string `json:"Model"`
	Category       string `json:"Category"`
	Status         string `json:"Status"`
	Location       string `json:"Location"` //facility in NetBox site object
	Specifications string `json:"Specifications"`
}

type AssetRepository interface {
	GetAllAssets() ([]Asset, error)
}
