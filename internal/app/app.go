package app

import (
	"github.com/tanveerprottoy/templates-go-gin/internal/app/module/auth"
	"github.com/tanveerprottoy/templates-go-gin/internal/app/module/content"
	"github.com/tanveerprottoy/templates-go-gin/internal/app/module/user"
	"github.com/tanveerprottoy/templates-go-gin/internal/pkg"
	"github.com/tanveerprottoy/templates-go-gin/internal/pkg/constant"
	"github.com/tanveerprottoy/templates-go-gin/internal/pkg/middleware"
	"github.com/tanveerprottoy/templates-go-gin/internal/pkg/router"
	routerPkg "github.com/tanveerprottoy/templates-go-gin/internal/pkg/router"
	"github.com/tanveerprottoy/templates-go-gin/pkg/data/sql/mysql"

	"github.com/go-playground/validator/v10"
	// "go.uber.org/zap"
)

// App struct
type App struct {
	DBClient      *mysql.Client
	gin           *pkg.Gin
	Middlewares   []any
	AuthModule    *auth.Module
	UserModule    *user.Module
	ContentModule *content.Module
	Validate      *validator.Validate
}

func NewApp() *App {
	a := new(App)
	a.initComponents()
	return a
}

func (a *App) initDB() {
	a.DBClient = mysql.GetInstance()
}

func (a *App) initMiddlewares() {
	authMiddleWare := middleware.NewAuthMiddleware(a.AuthModule.Service)
	a.Middlewares = append(a.Middlewares, authMiddleWare)
}

func (a *App) initModules() {
	a.UserModule = user.NewModule(a.MongoDBClient.DB, a.DBClient.DB, a.Validate)
	a.AuthModule = auth.NewModule(a.UserModule.Service)
	a.ContentModule = content.NewModule(a.DBClient.DB)
}

func (a *App) initModuleRouters() {
	m := a.Middlewares[0].(*middleware.AuthMiddleware)
	routerPkg.RegisterUserRoutes(a.gin.Engine, constant.V1, a.UserModule, m)
	routerPkg.RegisterContentRoutes(a.gin.Engine, constant.V1, a.ContentModule)
}

/* func (a *App) initLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"proxy.log",
	}
	cfg.Build()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	task := "taskName"
	logger.Info("failed to do task",
		// Structured context as strongly typed Field values.
		zap.String("url", task),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
} */

// Init app
func (a *App) initComponents() {
	a.initDB()
	a.gin = pkg.NewGin()
	a.initModules()
	a.initMiddlewares()
	a.initModuleRouters()
	// a.initLogger()
	// setup global middlewares
	router.RegisterGlobalMiddlewares(a.gin.Engine)
}

// Run app
func (a *App) Run() {
	a.gin.Engine.Run(":8080")
}

// Run app
func (a *App) RunTLS() {
	a.gin.Engine.Run(":443", "cert.crt", "key.key")
}
