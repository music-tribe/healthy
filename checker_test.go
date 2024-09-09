package healthy

import (
	"context"
	errs "errors"
	"testing"

	"github.com/hellofresh/health-go/v5"
	"github.com/music-tribe/errors"
)

func TestNewChecker(t *testing.T) {
	type args struct {
		name      string
		checkFunc health.CheckFunc
	}
	tests := []struct {
		name            string
		args            args
		wantName        string
		wantCheckErrMsg string
	}{
		{
			name: "when the name is not added, it should default to health-check",
			args: args{
				name:      "",
				checkFunc: NewMockChecker(nil),
			},
			wantName:        defaultName,
			wantCheckErrMsg: "",
		},
		{
			name: "when the name is added, it should return this name",
			args: args{
				name:      "hello",
				checkFunc: NewMockChecker(nil),
			},
			wantName:        "hello",
			wantCheckErrMsg: "",
		},
		{
			name: "when the checkFunc is missing, it should default to the stub checkFunc",
			args: args{
				name:      "hello",
				checkFunc: nil,
			},
			wantName:        "hello",
			wantCheckErrMsg: stubCheckerErrorMessage,
		},
		{
			name: "when the checkFunc is added, it should return any errors from that checkFunc",
			args: args{
				name:      "hello",
				checkFunc: NewMockChecker(errors.NewCloudError(400, "oops")),
			},
			wantName:        "hello",
			wantCheckErrMsg: "oops",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewChecker(tt.args.name, tt.args.checkFunc)
			if got.Name != tt.wantName {
				t.Errorf("wanted name to be %s but got %s\n", tt.wantName, got.Name)
			}

			if err := got.CheckFunc(context.Background()); err != nil {
				msg := err.Error()
				ce := new(errors.CloudError)
				if errs.As(err, &ce) {
					msg = ce.Message
				}

				if msg != tt.wantCheckErrMsg {
					t.Errorf("wanted err msg to be %s but got %s\n", tt.wantName, msg)
				}

			}
		})
	}
}
