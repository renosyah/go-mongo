package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/renosyah/go-mongo/auth"
	"github.com/renosyah/go-mongo/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbPool  *mongo.Database
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use: "app",
	PreRun: func(cmd *cobra.Command, args []string) {
		router.Init(dbPool)
		auth.Init(dbPool, &auth.Option{
			Enable: viper.GetBool("auth.enable"),
		})
	},
	Run: func(cmd *cobra.Command, args []string) {

		r := mux.NewRouter()

		apiRouter := r.PathPrefix("/api/v1").Subrouter()
		apiRouter.Use(auth.AuthenticationMiddleware)

		// users api
		apiRouter.Handle("/users", router.HandlerFunc(router.HandlerAddUser)).Methods(http.MethodPost)
		apiRouter.Handle("/users-list", router.HandlerFunc(router.HandlerAllUser)).Methods(http.MethodPost)
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerOneUser)).Methods(http.MethodGet)
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerUpdateUser)).Methods(http.MethodPut)
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerDeleteUser)).Methods(http.MethodDelete)

		r.PathPrefix("/").Handler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "max-age=-1")
				http.StripPrefix("/", http.FileServer(http.Dir("files"))).ServeHTTP(w, r)
			}),
		)

		port := viper.GetInt("app.port")
		p := os.Getenv("PORT")
		if p != "" {
			port, _ = strconv.Atoi(p)
		}

		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			ReadTimeout:  time.Duration(viper.GetInt("read_timeout")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt("write_timeout")) * time.Second,
			IdleTimeout:  time.Duration(viper.GetInt("idle_timeout")) * time.Second,
		}

		done := make(chan bool, 1)
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			log.Println("Server is shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.SetKeepAlivesEnabled(false)
			if err := server.Shutdown(ctx); err != nil {
				log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
			}
			close(done)
		}()

		log.Println("Server is ready to handle requests at", fmt.Sprintf(":%d", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", fmt.Sprintf(":%d", port), err)
		}

		<-done
		log.Println("Server stopped")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is github.com/renosyah/go-mysql/.server.toml)")
	cobra.OnInitialize(initConfig, initDB)
}

func initDB() {

	dbConfig := fmt.Sprintf("mongodb://%s:%d/?compressors=%s&gssapiServiceName=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.compressors"),
		viper.GetString("database.gssapi_service_name"))

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConfig))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error open DB: %v\n", err))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error connect to DB service: %v\n", err))
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error ping DB: %v\n", err))
		return
	}

	dbPool = client.Database(viper.GetString("database.name"))
}

func initConfig() {
	viper.SetConfigType("toml")
	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/go-mysql")
		viper.SetConfigName(".server")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
