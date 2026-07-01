package cli

import (
	"reflect"
	"testing"
)

func TestTestCommandParseArgs(t *testing.T) {
	cmd := &TestCommand{}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "defaults to all packages",
			want: []string{"test", "./..."},
		},
		{
			name: "coverage package and passthrough",
			args: []string{"--coverage", "./internal/...", "--", "-run", "TestConfig"},
			want: []string{"test", "-cover", "./internal/...", "-run", "TestConfig"},
		},
		{
			name: "common go test args without separator",
			args: []string{"./internal/...", "-run", "TestConfig", "-count=1"},
			want: []string{"test", "./internal/...", "-run", "TestConfig", "-count=1"},
		},
		{
			name: "go test args before package",
			args: []string{"-run", "TestConfig", "./internal/..."},
			want: []string{"test", "./...", "-run", "TestConfig", "./internal/..."},
		},
		{
			name: "race and bench",
			args: []string{"--race", "--bench", "."},
			want: []string{"test", "-race", "-bench", ".", "./..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.parseArgs(tt.args)
			if err != nil {
				t.Fatalf("parseArgs() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseArgs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTestCommandParseArgsHelp(t *testing.T) {
	got, err := (&TestCommand{}).parseArgs([]string{"--help"})
	if err != nil {
		t.Fatalf("parseArgs() unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("parseArgs() = %#v, want nil for help", got)
	}
}

func TestTestCommandParseArgsBenchRequiresValue(t *testing.T) {
	if _, err := (&TestCommand{}).parseArgs([]string{"--bench"}); err == nil {
		t.Fatal("expected --bench without value to fail")
	}
}
