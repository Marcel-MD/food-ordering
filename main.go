package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/Marcel-MD/food-ordering/domain"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config()
	domain.SetConfig(cfg)

	om := domain.OrderManager{
		Orders: make(map[int64]domain.Order),
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var restaurant domain.Restaurant
		err := json.NewDecoder(r.Body).Decode(&restaurant)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		domain.Restaurants[restaurant.RestaurantId] = restaurant

		log.Info().Str("restaurant", restaurant.Name).Msg("Registered restaurant")

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {

		restaurants := make([]domain.Restaurant, 0, len(domain.Restaurants))
		for _, value := range domain.Restaurants {
			restaurants = append(restaurants, value)
		}

		menu := domain.Menu{
			Restaurants:     len(restaurants),
			RestaurantsData: restaurants,
		}

		log.Debug().Msg("Sending menu")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(menu)
	}).Methods("GET")

	r.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Debug().Int("client_id", order.ClientId).Msg("Received order")

		orderResponse := om.ManageOrder(order)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderResponse)
	}).Methods("POST")

	http.ListenAndServe(":"+cfg.FoodOrderingPort, r)
}

func config() domain.Config {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	file, err := os.Open("config/cfg.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening menu.json")
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var cfg domain.Config
	json.Unmarshal(byteValue, &cfg)

	return cfg
}
