package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSinkConfig_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		sc              SinkConfig
		wantErrIs       error
		wantErrContains string
	}{
		{
			name: "missing-name",
			sc: SinkConfig{
				EventTypes: []Type{EveryType},
				Type:       FileSink,
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "missing sink name",
		},
		{
			name: "missing-EventType",
			sc: SinkConfig{
				Name: "sink-name",
				Type: FileSink,
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "missing event types",
		},
		{
			name: "invalid-EventType",
			sc: SinkConfig{
				Name:       "sink-name",
				EventTypes: []Type{"invalid"},
				Type:       FileSink,
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "not a valid event type",
		},
		{
			name: "missing-sink-type",
			sc: SinkConfig{
				Name:       "sink-name",
				EventTypes: []Type{EveryType},
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "not a valid sink type",
		},
		{
			name: "invalid-sink-type",
			sc: SinkConfig{
				Name:       "sink-name",
				EventTypes: []Type{EveryType},
				Type:       "invalid",
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "not a valid sink type",
		},
		{
			name: "missing-format",
			sc: SinkConfig{
				Name:       "sink-name",
				Type:       FileSink,
				EventTypes: []Type{EveryType},
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "not a valid sink format",
		},
		{
			name: "invalid-format",
			sc: SinkConfig{
				Name:       "sink-name",
				Format:     "invalid",
				Type:       FileSink,
				EventTypes: []Type{EveryType},
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "not a valid sink format",
		},
		{
			name: "file-sink-with-no-file-name",
			sc: SinkConfig{
				EventTypes: []Type{EveryType},
				Type:       FileSink,
				Format:     JSONSinkFormat,
				FileConfig: &FileSinkTypeConfig{},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: "missing file name",
		},
		{
			name: "type mismatch file type stderr config",
			sc: SinkConfig{
				EventTypes:   []Type{EveryType},
				Type:         FileSink,
				Format:       JSONSinkFormat,
				StderrConfig: &StderrSinkTypeConfig{},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: `missing "file" block`,
		},
		{
			name: "type mismatch stderr type file config",
			sc: SinkConfig{
				EventTypes: []Type{EveryType},
				Type:       StderrSink,
				Format:     JSONSinkFormat,
				FileConfig: &FileSinkTypeConfig{},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: `mismatch between sink type and sink configuration block`,
		},
		{
			name: "type mismatch both types file config",
			sc: SinkConfig{
				EventTypes:   []Type{EveryType},
				Type:         FileSink,
				Format:       JSONSinkFormat,
				StderrConfig: &StderrSinkTypeConfig{},
				FileConfig:   &FileSinkTypeConfig{},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: `too many sink type config blocks`,
		},
		{
			name: "type mismatch both types stderr config",
			sc: SinkConfig{
				EventTypes:   []Type{EveryType},
				Type:         StderrSink,
				Format:       JSONSinkFormat,
				StderrConfig: &StderrSinkTypeConfig{},
				FileConfig:   &FileSinkTypeConfig{},
			},
			wantErrIs:       ErrInvalidParameter,
			wantErrContains: `too many sink type config blocks`,
		},
		{
			name: "valid",
			sc: SinkConfig{
				Name:       "valid",
				EventTypes: []Type{EveryType},
				Type:       FileSink,
				FileConfig: &FileSinkTypeConfig{
					FileName: "tmp.file",
				},
				Format: JSONSinkFormat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			err := tt.sc.Validate()
			if tt.wantErrIs != nil {
				require.Error(err)
				assert.ErrorIs(err, tt.wantErrIs)
				if tt.wantErrContains != "" {
					assert.Contains(err.Error(), tt.wantErrContains)
				}
				return
			}
			assert.NoError(err)
		})
	}
}
