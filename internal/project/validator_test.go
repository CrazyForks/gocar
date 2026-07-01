package project

import "testing"

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "myapp"},
		{name: "my-app_1"},
		{name: "", wantErr: true},
		{name: "-myapp", wantErr: true},
		{name: ".myapp", wantErr: true},
		{name: "1myapp", wantErr: true},
		{name: "my app", wantErr: true},
		{name: "internal", wantErr: true},
		{name: "TEST", wantErr: true},
	}

	for _, tt := range tests {
		err := ValidateProjectName(tt.name)
		if tt.wantErr && err == nil {
			t.Fatalf("ValidateProjectName(%q) expected error", tt.name)
		}
		if !tt.wantErr && err != nil {
			t.Fatalf("ValidateProjectName(%q) unexpected error: %v", tt.name, err)
		}
	}
}
