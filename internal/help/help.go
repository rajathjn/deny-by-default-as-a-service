package help

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/utils"
)

type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type ContentNegotiation struct {
	Default     string   `json:"default"`
	JSONMethods []string `json:"json_methods"`
}

// Need this struct to ensure the JSON response has the correct format and fields
type HelpResponse struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Endpoints          []Endpoint         `json:"endpoints"`
	ContentNegotiation ContentNegotiation `json:"content_negotiation"`
}

func getHelpResponse() HelpResponse {
	response := HelpResponse{
		Name:        "Deny By Default as a Service",
		Description: "An API that returns random, generic, creative, and sometimes hilarious rejection reasons (and acceptances!).",
		Endpoints: []Endpoint{
			{Path: "/", Method: "GET", Description: "Returns a random \"no\" reason."},
			{Path: "/no", Method: "GET", Description: "Returns a random \"no\" reason."},
			{Path: "/yes", Method: "GET", Description: "Returns a random \"yes\" reason."},
			{Path: "/random", Method: "GET", Description: "Returns a random \"yes\" or \"no\" reason."},
			{Path: "/health", Method: "GET", Description: "Health check status."},
			{Path: "/help", Method: "GET", Description: "Shows this help information."},
			{Path: "/*", Method: "GET", Description: "Any other route returns a positive \"yes\" reason."},
		},
		ContentNegotiation: ContentNegotiation{
			Default: "plain text",
			JSONMethods: []string{
				"Add `?format=json` query parameter",
				"Set `Accept: application/json` header",
				"Set `Content-Type: application/json` header",
			},
		},
	}
	return response
}

func HelpText() string {
	response := getHelpResponse()
	var sb strings.Builder
	sb.WriteString(response.Name + "\n")
	sb.WriteString(response.Description + "\n\n")
	sb.WriteString("Endpoints:\n")
	for _, ep := range response.Endpoints {
		sb.WriteString(fmt.Sprintf("  %s %s - %s\n", ep.Method, ep.Path, ep.Description))
	}
	sb.WriteString("\nContent Negotiation:\n")
	sb.WriteString(fmt.Sprintf("  Default: %s\n", response.ContentNegotiation.Default))
	sb.WriteString("  JSON methods:\n")
	for _, m := range response.ContentNegotiation.JSONMethods {
		sb.WriteString(fmt.Sprintf("  - %s\n", m))
	}
	return sb.String()
}

func Help(c *gin.Context) {
	if utils.WantsJSON(c) {
		c.JSON(http.StatusOK, getHelpResponse())
		return
	}
	c.String(http.StatusOK, HelpText())
}
