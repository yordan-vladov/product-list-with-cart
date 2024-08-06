package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	. "product-list/product"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func Products() ([]Product, error) {
	jsonFile, err := os.Open("data.json")

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var products []Product

	json.Unmarshal(byteValue, &products)

	return products, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func cartItems(products []Product) []Product {
	var filteredProducts []Product
	for _, product := range products {
		if product.Quantity > 0 {
			filteredProducts = append(filteredProducts, product)
		}
	}
	return filteredProducts
}
func main() {
	funcMap := template.FuncMap{
		"mul": func(a float32, b int) float32 {
			return a * float32(b)
		},
		"cartItems": cartItems,
		"totalPayment": func(products []Product) float32 {
			sum := float32(0.0)
			for _, product := range products {
				sum += product.Price * float32(product.Quantity)
			}

			return sum
		},
		"orderedItemsQuantity": func(products []Product) int {
			totalQuantity := 0
			for _, product := range products {
				totalQuantity += product.Quantity
			}

			return totalQuantity
		},
	}
	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("public/views/*.html")),
	}

	products, err := Products()

	for i := 0; i < len(products); i++ {
		products[i].Id = i
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", products)
	})

	e.GET("/cart", func(c echo.Context) error {
		return c.Render(http.StatusOK, "cart", cartItems(products))
	})

	e.POST("/decrement/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			e.Logger.Panic("Invalid ID")
		}

		if products[id].Quantity > 0 {
			products[id].Quantity--
		}

		return c.Render(http.StatusOK, "product", products[id])
	})
	e.POST("/increment/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Panic("Invalid ID")
		}

		products[id].Quantity++

		return c.Render(http.StatusOK, "product", products[id])
	})
	e.DELETE("/remove-cart-item/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			e.Logger.Panic("Invalid ID")
		}

		products[id].Quantity = 0

		return c.Render(http.StatusOK, "cart", cartItems(products))
		//return c.NoContent(http.StatusOK)
	})

	e.POST("/confirmation-message", func(c echo.Context) error {

		defer clearProducts(products)

		return c.Render(http.StatusOK, "confirmation-message", cartItems(products))
	})

	e.Static("images", "public/images")
	e.Static("styles", "public/styles")
	e.Static("scripts", "public/scripts")

	e.Logger.Fatal(e.Start(":1323"))

}

func clearProducts(products []Product) {
	for i := 0; i < len(products); i++ {
		products[i].Quantity = 0
	}
}
