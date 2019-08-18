package models

// import (
// 	"time"

// 	tools "github.com/sjljrvis/deploynow/tools"

// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// )

// const (
// 	collection = "repository"
// )

// // Repository model schema
// type Repository struct {
// 	ID             bson.ObjectId          `bson:"_id" json:"_id"`
// 	RepositoryName string                 `bson:"repositoryName" json:"repositoryName"`
// 	UserID         bson.ObjectId          `bson:"userId" json:"userId"`
// 	Language       string                 `bson:"language" json:"language"`
// 	Path           string                 `bson:"path" json:"path"`
// 	PathDocker     string                 `bson:"pathDocker" json:"pathDocker"`
// 	Date           time.Time              `bson:"date" json:"date"`
// 	Description    string                 `bson:"description" json:"description"`
// 	State          string                 `bson:"state" json:"state"`
// 	Github         map[string]interface{} `bson:"github" json:"github"`
// }

// /*
// CreateIndex creates  unique index for a collection
// */
// func CreateIndex() error {
// 	for _, key := range []string{"repositoryName"} {
// 		index := mgo.Index{
// 			Key:    []string{key},
// 			Unique: true,
// 		}
// 		if err := tools.MongoDB.C(collection).EnsureIndex(index); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// /*
// FindAll runs empty search and returns all documents from collection
// */
// func FindAll(userID string) ([]Repository, error) {
// 	var results []Repository
// 	err := tools.MongoDB.C(collection).Find(bson.M{"userId": bson.ObjectIdHex(userID)}).All(&results)
// 	return results, err
// }

// /*
// FindOneByID runs query based on _id  and returns single document from collection
// */
// func FindOneByID(userID string, id string) (Repository, error) {
// 	var result Repository

// 	err := tools.MongoDB.C(collection).FindId(bson.M{
// 		"userId": bson.ObjectIdHex(userID),
// 		"_id":    bson.ObjectIdHex(id)}).One(&result)

// 	return result, err
// }

// /*
// FindOne runs query based on {query} and returns single document from collection
// */
// func FindOne(query map[string]interface{}) (Repository, error) {
// 	var result Repository
// 	err := tools.MongoDB.C(collection).Find(query).One(&result)
// 	return result, err
// }

// /*
// Create runs query based on _id  and returns single document from collection
// */
// func Create(repository Repository) error {
// 	err := tools.MongoDB.C(collection).Insert(&repository)
// 	return err
// }

// /*
// FindByQuery runs query based on query  and returns array of document from collection
// */
// func FindByQuery(query map[string]interface{}) ([]Repository, error) {
// 	var results []Repository
// 	err := tools.MongoDB.C(collection).Find(query).All(&results)
// 	return results, err
// }
