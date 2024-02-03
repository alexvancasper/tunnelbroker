package serviceHanlders

import (
	"time"

	db "github.com/alexvancasper/TunnelBroker/ipam/internal/database"

	"github.com/sirupsen/logrus"
)

// Handler set of variables for DB management
type Handler struct {
	//Database connection instance
	DB *db.Instance
	//Log instance
	Log *logrus.Logger
	//Timeout - database timeout
	Timeout time.Duration
}

var SHandler Handler
