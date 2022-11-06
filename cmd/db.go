package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chux0519/yeti/cmd/migrations"
	"github.com/chux0519/yeti/pkg/config"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var dbFlag = config.DBFlag{}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "yeti db helper",
	Run: func(cmd *cobra.Command, args []string) {
		println("hi db")
	},
}

func getDB() *bun.DB {
	sqldb, err := sql.Open(sqliteshim.ShimName, dbFlag.URI)
	if err != nil {
		panic(err)
	}
	return bun.NewDB(sqldb, sqlitedialect.New())
}

var dbInitCmd = &cobra.Command{
	Use:   "init",
	Short: "create migration tables",
	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		migrator := migrate.NewMigrator(db, migrations.Migrations)
		ctx := context.Background()
		err := migrator.Init(ctx)

		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		migrator := migrate.NewMigrator(db, migrations.Migrations)
		ctx := context.Background()

		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}

		if group.ID == 0 {
			fmt.Printf("there are no new migrations to run\n")
		}

		fmt.Printf("migrated to %s\n", group)
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback the last migration group",
	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		migrator := migrate.NewMigrator(db, migrations.Migrations)
		ctx := context.Background()

		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}

		if group.ID == 0 {
			fmt.Printf("there are no groups to roll back\n")
		}

		fmt.Printf("rolled back %s\n", group)
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create up and down SQL migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		migrator := migrate.NewMigrator(db, migrations.Migrations)
		ctx := context.Background()

		name := strings.Join(args, "_")
		files, err := migrator.CreateSQLMigrations(ctx, name)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, mf := range files {
			fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
		}
	},
}

func init() {
	iflags := dbInitCmd.Flags()
	iflags.StringVar(&dbFlag.URI, "uri", "", "sqlite3 uri")

	mflags := dbMigrateCmd.Flags()
	mflags.StringVar(&dbFlag.URI, "uri", "", "sqlite3 uri")

	rblags := dbRollbackCmd.Flags()
	rblags.StringVar(&dbFlag.URI, "uri", "", "sqlite3 uri")

	clags := dbCreateCmd.Flags()
	clags.StringVar(&dbFlag.URI, "uri", "", "sqlite3 uri")

	dbCmd.AddCommand(dbInitCmd)
	dbCmd.AddCommand(dbMigrateCmd)
	dbCmd.AddCommand(dbRollbackCmd)
	dbCmd.AddCommand(dbCreateCmd)
	rootCmd.AddCommand(dbCmd)
}
