package entities

type Burger struct {
	ID          int     `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Weight      int     `json:"weight,omitempty"`
	Calorie     int     `json:"calorie,omitempty"`
	Price       float64 `json:"price,omitempty"`
	ImagePath   string  `json:"imagePath"`
}