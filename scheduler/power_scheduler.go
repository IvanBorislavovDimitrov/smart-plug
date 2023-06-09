package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
	plugclient "github.com/IvanBorislavovDimitrov/smart-charger/plug_client"
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
	service_model "github.com/IvanBorislavovDimitrov/smart-charger/service/model"
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
	log.Println("Reconciliation of the plugs - START")
	page := 1
	plugs := ps.listPlugs(page)
	for len(plugs) != 0 {
		ps.reconcilePlugsBatch(plugs)
		page++
		plugs = ps.listPlugs(page)
	}
	log.Println("Reconciliation of the plugs - END")
}

func (ps *PowerScheduler) listPlugs(page int) []*model.Plug {
	ctx := context.Background()
	plugs, err := ps.plugService.ListPlugs(ctx, page, 50)
	if err != nil {
		log.Fatal("Error while listing plugs", err)
	}
	return service_model.FromServiceModelPlugsToPlugs(plugs)
}

func (ps *PowerScheduler) reconcilePlugsBatch(plugs []*model.Plug) {
	plugChannel := make(chan *model.Plug)
	go func() {
		wg := sync.WaitGroup{}
		for _, plug := range plugs {
			wg.Add(1)
			go ps.reconcilePlugState(plug, &wg, plugChannel)
		}
		wg.Wait()
		close(plugChannel)
	}()
	for plug := range plugChannel {
		log.Println(fmt.Sprintf("Plug with IP: %s and name: %s was updated, state: %s", plug.IPAddress, plug.Name, plug.State))
	}
}

func (ps *PowerScheduler) reconcilePlugState(plug *model.Plug, wg *sync.WaitGroup, plugChannel chan *model.Plug) {
	defer wg.Done()
	power, err := plugclient.GetCurrentlyUsedPowerInWatts(*plug)
	if err != nil {
		log.Println("Could not reconcile state of the plug. Error getting data!", err)
		return
	}
	log.Println(fmt.Sprintf("Plug with id: %s is consuming power of %f", plug.IPAddress, power))
	if power <= plug.PowerToTurnOff {
		log.Println(fmt.Sprintf("Turning off plug: %s", plug.Name))
		err = plugclient.TurnOffPlug(*plug)
		if err != nil {
			log.Println("Error while turning on plug: ", err)
		}
		plug.State = "OFF"
	} else {
		log.Println(fmt.Sprintf("Updating DB state. Plug in turned on: %s", plug.Name))
		plug.State = "ON"
	}
	serviceModelPlug, err := ps.plugService.UpdatePlug(context.Background(), service_model.ToServiceModelPlugFromPlug(*plug))
	if err != nil {
		log.Println("Could not reconcile state of the plug. Error updating plug!", err)
		return
	}
	updatedPlug := service_model.FromServiceModelPlugToPlug(*serviceModelPlug)
	plugChannel <- &updatedPlug
}

func (ps *PowerScheduler) TurnOnPlugs() {
	log.Println("Turning on plugs - START")
	page := 1
	plugs := ps.listPlugs(page)
	for len(plugs) != 0 {
		ps.turnOnPlugsBatch(plugs)
		page++
		plugs = ps.listPlugs(page)
	}
	log.Println("Turning on plugs - END")
}

func (ps *PowerScheduler) turnOnPlugsBatch(plugs []*model.Plug) {
	plugChannel := make(chan *model.Plug)
	go func() {
		wg := sync.WaitGroup{}
		for _, plug := range plugs {
			wg.Add(1)
			go ps.turnOnPlug(plug, &wg, plugChannel)
		}
		wg.Wait()
		close(plugChannel)
	}()
	for plug := range plugChannel {
		log.Println(fmt.Sprintf("Plug with IP: %s and name: %s was updated, state: %s", plug.IPAddress, plug.Name, plug.State))
	}
}

func (ps *PowerScheduler) turnOnPlug(plug *model.Plug, wg *sync.WaitGroup, plugChannel chan *model.Plug) {
	defer wg.Done()
	err := plugclient.TurnOnPlug(*plug)
	if err != nil {
		log.Println("Could not start plug", err)
		return
	}
	plug.State = "ON"
	serviceModelPlug, err := ps.plugService.UpdatePlug(context.Background(), service_model.ToServiceModelPlugFromPlug(*plug))
	if err != nil {
		log.Println("Could not update plug", err)
		return
	}
	updatedPlug := service_model.FromServiceModelPlugToPlug(*serviceModelPlug)
	plugChannel <- &updatedPlug
}
