package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"github.com/XXueTu/wise/internal/config"
	"github.com/XXueTu/wise/internal/handler"
	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
)

var configFile = flag.String("f", "etc/wise-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithFileServer("/", http.Dir("dist")))
	defer server.Stop()

	// 添加中间件处理根路径重定向
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.ServeFile(w, r, "dist/index.html")
				return
			}
			next(w, r)
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 初始化数据库
	model.InitDB()
	server.Start()
}
