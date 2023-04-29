package plugclient

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/IvanBorislavovDimitrov/smart-charger/entity"
	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
)

var HttpClient *http.Client = &http.Client{}

func GetCurrentlyUsedPowerInWatts(plug model.Plug) (float64, error) {
	resp, err := HttpClient.Get("http://" + plug.IPAddress + "/status")
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	plugEntity := &entity.PlugEntity{}
	err = json.Unmarshal(body, plugEntity)
	if err != nil {
		return -1, err
	}
	return plugEntity.Meters[0].Power, nil
}

func TurnOnPlug(plug model.Plug) error {
	_, err := HttpClient.Get("http://" + plug.IPAddress + "/relay/0?turn=on")
	if err != nil {
		return err
	}
	return nil
}

func TurnOffPlug(plug model.Plug) error {
	_, err := HttpClient.Get(plug.IPAddress + "/relay/0?turn=off")
	if err != nil {
		return err
	}
	return nil
}
