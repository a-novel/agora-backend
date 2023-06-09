package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/a-novel/agora-backend/framework/bunframework"
	"github.com/a-novel/agora-backend/framework/bunframework/pgconfig"
	"github.com/a-novel/agora-backend/migrations"
	"github.com/gookit/color"
	"io/fs"
	"os"
	"strings"
	"time"
)

var (
	dsn string

	steps = []string{
		"🔌 Acquiring connection",
		"🏃 Applying migrations",
	}
)

func init() {
	flag.StringVar(&dsn, "d", os.Getenv("POSTGRES_URL"), "database to rollback")
}

func requireConfirmation(txt string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(txt + " [Y/n] ")
	scanner.Scan()
	if strings.ToLower(scanner.Text()) != "y" {
		os.Exit(0)
	}
}

func printSteps(reprint bool, current int) {
	if reprint {
		for i := 0; i < len(steps); i++ {
			fmt.Printf("\r\033[1A\033[0K")
		}
	}

	for i, step := range steps {
		c := uint8(245)
		if i == current {
			c = 220
		} else if i < current {
			c = 40
		}

		color.C256(c).Printf("- %s\n", step)
	}
}

func quit(err string) {
	fmt.Println("")
	fmt.Println("")
	color.C256(9).Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	color.C256(45).Println("Rolling back local migrations.")
	fmt.Println("")

	fmt.Printf("The migrations will be applied to this instance: %s\n", color.C256(13).Sprint(dsn))
	requireConfirmation("Confirm ?")

	printSteps(false, 0)

	postgresClient, sqlClient, err := bunframework.NewClient(context.Background(), bunframework.Config{
		Driver: pgconfig.Driver{
			DSN:         dsn,
			DialTimeout: 120 * time.Second,
		},
		DiscardUnknownColumns: true,
	})
	if err != nil {
		quit(fmt.Sprintf("💥 failed to acquire connection to '%s': %s", dsn, err.Error()))
		return
	}

	defer postgresClient.Close()
	defer sqlClient.Close()

	steps[0] = "🔌 Connection acquired"
	printSteps(true, 1)

	migrationsConfig := &bunframework.MigrateConfig{
		Files: []fs.FS{migrations.Migrations},
	}

	if err := migrationsConfig.Rollback(context.Background(), postgresClient); err != nil {
		quit(fmt.Sprintf("💥 failed to rollback last migration group: %s", err.Error()))
		return
	}

	steps[0] = "🔌 Last migration group rolled back successfully"
	printSteps(true, 2)

	fmt.Println("")
	color.C256(45).Println("🚀 Migrations rolled back successfully!")
	color.C256(255).Println("The following migrations where rolled back:")
	for _, migration := range migrationsConfig.Report().Migrations {
		if migration.IsApplied() {
			color.C256(255).Printf("\t%s\n", migration)
		} else {
			color.C256(245).Printf("\t %s\n", migration)
		}
	}
}
