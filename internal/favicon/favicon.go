package favicon

import (
	"embed"
)

//go:embed favicon.ico
var faviconFile embed.FS

func GetFavicon() ([]byte, error) {
	return faviconFile.ReadFile("favicon.ico")
}
