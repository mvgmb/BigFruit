package main

import (
	"log"

	bigfruit "github.com/mvgmb/BigFruit/bigfruit"
)

func main() {
	client := bigfruit.NewBigFruitClient()

	err := client.DownloadBigFile("data.txt", "./")
	if err != nil {
		log.Println(err)
	}
}
