package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"webook/internal/bootstrap"
	"webook/internal/repository"
	articlerepo "webook/internal/repository/article"
	"webook/internal/repository/cache"
	userdao "webook/internal/repository/dao"
	articledao "webook/internal/repository/dao/article"
	"webook/internal/service"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/internal/web/middleware"
)

func main() {
	initConfig()

	l := bootstrap.InitLogger()
	db := bootstrap.InitDB(l)
	redisClient := bootstrap.InitRedis()

	// repository 层
	userDAO := userdao.NewUserDAO(db)
	userCache := cache.NewUserCache(redisClient)
	userRepo := repository.NewUserRepository(userDAO, userCache)

	articleDAO := articledao.NewArticleDAO(db)
	articleCache := cache.NewRedisArticleCache(redisClient)
	articleRepo := articlerepo.NewArticleRepositoryWithCache(articleDAO, db, l, userRepo, articleCache)
	interactiveSvc := service.NewDBInteractiveService(db)

	// service 层
	userSvc := service.NewUserService(userRepo)
	articleSvc := service.NewArticleService(articleRepo, l, nil)

	// handler & middleware
	jwtHandler := ijwt.NewRedisJWTHandler(redisClient)
	userHdl := web.NewUserHandler(userSvc, jwtHandler)
	articleHdl := web.NewArticleHandler(articleSvc, interactiveSvc, l)

	server := gin.Default()
	// 允许前端开发端口跨域访问
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
	}))
	server.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "msg": "ok"})
	})
	server.Use(middleware.NewLoginJwtMiddlewareBuilder(jwtHandler).
		IgnorePaths("/health").
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		IgnorePaths("/articles/pub").
		Build())

	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
