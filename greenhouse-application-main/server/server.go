package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kaua-matheus/greenhouse-application/database"
)

func main() {
	Serve()
}

func Serve() {
	router := gin.Default()
	allowedOrigins := []string{
        "*", 
    }
    
    var validOrigins []string
    for _, origin := range allowedOrigins {
        if origin != "" {
            validOrigins = append(validOrigins, origin)
        }
    }
    
    router.Use(cors.New(cors.Config{
        AllowOrigins:     validOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
    }))

	connection, err := database.NewConnection()
	if err != nil {
		fmt.Printf("An error occurs trying to create a connection: %s", err)
		return
	}

	// Parameters
	// Get
	router.GET("/parameters", func(c *gin.Context) {
		data, err := database.GetAllParameters(connection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "An error occurs trying to get all the parameters",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
		})
	})

	// Put
	router.PUT("/parameter/:id", func(c *gin.Context) {
		parameter := map[string]interface{}{}

		id_str := c.Param("id")
		id, err := strconv.Atoi(id_str)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Error: Invalid UINT format",
			})
			return
		}

		if err := c.BindJSON(&parameter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid JSON payload",
			})
			return
		}

		err = database.UpdateParameter(connection, uint(id), parameter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error: Couldn't update the parameter",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   "Parameter updated successfully",
		})

	})

	// Data
	// Get
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "A API está funcionando.",
		})
	})

	router.GET("/data", func(c *gin.Context) {
		data, err := database.GetAllData(connection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "An error occurs trying to get all the data",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
		})
	})

	// Post
	router.POST("/data", func(c *gin.Context) {
		data := database.GlpData{}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid JSON payload",
			})
			return
		}

		database.AddData(connection, data)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   "Data registered successfully",
		})
	})

	// Put
	router.PUT("/data/:id", func(c *gin.Context) {
		data := database.GlpData{}

		id_str := c.Param("id")
		id, err := uuid.Parse(id_str)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Error: Invalid UUID format",
			})
			return
		}

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid JSON payload",
			})
			return
		}

		err = database.UpdateData(connection, id, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error: Couldn't update the data",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   "Data updated successfully",
		})

	})

	// Delete
	router.DELETE("/data/:id", func(c *gin.Context) {
		id_str := c.Param("id")
		id, err := uuid.Parse(id_str)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Error: Invalid UUID format",
			})
			return
		}

		if err := database.DeleteData(connection, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error: Couldn't delete the data",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   "Data deleted successfully",
		})
	})

	router.Run("0.0.0.0:8080")
}
