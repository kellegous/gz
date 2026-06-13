package internal

import (
	_ "embed"

	"github.com/kellegous/gz/internal/store/schema"
)

//go:embed 0001.up.sql
var up0001 string

//go:embed 0001.dn.sql
var dn0001 string

var Migrations = []schema.Migration{
	{Up: up0001, Down: dn0001},
}
