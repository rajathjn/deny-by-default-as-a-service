package utils

import (
	"encoding/json"
	"os"
	"math/rand/v2"
	"log"
)

func ReadReasons(reasons_file string) (int,[]string) {

	jsonfile, err := os.ReadFile(reasons_file)
	if err != nil {
		log.Fatalf("Error in reading the file: %v\n",err)
		panic("Exiting the Program.")
	}

	var reasons_list []string
	err = json.Unmarshal(jsonfile, &reasons_list)
	if err != nil {
		log.Fatalf("Unable to read the reasons list: %v", err)
		panic("Exiting the program")
	}

	return len(reasons_list), reasons_list
}

func Get_random_reason(reason_len int, reason_list []string) string {
	return reason_list[ rand.IntN(reason_len+1) ] // Range of rand.IntN is [0,n)
}
