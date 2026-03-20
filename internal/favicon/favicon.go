package favicon

import (
	"embed"
)

//go:embed favicon.ico
var favion_file embed.FS

func Get_favicon() ([]byte, error) {
	return favion_file.ReadFile("favicon.ico")
}