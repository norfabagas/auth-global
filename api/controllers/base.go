package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DBDriver, DBUser, DBPassword, DBPort, DBHost, DBName string) {
	var err error

	if DBDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)
		server.DB, err = gorm.Open(DBDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", DBName)
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("Connected to database %s\n", DBName)
		}
	}

	server.Router = mux.NewRouter()

	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to address: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
