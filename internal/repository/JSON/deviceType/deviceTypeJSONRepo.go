package deviceType

import (
	"encoding/json"
	"os"

	"github.com/timohahaa/report-maker/internal/models"
)

type DeviceTypeJSONRepo struct {
}

func (rep *DeviceTypeJSONRepo) GetAllDeviceTypes() ([]models.DeviceType, error) {
	//just read from the json file
	data, err := os.ReadFile("rawJson/Netbox/device_types.json")
	if err != nil {
		return nil, err
	}
	deviceTypePage := &models.DeviceTypePage{}
	err = json.Unmarshal(data, deviceTypePage)
	if err != nil {
		return nil, err
	}
	return deviceTypePage.Results, nil
}
