package main

import (
	"log"
	"net/http"

	"git.carried.ru/opcup23/backend/api"
	"git.carried.ru/opcup23/backend/backend/gorm"
)

func main() {
	var err error

	conf := Config{}
	if err = conf.Init(); err != nil {
		log.Fatal(err)
	}

	s := &api.Server{}
	switch conf.DB.Backend {
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite":
		s.DB, err = gorm.Open(conf.DSN, conf.DB.Backend)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unsupported backend: %s", conf.DB.Backend)
	}

	s.DB.Setup()

	s.DevMode = conf.Server.DevMode
	s.TokenSecret = []byte(conf.Server.Secret)

	s.RegisterAPI()

	log.Fatal(http.ListenAndServe(":"+conf.Server.Port, s.Mux))
}
