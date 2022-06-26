package middlewares

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"time"

	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/fatih/color"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/compress/zlib"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

type (
	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
		Priority() int
	}

	Factory struct {
		Debug       bool
		Middlewares []Middleware
		Logger      *loggers.Logger
	}

	responseWrapper struct {
		http.ResponseWriter
		statusCode int
	}
)

func (w *responseWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWrapper) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}

	return w.statusCode
}

func (m *Factory) Register(middlewares []Middleware) {
	for _, v := range middlewares {
		m.Add(v)
	}
}

func (m *Factory) Add(middlware Middleware) {
	m.Middlewares = append(m.Middlewares, middlware)
}

func (m *Factory) Sort() {
	sort.Slice(m.Middlewares, func(i, j int) bool {
		return m.Middlewares[i].Priority() > m.Middlewares[j].Priority()
	})
}

func (m *Factory) Attach(handler http.Handler) http.Handler {
	ctx := context.WithValue(context.Background(), "scope", "middleware")
	internal := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()
		if !m.Debug {
			for _, middleware := range m.Middlewares {
				if stop := middleware.Attach(request, response); stop {
					return
				}
			}

			handler.ServeHTTP(response, request)

			elapsed := time.Since(start)

			var execution bytes.Buffer
			execution.WriteString("Execution time: ")
			execution.WriteString(elapsed.String())

			m.Logger.Info(ctx, execution.String())

			return
		}

		wrapper := responseWrapper{ResponseWriter: response}
		for _, middleware := range m.Middlewares {
			if stop := middleware.Attach(request, response); stop {
				var stopper bytes.Buffer
				stopper.WriteString("Middleware stopped by: ")
				stopper.WriteString(reflect.TypeOf(middleware).Name())

				m.Logger.Debug(ctx, stopper.String())

				return
			}
		}

		handler.ServeHTTP(&wrapper, request)

		elapsed := time.Since(start)

		var statusCode string
		uri, _ := url.QueryUnescape(request.RequestURI)
		mGet := color.New(color.BgHiGreen, color.FgBlack)
		mPost := color.New(color.BgYellow, color.FgBlack)
		mPut := color.New(color.BgCyan, color.FgBlack)
		mDelete := color.New(color.BgRed, color.FgBlack)

		switch request.Method {
		case http.MethodPost:
			mPost.Print("[POST]")
		case http.MethodPatch:
			mPost.Print("[PATCH]")
		case http.MethodPut:
			mPut.Print("[PUT]")
		case http.MethodDelete:
			mDelete.Print("[DELETE]")
		default:
			mGet.Print("[GET]")
		}

		switch {
		case wrapper.StatusCode() < 300:
			statusCode = color.New(color.FgGreen, color.Bold).Sprintf("%d", wrapper.StatusCode())
		case wrapper.StatusCode() < 400:
			statusCode = color.New(color.FgYellow, color.Bold).Sprintf("%d", wrapper.StatusCode())
		default:
			statusCode = color.New(color.FgRed, color.Bold).Sprintf("%d", wrapper.StatusCode())
		}

		var elapsedString string
		switch {
		case elapsed.Seconds() < 1.00:
			elapsedString = color.New(color.FgGreen, color.Bold).Sprint(elapsed)
		case elapsed.Seconds() < 5.00:
			elapsedString = color.New(color.FgYellow, color.Bold).Sprint(elapsed)
		case elapsed.Seconds() > 5.00:
			elapsedString = color.New(color.FgRed, color.Bold).Sprint(elapsed)
		}

		var stdLog bytes.Buffer
		stdLog.WriteString("\t")
		stdLog.WriteString(statusCode)
		stdLog.WriteString("\t")
		stdLog.WriteString(elapsedString)
		stdLog.WriteString("\t")
		stdLog.WriteString(uri)

		fmt.Println(stdLog.String())
	})

	deflateEncoder, _ := zlib.New(zlib.Options{})
	brotliEncoder, _ := brotli.New(brotli.Options{})
	gzipEncoder, _ := pgzip.New(pgzip.Options{
		Level:     pgzip.DefaultCompression,
		BlockSize: 1 << 20,
		Blocks:    4,
	})

	compress, _ := httpcompression.Adapter(
		httpcompression.Compressor(brotli.Encoding, 2, brotliEncoder),
		httpcompression.Compressor(pgzip.Encoding, 1, gzipEncoder),
		httpcompression.Compressor(zlib.Encoding, 0, deflateEncoder),
		httpcompression.Prefer(httpcompression.PreferClient),
		httpcompression.MinSize(165),
	)

	return compress(internal)
}
