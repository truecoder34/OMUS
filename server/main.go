package main

import (
	"OMUS/server/helper"
	"log"
)

func main() {
	//app.StartApp()
	var data2Encode uint64 = 0
	encoded_number := helper.Encode(data2Encode)
	log.Printf("[RESULT] Enocde() - %s", encoded_number)

	decode, err := helper.Decode("godevblogusinggomodules")
	if err != nil {
		log.Fatalln("")
	}
	log.Printf("[RESULT] Decode() - %d", decode)

}
