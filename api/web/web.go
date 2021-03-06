package web

import (
	"errors"
	"fmt"
	"time"

	"github.com/JSTOR-Labs/pep/api/discovery"
	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/web/routes"
	"github.com/JSTOR-Labs/pep/api/web/routes/admin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Listen(port int) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msgf("Starting PEP API")
	app := echo.New()
	app.HideBanner = true
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	err := errors.New("elasticsearch not connected")
	tries := 0
	for err != nil && tries < 10 {
		err = elasticsearch.Connect()
		if err != nil {
			log.Error().Err(err).Msg("Failed to connect to Elasticsearch")
		}
		time.Sleep(time.Second * (5 * time.Duration(tries)))
		tries++
	}

	app.POST("/search", routes.Search)
	// TODO: healthcheck
	app.POST("/request", routes.SubmitRequests)
	app.POST("/login", routes.Login)
	app.GET("/version", routes.VersionInfo)
	adminGrp := app.Group("/admin")
	adminGrp.Use(middleware.JWT([]byte(viper.GetString("auth.signing_key"))))
	adminGrp.GET("/request", admin.AdminGetRequests)
	adminGrp.PATCH("/request", admin.AdminUpdateRequest)
	adminGrp.POST("/pdf/check", admin.CheckPDFs)
	adminGrp.GET("/pdf/:doi/:pdf", admin.GetPDF)
	adminGrp.GET("/usb", admin.GetUSBDevices)
	adminGrp.POST("/usb", admin.FormatUSBDevice)
	adminGrp.POST("/usb/:name", admin.BuildFlashDrive)
	adminGrp.POST("/snapshot", admin.SnapshotStatus)
	adminGrp.GET("/snapshot", admin.GetRestoreStatus)
	adminGrp.GET("/indices", admin.GetIndexData)
	app.Static("/", viper.GetString("web.root"))

	if !viper.GetBool("runtime.flash_drive_mode") {
		svc, err := discovery.SetupDiscovery(port)
		if err != nil {
			app.Logger.Warn("discovery setup failed")
		} else {
			defer discovery.ShutdownDiscovery(svc)
		}
	}
	log.Fatal().Err(app.Start(fmt.Sprintf(":%d", port))).Int("port", port).Msg("Failed to listen")
}
