package migrations

import "embed"

//go:embed migrations/*
var Migrations embed.FS
