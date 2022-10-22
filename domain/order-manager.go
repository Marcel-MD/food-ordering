package domain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

var orderId int64

type OrderManager struct {
	Orders map[int64]Order
}

func (om *OrderManager) ManageOrder(order Order) OrderResponse {
	order.OrderId = atomic.AddInt64(&orderId, 1)
	om.Orders[order.OrderId] = order

	response := OrderResponse{
		OrderId: order.OrderId,
		Orders:  make([]OrderResponseData, 0, len(order.Orders)),
	}

	for _, orderData := range order.Orders {
		restaurant, ok := Restaurants[orderData.RestaurantId]
		if !ok {
			continue
		}

		jsonBody, err := json.Marshal(orderData)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling order")
		}
		contentType := "application/json"

		r, err := http.Post(restaurant.Address+"/v2/order", contentType, bytes.NewReader(jsonBody))
		if err != nil {
			log.Fatal().Err(err).Msg("Error sending order to restaurant")
		}

		var orderResponse OrderResponseData
		err = json.NewDecoder(r.Body).Decode(&orderResponse)
		if err != nil {
			log.Fatal().Err(err).Msg("Error decoding order response")
		}

		orderResponse.RestaurantAddress = restaurant.Address
		response.Orders = append(response.Orders, orderResponse)
		log.Debug().Int64("order_id", order.OrderId).Str("restaurant", restaurant.Name).Msg("Order sent to restaurant")
	}

	log.Info().Int64("order_id", order.OrderId).Msg("Order sent to all restaurants")
	return response
}

func (om *OrderManager) ManageRating(rating Rating) {

	for _, orderRating := range rating.Orders {

		restaurant, ok := Restaurants[orderRating.RestaurantId]
		if !ok {
			continue
		}

		jsonBody, err := json.Marshal(orderRating)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling rating")
		}
		contentType := "application/json"

		r, err := http.Post(restaurant.Address+"/v2/rating", contentType, bytes.NewReader(jsonBody))
		if err != nil {
			log.Fatal().Err(err).Msg("Error sending rating to restaurant")
		}

		var ratingResponse RatingResponse
		err = json.NewDecoder(r.Body).Decode(&ratingResponse)
		if err != nil {
			log.Fatal().Err(err).Msg("Error decoding rating response")
		}

		restaurant.Rating = ratingResponse.RestaurantAvgRating
		Restaurants[ratingResponse.RestaurantId] = restaurant
	}

	totalRating := 0.0
	for _, restaurant := range Restaurants {
		totalRating += restaurant.Rating
	}

	log.Info().Float64("avg_rating", totalRating/float64(len(Restaurants))).Msg("Average simulation rating")
}
