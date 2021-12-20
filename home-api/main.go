package main

import (
	"fmt"
	"net/http"

	"github.com/JSTOR-Labs/pep/home-api/models"
	"github.com/JSTOR-Labs/pep/home-api/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/pep")
	viper.SetConfigName("home")
	viper.SetConfigType("toml")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "db.sqlite")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	driver := viper.GetString("database.driver")
	var dialector gorm.Dialector
	switch driver {
	case "sqlite":
		dialector = sqlite.Open(viper.GetString("database.dsn"))
	case "postgres":
		dialector = postgres.Open(viper.GetString("database.dsn"))
	default:
		log.Fatal().Msgf("unsupported database driver: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	migrate(db)

	w := web.New(db)

	r.Mount("/api/v1", w.V1())
	port := viper.GetInt("http.port")
	log.Info().Int("port", port).Msg("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Device{}, &models.Asset{}, &models.Version{})
}
