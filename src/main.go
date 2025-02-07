package main

import (
	"deduper/src/internal"
	"fmt"
	"log"
)

func main() {

	service := internal.NewDeduperService()

	outputPath, err := service.FindDuplicates()

	if err != nil {
		log.Fatalf("Error processing contacts: %v", err)
	}

	fmt.Printf("Results have been saved to: %s\n", outputPath)
}
