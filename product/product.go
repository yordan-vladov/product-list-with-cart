package product

type Product struct {
	Id       int
	Quantity int
	Image    struct {
		Thumbnail string `json:"thumbnail"`
		Mobile    string `json:"mobile"`
		Tablet    string `json:"tablet"`
		Desktop   string `json:"Desktop"`
	} `json:"image"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float32 `json:"price"`
}
