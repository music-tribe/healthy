package healthy

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		serviceName string
		version     string
		checkers    []Checker
	}
	tests := []struct {
		name        string
		args        args
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{
			name: "when the name is missing we should get an error",
			args: args{
				serviceName: "",
				version:     "1.2.3",
				checkers:    []Checker{NewChecker("hello", &MockPinger{})},
			},
			wantErr:     true,
			wantVersion: "",
		},
		{
			name: "when the version is missing we should get an error",
			args: args{
				serviceName: "hello",
				version:     "",
				checkers:    []Checker{NewChecker("hello", &MockPinger{})},
			},
			wantErr:     true,
			wantVersion: "",
			wantName:    "",
		},
		{
			name: "when the checkers is passed as nil, it should be populated",
			args: args{
				serviceName: "hello",
				version:     "1.2.3",
				checkers:    []Checker{NewChecker("hello", &MockPinger{})},
			},
			wantErr:     false,
			wantName:    "hello",
			wantVersion: "1.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.serviceName, tt.args.version, tt.args.checkers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error to be %v but got %v", tt.wantErr, err)
			}

			if got.Name() != tt.wantName {
				t.Errorf("wanted name to be %s but got %s\n", tt.wantName, got.Name())
				return
			}

			if got.Version() != tt.wantVersion {
				t.Errorf("wanted version to be %s but got %s\n", tt.wantVersion, got.Version())
				return
			}

			if got.Checkers() == nil {
				t.Error("we never want checkers to be nil - always empty array")
				return
			}

		})
	}
}
