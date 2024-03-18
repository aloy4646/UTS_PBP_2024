package main

import (
	"fmt"
	"log"
	"net/http"
	"uts/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rooms/{id_game}", controller.GetAllRoomsByGame).Methods("GET")
	router.HandleFunc("/rooms/detail/{id}", controller.GetDetailRoom).Methods("GET")
	router.HandleFunc("/rooms", controller.InsertRoom).Methods("POST")
	router.HandleFunc("/rooms", controller.LeaveRoom).Methods("DELETE")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://ithb.ac.id"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
