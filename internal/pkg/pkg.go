package pkg

import (
	"flag"
	"yandex-go-advanced/internal/config"
)

func ParseFlag() {
	config.Init()
	flag.Parse()
}
