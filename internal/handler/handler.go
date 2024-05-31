package handler

import (
	"encoding/json"
	"log"
	"multithreaded-data-processor/internal/entities"
	"multithreaded-data-processor/internal/resultData"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Build(router *chi.Mux, store *resultData.ResultDataStorage) {
	router.Use(middleware.Recoverer)

	controller := NewController(store)

	router.Get("/", controller.GetData)

}

type Controller struct {
	storage *resultData.ResultDataStorage
}

func NewController(storage *resultData.ResultDataStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (c *Controller) GetData(w http.ResponseWriter, r *http.Request) {
	var result entities.ResultT
	resultSetT, err := c.storage.GetResultData()
	if err != nil {
		result.Status = false
		result.Error = "Error on collect data"
	} else {
		checkFull := c.storage.IsFull()
		switch checkFull {
		case true:
			result.Status = true
			result.Data = resultSetT
		case false:
			result.Status = false
			result.Error = "Error on collect data"
		}
	}

	if result.Error != "" {
		log.Printf("Error: %s", result.Error)
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error converting ResulT to json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(res)
}
