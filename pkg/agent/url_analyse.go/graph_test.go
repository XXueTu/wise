package url_analyse

import (
	"context"
	"testing"
	"time"
)

func TestBuildAnalysisGraph(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test-compile",
			args: args{
				ctx: context.Background(),
				url: "https://mp.weixin.qq.com/s/1234567890",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BuildAnalysisGraph(nil)
			if err != nil {
				t.Errorf("BuildAnalysisGraph() error = %v", err)
				return
			}
			RunUrlAnalyseAgent("123", tt.args.url)
			time.Sleep(10 * time.Second)
		})
	}
}
