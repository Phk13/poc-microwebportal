package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type animal struct {
	id         int
	animalType string
	nickname   string
	zone       int
	age        int
}

func main() {
	db, err := sql.Open("mysql", "micro:micro123..@tcp(172.17.0.8:3306)/micro")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM micro.animals WHERE ID=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	animals := []animal{}
	for rows.Next() {
		a := animal{}
		err := rows.Scan(&a.id, &a.animalType, &a.nickname, &a.zone, &a.age)
		if err != nil {
			log.Println(err)
			continue
		}
		animals = append(animals, a)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(animals)

	row := db.QueryRow("SELECT * FROM micro.animals WHERE AGE > ?", 10)
	a := animal{}
	err = row.Scan(&a.id, &a.animalType, &a.nickname, &a.zone, &a.age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)

	result, err := db.Exec("INSERT INTO micro.animals (ANIMAL_TYPE, NICKNAME, ZONE, AGE) VALUES ('Carnotaurus', 'Carno', 3, 22)")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
}
