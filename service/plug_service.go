package service

import (
	"context"
	"fmt"
	"log"
	"time"

	service_model "github.com/IvanBorislavovDimitrov/smart-charger/service/model"
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

func (ps *PlugService) AddPlug(ctx context.Context, plug service_model.Plug) (*service_model.Plug, error) {
	poolCon, err := ps.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(ctx)
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
	return &service_model.Plug{
		ID:             uuid,
		Name:           plug.Name,
		IPAddress:      plug.IPAddress,
		PowerToTurnOff: plug.PowerToTurnOff,
		CreatedAt:      fmt.Sprint(now),
	}, nil
}

func (ps *PlugService) UpdatePlug(ctx context.Context, plug service_model.Plug) (*service_model.Plug, error) {
	poolCon, err := ps.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(ctx)
	transaction, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	plugForUpdate, err := ps.getPlugForUpdate(ctx, conn, plug)
	if err != nil {
		return nil, err
	}
	_, err = transaction.Exec(ctx, "UPDATE plugs SET name = $1, ip_address = $2, power_to_turn_off = $3, state = $4 WHERE id = $5",
		plugForUpdate.Name, plugForUpdate.IPAddress, plugForUpdate.PowerToTurnOff, plugForUpdate.State, plugForUpdate.ID)
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
	return ps.GetPlug(ctx, plug.ID)
}

func (ps *PlugService) getPlugForUpdate(ctx context.Context, conn *pgx.Conn, plug service_model.Plug) (*service_model.Plug, error) {
	row := conn.QueryRow(ctx, "SELECT id, name, ip_address, power_to_turn_off, created_at, state FROM plugs where id = $1", plug.ID)
	servicePlugModel, err := readPlug(row)
	if err != nil {
		return nil, err
	}
	if plug.IPAddress != "" {
		servicePlugModel.IPAddress = plug.IPAddress
	}
	if plug.Name != "" {
		servicePlugModel.Name = plug.Name
	}
	if plug.PowerToTurnOff >= 0 {
		servicePlugModel.PowerToTurnOff = plug.PowerToTurnOff
	}
	if plug.State != "" {
		servicePlugModel.State = plug.State
	}
	return servicePlugModel, nil
}

func (ps *PlugService) ListPlugs(ctx context.Context, page, perPage int) ([]*service_model.Plug, error) {
	poolCon, err := ps.pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(ctx)
	offset := (page - 1) * perPage
	rows, err := conn.Query(ctx, "SELECT id, name, ip_address, power_to_turn_off, created_at, state FROM plugs offset $1 limit $2 ", offset, perPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var plugs []*service_model.Plug
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

func (ps *PlugService) GetPlug(ctx context.Context, id string) (*service_model.Plug, error) {
	poolCon, err := ps.pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer poolCon.Release()
	conn := poolCon.Conn()
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT id, name, ip_address, power_to_turn_off, created_at, state FROM plugs where id = $1", id)
	return readPlug(row)
}

func readPlug(row pgx.Row) (*service_model.Plug, error) {
	var id string
	var name string
	var ipAddress string
	var createdAt time.Time
	var powerToTurnOff float64
	var state string
	err := row.Scan(&id, &name, &ipAddress, &powerToTurnOff, &createdAt, &state)
	if err != nil {
		return nil, err
	}
	return &service_model.Plug{
		ID:             id,
		Name:           name,
		IPAddress:      ipAddress,
		PowerToTurnOff: powerToTurnOff,
		CreatedAt:      fmt.Sprint(createdAt),
		State:          state,
	}, nil
}
