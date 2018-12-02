package mid

import (
	"context"
	"errors"
	"github.com/abbi-gaurav/go-learning-projects/go-web-services/internal/platform/web"
	errors2 "github.com/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
)

func ErrorHandler(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s: Error : Panic caught %s\n", v.TraceID, r)
				web.RespondError(ctx, w, errors.New("unhandled"), http.StatusInternalServerError)
				log.Printf("%s : ERROR : StackTrace %s\n", v.TraceID, debug.Stack())
			}
		}()

		if err := next(ctx, w, r, params); err != nil {
			if errors2.Cause(err) != web.ErrNotFound {
				log.Printf("%s : ERROR : %v\n", v.TraceID, err)
			}

			web.Error(ctx, w, errors2.Cause(err))

			return nil
		}

		return nil
	}

	return h
}
