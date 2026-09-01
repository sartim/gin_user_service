package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gin-shop-api/internal/config"
	"gin-shop-api/internal/controllers"
	"gin-shop-api/internal/helpers/crypto"
	"gin-shop-api/internal/middleware"
	"gin-shop-api/internal/migrations"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	actionRunServer       = "run-server"
	actionMigrate         = "migrate"
	actionRollback        = "rollback"
	actionDropTables      = "drop-tables"
	actionCreateSuperUser = "create-super-user"
)

func main() {
	action := flag.String("action", actionRunServer, "run-server, migrate, drop-tables, or create-super-user")
	flag.Parse()

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("load .env: %v", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}
	db, err := repository.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("database pool unavailable: %v", err)
	}
	defer sqlDB.Close()

	switch *action {
	case actionRunServer:
		err = runServer(db, cfg)
	case actionMigrate:
		err = migrate(db)
	case actionRollback:
		err = migrations.Down(context.Background(), db)
	case actionDropTables:
		err = dropTables(db)
	case actionCreateSuperUser:
		err = createSuperUser(db)
	default:
		err = fmt.Errorf("unknown action %q", *action)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func newRouter(db *gorm.DB, cfg config.App) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	metrics := middleware.NewMetrics()
	router.Use(gin.Logger(), gin.Recovery(), metrics.Middleware(), middleware.CORS(cfg.AllowedOrigins))
	router.GET("/metrics", metrics.Handler(cfg.MetricsToken))

	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/health/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	api := router.Group("/api/v1")
	authController := controllers.NewAuthController(db, cfg.SecretKey, cfg.AccessTokenTTL)
	authController.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.RequireAuth(db, cfg.SecretKey))
	admin := protected.Group("")
	admin.Use(middleware.RequireAdmin())
	controllers.NewUserController(db).RegisterRoutes(admin)
	controllers.NewRoleController(db).RegisterRoleRoutes(admin)
	controllers.NewPermissionController(db).RegisterPermissionRoutes(admin)

	return router
}

func runServer(db *gorm.DB, cfg config.App) error {
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           newRouter(db, cfg),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("user service listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-shutdownSignal.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}

func migrate(db *gorm.DB) error {
	return migrations.Up(context.Background(), db)
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.RolePermission{}, &models.UserPermission{},
		&models.Permission{}, &models.Role{}, &models.User{},
	)
}

func createSuperUser(db *gorm.DB) error {
	firstName := prompt("first_name:")
	lastName := prompt("last_name:")
	email := strings.ToLower(prompt("email:"))
	password := prompt("password:")
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
		IsActive:  true,
		IsAdmin:   true,
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("create super user: %w", err)
	}
	log.Printf("created super user %s", user.Email)
	return nil
}

func prompt(label string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		value, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("read input: %v", err)
		}
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
}
