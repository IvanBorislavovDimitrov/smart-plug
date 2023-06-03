package model

import (
	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
)

func ToServiceModelPlugFromNewPlug(newPlug model.NewPlug) Plug {
	return Plug{
		Name:           newPlug.Name,
		IPAddress:      newPlug.IPAddress,
		PowerToTurnOff: newPlug.PowerToTurnOff,
	}
}

func ToServiceModelPlugFromPlug(plug model.Plug) Plug {
	return Plug{
		ID:             plug.ID,
		Name:           plug.Name,
		IPAddress:      plug.IPAddress,
		PowerToTurnOff: plug.PowerToTurnOff,
		CreatedAt:      plug.CreatedAt,
		State:          plug.State,
	}
}

func ToServiceModelPlugFromUpdatedPlug(updatedPlug model.UpdatedPlug) Plug {
	return Plug{
		ID:             updatedPlug.ID,
		Name:           getStringOrReturnEmpty(updatedPlug.Name),
		IPAddress:      getStringOrReturnEmpty(updatedPlug.IPAddress),
		PowerToTurnOff: getFloatOrReturnNegative(updatedPlug.PowerToTurnOff),
	}
}

func FromServiceModelPlugToPlug(serviceModelPlug Plug) model.Plug {
	return model.Plug{
		ID:             serviceModelPlug.ID,
		Name:           serviceModelPlug.Name,
		IPAddress:      serviceModelPlug.IPAddress,
		PowerToTurnOff: serviceModelPlug.PowerToTurnOff,
		CreatedAt:      serviceModelPlug.CreatedAt,
		State:          serviceModelPlug.State,
	}
}

func FromServiceModelPlugsToPlugs(serviceModelPlugs []*Plug) []*model.Plug {
	var plugs []*model.Plug
	for _, serviceModelPlug := range serviceModelPlugs {
		plug := FromServiceModelPlugToPlug(*serviceModelPlug)
		plugs = append(plugs, &plug)
	}
	return plugs
}

func getStringOrReturnEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func getFloatOrReturnNegative(fl *float64) float64 {
	if fl == nil {
		return -1
	}
	return *fl
}
