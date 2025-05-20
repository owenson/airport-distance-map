package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Airport struct {
	ID           string   `json:"id"`
	Ident        string   `json:"ident"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	LatitudeDeg  *float64 `json:"latitude_deg"`
	LongitudeDeg *float64 `json:"longitude_deg"`
	ElevationFt  *float64 `json:"elevation_ft"`
	Continent    *string  `json:"continent"`
	IsoCountry   *string  `json:"iso_country"`
	IsoRegion    *string  `json:"iso_region"`
	Municipality *string  `json:"municipality"`
	IcaoCode     *string  `json:"icao_code"`
	IataCode     *string  `json:"iata_code"`
}

func main() {
	// Initialize database connection
	db, err := sql.Open("sqlite3", "./airports.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static files
	r.Static("/static", "./static")

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Search endpoint
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
			return
		}

		rows, err := db.Query(`
			SELECT id, ident, type, name, latitude_deg, longitude_deg, elevation_ft,
				   continent, iso_country, iso_region, municipality, icao_code, iata_code
			FROM airports
			WHERE name LIKE ? OR ident LIKE ?
			LIMIT 50
		`, "%"+query+"%", "%"+query+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var airports []Airport
		for rows.Next() {
			var a Airport
			err := rows.Scan(
				&a.ID, &a.Ident, &a.Type, &a.Name, &a.LatitudeDeg, &a.LongitudeDeg,
				&a.ElevationFt, &a.Continent, &a.IsoCountry, &a.IsoRegion,
				&a.Municipality, &a.IcaoCode, &a.IataCode,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			airports = append(airports, a)
		}

		c.JSON(http.StatusOK, airports)
	})

	// Start server
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
} 