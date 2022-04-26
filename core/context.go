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
	handler  ControllerHandler

	// 标记 handler 处理是否超时
	hasStopped int32
	writeMutex *sync.Mutex
}

func NewContext(request *http.Request, response http.ResponseWriter) *Context {
	return &Context{
		request:    request,
		response:   response,
		writeMutex: &sync.Mutex{},
	}
}

func (c *Context) WriteMutex() *sync.Mutex {
	return c.writeMutex
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.response
}

func (c *Context) SetHasStopped() {
	atomic.AddInt32(&c.hasStopped, 1)
}

func (c *Context) HasStopped() int32 {
	return atomic.LoadInt32(&c.hasStopped)
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) QueryInt(key string, def int) int {
	params := c.QueryAll()
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

func (c *Context) QueryString(key string, def string) string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		l := len(vals)
		if l > 0 {
			return vals[l-1]
		}
	}
	return def
}

func (c *Context) QueryArray(key string, def []string) []string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		c.request.ParseForm()
		return c.request.PostForm
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) int {
	params := c.FormAll()
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

func (c *Context) FormString(key string, def string) string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		l := len(params)
		if l > 0 {
			return vals[l-1]
		}
	}
	return def
}

func (c *Context) FormArray(key string, def []string) []string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (c *Context) BindJSON(obj interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}
		//c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx request empty")
	}
	return nil
}

func (c *Context) JSON(status int, obj interface{}) error {
	c.WriteMutex().Lock()
	defer c.WriteMutex().Unlock()
	if c.HasStopped() != 0 {
		return nil
	}
	byt, err := json.Marshal(obj)
	if err != nil {
		c.response.WriteHeader(http.StatusInternalServerError)
		return err
	}
	// Header 必须在 WriteHeader 方法之前设置
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status)
	c.response.Write(byt)
	return nil
}

func (c *Context) ExecTimeout() {
	c.WriteMutex().Lock()
	defer c.writeMutex.Unlock()
	c.SetHasStopped()
	c.response.WriteHeader(http.StatusInternalServerError)
	c.response.Write([]byte("time out"))
}

func (c *Context) ExecPanic() {
	c.WriteMutex().Lock()
	defer c.WriteMutex().Unlock()
	c.SetHasStopped()
	c.response.WriteHeader(http.StatusInternalServerError)
	c.response.Write([]byte("panic"))
}

func (c *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (c *Context) Text(status int, obj string) error {
	return nil
}
