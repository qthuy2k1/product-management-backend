package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/joho/godotenv"
)

const PORT = 3000

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	database, err := repositories.Initialize(dbUrl)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer database.Close()
	boil.SetDB(database)

	// db cache
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASSWORD")
	redis := repositories.RedisInitialize(redisPort, redisPass)

	routerHandlers := InitRoutes(database, redis)

	log.Printf("Started server on %d", PORT)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), routerHandlers); err != nil {
		log.Fatal(err)
	}
}
