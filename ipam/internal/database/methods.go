package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

var (
	EmptySelect             error = errors.New("no select rows")
	PrefixTypeNotRecognized error = errors.New("prefix type not recognized")
	DeadContext             error = errors.New("context is dead")
)

const (
	P2P int = iota
	DELEGATED
)

func (i *Instance) AcquirePrefix(ctx context.Context, prefixLen int64) (string, error) {

	select {
	case <-ctx.Done():
		return "", DeadContext
	default:
		var table string
		query := "SELECT prefix FROM %s WHERE released=true LIMIT 1;"
		switch prefixLen {
		case 127:
			table = "p2p"
		case 64:
			table = "delegated"
		}
		val, err := i.Conn.Query(ctx, fmt.Sprintf(query, table))
		if err == pgx.ErrNoRows {
			return "", EmptySelect
		}
		if err != nil {
			return "", fmt.Errorf("select error: %v", err)
		}
		defer val.Close()
		var prefix string
		if val.Next() {
			err = val.Scan(&prefix)
			if err != nil {
				return "", fmt.Errorf("not able to Scan result: %v", err)
			}
		}

		val, err = i.Conn.Query(ctx, fmt.Sprintf("UPDATE %s SET released=false WHERE prefix=$1;", table), prefix)
		defer val.Close()
		if err != nil {
			return "", fmt.Errorf("select/update error: %v", err)
		}
		return prefix, nil
	}
}

func (i *Instance) ReleasePrefix(ctx context.Context, prefix string, prefixLen int64) error {
	select {
	case <-ctx.Done():
		return DeadContext
	default:
		var table string
		query := "UPDATE %s SET released=true WHERE prefix=$1;"
		switch prefixLen {
		case 127:
			table = "p2p"
		case 64:
			table = "delegated"
		}
		val, err := i.Conn.Query(ctx, fmt.Sprintf(query, table), fmt.Sprintf("%s/%d", prefix, prefixLen))
		defer val.Close()
		if err == pgx.ErrNoRows {
			return EmptySelect
		}
		if err != nil {
			return fmt.Errorf("select error: %v", err)
		}
		return nil
	}

}
