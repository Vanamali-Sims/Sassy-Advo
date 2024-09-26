package main

import (
	"html/template"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Define sarcastic responses with sass levels
var sarcasticResponses = map[string]map[string][]string{
	"coffee": {
		"mild": {
			"Coffee's okay, I guess. Just keep it reasonable.",
		},
		"medium": {
			"Coffee? You're basically a walking espresso shot.",
		},
		"extra spicy": {
			"Another cup of coffee? Why stop when you can levitate?",
		},
	},
	"exercise": {
		"mild": {
			"A little exercise never hurt anyone... I think.",
		},
		"medium": {
			"Exercise? Sure, walking to the fridge counts.",
		},
		"extra spicy": {
			"No pain, no gain, but isn't sitting pain-free?",
		},
	},
	"dessert": {
		"mild": {
			"Dessert for breakfast? Maybe once in a while...",
		},
		"medium": {
			"Dessert for breakfast? You’re living your best life.",
		},
		"extra spicy": {
			"Sugar in the morning? Why not go all-in with some cake?",
		},
	},
}

// Struct to hold dynamic values
type AdviceData struct {
	Category  string
	SassLevel string
	Response  string
}

func getSarcasticResponse(category, sassLevel string) string {
	responses := sarcasticResponses[category][sassLevel]
	rand.Seed(time.Now().UnixNano()) // Seed to ensure random responses
	return responses[rand.Intn(len(responses))]
}

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Health advice endpoint by category and sass level
	app.Get("/:category", func(c *fiber.Ctx) error {
		category := c.Params("category")
		sassLevel := c.Query("sass", "medium") // Default sass level is "medium"

		// Validate if category and sass level exist
		if _, categoryExists := sarcasticResponses[category]; categoryExists {
			if _, sassExists := sarcasticResponses[category][sassLevel]; sassExists {
				response := getSarcasticResponse(category, sassLevel)

				// Load the HTML template file
				tmpl, err := template.ParseFiles("../templates/index.html")
				if err != nil {
					return c.Status(500).SendString("Template loading error!")
				}

				// Prepare the data to inject into the template
				data := AdviceData{
					Category:  category,
					SassLevel: sassLevel,
					Response:  response,
				}

				// Set the content type to HTML and render the template
				c.Type("html")
				err = tmpl.Execute(c.Response().BodyWriter(), data)
				if err != nil {
					return c.Status(500).SendString("Template execution error!")
				}
				return nil
			}
			// If sass level is invalid, return error
			return c.Status(400).SendString("Invalid sass level. Try 'mild', 'medium', or 'extra spicy'.")
		}

		// Return a 404 error if category not found
		return c.Status(404).SendString("No sarcasm found for that category.")
	})

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}
