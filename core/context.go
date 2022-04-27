package core

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (ctx *Context) HasStopped() int32 {
	return atomic.LoadInt32(&ctx.hasStopped)
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

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		l := len(vals)
		if l > 0 {
			intval, err := strconv.Atoi(vals[l-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		l := len(vals)
		if l > 0 {
			return vals[l-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		ctx.request.ParseForm()
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		l := len(params)
		if l > 0 {
			intval, err := strconv.Atoi(vals[l-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		l := len(params)
		if l > 0 {
			return vals[l-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
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

func (ctx *Context) JSON(status int, obj interface{}) error {
	ctx.WriteMutex().Lock()
	defer ctx.WriteMutex().Unlock()
	if ctx.HasStopped() != 0 {
		return nil
	}
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.response.WriteHeader(http.StatusInternalServerError)
		return err
	}
	// Header 必须在 WriteHeader 方法之前设置
	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.response.WriteHeader(status)
	ctx.response.Write(byt)
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

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
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
