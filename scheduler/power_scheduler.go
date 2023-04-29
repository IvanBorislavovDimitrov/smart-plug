package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
	plugclient "github.com/IvanBorislavovDimitrov/smart-charger/plug_client"
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
)

type PowerScheduler struct {
	plugService *service.PlugService
}

func NewPowerScheduler(plugService *service.PlugService) *PowerScheduler {
	return &PowerScheduler{
		plugService: plugService,
	}
}

func (ps *PowerScheduler) ReconcilePlugsStates() {
	page := 1
	plugs := ps.listPlugs(page)
	for len(plugs) != 0 {
		ps.reconcilePlugsBatch(plugs)
		page++
		plugs = ps.listPlugs(page)
	}
}

func (ps *PowerScheduler) listPlugs(page int) []*model.Plug {
	ctx := context.Background()
	plugs, err := ps.plugService.ListPlugs(ctx, page, 50)
	if err != nil {
		log.Fatal("Error while listing plugs", err)
	}
	return plugs
}

func (ps *PowerScheduler) reconcilePlugsBatch(plugs []*model.Plug) {
	plugChannel := make(chan *model.Plug)
	for _, plug := range plugs {
		go ps.reconcilePlugState(plug, plugChannel)
	}
	for plug := range plugChannel {
		log.Println("Plug with IP was updated: " + plug.IPAddress)
	}
}

func (ps *PowerScheduler) reconcilePlugState(plug *model.Plug, plugChannel chan *model.Plug) {
	power, err := plugclient.GetCurrentlyUsedPowerInWatts(*plug)
	if err != nil {
		log.Println("Could not reconcile state of the plug", err)
	}
	log.Println(fmt.Sprintf("Plug with id: %s is consumes power of %f", plug.IPAddress, power))
	if power <= plug.PowerToTurnOff {
		err = plugclient.TurnOffPlug(*plug)
		plug.State = "OFF"
		ps.plugService.UpdatePlug(context.Background(), plug)
		if err != nil {
			log.Println("Could not reconcile state of the plug", err)
		}
	}
	plugChannel <- plug
}

func (ps *PowerScheduler) TurnOnPlugs() {
	page := 1
	plugs := ps.listPlugs(page)
	for len(plugs) != 0 {

		page++
		plugs = ps.listPlugs(page)
	}
}
