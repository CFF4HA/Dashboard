package main

import (
	"fmt"
	"github.com/CFF4HA/Dashboard/internal/handlers/ingredient"
)

func main() {
	i, err := ingredient.RetrieveIngredient("water")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ingredient: %v\n", i)
}
