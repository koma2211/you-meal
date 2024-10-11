package entities

type ClientInfo struct {
	Name        string   `json:"name" binding:"required"`
	PhoneNumber string   `json:"phoneNumber" binding:"required"`
	IsDelivery  bool     `json:"isDelivery" binding:"required"`
	Address     string   `json:"address,omitempty"`
	Floor       int      `json:"floor,omitempty"`
	Burgers     []Orders `json:"burgers,omitempty"`
	Snacks      []Orders `json:"snacks,omitempty"`
}

type Orders struct {
	ID       int `json:"id" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}
