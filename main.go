package main

import (
	"database/sql"
	"github.com/JohnKucharsky/golang-sqlc/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

var validate *validator.Validate

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := godotenv.Load(); err != nil {
		log.Println("Can't load env")
	}

	port := "8080"
	portEnv := os.Getenv("PORT")
	var portString = portEnv
	if portEnv == "" {
		portString = port
	}

	dbAddress := "postgres://postgres:pass@db:5432/data?sslmode=disable"
	dbAddressEnv := os.Getenv("DB_URL")
	var dbAddressString = dbAddressEnv
	if dbAddressEnv == "" {
		dbAddressString = dbAddress
	}

	conn, err := sql.Open("postgres", dbAddressString)
	if err != nil {
		log.Fatal("Can't connect to db", err.Error())
	}
	log.Print("Connected to db")

	apiCfg := apiConfig{DB: database.New(conn)}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Println("With instance", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/schema",
		"postgres", driver,
	)
	if err != nil {
		log.Println("New with database instance", err.Error())
	}
	err = m.Up()
	if err != nil {
		log.Println("Unable to migrate", err.Error())
	}
	log.Println("Migration succeeded")

	router := chi.NewRouter()

	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins: []string{"https://*", "http://*"},
				AllowedMethods: []string{
					"GET",
					"POST",
					"PUT",
					"DELETE",
					"OPTIONS",
				},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("Can't start the server")
	}
}
