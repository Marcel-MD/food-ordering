package domain

type Config struct {
	TimeUnit         int    `json:"time_unit"`
	FoodOrderingPort string `json:"food_ordering_port"`
}

var cfg Config = Config{
	TimeUnit:         250,
	FoodOrderingPort: "8090",
}

func SetConfig(c Config) {
	cfg = c
}
