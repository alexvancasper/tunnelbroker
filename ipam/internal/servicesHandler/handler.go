package serviceHanlders

import (
	"context"
	"time"

	psql "github.com/alexvancasper/TunnelBroker/ipam/internal/database"
	"github.com/alexvancasper/TunnelBroker/ipam/internal/handler/operations"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
)

func GetAcquirePrefixlenHandler(params operations.GetAcquireParams) middleware.Responder {

	l := SHandler.Log.WithFields(logrus.Fields{
		"function": "GetAcquirePrefixlenHandler",
	})
	l.Debugf("acquired prefix length: %d", params.Prefixlen)
	l.Info("data received")

	ctx, ctxCancelFunc := context.WithTimeout(context.Background(), SHandler.Timeout*time.Millisecond)
	defer ctxCancelFunc()
	prefix, err := SHandler.DB.AcquirePrefix(ctx, params.Prefixlen)
	if err != nil && err != psql.EmptySelect {
		l.Errorf("upredictable error: %s", err)
		return &operations.GetAcquireInternalServerError{}
	}
	if err == psql.EmptySelect {
		l.Errorf("The DB does not have available prefixes %s", err)
		return &operations.GetAcquireInternalServerError{}
	}

	return &operations.GetAcquireCreated{Payload: prefix}
}

func DeleteReleasePrefixHandler(params operations.DeleteReleaseParams) middleware.Responder {
	l := SHandler.Log.WithFields(logrus.Fields{
		"function": "DeleteReleasePrefixHandler",
	})
	l.Debugf("releasing prefix: %s length: %d", params.Prefix, params.Prefixlen)
	l.Info("data received")

	ctx, ctxCancelFunc := context.WithTimeout(context.Background(), SHandler.Timeout*time.Millisecond)
	defer ctxCancelFunc()
	err := SHandler.DB.ReleasePrefix(ctx, params.Prefix, params.Prefixlen)
	if err != nil && err != psql.EmptySelect {
		l.Errorf("upredictable error: %s", err)
		return &operations.DeleteReleaseBadRequest{}
	}
	if err == psql.EmptySelect {
		l.Errorf("The DB does not have the prefix %s", err)
		return &operations.DeleteReleaseInternalServerError{}
	}

	return &operations.DeleteReleaseOK{}
}
