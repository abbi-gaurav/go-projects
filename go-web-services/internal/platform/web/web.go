package web

import (
	"time"
	"context"
	"net/http"
	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

type ctxKey int

const KeyValues ctxKey = 1
const TraceIDHeader = "X-Trace-Id"

type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

type App struct {
	*httptreemux.TreeMux
	mw []Middelware
}

func New(mw ...Middelware) *App {
	return &App{
		TreeMux: httptreemux.New(),
		mw:      mw,
	}
}

func (a *App) Handle(verb, path string, handler Handler) {
	handler = wrapMiddleware(handler, a.mw)

	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		v := Values{
			TraceID: uuid.New(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		w.Header().Set(TraceIDHeader, v.TraceID)

		handler(ctx, w, r, params)
	}

	a.TreeMux.Handle(verb, path, h)
}
