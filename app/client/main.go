package main

import (
	"log"

	bigfruit "github.com/mvgmb/BigFruit/bigfruit"
)

func main() {
	client, err := bigfruit.NewBigFruitClient()
	if err != nil {
		log.Fatal(err)
	}

	err = client.DownloadBigFile("data.txt", "./")
	if err != nil {
		log.Println(err)
	}
}
