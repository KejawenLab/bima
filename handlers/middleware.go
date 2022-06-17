package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/fatih/color"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/compress/zlib"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

type (
	Middleware struct {
		Debug       bool
		Dispatcher  *events.Dispatcher
		Middlewares []configs.Middleware
		Logger      *Logger
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

func (m *Middleware) Register(middlewares []configs.Middleware) {
	for _, v := range middlewares {
		m.Add(v)
	}
}

func (m *Middleware) Add(middlware configs.Middleware) {
	m.Middlewares = append(m.Middlewares, middlware)
}

func (m *Middleware) Sort() {
	sort.Slice(m.Middlewares, func(i, j int) bool {
		return m.Middlewares[i].Priority() > m.Middlewares[j].Priority()
	})
}

func (m *Middleware) Attach(handler http.Handler) http.Handler {
	ctx := context.WithValue(context.Background(), "scope", "middleware")
	start := time.Now()
	internal := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if !m.Debug {
			for _, middleware := range m.Middlewares {
				if stop := middleware.Attach(request, response); stop {
					return
				}
			}

			handler.ServeHTTP(response, request)

			elapsed := time.Since(start)

			m.Logger.Info(ctx, fmt.Sprintf("Execution time: %s", elapsed))

			return
		}

		wrapper := responseWrapper{ResponseWriter: response}
		for _, middleware := range m.Middlewares {
			if stop := middleware.Attach(request, response); stop {
				m.Logger.Debug(ctx, fmt.Sprintf("Middleware stopped by: %s", reflect.TypeOf(middleware).Name()))

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

		fmt.Printf("\t%s\t%s\t%s\n", statusCode, elapsed, uri)
	})

	deflateEncoder, err := zlib.New(zlib.Options{})
	if err != nil {
		m.Logger.Fatal(ctx, err.Error())
	}

	brotliEncoder, err := brotli.New(brotli.Options{})
	if err != nil {
		m.Logger.Fatal(ctx, err.Error())
	}

	gzipEncoder, err := pgzip.New(pgzip.Options{
		Level:     pgzip.DefaultCompression,
		BlockSize: 1 << 20,
		Blocks:    4,
	})
	if err != nil {
		m.Logger.Fatal(ctx, err.Error())
	}

	compress, err := httpcompression.Adapter(
		httpcompression.Compressor(brotli.Encoding, 2, brotliEncoder),
		httpcompression.Compressor(pgzip.Encoding, 1, gzipEncoder),
		httpcompression.Compressor(zlib.Encoding, 0, deflateEncoder),
		httpcompression.Prefer(httpcompression.PreferServer),
		httpcompression.MinSize(165),
	)

	if err != nil {
		m.Logger.Fatal(ctx, err.Error())
	}

	return compress(internal)
}
