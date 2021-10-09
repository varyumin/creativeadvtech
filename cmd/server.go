/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Conn *redis.Client

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Conn = red.Connect()
		log := logrus.New()
		log.SetOutput(os.Stdout)
		prometheus.MustRegister(opsProcessed)
		recordMetrics()

		log.Info("Starting the app...")

		router := mux.NewRouter()
		router.HandleFunc("/health", Health)
		router.Handle("/metrics", promhttp.Handler())
		router.HandleFunc("/", GetCounter).Methods("GET")
		router.HandleFunc("/", UpdateCounter).Methods("POST")
		router.HandleFunc("/", IncrementCounter).Methods("PUT")
		router.HandleFunc("/", DeleteCounter).Methods("DELETE")

		serv := http.Server{
			Addr:    net.JoinHostPort(web.Host, web.Port),
			Handler: router,
		}

		go serv.ListenAndServe()

		log.Infof("Started the app...: %s:%s", web.Host, web.Port)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		<-interrupt

		log.Info("Stopping app...")

		timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		err := serv.Shutdown(timeout)
		if err != nil {
			log.Error("Error when shutdown app: %v", err)
		}

		log.Info("The app stopped")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&web.Host, "web-server", "s", "0.0.0.0", "Host web server")
	serverCmd.PersistentFlags().StringVarP(&web.Port, "web-port", "p", "8080", "Port web server")
	serverCmd.PersistentFlags().StringVarP(&red.Host, "redis-server", "r", "127.0.0.1", "Redis server")
	serverCmd.PersistentFlags().StringVarP(&red.Port, "redis-port", "d", "6379", "Redis port")
	serverCmd.PersistentFlags().StringVarP(&red.Password, "redis-password", "P", "", "Redis password")
	serverCmd.PersistentFlags().IntVarP(&red.DB, "redis-db", "D", 0, "Redis DB")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
