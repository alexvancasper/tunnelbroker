package psql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	Conn       *pgxpool.Pool
	PoolConfig *pgxpool.Config
}

// var ErrConnection error = errors.New("Database connection error")

// New - makes DB connection
// dsn - database string
// number - how many connections in pool
func New(dsn string, number int) (*Instance, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("Connect to DB error: %v", err)
	}
	poolConfig.MaxConns = int32(number)

	c, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	_, err = c.Exec(context.Background(), ";")
	if err != nil {
		return nil, fmt.Errorf("Ping failed: %v\n", err)
	}

	repo := Instance{Conn: c, PoolConfig: poolConfig}

	return &repo, nil
}

func (i *Instance) CloseConnection() {
	defer i.Conn.Close()
}
