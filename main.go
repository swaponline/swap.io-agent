package main

import (
	"log"
	"swap.io-agent/src/configLoader"
	"swap.io-agent/src/httpHandler"
	"swap.io-agent/src/httpServer"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/socketServer"
)

func main() {
	registry := serviceRegistry.NewServiceRegistry()

	err := configLoader.InitializeConfig()
	if err != nil {panic(err)}

	db, err := redisStore.InitializeDB()
	if err != nil {
		log.Panicf("redisStore not initialize, err: %v", err)
	}

	err = registry.RegisterService(&db)
	if err != nil {
		log.Panicln(err.Error())
	}

	socketServerEntity := socketServer.InitializeServer()
	err = registry.RegisterService(socketServerEntity)
	if err != nil {
		log.Panicln(err.Error())
	}

	httpServerEntity := httpServer.InitializeServer()
	err = registry.RegisterService(httpServerEntity)
	if err != nil {
		log.Panicln(err.Error())
	}

	httpHandlerEntity := httpHandler.InitializeServer()
	err = registry.RegisterService(httpHandlerEntity)
	if err != nil {
		log.Panicln(err.Error())
	}

	registry.StartAll()

	<-make(chan struct{})
}