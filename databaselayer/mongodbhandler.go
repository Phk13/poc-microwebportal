package databaselayer

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBHandler struct {
	*mgo.Session
}

func NewMongoDBHandler(connection string) (*MongoDBHandler, error) {
	s, err := mgo.Dial(connection)
	return &MongoDBHandler{
		Session: s,
	}, err
}

func (handler *MongoDBHandler) GetAvailableDinos() ([]Animal, error) {
	s := handler.getFreshSession()
	defer s.Close()
	animals := []Animal{}
	err := s.DB("micro").C("animals").Find(nil).All(&animals)
	return animals, err
}

func (handler *MongoDBHandler) GetDinoByNickname(nickname string) (Animal, error) {
	s := handler.getFreshSession()
	defer s.Close()
	a := Animal{}
	err := s.DB("micro").C("animals").Find(bson.M{"nickname": nickname}).One(&a)
	return a, err
}
func (handler *MongoDBHandler) GetDinosByType(dinoType string) ([]Animal, error) {
	s := handler.getFreshSession()
	defer s.Close()
	animals := []Animal{}
	err := s.DB("micro").C("animals").Find(bson.M{"animal_type": dinoType}).All(&animals)
	return animals, err
}
func (handler *MongoDBHandler) AddAnimal(a Animal) error {
	s := handler.getFreshSession()
	defer s.Close()
	return s.DB("micro").C("animals").Insert(a)
}
func (handler *MongoDBHandler) UpdateAnimal(a Animal, nickname string) error {
	s := handler.getFreshSession()
	defer s.Close()
	return s.DB("micro").C("animals").Update(bson.M{"nickname": nickname}, a)
}

func (handler *MongoDBHandler) getFreshSession() *mgo.Session {
	return handler.Session.Copy()
}
