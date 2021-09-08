package main

import (
	moduleServer "github.com/ProjectAthenaa/footsites/module"
	"github.com/ProjectAthenaa/newbalance/config"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"time"
)

func init() {
	if err := sonic.RegisterModule(config.Module); err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln("start listener: ", err)
	}

	server := grpc.NewServer()

	module.RegisterModuleServer(server, moduleServer.Server{})

	log.Info("NewBalance Module Initialized")
	if err = server.Serve(listener); err != nil {
		log.Fatalln("start server: ", err)
	}
}
