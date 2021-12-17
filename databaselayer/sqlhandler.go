package databaselayer

import (
	"database/sql"
	"fmt"
	"log"
)

type SQLHandler struct {
	*sql.DB
}

func (handler *SQLHandler) GetAvailableDinos() ([]Animal, error) {
	return handler.sendQuery("select id, animal_type, nickname, zone, age from animals")
}

func (handler *SQLHandler) GetDinoByNickname(nickname string) (Animal, error) {
	row := handler.QueryRow(fmt.Sprintf("select id, animal_type, nickname, zone, age from animals where nickname = '%s'", nickname))

	a := Animal{}
	err := row.Scan(&a.ID, &a.AnimalType, &a.Nickname, &a.Zone, &a.Age)
	return a, err
}

func (handler *SQLHandler) GetDinosByType(dinoType string) ([]Animal, error) {
	return handler.sendQuery(fmt.Sprintf("select id, animal_type, nickname, zone, age from animals where animal_type = '%s'", dinoType))
}

func (handler *SQLHandler) AddAnimal(a Animal) error {
	_, err := handler.Exec(fmt.Sprintf("insert into animals (animal_type, nickname, zone, age) values ('%s', '%s', %d, %d)", a.AnimalType, a.Nickname, a.Zone, a.Age))
	return err
}

func (handler *SQLHandler) UpdateAnimal(a Animal, nickname string) error {
	_, err := handler.Exec(fmt.Sprintf("update animals set animal_type = '%s', nickname = '%s', zone = %d, age=%d where nickname = '%s'", a.AnimalType, a.Nickname, a.Zone, a.Age, nickname))
	return err
}

func (handler *SQLHandler) sendQuery(q string) ([]Animal, error) {
	Animals := []Animal{}
	rows, err := handler.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		a := Animal{}
		err := rows.Scan(&a.ID, &a.AnimalType, &a.Nickname, &a.Zone, &a.Age)
		if err != nil {
			log.Println(err)
			continue
		}
		Animals = append(Animals, a)
	}
	return Animals, rows.Err()
}
