package url_analyse

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
)

func EndNodeHandler(ctx context.Context, param map[string]any) (map[string]any, error) {
	logx.Infof("end node handler: %+v", param)
	return param, nil
}

var url_analyse_runnable compose.Runnable[map[string]any, any]
var traceHandler *callbacks.HandlerBuilder
var svcCtx *svc.ServiceContext

type UrlAnalyseStep struct {
	Step int64
}

var UrlAnalyseSteps = map[string]UrlAnalyseStep{
	"start":  {Step: 0},
	"check":  {Step: 1},
	"read":   {Step: 2},
	"split":  {Step: 3},
	"mark":   {Step: 4},
	"vector": {Step: 5},
	"index":  {Step: 6},
}

const (
	nodeOfCheck  = "check"
	nodeOfRead   = "read"
	nodeOfSplit  = "split"
	nodeOfMark   = "mark"
	nodeOfVector = "vector"
	nodeOfIndex  = "index"
)

func BuildAnalysisGraph(svct *svc.ServiceContext) error {
	logx.Infof("build analysis graph")
	svcCtx = svct
	ctx := context.Background()
	wf := compose.NewWorkflow[map[string]any, any]()
	wf.AddLambdaNode(nodeOfCheck,
		compose.InvokableLambda(CheckNodeHandler), compose.WithNodeName(nodeOfCheck)).AddInput(compose.START)

	wf.AddLambdaNode(nodeOfRead,
		compose.InvokableLambda(ReadNodeHandler), compose.WithNodeName(nodeOfRead)).
		AddInput(nodeOfCheck)

	wf.AddLambdaNode(nodeOfSplit,
		compose.InvokableLambda(SplitNodeHandler), compose.WithNodeName(nodeOfSplit)).
		AddInput(nodeOfRead)

	wf.AddLambdaNode(nodeOfMark,
		compose.InvokableLambda(MarkNodeHandler), compose.WithNodeName(nodeOfMark)).
		AddInput(nodeOfSplit)

	wf.End().AddInput(nodeOfMark)
	runnable, err := wf.Compile(ctx)
	if err != nil {
		return err
	}
	traceHandler = callbacks.NewHandlerBuilder()
	traceHandler.OnStartFn(
		func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			logx.Infof("onStart, runInfo: %v, input: %v", info, input)
			value := ctx.Value(urlAnalyseContextKey).(map[string]string)
			tid := value["tid"]
			pid := model.GenUid()
			name := tenl(info)
			if name == "" {
				name = pid
			}
			ctx = context.WithValue(ctx, TraceId, pid)
			// 创建task_plans
			jsonInput, _ := json.Marshal(input)
			_ = svcCtx.TaskPlansModel.Create(ctx, &model.TaskPlans{
				Tid:       tid,
				Pid:       pid,
				BeforePid: "",
				Next:      "",
				Types:     info.Type,
				Name:      name,
				Index:     0,
				Status:    model.TaskPlanStatusInit,
				Params:    string(jsonInput),
				Result:    "{}",
				Duration:  0,
			})
			// 更新 task 当前步骤
			_ = svcCtx.TasksModel.UpdateStateAndStep(ctx, tid, name, UrlAnalyseSteps[name].Step, "{}")
			return ctx
		}).OnEndFn(
		func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			logx.Infof("onEnd, runInfo: %v, out: %v", info, output)
			pid := ctx.Value(TraceId).(string)
			taskPlan, _ := svcCtx.TaskPlansModel.GetByPid(ctx, pid)
			jsonOutput, _ := json.Marshal(output)
			taskPlan.Result = string(jsonOutput)
			taskPlan.Status = model.TaskPlanStatusSuccess
			_ = svcCtx.TaskPlansModel.Update(ctx, taskPlan)
			return ctx
		}).OnErrorFn(
		func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
			logx.Errorf("onError, runInfo: %v, err: %v", info, err)
			pid := ctx.Value(TraceId).(string)
			taskPlan, _ := svcCtx.TaskPlansModel.GetByPid(ctx, pid)
			taskPlan.Status = model.TaskPlanStatusFailed
			if err != nil {
				taskPlan.Error = err.Error()
			}
			_ = svcCtx.TaskPlansModel.Update(ctx, taskPlan)
			return ctx
		})
	url_analyse_runnable = runnable
	logx.Infof("build analysis graph success")
	return nil
}

func tenl(info *callbacks.RunInfo) string {
	if info.Name != "" {
		return info.Name
	}
	return ""
}

// contextKey 是一个自定义类型，用作 context 的 key
type contextKey string

const (
	urlAnalyseContextKey contextKey = "url_analyse_context"
	TraceId              contextKey = "trace_id"
)

func RunUrlAnalyseAgent(tid string, url string) error {
	start := map[string]any{
		"tid": tid,
		"url": url,
	}
	value := map[string]string{
		"tid": tid,
	}
	ctx := context.WithValue(context.Background(), urlAnalyseContextKey, value)
	_, err := url_analyse_runnable.Invoke(ctx, start, compose.WithCallbacks(traceHandler.Build()))
	if err != nil {
		logx.Errorf("run url analyse agent error: %v", err)
		return err
	}
	return nil
}
