package main

import (
	// "github.com/rajathjn/deny-by-default-as-a-service/cmd"
	"github.com/rajathjn/deny-by-default-as-a-service/utils"
	"log"
)

var positive_reasons string = "assets/reasons_y.json"
var negative_reasons string = "assets/reasons_n.json"

func main() {
	reason_len, reasons_list := utils.ReadReasons(positive_reasons)
	log.Printf("The loaded reasons list of length: %d and has a value: %s", reason_len, utils.Get_random_reason(reason_len, reasons_list))
}
