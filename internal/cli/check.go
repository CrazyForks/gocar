package cli

import (
	"fmt"

	"gocar/internal/project"
	"gocar/internal/util"
)

// CheckCommand check 命令
type CheckCommand struct{}

// Run 执行 check 命令
func (c *CheckCommand) Run(args []string) error {
	runTests := true
	race := false

	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h":
			fmt.Print(c.Help())
			return nil
		case "--no-test":
			runTests = false
		case "--race":
			race = true
		default:
			return fmt.Errorf("unknown option '%s' (run 'gocar check --help' for usage)", arg)
		}
	}

	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("Checking '%s'...\n", appName)

	steps := []struct {
		name string
		args []string
	}{
		{name: "fmt", args: []string{"fmt", "./..."}},
		{name: "vet", args: []string{"vet", "./..."}},
	}

	if runTests {
		testArgs := []string{"test", "./..."}
		if race {
			testArgs = []string{"test", "-race", "./..."}
		}
		steps = append(steps, struct {
			name string
			args []string
		}{name: "test", args: testArgs})
	}

	for _, step := range steps {
		fmt.Printf("Running go %s...\n", step.name)
		if err := util.RunCommand(projectRoot, "go", step.args...); err != nil {
			return fmt.Errorf("go %s failed: %w", step.name, err)
		}
	}

	fmt.Println("Check passed")
	return nil
}

// Help 返回帮助信息
func (c *CheckCommand) Help() string {
	return `gocar check - Run fmt, vet, and tests

USAGE:
    gocar check [OPTIONS]

OPTIONS:
    --no-test      Skip go test ./...
    --race         Run tests with the race detector
    --help         Show this help message

EXAMPLES:
    gocar check            Run go fmt, go vet, and go test
    gocar check --race     Run checks and race-enabled tests
    gocar check --no-test  Run only go fmt and go vet
`
}
