package gut

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FiberResponseWriter struct {
	StatusCode int
	Headers    http.Header
	Body       *bytes.Buffer
}

func NewFiberResponseWriter() *FiberResponseWriter {
	return &FiberResponseWriter{
		StatusCode: http.StatusOK,
		Headers:    make(http.Header),
		Body:       bytes.NewBuffer(nil),
	}
}

func (w *FiberResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *FiberResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *FiberResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func FiberRequestAdapter(c *fiber.Ctx, prefix string) (*http.Request, error) {
	url := string(c.Request().URI().Path())
	url = strings.TrimPrefix(url, prefix) + "?" + string(c.Request().URI().QueryString())
	method := string(c.Request().Header.Method())
	body := bytes.NewReader(c.Body())

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})

	return req, nil
}

func FiberResponseAdapter(w *FiberResponseWriter, c *fiber.Ctx) error {
	for key, values := range w.Headers {
		for _, value := range values {
			c.Response().Header.Add(key, value)
		}
	}

	c.Status(w.StatusCode)
	return c.Send(w.Body.Bytes())
}

func FiberAdapter(handler http.Handler, prefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := FiberRequestAdapter(c, prefix)
		if err != nil {
			return err
		}

		w := NewFiberResponseWriter()
		handler.ServeHTTP(w, req)

		return FiberResponseAdapter(w, c)
	}
}
