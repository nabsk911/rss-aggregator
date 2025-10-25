package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	//Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file.")
	}
	dbURL := os.Getenv("DB_URL")

	//Connnect to database
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer conn.Close()

	apiCfg := apiConfig{
		DB: db.New(conn),
	}

	go startScrapping(apiCfg.DB, 10, time.Minute)

	//Setup chi router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/users", apiCfg.handleCreateUser)
	router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	router.Get("/feeds", apiCfg.handleGetFeeds)
	router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))
	router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollows))
	router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))
	serve := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is running at port 8080...")
	serve.ListenAndServe()
}
