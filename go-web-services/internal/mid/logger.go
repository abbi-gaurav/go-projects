package mid

import (
	"context"
	"github.com/abbi-gaurav/go-projects/go-web-services/internal/platform/web"
	"log"
	"net/http"
	"time"
)

func RequestLogger(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)

		next(ctx, w, r, params)
		log.Printf("%s : (%d) : %s %s -> %s (%s)",
			v.TraceID,
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now),
		)
		return nil
	}

	return h
}
