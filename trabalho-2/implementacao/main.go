package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"trabalho-2/m/routes"

	"cloud.google.com/go/firestore"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "block-firebase-8c6e4")

	if err != nil {
		log.Fatalf("Error at get a firestore client: %s\n", err.Error())
	}

	defer client.Close()

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))

	router.Handle("/static/", http.StripPrefix("/static", fs))
	router.HandleFunc("/", routes.IndexRoute)
	router.HandleFunc("/add-block", routes.AddBlockRoute(ctx, client))

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
