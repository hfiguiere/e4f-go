package main

import (
	"os"
	"fmt"
	"log"
)

func (db *E4fDb) Print(roll *ExposedRoll) {
	fmt.Printf("%s %d\n", roll.FilmType, roll.Iso)

}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Not enough arg")
	}

	path := os.Args[1]
	e4fDb := Parse(path)
	e4fDb.buildMaps()

	for _, roll := range e4fDb.ExposedRolls {
		id := roll.Id
		fmt.Println("Roll:")
		e4fDb.Print(roll)
		exps := e4fDb.exposuresForRoll(id)
		for _, exp := range exps {
			fmt.Printf("Exposure %d: ", exp.Number)
			fmt.Println(exp)
		}
	}
}
