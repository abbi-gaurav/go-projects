package integration

import (
	"context"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/app"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	options := opts.DefaultOptions
	informer := newFakeInformer(ctx)
	application := app.New(&options, informer)
	httptest.NewServer(application.ServerMux)
	err := application.Run(ctx.Done())

	println("Set up done...", err)

	retCode := m.Run()
	os.Exit(retCode)
}

func Test_Create(t *testing.T) {
	cake, err := createNewCake("black-forest", "bavaria", "choclate")
	print(cake, err)
}
