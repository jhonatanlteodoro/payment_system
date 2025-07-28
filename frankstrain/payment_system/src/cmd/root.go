package cmd

import (
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/api"
	"github.com/jhonatanlteodoro/payment_system/src/workers"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:   "system-run",
	Short: "Payment system",
}

func getCmdServer(shutdown chan os.Signal) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Args: ", args)
			fmt.Println("server")
			api.StartApiServer(shutdown)
		},
	}
}

func getCmdFullServer(shutdown chan os.Signal) *cobra.Command {
	return &cobra.Command{
		Use:   "full-server",
		Short: "Run full server - Api Server + All workers",
		Run: func(cmd *cobra.Command, args []string) {
			wg := &sync.WaitGroup{}
			wg.Add(2)

			log.Println("Stating full server...")
			go func() {
				defer wg.Done()
				api.StartApiServer(shutdown)
			}()

			go func() {
				defer wg.Done()
				workers.RunAllWorkers(shutdown)
			}()
			wg.Wait()
		},
	}
}

func getCmdWorker(shutdown chan os.Signal) *cobra.Command {
	var runWorker = &cobra.Command{
		Use:     "worker",
		Short:   "Run worker by name",
		Long:    "Available workers: start-payment, process-payment, notify, ALL",
		Example: "worker --name workerName",
		Run: func(cmd *cobra.Command, args []string) {
			workerName, _ := cmd.Flags().GetString("name")
			if workerName == "" {
				fmt.Println("Please provide a valid worker name")
				cmd.Usage()
				os.Exit(1)
			}

			switch workerName {
			case "ALL":
				workers.RunAllWorkers(shutdown)
			case "start-payment":
				workers.RunStartPaymentWorker(shutdown)
			case "process-payment":
				workers.RunProcessPaymentWorker(shutdown)
			case "notify":
				workers.RunNotifyWorker(shutdown)
			default:
				fmt.Println("Please provide a valid worker name")
				cmd.Usage()
				os.Exit(1)
			}
		},
	}

	runWorker.Flags().StringP(
		"name", "n", "", `Name of the worker. Currently available options are: - w1 - w2`,
	)
	return runWorker
}

func Execute(shutdown chan os.Signal) {
	rootCmd.AddCommand(getCmdServer(shutdown))
	rootCmd.AddCommand(getCmdWorker(shutdown))
	rootCmd.AddCommand(getCmdFullServer(shutdown))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
