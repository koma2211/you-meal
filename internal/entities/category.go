package entities

type Menu struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID    string `json:"categoryId" binding:"required"`
	Title string `json:"title" binding:"required"`
}

type MealCategory struct {
	ID    int    `json:"id,omitempty"`
	Meals []Meal `json:"meals,omitempty"`
}

type Meal struct {
	ID          int          `json:"id,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Weight      int          `json:"weight,omitempty"`
	Calorie     int          `json:"calorie,omitempty"`
	Price       float64      `json:"price,omitempty"`
	ImagePath   string       `json:"imagePath,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

type Ingredient struct {
	Id    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}
