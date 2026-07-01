package cli

import (
	"fmt"
	"strings"
)

// CommandInfo describes a built-in command in one place.
type CommandInfo struct {
	Name        string
	Usage       string
	Description string
	Example     string
}

var builtInCommands = []CommandInfo{
	{Name: "new", Usage: "new <name>", Description: "Create a standard Go application", Example: "gocar new myapp"},
	{Name: "init", Usage: "init", Description: "Initialize .gocar.toml in current project", Example: "gocar init"},
	{Name: "build", Usage: "build [OPTIONS]", Description: "Build the project", Example: "gocar build --release"},
	{Name: "run", Usage: "run [args...]", Description: "Run the project", Example: "gocar run"},
	{Name: "clean", Usage: "clean", Description: "Clean build artifacts", Example: "gocar clean"},
	{Name: "fmt", Usage: "fmt [packages...]", Description: "Format Go code", Example: "gocar fmt"},
	{Name: "vet", Usage: "vet [packages...]", Description: "Run go vet", Example: "gocar vet"},
	{Name: "test", Usage: "test [OPTIONS] [packages...]", Description: "Run tests", Example: "gocar test --coverage"},
	{Name: "check", Usage: "check [OPTIONS]", Description: "Run vet and tests", Example: "gocar check"},
	{Name: "add", Usage: "add <package>...", Description: "Add dependencies to go.mod", Example: "gocar add github.com/gin-gonic/gin"},
	{Name: "update", Usage: "update [package]...", Description: "Update dependencies", Example: "gocar update"},
	{Name: "tidy", Usage: "tidy", Description: "Tidy up go.mod and go.sum", Example: "gocar tidy"},
	{Name: "commands", Usage: "commands", Description: "List built-in and custom commands", Example: "gocar commands"},
	{Name: "doctor", Usage: "doctor", Description: "Check project and toolchain setup", Example: "gocar doctor"},
	{Name: "help", Usage: "help", Description: "Print this help message", Example: "gocar help"},
	{Name: "version", Usage: "version", Description: "Print version info", Example: "gocar version"},
}

func builtInCommandNames() []string {
	names := make([]string, 0, len(builtInCommands))
	for _, info := range builtInCommands {
		names = append(names, info.Name)
	}
	return names
}

func builtInCommandInfo(name string) (CommandInfo, bool) {
	for _, info := range builtInCommands {
		if info.Name == name {
			return info, true
		}
	}
	return CommandInfo{}, false
}

func formatBuiltInCommands() string {
	var b strings.Builder
	for _, info := range builtInCommands {
		fmt.Fprintf(&b, "    %-38s %s\n", info.Usage, info.Description)
	}
	return b.String()
}

func formatExamples() string {
	var b strings.Builder
	seen := map[string]bool{}
	for _, info := range builtInCommands {
		if info.Example == "" || seen[info.Example] {
			continue
		}
		seen[info.Example] = true
		fmt.Fprintf(&b, "    %-38s %s\n", info.Example, info.Description)
	}
	return b.String()
}
