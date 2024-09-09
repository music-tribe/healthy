package healthy

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/hellofresh/health-go/v5"
	"github.com/music-tribe/errors"
)

func mongoDBChecker(connString string) health.CheckFunc {
	return func(ctx context.Context) error {
		session, err := mgo.Dial(connString)
		if session != nil {
			defer session.Close()
		}
		if err != nil {
			return errors.NewCloudError(500, err.Error())
		}

		return nil
	}
}

func NewMongoDbCheckerWithConnectionString(name, dbConnString string) Checker {
	return NewChecker(name, mongoDBChecker(dbConnString))
}
