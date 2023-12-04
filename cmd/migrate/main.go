package migrations

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/spf13/cobra"
)

type migration struct {
	name string
	sql  string
}

//go:embed migrations/2023_12_04_041700_base_schema_n_fn.sql
var _2023_12_04_041700_base_schema_n_fn string

//go:embed migrations/2023_12_04_041701_base_tables.sql
var _2023_12_04_041701_base_tables string

//go:embed migrations/2023_12_04_041702_base_roles.sql
var _2023_12_04_041702_base_roles string

//go:embed migrations/2023_12_04_041703_add_dev_test_user.sql
var _2023_12_04_041703_add_dev_test_user string

func migrationList() []migration {
	migrations := []migration{}
	migrations = append(migrations, migration{"_2023_12_04_041700_base_schema_n_fn", _2023_12_04_041700_base_schema_n_fn})
	migrations = append(migrations, migration{"_2023_12_04_041701_base_tables", _2023_12_04_041701_base_tables})
	migrations = append(migrations, migration{"_2023_12_04_041702_base_roles", _2023_12_04_041702_base_roles})

	// only if developing/testing
	if os.Getenv(config.RUOK_ENVIRONMENT) != config.ProdRuokEnvironment {
		migrations = append(migrations, migration{"_2023_12_04_041703_add_dev_test_user", _2023_12_04_041703_add_dev_test_user})
	}

	return migrations
}

func migrate() {
	log.Println("Starting Migration Process")
	cfg := config.FromEnvs()
	s, close := storage.NewStorage(&cfg)
	defer close()

	ctx := context.Background()
	log.Printf("Starting Transactions\n")
	tx, err := s.GetClient().Begin(ctx)
	if err != nil {
		log.Fatalf("Could not start transaction to setup db: %q\n", err.Error())
	}
	defer tx.Rollback(ctx)

	migrations := migrationList()

	for i := 0; i < len(migrations); i++ {
		log.Printf("Starting migration number %d with name:%q\n\n", i, migrations[i].name)
		_, err := tx.Exec(ctx, migrations[i].sql)
		if err != nil {
			log.Fatalf("rolling back, couldn't exec migration: with name: %q. Error: %q\n\n", migrations[i].name, err.Error())
		}
		log.Printf("Success! Done with migration number %d with name : %q\n\n", i, migrations[i].name)

	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("couldn't commit transaction to generate db. error %q.\n", err.Error())
	}

}

var SetupDB = &cobra.Command{
	Use:   "setupdb",
	Short: "Runs all migrations needed to setup postgres to work with ruok",
	Long: `Runs all migrations needed to setup postgres to work with ruok. It will create:
  * the ruok schema
  * some utility funcions
  * all tables needed
  * couple roles 
`,
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}
