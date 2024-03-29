package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/KejawenLab/bima/v4/mocks/middlewares"
	"github.com/stretchr/testify/mock"
)

type SlowMiddleware struct {
	hold int
}

func (a *SlowMiddleware) Attach(_ *http.Request, response http.ResponseWriter) bool {
	time.Sleep(time.Duration(a.hold) * time.Second)

	return false
}

func (a *SlowMiddleware) Priority() int {
	return 0
}

func Test_Middleware_Debug_True_Without_Stop(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("PATCH", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_Without_Stop_Method_Head(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("HEAD", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_With_Stop(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(true).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_False_Without_Stop(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: false,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("PATCH", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_False_With_Stop(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(true).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: false,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)

	w.WriteHeader(http.StatusPermanentRedirect)
}

func Test_Middleware_Debug_True_Without_Stop_Return_3XX(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPermanentRedirect)
	})

	req := httptest.NewRequest("POST", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_Without_Stop_Return_4XX(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Once()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Once()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	req := httptest.NewRequest("PUT", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_Without_Stop_Return_5XX(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Maybe()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Maybe()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})

	req := httptest.NewRequest("DELETE", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_Slow_Middleware(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Maybe()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Maybe()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
		&SlowMiddleware{hold: 3},
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})

	req := httptest.NewRequest("DELETE", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}

func Test_Middleware_Debug_True_Slowest_Middleware(t *testing.T) {
	middleware1 := mocks.NewMiddleware(t)
	middleware1.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware1.On("Priority").Return(1).Maybe()

	middleware2 := mocks.NewMiddleware(t)
	middleware2.On("Attach", mock.Anything, mock.Anything).Return(false).Once()
	middleware2.On("Priority").Return(2).Maybe()

	factory := Factory{
		Debug: true,
	}
	factory.Register([]Middleware{
		middleware1,
		middleware2,
		&SlowMiddleware{hold: 5},
	})

	factory.Sort()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})

	req := httptest.NewRequest("DELETE", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	factory.Attach(handler).ServeHTTP(w, req)

	middleware1.AssertExpectations(t)
	middleware2.AssertExpectations(t)
}
