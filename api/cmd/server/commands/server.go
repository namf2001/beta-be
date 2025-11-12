package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"beta-be/cmd/server/commands/router"
	"beta-be/internal/controller/user"
	"beta-be/internal/pkg/config"
	"beta-be/internal/pkg/logger"
	"beta-be/internal/pkg/middleware/auth"
	"beta-be/internal/repository"
	"beta-be/internal/repository/ent"

	_ "entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configEnv string
)

func init() {
	serverCmd.PersistentFlags().StringVarP(&configEnv, "config", "c", "api/configs/.env", "Start server with provided configuration file")
	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Start API server",
	Example:      "beta-be server -c configs/.env",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		setup()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func setup() {
	if err := config.InitializeAppConfig(configEnv); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{})
	}

	log.Println("Setting up server...")
}

func run() error {
	ctx := context.Background()
	// setup databases
	entClient, err := setupPostgresConnection(ctx)
	if err != nil {
		return err
	}

	repo := repository.New(entClient)

	handler, err := initRouter(ctx, repo)
	if err != nil {
		panic(err)
	}

	// setup http server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Gracefull Shutdown
	go func() {
		logger.InfoF("success to listen and serve on :%d", logrus.Fields{}, config.AppConfig.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// make blocking channel and waiting for a signal
	<-quit
	logger.Info("shutdown server ...", logrus.Fields{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	logger.Info("timeout of 5 seconds.", logrus.Fields{})
	logger.Info("server exiting", logrus.Fields{})
	return nil
}

func initRouter(ctx context.Context, repo repository.Registry) (http.Handler, error) {
	gin.SetMode(gin.ReleaseMode)
	if config.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	}
	// create a new router instance
	g := gin.New()

	userCtrl := user.New(repo)
	authMiddleware, err := auth.NewJWTAuth()
	if err != nil {
		return nil, err
	}
	r := router.New(ctx, userCtrl, authMiddleware)

	return r.Handler(g)
}

func setupPostgresConnection(ctx context.Context) (*ent.Client, error) {
	// Create an ent.Driver from the given data source name (dsn).
	drv, err := sql.Open(dialect.Postgres, config.AppConfig.PgUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to postgres: %v", err)
	}

	// Create an ent.Client with the driver.
	client := ent.NewClient(ent.Driver(drv))

	// Optional: Set connection limits on the underlying sql.DB.
	db := drv.DB()
	db.SetMaxOpenConns(config.AppConfig.PgPoolMaxOpenConns)
	db.SetMaxIdleConns(config.AppConfig.PgPoolMaxIdleConns)
	db.SetConnMaxLifetime(15 * time.Minute)

	// Optional: Ping the database to ensure connection is established.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %v", err)
	}

	return client, nil
}
