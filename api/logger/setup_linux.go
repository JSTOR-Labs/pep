package logger

import (
	"log/syslog"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	syslogWriter, err := syslog.New(syslog.LOG_INFO, "pepapi")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create syslog writer")
	} else {
		log.Logger = log.Output(zerolog.SyslogLevelWriter(syslogWriter))
	}
}
