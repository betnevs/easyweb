package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

var _ IResponse = (*Context)(nil)

type IResponse interface {
	JSON(obj interface{}) IResponse

	JSONP(obj interface{}) IResponse

	XML(obj interface{}) IResponse

	HTML(tpl string, obj interface{}) IResponse

	Text(format string, values ...interface{}) IResponse

	Redirect(path string) IResponse

	SetHeader(key string, val string) IResponse

	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	SetStatus(code int) IResponse

	SetOkStatus() IResponse
}

func (ctx *Context) JSON(obj interface{}) IResponse {
	ctx.WriteMutex().Lock()
	defer ctx.WriteMutex().Unlock()

	if ctx.HasStopped() {
		return ctx
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	ctx.SetHeader("Content-Type", "application/json")
	ctx.response.Write(b)
	return ctx
}

func (ctx *Context) JSONP(obj interface{}) IResponse {
	callBack, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Content-type", "application/javascript")
	callBack = template.JSEscapeString(callBack)

	_, err := ctx.response.Write([]byte(callBack))
	if err != nil {
		return ctx
	}

	_, err = ctx.response.Write([]byte("("))
	if err != nil {
		return ctx
	}

	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}

	_, err = ctx.response.Write(ret)
	if err != nil {
		return ctx
	}

	_, err = ctx.response.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx
}

func (ctx *Context) XML(obj interface{}) IResponse {
	ctx.WriteMutex().Lock()
	defer ctx.WriteMutex().Unlock()

	if ctx.HasStopped() {
		return ctx
	}

	b, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	ctx.SetHeader("Content-Type", "application/xml")
	ctx.response.Write(b)
	return ctx
}

func (ctx *Context) HTML(tpl string, obj interface{}) IResponse {
	t, err := template.New("").ParseFiles(tpl)
	if err != nil {
		return ctx
	}
	ctx.SetHeader("Content-Type", "text/html")
	if err = t.Execute(ctx.response, obj); err != nil {
		return ctx
	}
	return ctx
}

func (ctx *Context) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.response.Write([]byte(out))
	return ctx
}

func (ctx *Context) Redirect(path string) IResponse {
	http.Redirect(ctx.response, ctx.request, path, http.StatusFound)
	return ctx
}

func (ctx *Context) SetHeader(key string, val string) IResponse {
	ctx.response.Header().Add(key, val)
	return ctx
}

func (ctx *Context) SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.response, &http.Cookie{
		Name:     key,
		Value:    val,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})

	return ctx
}

func (ctx *Context) SetStatus(code int) IResponse {
	ctx.response.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetOkStatus() IResponse {
	return ctx.SetStatus(http.StatusOK)
}
