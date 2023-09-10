package main

import (
	"bufio"
	"flag"
	"fmt"
	"gin-shop-api/internal/config"
	"gin-shop-api/internal/controllers"
	"gin-shop-api/internal/helpers/crypto"
	"gin-shop-api/internal/helpers/env"
	"gin-shop-api/internal/middleware"
	"gin-shop-api/internal/models"
	"gin-shop-api/internal/repository"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var action string

type Action struct {
	Env             string
	RunServer       string
	CreateTables    string
	DropTables      string
	CreateSuperUser string
	SetupService    string
}

var actions = Action{
	Env:             "env",
	RunServer:       "run-server",
	CreateTables:    "create-tables",
	DropTables:      "drop-tables",
	CreateSuperUser: "create-super-user",
	SetupService:    "setup-service",
}

type EnvVars struct {
	ENV    string
	PORT   string
	DB_URL string
}

func init() {
	gin.ForceConsoleColor()

	// Override logging
	log.SetPrefix("\u001b[31mERROR: \u001b[0m")
	log.SetFlags(log.LstdFlags | log.Ldate | log.Lmicroseconds | log.Llongfile)

	// Load environment variables
	// Check if .env file exists
	if _, err := os.Stat(".env"); err == nil {
		env.LoadEnvVars()
	}

	repository.ConnectToDb()
}

func healthCheckRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func RegisterRoutes(r *gin.Engine) {
	APIVersion := "/api/v1"

	// Auth API
	controllers.RegisterAuthRoutes(r.Group(APIVersion))

	// User API
	userCtrl := controllers.NewUserController(repository.DB)
	userCtrl.RegisterUserRoutes(r.Group(APIVersion))
}

func runServer() {
	if action == "run-server" {
		r := gin.Default()
		r.Use(gin.LoggerWithFormatter(
			func(param gin.LogFormatterParams) string {
				// your custom format
				return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
					param.ClientIP,
					param.TimeStamp.Format(time.RFC1123),
					param.Method,
					param.Path,
					param.Request.Proto,
					param.StatusCode,
					param.Latency,
					param.Request.UserAgent(),
					param.ErrorMessage,
				)
			}))
		r.Use(gin.Recovery())
		r.Use(middleware.CORSMiddleware())
		healthCheckRoute(r)
		RegisterRoutes(r)

		r.Run()
	}
}

func createUserTable() {
	err := repository.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
}

func createRoleTable() {
	err := repository.DB.AutoMigrate(&models.Role{})
	if err != nil {
		panic(err)
	}
}

func createUserPermissionTable() {
	err := repository.DB.AutoMigrate(&models.UserPermission{})
	if err != nil {
		panic(err)
	}
}

func dropUserTable() {
	err := repository.DB.Migrator().DropTable(&models.User{})
	if err != nil {
		panic(err)
	}
}

func makeMigrations() {
	if action == "create-tables" {
		createTables()
		fmt.Println("Finished running migrations")
	}
}

func createTables() {
	repository.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	createUserTable()
	createRoleTable()
	createUserPermissionTable()
}

func dropTables() {
	if action == "drop-tables" {
		dropUserTable()
		fmt.Println("Finished dropping tables")
	}
}

func createSuperUser() {
	firstName := StringPrompt("first_name:")
	lastName := StringPrompt("last_name:")
	email := StringPrompt("email:")
	password := StringPrompt("password:")

	user := models.User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  crypto.HashPassword(password),
		IsActive:  true,
	}
	result := repository.DB.Create(&user)

	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("Finished creating super user record")
}

func SetupService() {
	// TODO setup service
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func setupEnvVars() {
	envConfig := config.Config{EnvVar: config.Environment}
	env := envConfig.Get()

	portConfig := config.Config{EnvVar: config.Port}
	port := portConfig.Get()

	dbUrlConfig := config.Config{EnvVar: config.DbUrl}
	dbUrl := dbUrlConfig.Get()

	service := EnvVars{
		ENV:    env,
		PORT:   port,
		DB_URL: dbUrl,
	}
	// Scaffold from template
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/templates/.env", config.CurrDir()))
	if err != nil {
		log.Panic(err)
	}

	// Path to env file name
	envFileName := fmt.Sprint(".env")
	path := fmt.Sprintf("%s/%s", config.CurrDir(), envFileName)

	// Generate file to systemd path
	outputFile, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, service)
	if err != nil {
		log.Panic(err)
	}
}

func launchAction() {
	switch action {
	case actions.Env:
		setupEnvVars()
	case actions.RunServer:
		runServer()
	case actions.DropTables:
		dropTables()
	case actions.CreateTables:
		createTables()
	case actions.DropTables:
		dropTables()
	case actions.CreateSuperUser:
		createSuperUser()
	}
}

func main() {
	flag.StringVar(&action,
		"action", "",
		"action e.g. run-server, create-tables, drop-tables, create-super-user")
	flag.Parse()
	if action != "" {
		launchAction()
	}
}
