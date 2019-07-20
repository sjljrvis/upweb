package db

import (
	"github.com/sjljrvis/deploynow/log"
	mgo "gopkg.in/mgo.v2"
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
		loadIndexes()
		log.Info().Msgf("loading Indexes -finished")
	}
}

func loadIndexes() {
	MongoDB.C("user").EnsureIndex(userModelIndex())
	MongoDB.C("repository").EnsureIndex(userModelIndex())
}

func userModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"userName", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func repositoryModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"repositoryName"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}
