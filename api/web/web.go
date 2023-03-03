package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/JSTOR-Labs/pep/api/elasticsearch"
	"github.com/JSTOR-Labs/pep/api/pdfs"
	"github.com/JSTOR-Labs/pep/api/utils"
	"github.com/JSTOR-Labs/pep/api/web/routes"
	"github.com/JSTOR-Labs/pep/api/web/routes/admin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Route struct {
	Handler func(echo.Context) error
	Type    string
	IsAdmin bool
}

func PathHasMethod(rts []Route, method string) bool {
	for _, rt := range rts {
		if rt.Type == method {
			return true
		}
	}
	return false
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	log.Error().Err(err).Msg("HTML Error")
	root, err := utils.GetRoot()
	if err != nil {
		log.Error().Err(err).Msg("Failed to find root path")
	}

	errorPage := filepath.Join(root, fmt.Sprint(code)+".html")
	if err := c.File(errorPage); err != nil {
		log.Error().Err(err).Msg("Failed to find to error page")
	}
}

func Listen(port int) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Starting PEP API")

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

	rtPaths := map[string][]Route{
		"/search": {
			{
				Handler: routes.Search,
				Type:    http.MethodPost,
				IsAdmin: false,
			},
		},
		"/request": {
			{
				Handler: routes.SubmitRequests,
				Type:    http.MethodPost,
				IsAdmin: false,
			},
			{
				Handler: admin.AdminGetRequests,
				Type:    http.MethodGet,
				IsAdmin: true,
			},
			{
				Handler: admin.AdminUpdateRequest,
				Type:    http.MethodPatch,
				IsAdmin: true,
			},
		},
		"/login": {
			{
				Handler: routes.Login,
				Type:    http.MethodPost,
				IsAdmin: false,
			},
		},
		"/version": {
			{
				Handler: routes.VersionInfo,
				Type:    http.MethodPost,
				IsAdmin: false,
			},
		},
		"/pdf/check": {
			{
				Handler: admin.CheckPDFs,
				Type:    http.MethodPost,
				IsAdmin: true,
			},
		},
		"/pdf/:doi/:pdf": {
			{
				Handler: routes.GetPDF,
				Type:    http.MethodGet,
				IsAdmin: false,
			},
			{
				Handler: routes.GetPDF,
				Type:    http.MethodGet,
				IsAdmin: true,
			},
		},
		"/snapshot": {
			{
				Handler: admin.SnapshotStatus,
				Type:    http.MethodGet,
				IsAdmin: true,
			},
			{
				Handler: admin.GetRestoreStatus,
				Type:    http.MethodPost,
				IsAdmin: true,
			},
		},
		"/indices": {
			{
				Handler: admin.GetIndexData,
				Type:    http.MethodGet,
				IsAdmin: true,
			},
		},
	}

	adminGrp := app.Group("/admin")
	adminGrp.Use(middleware.JWT([]byte(viper.GetString("auth.signing_key"))))

	for path, rts := range rtPaths {
		for _, rt := range rts {
			if !rt.IsAdmin {
				switch rt.Type {
				case http.MethodGet:
					app.GET(path, rt.Handler)
				case http.MethodPost:
					app.POST(path, rt.Handler)
				case http.MethodPatch:
					app.PATCH(path, rt.Handler)
				}
			} else {
				switch rt.Type {
				case http.MethodGet:
					adminGrp.GET(path, rt.Handler)
				case http.MethodPost:
					adminGrp.POST(path, rt.Handler)
				case http.MethodPatch:
					adminGrp.PATCH(path, rt.Handler)
				}
			}
		}
	}

	exPath, err := utils.GetExecutablePath()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find executable path")
		return
	}
	root, err := utils.GetRoot()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find root")
		return
	}
	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			val, ok := rtPaths[c.Request().URL.Path]
			if ok && PathHasMethod(val, c.Request().Method) {
				return true
			}
			return false
		},
		Root:  root,
		HTML5: true,
	}))

	app.HTTPErrorHandler = customHTTPErrorHandler

	if _, err := os.Stat(filepath.Join(exPath, "content", "pdfindex.dat")); err != nil {
		log.Info().Msg("Generating PDF Index. This may take some time.")
		pdfPath, err := utils.GetPDFPath()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get pdf path")
			return
		}
		pdfs.GenerateIndex(pdfPath)
	}

	log.Fatal().Err(app.Start(fmt.Sprintf(":%d", port))).Int("port", port).Msg("Failed to listen")
}
