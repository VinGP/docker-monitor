package repo

import (
	"backend/internal/model"
	"backend/pkg/postgres"
	"log"
)

type ContainerStatusRepo struct {
	*postgres.Postgres
}

const containerStatusTableName = "container_status"

func createContainerStatusTable(pg *postgres.Postgres) {
	sql := `
	CREATE TABLE IF NOT EXISTS container_status (
	    ip VARCHAR(45) PRIMARY KEY,
	    ping_time FLOAT,
	    last_success TIMESTAMP NULL
	)
	    `

	_, err := pg.DB.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

func New(pg *postgres.Postgres) *ContainerStatusRepo {
	createContainerStatusTable(pg)
	return &ContainerStatusRepo{pg}
}

func (r *ContainerStatusRepo) UpsertContainerStatus(status *model.ContainerStatus) error {
	sb := r.Builder.
		Insert(containerStatusTableName).
		Columns("ip", "ping_time", "last_success").
		Values(status.IPAddress, status.PingTime, status.LastSuccess).
		Suffix(`
			ON CONFLICT (ip) 
			DO UPDATE SET 
				ping_time = EXCLUDED.ping_time,
				last_success = EXCLUDED.last_success
		`)

	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ContainerStatusRepo) GetAll() ([]model.ContainerStatus, error) {
	sb := r.Builder.Select("*").From(containerStatusTableName)
	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var containers = make([]model.ContainerStatus, 0)
	err = r.DB.Select(&containers, query, args...)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (r *ContainerStatusRepo) DeleteAll() error {
	sb := r.Builder.Delete(containerStatusTableName)
	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
