package routes

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"trabalho-2/m/database"
	"trabalho-2/m/model"

	"cloud.google.com/go/firestore"
)

func enableCORS(rw *http.ResponseWriter, r *http.Request) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "https://evening-cove-12029.herokuapp.com")
	(*rw).Header().Set("Access-Control-Allow-Methods",
		fmt.Sprintf("%s, %s", http.MethodOptions, http.MethodPost))
	(*rw).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// IndexRoute -> Index route
func IndexRoute(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw,
			"This route only accept Get requests.", http.StatusMethodNotAllowed)
		return
	}

	if err := template.Must(
		template.ParseFiles("templates/index.html")).Execute(rw, nil); err != nil {
		log.Println(err)
		http.Error(rw,
			"Could not render the index page.", http.StatusInternalServerError)
		return
	}
}

// AddBlockRoute -> Route to add a new Block to firestore.
func AddBlockRoute(ctx context.Context, client *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		enableCORS(&rw, r)

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(rw,
				"This route only accept Post requests.", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		block := model.NewBlock()

		if err := block.Deserialize(r.Body); err != nil {
			log.Println(err)
			http.Error(rw, "Could not parse the payload.", http.StatusBadRequest)
			return
		}

		if err := database.AddBlock(ctx, client, model.Collection, block); err != nil {
			log.Println(err)
			http.Error(rw, "Could not store the block", http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, "O bloco (%s, %s, %s) foi cadastrado com sucesso.", block.ID,
			block.Name, block.Timestamp.String())
	}
}
