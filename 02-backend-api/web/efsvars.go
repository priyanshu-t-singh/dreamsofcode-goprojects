package web

import "embed"

// staticFiles embeds the entire web directory contents into the binary.
//
//go:embed static/*
var staticFiles embed.FS

//go:embed templates/index.html
var indexHTML []byte
