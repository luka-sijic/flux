package server

import (
	"log"
	"os"

	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/routes"
	"github.com/luka-sijic/flux/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file", err)
	}
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS")},
		AllowMethods: []string{echo.OPTIONS, echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderCookie,
		},
		ExposeHeaders:    []string{echo.HeaderSetCookie},
		AllowCredentials: true,
	}))

	//e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app, err := database.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Init(); err != nil {
		log.Fatalf("could not create users table: %v", err)
	}

	userSvc := service.NewService(app)
	friendSvc := service.NewService(app)
	routes.UserRoutes(e, userSvc)
	routes.FriendRoutes(e, friendSvc)

	e.Logger.Fatal(e.Start(":8015"))
}
