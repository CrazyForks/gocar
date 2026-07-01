package cli

import "testing"

func TestNewAppRegistersCoreCommands(t *testing.T) {
	app := NewApp()

	for _, name := range []string{"new", "init", "build", "run", "clean", "add", "update", "tidy", "test", "check", "commands", "doctor"} {
		if app.commands[name] == nil {
			t.Fatalf("command %q was not registered", name)
		}
	}
}

func TestProtectedCommands(t *testing.T) {
	if !isProtectedCommand("new") || !isProtectedCommand("init") {
		t.Fatal("new and init should be protected")
	}
	if isProtectedCommand("test") || isProtectedCommand("check") {
		t.Fatal("test and check should be overrideable project commands")
	}
}
