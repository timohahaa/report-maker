package device

import (
	"encoding/json"
	"os"

	"github.com/timohahaa/report-maker/internal/models"
)

type DeviceJSONRepo struct {
}

func (rep *DeviceJSONRepo) GetAllDevices() ([]models.Device, error) {
	//just read from the json file
	data, err := os.ReadFile("rawJson/Netbox/devices.json")
	if err != nil {
		return nil, err
	}
	devicePage := &models.DevicePage{}
	err = json.Unmarshal(data, devicePage)
	if err != nil {
		return nil, err
	}
	return devicePage.Results, nil
}
