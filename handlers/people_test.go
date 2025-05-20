package handlers_test

import (
	"bytes"
	"effective-mobile-test/config"
	"effective-mobile-test/handlers"
	"effective-mobile-test/logger"
	"effective-mobile-test/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "modernc.org/sqlite" // вот этот драйвер вместо mattn/go-sqlite3
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open gorm DB: %v", err)
	}

	if err := db.AutoMigrate(&models.Person{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	config.DB = db
	return db
}

func setupRouter() *gin.Engine {
	logger.Init()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handlers.RegisterRoutes(r)
	return r
}
func TestCreatePerson(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	person := map[string]string{
		"name":    "Ivan",
		"surname": "Petrov",
	}
	body, err := json.Marshal(person)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/people/", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListPeople(t *testing.T) {
	db := setupTestDB()
	db.Create(&models.Person{Name: "Ivan", Surname: "Petrov", Nationality: "RU"})

	router := setupRouter()
	req, err := http.NewRequest("GET", "/people/?name=Ivan", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Ivan")
}

func TestUpdatePerson(t *testing.T) {
	db := setupTestDB()
	person := models.Person{Name: "Old", Surname: "Name"}
	db.Create(&person)

	router := setupRouter()

	update := map[string]string{"name": "New"}
	body, err := json.Marshal(update)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/people/"+strconv.Itoa(int(person.ID)), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New")
}

func TestDeletePerson(t *testing.T) {
	db := setupTestDB()
	person := models.Person{Name: "Delete", Surname: "Me"}
	db.Create(&person)

	router := setupRouter()

	req, err := http.NewRequest("DELETE", "/people/"+strconv.Itoa(int(person.ID)), nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
