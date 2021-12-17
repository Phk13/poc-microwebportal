package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type animal struct {
	id         int
	animalType string
	nickname   string
	zone       int
	age        int
}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=abc123.. host=172.17.0.9 port=5432 dbname=micro sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM animals WHERE ID=$1", 1)
	handlerows(rows, err)

	row := db.QueryRow("SELECT * FROM animals WHERE AGE > $1", 10)
	a := animal{}
	err = row.Scan(&a.id, &a.animalType, &a.nickname, &a.zone, &a.age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)

	// result, err := db.Exec("INSERT INTO animals (ANIMAL_TYPE, NICKNAME, ZONE, AGE) VALUES ('Carnotaurus', 'Carno', 3, 22)")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result.LastInsertId()) //not supported
	// fmt.Println(result.RowsAffected())

	// var id int

	// db.QueryRow("UPDATE animals SET AGE = $1 WHERE ID = $2 RETURNING ID", 16, 2).Scan(&id)
	// fmt.Println(id)

	fmt.Println("Statements")
	stmt, err := db.Prepare("select * from animals where age > $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err = stmt.Query(5)
	handlerows(rows, err)

	rows, err = stmt.Query(10)
	handlerows(rows, err)

	testTransaction(db)
}

func handlerows(rows *sql.Rows, err error) {
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
}

func testTransaction(db *sql.DB) {
	fmt.Println("Transactions")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("select * from animals where age > $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(15)
	handlerows(rows, err)

	rows, err = stmt.Query(10)
	handlerows(rows, err)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
