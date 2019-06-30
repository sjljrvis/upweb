package db

import (
	mgo "gopkg.in/mgo.v2"
	"github.com/sjljrvis/deploynow/log"
)

//MongoDB -Connecting to DB
var MongoDB *mgo.Database

//MongoConnect -Connecting to DB
func MongoConnect(mongoURI, dbName string) {
	log.Info().Msg("Connecting to mongo")
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		MongoDB = session.DB(dbName)
		log.Info().Msgf("Connected to mongo :" + dbName)

		log.Info().Msgf("loading Indexes - started")
  	MongoDB.C("user").EnsureIndex(userModelIndex())
		log.Info().Msgf("loading Indexes -finished")
	}
}

func userModelIndex() mgo.Index {
  return mgo.Index{
    Key:        []string{"userName" , "email"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }
}
