package main

import (
	"context"
	"github.com/talon-one/assignment-props/access_logger"
	"net/http"

	"github.com/talon-one/assignment-props/dat"
	"github.com/talon-one/assignment-props/lib"
	"github.com/talon-one/assignment-props/logger"
	"go.uber.org/zap"
)

func main() {
	logger := logger.InitLogger()
	ctx := context.Background()

	db, err := dat.InitDB(ctx, logger)
	if err != nil {
		logger.Fatal("failed to init DB", zap.Error(err))
	}
	defer dat.ExitDb(logger)

	_, err = db.Exec(ctx, dat.GetSchemaSQL())
	if err != nil {
		logger.Fatal("failed to execute SQL schema", zap.Error(err))
	}

	access_logger.Init()
	// TODO Good to create a number of workers inside
	go access_logger.Run(logger)

	address := ":8080"
	server := http.Server{
		Addr:    address,
		Handler: lib.NewRouter(),
	}

	logger.Info("Listening...", zap.String("address", address))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("ListenAndServe failed", zap.Error(err))
	}
}
