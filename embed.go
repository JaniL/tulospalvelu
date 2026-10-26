package tulospalvelu

import "embed"

//go:embed assets vendor index.html
var StaticFS embed.FS
