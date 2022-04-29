package core

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	handlers []ControllerHandler
	index    int

	// 标记 handler 处理是否超时
	hasStopped int32
	writeMutex *sync.Mutex

	params map[string]string
}

func NewContext(request *http.Request, response http.ResponseWriter) *Context {
	return &Context{
		request:    request,
		response:   response,
		writeMutex: &sync.Mutex{},
		index:      -1,
	}
}

func (ctx *Context) WriteMutex() *sync.Mutex {
	return ctx.writeMutex
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.response
}

func (ctx *Context) SetHasStopped() {
	atomic.AddInt32(&ctx.hasStopped, 1)
}

func (ctx *Context) HasStopped() bool {
	return atomic.LoadInt32(&ctx.hasStopped) == 1
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) BindJSON(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		//ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx request empty")
	}
	return nil
}

func (ctx *Context) ExecTimeout() {
	ctx.WriteMutex().Lock()
	defer ctx.writeMutex.Unlock()
	ctx.SetHasStopped()
	ctx.response.WriteHeader(http.StatusInternalServerError)
	ctx.response.Write([]byte("time out"))
}

func (ctx *Context) ExecPanic() {
	ctx.WriteMutex().Lock()
	defer ctx.WriteMutex().Unlock()
	ctx.SetHasStopped()
	ctx.response.WriteHeader(http.StatusInternalServerError)
	ctx.response.Write([]byte("panic"))
}

func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}
