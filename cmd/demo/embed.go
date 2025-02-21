//go:build embed

package main

import (
	"embed"

	"github.com/euiko/webapp/module/static"
)

//go:embed ui/dist
var embedFs embed.FS

func init() {
	// inject static files into webapp
	static.EmbedFS = embedFs
}
