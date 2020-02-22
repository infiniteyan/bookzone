package main

import (
	"bookzone/controllers"
	"bookzone/models"
	_ "bookzone/models"
	_ "bookzone/sysinit"
	"bookzone/util"
	"bookzone/util/log"
	"context"
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
	template := iris.HTML("./views", ".html")
	template.AddFunc("doesCollection", models.NewCollection().DoesCollection)
	template.AddFunc("inMap", util.InMap)
	application.RegisterView(template)
	application.StaticWeb("/static", "./static")
	application.StaticWeb("/uploads", "./uploads")
	application.Logger().SetLevel("debug")

	mvc.New(application.Party("/user")).Handle(new(controllers.UserController))
	mvc.New(application.Party("/setting")).Handle(new(controllers.SettingController))
	mvc.New(application.Party("/explore")).Handle(new(controllers.ExploreController))
	mvc.New(application.Party("/account")).Handle(new(controllers.AccountController))
	mvc.New(application.Party("/books")).Handle(new(controllers.DocumentController))
	mvc.New(application.Party("/read")).Handle(new(controllers.DocumentController))
	mvc.New(application.Party("/api")).Handle(new(controllers.DocumentController))
	mvc.New(application.Party("/book")).Handle(new(controllers.BookController))
	mvc.New(application.Party("/manager")).Handle(new(controllers.ManagerController))
	mvc.New(application).Handle(new(controllers.HomeController))
	return application
}

func main() {
	log.SetLogLevel(log.DebugLevel)
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
			log.Infof("shutdown...")
			timeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			app.Shutdown(ctx)
		}
	}()

	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}