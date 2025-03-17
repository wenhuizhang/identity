package assets

import (
	"embed"
	"io/fs"
)

//go:embed web/*
var embeddedFiles embed.FS

// Ensure only the `static/` subdirectory is served correctly
var StaticFiles, _ = fs.Sub(embeddedFiles, "web")
