package static

import "embed"

//go:embed assets/*  index.html
var Static embed.FS
