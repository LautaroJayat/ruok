package main

import (
	m20231117133101 "github.com/back-end-labs/ruok/cmd/migrate/migrations/20231117133101"
	m20231117172601 "github.com/back-end-labs/ruok/cmd/migrate/migrations/20231117172601"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/storage"
)

func main() {

	migrations := []func(storage.Storage){
		m20231117133101.Migrate20231117133101,
		m20231117172601.Migrate20231117172601,
	}
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()
	for i := 0; i < len(migrations); i++ {
		migrations[i](s)
	}

}
