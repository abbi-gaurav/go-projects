package web

type Middelware func(Handler) Handler

func wrapMiddleware(handler Handler, mw []Middelware) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		if mw[i] != nil {
			handler = mw[i](handler)
		}
	}

	return handler
}
