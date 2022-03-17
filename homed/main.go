package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("toml")
	viper.SetConfigName("homed")
	viper.AddConfigPath("/etc/pep")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func main() {
	cfg := flag.String("url", "", "Home URL")
	version := flag.String("version", "", "Write a version to the config file")
	flag.Parse()
	if *cfg != "" {
		viper.Set("system.version", *version)
		viper.Set("system.buildtime", time.Now())
		viper.Set("homed.url", *cfg)
		viper.WriteConfigAs("homed.toml")
		return
	}

	done := make(chan bool)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		done <- true
	}()

	d := NewDaemon(&http.Client{
		Timeout: time.Second * 10,
	}, done)
	d.Run()
}
