package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/http"
  "time"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
)

type Place struct {
	Id   string
	Name string `json:"name"`
	Ctry string `json:"country"`
	Desc string `json:"description"`
	Lat  string `json:"latitude"`
	Lon  string `json:"longitude"`
}

func main() {
	router := SetupRouter()
	router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		v1.POST("/place", CreatePlace)
		v1.GET("/place/:id", GetPlace)
		v1.GET("/places", GetAllPlaces)
		v1.DELETE("/place", DeletePlace)
		v1.GET("/health", HealthCheck)
		v1.GET("/url/:id", GetPlaceUrl)
	}
	return router
}

func CreatePlace(c *gin.Context) {
	var newPlace Place
	c.Bind(&newPlace)

	stmt, err := Init().Prepare("INSERT INTO places (name, country, description, latitude, longitude) VALUES (?,?,?,?,?);")

	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(newPlace.Name, newPlace.Ctry, newPlace.Desc, newPlace.Lat, newPlace.Lon)

	if err != nil {
		fmt.Print(err.Error())
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(" %s successfully created.", newPlace.Name),
	})
}

func GetAllPlaces(c *gin.Context) {
	var (
		place  Place
		places []Place
	)
	rows, err := Init().Query("SELECT * FROM places;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&place.Id, &place.Name, &place.Ctry, &place.Desc, &place.Lat, &place.Lon)
		places = append(places, place)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

  //  Properly close resultset
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, gin.H{
		"result": places,
	})
}

func GetPlace(c *gin.Context) {
	var (
		place  Place
		result gin.H
	)
	id := c.Param("id")
	err := Init().QueryRow("SELECT * FROM places WHERE id = ?;", id).Scan(&place.Id, &place.Name, &place.Ctry, &place.Desc, &place.Lat, &place.Lon)

	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": place,
		}
	}

	c.IndentedJSON(http.StatusOK, result)
}

func DeletePlace(c *gin.Context) {
	id := c.Query("id")
	stmt, err := Init().Prepare("DELETE FROM places WHERE id= ?;")

	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted place: %s", id),
	})
}

func GetPlaceUrl(c *gin.Context) {
	var (
		place  Place
		result gin.H
	)
	id := c.Param("id")
	err := Init().QueryRow("SELECT * FROM places WHERE id = ?;", id).Scan(&place.Id, &place.Name, &place.Ctry, &place.Desc, &place.Lat, &place.Lon)

	if err != nil {
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": place,
		}
	}

  // Google Maps format
	//Url := "https://maps.google.com/maps?q=" + place.Lat + "," + place.Lon + "&t=k"

  // Openstreetmap format
  Url := "https://www.openstreetmap.org/#map=15/" + place.Lat + "/" + place.Lon

	result = gin.H{
		"result": Url,
	}

	c.PureJSON(http.StatusOK, result)
}

func HealthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"health": "passed"})
}

func Init() *sql.DB {
	mysqlHost, provided := os.LookupEnv("MARIADB_HOST")
	if !provided {
		mysqlHost = "localhost"
	}
	mysqlUser, provided := os.LookupEnv("MARIADB_USER")
	if !provided {
		log.Fatalf("Environment variable %s is not set", "MARIADB_USER")
	}
	mysqlPass, provided := os.LookupEnv("MARIADB_PASS")
	if !provided {
		log.Fatalf("Environment variable %s is not set", "MARIADB_PASS")
	}
	config := mysql.Config{
		User:                 mysqlUser,
		Passwd:               mysqlPass,
		Net:                  "tcp",
		Addr:                 mysqlHost,
		AllowNativePasswords: true,
		DBName:               "placesdb",
	}
	db, err := sql.Open("mysql", config.FormatDSN())

	checkErr(err)

	err = db.Ping()
	checkErr(err)
	fmt.Printf("DB Connection successful.")

  // Set the maximum number of concurrently open connections (in-use + idle)
  // to 5. Setting this to less than or equal to 0 will mean there is no
  // maximum limit (which is also the default setting).
//  db.SetMaxOpenConns(5)

  // Set the maximum number of concurrently idle connections to 5. Setting this
  // to less than or equal to 0 will mean that no idle connections are retained.
//  db.SetMaxIdleConns(5)

  // Set the maximum lifetime of a connection to 1 hour. Setting it to 0
  // means that there is no maximum lifetime and the connection is reused
  // forever (which is the default behavior)
//  db.SetConnMaxLifetime(5*time.Minute)

	return db
}

func checkErr(err error) {
	if err != nil {
		fmt.Print(err.Error())
	}
}
