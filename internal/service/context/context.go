package context

import (
	"context"

	"github.com/ctwj/aavirus_helper/internal/pkg/log"
	"github.com/ctwj/aavirus_helper/internal/pkg/util"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var global *context.Context

func Set(ctx *context.Context) {
	global = ctx
}

func Get() *context.Context {
	return global
}

type Context struct {
	context.Context
}

func New() *Context {
	return &Context{log.NewTraceIDContext(context.Background(), util.NewTraceID())}
}

// func (ctx *Context) Success(data interface{}) string {
// s, _ := util.JsonEncode(OutputFormat(codes.OK.Code, codes.OK.Message, data))
// return s
// }

// func (ctx *Context) Error(err error, additional ...interface{}) string {
// code, message := codes.DecodeErr(err)
// m := OutputFormat(code, message)
// if config.RunReleaseMode != config.App.Mode {
// 	m["detail"] = additional
// }
// s, _ := util.JsonEncode(m)
// return s
// }

func (ctx *Context) EventsEmit(eventName string, optionalData ...interface{}) {
	EventsEmit(ctx, eventName, optionalData...)
}

func OutputFormat(code int, msg string, data ...interface{}) map[string]interface{} {
	m := map[string]interface{}{
		"errno":  code,
		"errmsg": msg,
	}
	if len(data) > 0 {
		m["data"] = data[0]
	}
	return m
}

func EventsEmit(ctx context.Context, eventName string, optionalData ...interface{}) {
	runtime.EventsEmit(ctx, eventName, optionalData...)
}
