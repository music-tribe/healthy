package healthy

import (
	"github.com/globalsign/mgo"
	"github.com/music-tribe/errors"
)

type Pinger interface {
	Ping() error
}

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

const stubPingerErrorMessage = "stubPinger: if you're seeing this message, you've not added a `Pinger` to your `Checker`!"

type stubPinger struct{}

func (sp *stubPinger) Ping() error {
	return errors.NewCloudError(400, stubPingerErrorMessage)
}

type MockPinger struct {
	Err error
}

func (mp *MockPinger) Ping() error {
	return mp.Err
}
