// Package assets contains the static files served by the web application.
package assets

import "embed"

// Files includes the Mini App and website assets in the production binary.
//
//go:embed miniapp/* website/*
var Files embed.FS
