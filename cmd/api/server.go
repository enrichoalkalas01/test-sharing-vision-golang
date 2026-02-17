package main

import (
	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/modules"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		modules.LoggerModule,
		modules.ConfigModule,
		modules.DatabaseModule,
		modules.ArticleModule,
		modules.ServerModule,
		modules.RouteModule,
	).Run()
}
