package entity

type PlugEntity struct {
	WifiStatus WifiStatus `json:"wifi_sta,omitempty"`
	Cloud      Cloud      `json:"cloud,omitempty"`
	Mqtt       Mqtt       `json:"mqtt,omitempty"`
	Time       string     `json:"time,omitempty"`
	Unixtime   string     `json:"unixtime,omitempty"`
	Meters     []Meter    `json:"meters,omitempty"` // Should always be one
}

type WifiStatus struct {
	Connected bool   `json:"connected,omitempty"`
	SSID      string `json:"ssid,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Rssi      int    `json:"rssi,omitempty"`
}

type Cloud struct {
	Enabled   bool `json:"enabled,omitempty"`
	Connected bool `json:"connected,omitempty"`
}

type Mqtt struct {
	Connected bool `json:"connected,omitempty"`
}

type Meter struct {
	Power     float64 `json:"power,omitempty"`
	Overpower float64 `json:"overpower,omitempty"`
}
