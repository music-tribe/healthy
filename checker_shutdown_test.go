package healthy

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestShutdownChecker_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := NewMockLogger(ctrl)
	type fields struct {
		isShuttingDown bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should return nil when not shutting down",
			fields: fields{
				isShuttingDown: false,
			},
			wantErr: false,
		},
		{
			name: "should return error when shutting down",
			fields: fields{
				isShuttingDown: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component := "test ShutdownChecker"
			shutdownLog := "test service shutting down"
			p, c := NewShutdownChecker(tt.name, logger, component, shutdownLog)
			if tt.fields.isShuttingDown {
				logger.EXPECT().Info(component, "SetShutdown", shutdownLog)
				p.SetShutdown()
			}
			if err := c.CheckFunc(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ShutdownChecker.CheckFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
