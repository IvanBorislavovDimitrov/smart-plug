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
	now := time.Now()
	res, err := transaction.Exec("INSERT INTO plugs VALUES ($1, $2, $3, $4, $5)", uuid, plug.Name, plug.IPAddress, plug.PowerToTurnOff, now)
	if err != nil {
		transaction.Rollback()
		log.Println(err)
		return nil, err
	}
	log.Println("Last insert ID" + fmt.Sprint(res.LastInsertId()))
	transaction.Commit()
	return &model.Plug{
		ID:             uuid,
		Name:           plug.Name,
		IPAddress:      plug.IPAddress,
		PowerToTurnOff: plug.PowerToTurnOff,
		CreatedAt:      fmt.Sprint(now),
	}, nil
}

func (ps *PlugService) ListPlugs(page, perPage int) ([]*model.Plug, error) {
	offset := (page - 1) * perPage
	rows, err := ps.db.Query("SELECT id, name, ip_address, power_to_turn_off, created_at FROM plugs offset $1 limit $2 ", offset, perPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var plugs []*model.Plug
	for rows.Next() {
		plug, err := readPlug(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		plugs = append(plugs, plug)
	}
	return plugs, nil
}

func readPlug(rows *sql.Rows) (*model.Plug, error) {
	var id string
	var name string
	var ipAddress string
	var createdAt time.Time
	var powerToTurnOff float64
	err := rows.Scan(&id, &name, &ipAddress, &powerToTurnOff, &createdAt)
	if err != nil {
		return nil, err
	}
	return &model.Plug{
		ID:             id,
		Name:           name,
		IPAddress:      ipAddress,
		PowerToTurnOff: powerToTurnOff,
		CreatedAt:      fmt.Sprint(createdAt),
	}, nil
}
