package url_analyse

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

func Test_llmMark(t *testing.T) {
	type args struct {
		ctx      context.Context
		segments []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "test-item-1",
			args: args{
				ctx:      context.Background(),
				segments: []string{"逻辑处理器，对 G 来说，P 相当于 CPU 核，G 只有绑定到 P 才能被调度。对 M 来说，P 提供了相关的执行环境(Context)"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := llmMark(tt.args.ctx, tt.args.segments)
			logx.Infof("got: %v, err: %v", got, err)
		})
	}
}
