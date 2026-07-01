package cli

import (
	"fmt"

	"gocar/internal/project"
	"gocar/internal/util"
)

// TestCommand test 命令
type TestCommand struct{}

// Run 执行 test 命令
func (c *TestCommand) Run(args []string) error {
	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	testArgs, err := c.parseArgs(args)
	if err != nil {
		return err
	}
	if testArgs == nil {
		fmt.Print(c.Help())
		return nil
	}

	fmt.Printf("Testing '%s'...\n", appName)
	if err := util.RunCommand(projectRoot, "go", testArgs...); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("Tests passed")
	return nil
}

func (c *TestCommand) parseArgs(args []string) ([]string, error) {
	testArgs := []string{"test"}
	packages := []string{}
	passThrough := []string{}
	forwardingGoTestArgs := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if forwardingGoTestArgs {
			passThrough = append(passThrough, arg)
			continue
		}
		switch arg {
		case "help", "--help", "-h":
			return nil, nil
		case "--coverage":
			testArgs = append(testArgs, "-cover")
		case "--race":
			testArgs = append(testArgs, "-race")
		case "--bench":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--bench requires a value")
			}
			testArgs = append(testArgs, "-bench", args[i+1])
			i++
		case "--":
			passThrough = append(passThrough, args[i+1:]...)
			i = len(args)
		default:
			if len(arg) > 0 && arg[0] == '-' {
				passThrough = append(passThrough, arg)
				forwardingGoTestArgs = true
			} else {
				packages = append(packages, arg)
			}
		}
	}

	if len(packages) == 0 {
		packages = append(packages, "./...")
	}

	testArgs = append(testArgs, packages...)
	testArgs = append(testArgs, passThrough...)
	return testArgs, nil
}

// Help 返回帮助信息
func (c *TestCommand) Help() string {
	return `gocar test - Run project tests

USAGE:
    gocar test [OPTIONS] [packages...] [go test args...]

OPTIONS:
    --coverage          Run tests with coverage (-cover)
    --race              Enable the race detector
    --bench <pattern>   Run benchmarks matching pattern
    --help              Show this help message

EXAMPLES:
    gocar test                      Run all tests
    gocar test ./internal/...       Run tests for selected packages
    gocar test --coverage           Run all tests with coverage
    gocar test --bench .            Run all benchmarks
    gocar test -run TestConfig      Pass extra arguments to go test
    gocar test -- -run TestConfig   Explicitly separate gocar and go test args
`
}
