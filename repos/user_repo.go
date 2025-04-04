package repo

import (
	"app/types"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// collection name for user
const (
	ACCOUNT_COLLECTION = "account"
)

// interface
// use the interface as a dependency for the next layer
type IUserRepo interface {
	LoginUser() (error, map[string]string)
	FindUser(email string) (*types.UserType, error)
	CreateUser(d types.UserRequestType) (string, error)
}

type MongoUser struct {
	dataBase *Mongo
}

func NewMongoUser(mongo *Mongo) *MongoUser {
	return &MongoUser{
		dataBase: mongo,
	}
}
func (m *MongoUser) FindUser(email string) (*types.UserType, error) {
	ctx, ctn := NewCTX()
	defer ctn()
	var UserBase types.UserType
	if err := m.dataBase.fromCollection(ACCOUNT_COLLECTION).FindOne(ctx, bson.M{"email": email}).Decode(&UserBase); err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {

			return nil, nil
		}
		return nil, err
	}
	return &UserBase, nil
}
func (m *MongoUser) CreateUser(d types.UserRequestType) (string, error) {
	ctx, ctn := NewCTX()
	defer ctn()

	res, err := m.dataBase.fromCollection(ACCOUNT_COLLECTION).InsertOne(ctx, d)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res.InsertedID), nil

}

func (m *MongoUser) LoginUser() (error, map[string]string) {
	coll := m.dataBase.fromCollection("")
	ctx, ctn := NewCTX()
	defer ctn()
	coll.FindOne(ctx, bson.M{
		"id":   "root",
		"role": "admin",
	})
	return nil, map[string]string{
		"id": "22",
	}
}
