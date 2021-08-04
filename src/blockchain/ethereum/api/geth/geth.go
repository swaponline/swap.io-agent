package geth

import (
	"log"
	"os"
)

type Geth struct {
	baseUrl string
}

func InitializeGeth() *Geth {
	gethBaseUrl := os.Getenv("GETH_BASE_URL")
	if len(gethBaseUrl) == 0 {
		log.Panicln("SET GETH_BASE_URL IN ENV")
	}

	return &Geth{
		baseUrl: gethBaseUrl,
	}
}

func (*Geth) Start() {}
func (*Geth) Stop() error {
	return nil
}
func (*Geth) Status() error {
	return nil
}