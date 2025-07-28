package main

import (
	"github.com/jhonatanlteodoro/payment_system/src/cmd"
	"github.com/jhonatanlteodoro/payment_system/src/db_scripts"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	shared_deps.NewSharedDependencies(shutdown)

	db_scripts.CreateSchema(shared_deps.GetSharedDependencies().PaymentsDBConn)

	cmd.Execute(shutdown)
}
