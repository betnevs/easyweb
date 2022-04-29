package core

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net"
	"strings"

	"github.com/spf13/cast"
)

var _ IRequest = (*Context)(nil)

type IRequest interface {
	// QueryXXX get query params, such as bar.com?name=nevs&age=18
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// ParamXXX get router params, such as /user/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	// FormXXX get form params
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)
	Form(key string) interface{}

	BindJson(obj interface{}) error

	BindXml(obj interface{}) error

	GetRawData() ([]byte, error)

	Uri() string
	Method() string
	Host() string
	ClintIP() string

	Headers() map[string][]string
	Header(key string) string

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToInt(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToInt64(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat32(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat64(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToBool(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return v[0], true
		}
	}
	return def, false
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		return v, true
	}
	return def, false
}

func (ctx *Context) Query(key string) interface{} {
	params := ctx.QueryAll()
	if v, ok := params[key]; ok {
		return v[0]
	}
	return nil
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToInt(v), true
	}
	return def, false
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToInt64(v), true
	}
	return def, false
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToFloat32(v), true
	}
	return def, false
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToFloat64(v), true
	}
	return def, false
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToBool(v), true
	}
	return def, false
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	if v := ctx.Param(key); v != nil {
		return cast.ToString(v), true
	}
	return def, false
}

func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if v, ok := ctx.params[key]; ok {
			return v
		}
	}
	return nil
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		ctx.request.ParseForm()
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToInt(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToInt64(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat32(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat64(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return cast.ToBool(v[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormString(key string, def string) (string, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return v[0], true
		}
	}
	return def, false
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		return v, true
	}
	return def, false
}

func (ctx *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.request.FormFile(key)
}

func (ctx *Context) Form(key string) interface{} {
	params := ctx.FormAll()
	if v, ok := params[key]; ok {
		if len(v) > 0 {
			return v[0]
		}
	}
	return nil
}

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request == nil {
		return errors.New("ctx.request empty")
	}

	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}

	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) BindXml(obj interface{}) error {
	if ctx.request == nil {
		return errors.New("ctx.request empty")
	}

	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}

	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = xml.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) GetRawData() ([]byte, error) {
	if ctx.request == nil {
		return nil, errors.New("ctx.request empty")
	}

	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return nil, err
	}

	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	return ctx.request.Host
}

func (ctx *Context) ClintIP() string {
	if ctx.request == nil {
		return ""
	}
	r := ctx.request
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	IP := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if IP != "" {
		return IP
	}

	IP = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if IP != "" {
		return IP
	}
	IP, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	return IP
}

func (ctx *Context) Headers() map[string][]string {
	return ctx.request.Header
}

func (ctx *Context) Header(key string) string {
	return ctx.request.Header.Get(key)
}

func (ctx *Context) Cookies() map[string]string {
	ret := make(map[string]string)
	cookies := ctx.request.Cookies()
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookies := ctx.Cookies()
	if v, ok := cookies[key]; ok {
		return v, true
	}

	return "", false
}
