package healthy

import (
	errs "errors"
	"testing"

	"github.com/music-tribe/errors"
)

func TestNewChecker(t *testing.T) {
	type args struct {
		name   string
		pinger Pinger
	}
	tests := []struct {
		name           string
		args           args
		wantName       string
		wantPingErrMsg string
	}{
		{
			name: "when the name is not added, it should default to health-check",
			args: args{
				name:   "",
				pinger: &MockPinger{},
			},
			wantName:       defaultName,
			wantPingErrMsg: "",
		},
		{
			name: "when the name is added, it should return this name",
			args: args{
				name:   "hello",
				pinger: &MockPinger{},
			},
			wantName:       "hello",
			wantPingErrMsg: "",
		},
		{
			name: "when the pinger is missing, it should default to the stub pinger",
			args: args{
				name:   "hello",
				pinger: nil,
			},
			wantName:       "hello",
			wantPingErrMsg: stubPingerErrorMessage,
		},
		{
			name: "when the pinger is added, it should return any errors from that pinger",
			args: args{
				name:   "hello",
				pinger: &MockPinger{errors.NewCloudError(400, "oops")},
			},
			wantName:       "hello",
			wantPingErrMsg: "oops",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewChecker(tt.args.name, tt.args.pinger)
			if got.Name() != tt.wantName {
				t.Errorf("wanted name to be %s but got %s\n", tt.wantName, got.Name())
			}

			if err := got.Check(); err != nil {
				msg := err.Error()
				ce := new(errors.CloudError)
				if errs.As(err, &ce) {
					msg = ce.Message
				}

				if msg != tt.wantPingErrMsg {
					t.Errorf("wanted err msg to be %s but got %s\n", tt.wantName, msg)
				}

			}
		})
	}
}
