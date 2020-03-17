package main

import (
	"log"
	// "github.com/ChiragSehra/custom-logistic-regression-in-go/logisitcregression"
)

func main() {
	if err := logisticregression.Run(); err != nil {
		log.Fatal(err)
	}
}
