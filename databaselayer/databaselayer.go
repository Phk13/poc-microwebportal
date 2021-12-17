package databaselayer

import "errors"

const (
	MYSQL uint8 = iota
	POSTGRESQL
	MONGODB
)

type DinoDBHandler interface {
	GetAvailableDinos() ([]Animal, error)
	GetDinoByNickname(string) (Animal, error)
	GetDinosByType(string) ([]Animal, error)
	AddAnimal(Animal) error
	UpdateAnimal(Animal, string) error
}

type Animal struct {
	ID         int    `bson:"-"`
	AnimalType string `bson:"animal_type"`
	Nickname   string `bson:"nickname"`
	Zone       int    `bson:"zone"`
	Age        int    `bson:"age"`
}

var ErrorDBTypeNotSupported = errors.New("the database type provided is not supported")

func GetDatabaseHandler(dbtype uint8, connection string) (DinoDBHandler, error) {
	switch dbtype {
	case MYSQL:
		return NewMySQLHandler(connection)
	case MONGODB:
		return NewMongoDBHandler(connection)
	case POSTGRESQL:
		return NewPQHandler(connection)
	}
	return nil, ErrorDBTypeNotSupported
}
