package entities

type Client struct {
	Name        string         `json:"name" binding:"required"`
	PhoneNumber string         `json:"phoneNumber" binding:"required"`
	IsDelivery  *bool          `json:"isDelivery" binding:"required"`
	Address     *string        `json:"address,omitempty"`
	Floor       *int           `json:"floor,omitempty"`
	Orders      []OrderedMeals `json:"orders" binding:"required"`
}

type OrderedMeals struct {
	ID       int `json:"id" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}
