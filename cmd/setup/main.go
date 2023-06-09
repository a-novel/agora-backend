package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	envrc = `export GOPRIVATE="github.com/a-novel/*"
export POSTGRES_URL="postgres://postgres@localhost:5432/agora?sslmode=disable"
export POSTGRES_URL_TEST="postgres://test@localhost:5432/agora_test?sslmode=disable"`

	steps = []string{
		"📥 Installing packages",
		"📄 Creating '.envrc' file",
		"📁 Creating '.secrets' directory",
		"🔌 Sourcing environment",
	}
)

func info(txt string) {
	color.C256(245).Printf("\n%s\n\n", txt)
}

func quit(err string) {
	fmt.Println("")
	fmt.Println("")
	color.C256(9).Println(err)
	os.Exit(1)
}

func cmd(cmd string, args ...string) {
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		quit(fmt.Sprintf(
			"💥 failed to execute command '%s %s': %s",
			cmd, strings.Join(args, " "), err.Error(),
		))
		return
	}
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

func writeENVRC(wd string) (archived bool) {
	data, err := os.ReadFile(path.Join(wd, ".envrc"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		quit(fmt.Sprintf("💥 failed to check for '.envrc' file: %s", err.Error()))
	}
	if data != nil {
		err = os.WriteFile(path.Join(wd, ".envrc.old"), data, 0644)
		if err != nil {
			quit(fmt.Sprintf("💥 failed to write '.envrc' file content to '.envrc.old': %s", err.Error()))
			return
		}
		archived = true
	}

	// Try to fetch secret keys from environment variables, since they can't be pushed into the script.
	if sendgridKey := os.Getenv("SENDGRID_API_KEY"); sendgridKey != "" {
		envrc += fmt.Sprintf("\nexport SENDGRID_API_KEY=%q", sendgridKey)
	}

	err = os.WriteFile(path.Join(wd, ".envrc"), []byte(envrc), 0644)
	if err != nil {
		quit(fmt.Sprintf("💥 failed to write '.envrc' file: %s", err.Error()))
		return
	}

	return
}

func main() {
	color.C256(45).Println("Running setup for 'Agora des Écrivains'.")
	color.C256(245).Println(
		"This script will prepare your local repository. You may rerun it any time if you encounter an issue.",
	)
	fmt.Println("")

	wd := getWD()
	fmt.Printf("We detected the following working directory: %s\n", color.C256(13).Sprint(wd))
	requireConfirmation("Do you confirm this path?")
	fmt.Println("")

	printSteps(false, 0)

	cmd("go", "install", "gotest.tools/gotestsum@latest")
	steps[0] = "📥 Packages installed"
	printSteps(true, 1)

	archived := writeENVRC(wd)
	if archived {
		steps[1] = "📄 Overwritten content of previous '.envrc' file"
	} else {
		steps[1] = "📄 Created '.envrc' file in the root directory"
	}

	printSteps(true, 2)
	if err := os.MkdirAll(path.Join(wd, ".secrets"), os.ModePerm); err != nil {
		quit(fmt.Sprintf("💥 failed to create '.secrets' directory: %s", err.Error()))
	}
	steps[2] = "📁 '.secrets' directory created"

	printSteps(true, 3)
	cmd("direnv", "allow", wd)
	steps[3] = "🔌 Environment sourced"

	printSteps(true, 4)
	fmt.Println("")
	color.C256(45).Println("🚀 Setup completed successfully!")
	color.C256(255).Println("You can now run and test the application locally.")

	if archived {
		info(
			"ⓘ The previous '.envrc' file has been renamed to '.envrc.old'. You can restore it if you wish. " +
				"Note this file is ignored by github and won't be pushed to the repository. You can delete it " +
				"if you don't need it anymore.",
		)
	}
}

func getWD() string {
	ex, err := os.Getwd()
	if err != nil {
		quit(fmt.Sprintf("💥 failed to retrieve working directory: %s", err.Error()))
		return ""
	}

	return ex
}
