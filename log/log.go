package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	file, _ := os.OpenFile("out.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	multi := zerolog.MultiLevelWriter(
		// 2023-07-24 11:11:48 | 200 |  1.001s |       127.0.0.1 | GET     | /
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.DateTime,
		},
		diode.NewWriter(file, 1000, 10*time.Millisecond, func(missed int) {
			log.Error().Msgf("Dropped %d messages", missed)
		}),
	)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

}
