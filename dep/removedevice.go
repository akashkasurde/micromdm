package dep

import (
	"github.com/pkg/errors"
)

const (
	devicedisownPath = "devices/disown"
)

type DeviceDisownRequest struct {
	Devices []string `json:"devices"`
}

// type DeviceDisownResponse struct {
// 	Devices map[string]string `json:"devices"`
// }

func (c *Client) DisownDevice(ddr *DeviceDisownRequest) (interface{}, error) {
	req, err := c.newRequest("POST", devicedisownPath, &ddr)
	if err != nil {
		return nil, errors.Wrap(err, "DeviceDisown request")
	}

	var response interface{}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "DeviceDisown")
}
