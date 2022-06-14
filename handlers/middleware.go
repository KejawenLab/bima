package handlers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/compress/zlib"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

type Middleware struct {
	Dispatcher     *events.Dispatcher
	Middlewares    []configs.Middleware
	MuxMiddlewares []configs.MuxMiddlewares
	Logger         *Logger
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

	internal := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, middleware := range m.Middlewares {
			if stop := middleware.Attach(request, response); stop {
				m.Logger.Info(ctx, fmt.Sprintf("Middleware stopped by: %s", reflect.TypeOf(middleware).Name()))

				return
			}
		}

		m.Logger.Info(ctx, "Dispatching request event")
		m.Dispatcher.Dispatch(events.REQUEST_EVENT, &events.Request{
			HttpRequest: request,
		})

		m.Logger.Info(ctx, "Dispatching response event")
		m.Dispatcher.Dispatch(events.RESPONSE_EVENT, &events.Response{
			ResponseWriter: response,
		})

		handler.ServeHTTP(response, request)
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

	m.Logger.Info(ctx, "Attach compression middleware")
	last := compress(internal)
	for _, middleware := range m.MuxMiddlewares {
		last = middleware(last)
	}

	m.Logger.Info(ctx, fmt.Sprintf("Total middlewares: %d", len(m.MuxMiddlewares)+1))

	return last
}
