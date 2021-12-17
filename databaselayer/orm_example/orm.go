package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type animal struct {
	gorm.Model
	// ID         int    `gorm:"primary_key;not null;unique;AUTO_INCREMENT"`
	AnimalType string `gorm:"type:TEXT"`
	Nickname   string `gorm:"type:TEXT"`
	Zone       int    `gorm:"type:INTEGER"`
	Age        int    `gorm:"type:INTEGER"`
}

func main() {
	db, err := gorm.Open("postgres", "user=postgres password=abc123.. host=172.17.0.9 port=5432 dbname=micro sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.DropTableIfExists(&animal{})
	db.AutoMigrate(&animal{})

	a := animal{
		AnimalType: "Tyrannosaurus rex",
		Nickname:   "Rex",
		Zone:       1,
		Age:        11,
	}
	db.Table("animals").Create(&a)

	db.Save(&a)

	// db.Table("animal").Where("nickname = ? and zone = ?", "rapto", 2).Update("age", 16)
	animals := []animal{}
	db.Table("animals").Find(&animals, "age > ?", 10)
	fmt.Println(animals)
}
