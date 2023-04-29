package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IvanBorislavovDimitrov/smart-charger/graph/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlugService struct {
	pool *pgxpool.Pool
}

func NewPlugService(pool *pgxpool.Pool) *PlugService {
	return &PlugService{
		pool: pool,
	}
}

func (ps *PlugService) AddPlug(ctx context.Context, plug model.NewPlug) (*model.Plug, error) {
	poolCon, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(context.Background())
	transaction, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	uuid := fmt.Sprint(uuid.New())
	now := time.Now().UTC()
	plugState := "OFF"
	_, err = transaction.Exec(ctx, "INSERT INTO plugs VALUES ($1, $2, $3, $4, $5, $6)", uuid, plug.Name, plug.IPAddress, plug.PowerToTurnOff, now, plugState)
	if err != nil {
		transaction.Rollback(ctx)
		log.Println(err)
		return nil, err
	}
	err = transaction.Commit(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &model.Plug{
		ID:             uuid,
		Name:           plug.Name,
		IPAddress:      plug.IPAddress,
		PowerToTurnOff: plug.PowerToTurnOff,
		CreatedAt:      fmt.Sprint(now),
	}, nil
}

func (ps *PlugService) UpdatePlug(ctx context.Context, plug *model.Plug) error {
	poolCon, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(context.Background())
	transaction, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = transaction.Exec(ctx, "UPDATE plugs SET name = $1, ip_address = $2, power_to_turn_off = $3, state = $4 WHERE id = $5",
		plug.Name, plug.IPAddress, plug.PowerToTurnOff, plug.State, plug.ID)
	if err != nil {
		transaction.Rollback(ctx)
		log.Println(err)
		return err
	}
	err = transaction.Commit(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ps *PlugService) ListPlugs(ctx context.Context, page, perPage int) ([]*model.Plug, error) {
	poolCon, err := ps.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(context.Background())
	offset := (page - 1) * perPage
	rows, err := conn.Query(ctx, "SELECT id, name, ip_address, power_to_turn_off, created_at, state FROM plugs offset $1 limit $2 ", offset, perPage)
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

func readPlug(rows pgx.Rows) (*model.Plug, error) {
	var id string
	var name string
	var ipAddress string
	var createdAt time.Time
	var powerToTurnOff float64
	var state string
	err := rows.Scan(&id, &name, &ipAddress, &powerToTurnOff, &createdAt, &state)
	if err != nil {
		return nil, err
	}
	return &model.Plug{
		ID:             id,
		Name:           name,
		IPAddress:      ipAddress,
		PowerToTurnOff: powerToTurnOff,
		CreatedAt:      fmt.Sprint(createdAt),
		State:          state,
	}, nil
}
