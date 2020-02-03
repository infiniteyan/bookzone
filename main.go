package main

import (
	_ "bookzone/sysinit"
	"bookzone/controllers"
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func newApplication() *iris.Application {
	application := iris.New()
	application.RegisterView(iris.HTML("./views", ".html"))
	application.StaticWeb("/static", "./static")
	application.Logger().SetLevel("debug")

	accountParty := application.Party("/account")
	mvc.New(accountParty).Handle(new(controllers.AccountController))
	mvc.New(application).Handle(new(controllers.HomeController))

	return application
}

func main() {
	app := newApplication()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			)
		select {
		case <- ch:
			fmt.Println("shutdown...")
			timeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			app.Shutdown(ctx)
		}
	}()

	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}