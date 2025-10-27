package gut

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"
)

func ConvertRequest(ctx *fasthttp.RequestCtx) *http.Request {
	defer func() {
		if err := recover(); err != nil {
			// TODO: use sentry logger
			Debug(err)
		}
	}()

	r := new(http.Request)

	r.Method = string(ctx.Method())
	uri := ctx.URI()
	// * ignore error
	r.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path()))

	// * headers
	r.Header = make(http.Header)
	r.Header.Add("Host", string(ctx.Host()))
	headers := ctx.Request.Header.All()
	for key, values := range headers {
		for _, value := range values {
			r.Header.Add(string(key), string(value))
		}
	}
	r.Host = string(ctx.Host())

	// * cookies
	cookies := ctx.Request.Header.Cookies()
	for key, value := range cookies {
		r.AddCookie(&http.Cookie{Name: string(key), Value: string(value)})
	}

	// * env
	r.RemoteAddr = ctx.RemoteAddr().String()

	// * query string
	r.URL.RawQuery = string(ctx.URI().QueryString())

	// * body
	r.Body = io.NopCloser(bytes.NewReader(ctx.Request.Body()))

	return r
}
