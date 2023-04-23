package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
	"github.com/google/uuid"
)

type PlugService struct {
	db *sql.DB
}

func NewPlugService(db *sql.DB) *PlugService {
	return &PlugService{
		db: db,
	}
}

func (ps *PlugService) AddPlug(plug model.NewPlug) (*model.Plug, error) {
	transaction, err := ps.db.Begin()
	if err != nil {
		return nil, err
	}
	uuid := fmt.Sprint(uuid.New())
	res, err := transaction.Exec("INSERT INTO plugs VALUES ($1, $2, $3, $4)", uuid, plug.Name, plug.IPAddress, time.Now())
	if err != nil {
		transaction.Rollback()
		log.Println(err)
		return nil, err
	}
	log.Println("Last insert ID" + fmt.Sprint(res.LastInsertId()))
	transaction.Commit()
	createdPlug := &model.Plug{
		ID:             uuid,
		Name:           plug.Name,
		IPAddress:      plug.IPAddress,
		PowerToTurnOff: plug.PowerToTurnOff,
	}
	return createdPlug, nil
}
