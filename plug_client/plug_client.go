package plugclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IvanBorislavovDimitrov/smart-charger/entity"
	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
)

var HttpClient *http.Client = &http.Client{}

func GetCurrentlyUsedPowerInWatts(plug model.Plug) (float64, error) {
	resp, err := HttpClient.Get("http://" + plug.IPAddress + "/status")
	if err != nil {
		log.Println("Error during HTTP call")
		return -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Println("Status code: " + fmt.Sprint(resp.StatusCode))
		return -1, errors.New("Request failure")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading body")
		return -1, err
	}
	plugEntity := &entity.PlugEntity{}
	err = json.Unmarshal(body, plugEntity)
	if err != nil {
		log.Println("Error while unmarshalling")
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
	_, err := HttpClient.Get("http://" + plug.IPAddress + "/relay/0?turn=off")
	if err != nil {
		return err
	}
	return nil
}
