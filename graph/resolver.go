package graph

import (
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	plugService *service.PlugService
}

func NewResolver(plugService *service.PlugService) *Resolver {
	return &Resolver{
		plugService: plugService,
	}
}
