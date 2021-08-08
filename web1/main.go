package main

import (
	"github.com/iitheogotut/web/web1/models"
	"github.com/iitheogotut/web/web1/routes"
	"github.com/iitheogotut/web/web1/utils"
	"log"
	"net/http"
)

func main(){
	models.Init()
	utils.LoadTemplates("templates/*.html")
	router := routes.NewRouter()
	http.Handle("/", router)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}


