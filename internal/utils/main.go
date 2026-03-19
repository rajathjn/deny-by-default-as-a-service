package utils

import (
	"encoding/json"
	"math/rand/v2"
	"log"
	"embed"
)

type Reasons struct {
	No 		[]string 	`json:"no"`
	Yes		[]string	`json:"yes"`
}

//go:embed reasons.json
var reasons_file embed.FS

var (
	reasons Reasons
	len_no_reasons int
	len_yes_reasons int
)

func init() {
	jsonfile, err := reasons_file.ReadFile("reasons.json")
	if err != nil {
		log.Fatalf("Error in reading the file: %v\n",err)
		panic("Exiting the Program.")
	}

	
	err = json.Unmarshal(jsonfile, &reasons)
	if err != nil {
		log.Fatalf("Unable to read the reasons list: %v", err)
		panic("Exiting the program")
	}

	len_no_reasons = len(reasons.No)
	len_yes_reasons = len(reasons.Yes)

} 

func Get_negative_reason() string {
	return reasons.No[ rand.IntN(len_no_reasons) ]
}

func Get_positive_reason() string {
	return reasons.Yes[ rand.IntN(len_yes_reasons) ]
}