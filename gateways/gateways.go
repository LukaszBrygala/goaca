package gateways

import (
	"errors"

	"database/sql"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db}
}

func (r Repository) Find(gatewayID, brokerType string) (Gateway, error) {
	var gateway Gateway

	row := r.db.QueryRow(
		`SELECT * FROM gateways
		WHERE gatewayid = $1 AND brokertype = $2`,
		gatewayID, brokerType,
	)

	err := row.Scan(
		&gateway.ID,
		&gateway.Password,
		&gateway.BrokerType,
		&gateway.Environment,
		&gateway.BrokerKey,
		&gateway.Meta,
	)

	if err == sql.ErrNoRows {
		return gateway, ErrNotFound
	}
	if err != nil {
		return gateway, err
	}

	return gateway, nil
}

type Gateway struct {
	ID          string `json:"gatewayId"`
	Password    string `json:"password"`
	BrokerType  string `json:"brokerType"`
	Environment string `json:"environment"`
	BrokerKey   string `json:"brokerKey"`
	Meta        string `json:"meta"`
}
