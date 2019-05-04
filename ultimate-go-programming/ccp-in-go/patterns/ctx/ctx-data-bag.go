package main

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func processRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	handleResponse(ctx)
}

func handleResponse(ctx context.Context) {
	fmt.Printf("Handling response for %v (auth %v)", UserID(ctx), AuthToken(ctx))
}

func main() {
	processRequest("jane", "abc123")
}
