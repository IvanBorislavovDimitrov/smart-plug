// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewPlug struct {
	Name           string  `json:"name"`
	IPAddress      string  `json:"ipAddress"`
	PowerToTurnOff float64 `json:"powerToTurnOff"`
}

type Plug struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	IPAddress      string  `json:"ipAddress"`
	PowerToTurnOff float64 `json:"powerToTurnOff"`
	CreatedAt      string  `json:"createdAt"`
	State          string  `json:"state"`
}

type UpdatedPlug struct {
	ID             string   `json:"id"`
	Name           *string  `json:"name,omitempty"`
	IPAddress      *string  `json:"ipAddress,omitempty"`
	PowerToTurnOff *float64 `json:"powerToTurnOff,omitempty"`
}
