package utils

import (
	"embed"
	"encoding/json"
	"log"
	"math/rand/v2"
	"strings"

	"github.com/gin-gonic/gin"
)

type Reasons struct {
	No  []string `json:"no"`
	Yes []string `json:"yes"`
}

//go:embed reasons.json
var reasonsFile embed.FS

var (
	reasons       Reasons
	lenNoReasons  int
	lenYesReasons int
)

func init() {
	jsonfile, err := reasonsFile.ReadFile("reasons.json")
	if err != nil {
		log.Fatalf("Error in reading the file: %v\n", err)
		panic("Exiting the Program.")
	}

	err = json.Unmarshal(jsonfile, &reasons)
	if err != nil {
		log.Fatalf("Unable to read the reasons list: %v", err)
		panic("Exiting the program")
	}

	lenNoReasons = len(reasons.No)
	lenYesReasons = len(reasons.Yes)
}

func GetNegativeReason() string {
	return reasons.No[rand.IntN(lenNoReasons)]
}

func GetPositiveReason() string {
	return reasons.Yes[rand.IntN(lenYesReasons)]
}

func WantsJSON(c *gin.Context) bool {
	if c.Query("format") == "json" {
		return true
	}
	accept := c.GetHeader("Accept")
	contentType := c.GetHeader("Content-Type")
	return strings.Contains(accept, "application/json") || strings.Contains(contentType, "application/json")
}
