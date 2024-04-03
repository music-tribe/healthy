package healthy

import (
	"github.com/globalsign/mgo"
	"github.com/music-tribe/errors"
)

// MongoDbPinger allows us to use mongoDB connection strings to setup a pinger
type MongoDbPinger struct {
	connString string
}

func (dbp *MongoDbPinger) Ping() error {
	session, err := mgo.Dial(dbp.connString)
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return errors.NewCloudError(500, err.Error())
	}

	return nil
}
