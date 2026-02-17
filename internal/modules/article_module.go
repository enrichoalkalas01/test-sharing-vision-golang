package modules

import (
	"github.com/enrichoalkalas01/test-sharing-vision-golang/internal/handler"
	repository "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/repository/article"
	usecase "github.com/enrichoalkalas01/test-sharing-vision-golang/internal/usecase/article"
	"go.uber.org/fx"
)

var ArticleModule = fx.Module("article",
	fx.Provide(repository.NewArticleRepository),
	fx.Provide(usecase.NewArticleUsecase),
	fx.Provide(handler.NewArticleHandler),
)
