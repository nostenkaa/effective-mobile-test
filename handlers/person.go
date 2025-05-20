package handlers

import (
	"effective-mobile-test/config"
	"effective-mobile-test/logger"
	"effective-mobile-test/models"
	"effective-mobile-test/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// RegisterRoutes регистрирует маршруты для работы с сущностью Person
// @Summary Register people routes
// @Description Register CRUD routes for Person entity
func RegisterRoutes(r *gin.Engine) {
	group := r.Group("/people")
	{
		group.POST("/", createPerson)
		group.GET("/", listPeople)
		group.PUT("/:id", updatePerson)
		group.DELETE("/:id", deletePerson)
	}
}

// createPerson создает нового человека
// @Summary Create a new person
// @Description Create a person with name, surname and optionally patronymic
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.Person true "Person data"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /people/ [post]
func createPerson(c *gin.Context) {
	log := logger.L().WithField("endpoint", "createPerson")
	log.Debug("Вызван endpoint createPerson")

	var input struct {
		Name       string `json:"name" binding:"required"`
		Surname    string `json:"surname" binding:"required"`
		Patronymic string `json:"patronymic"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Debugf("Ошибка биндинга JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Создание нового пользователя: %s %s", input.Name, input.Surname)

	enriched := services.EnrichPerson(input.Name)
	enriched.Name = input.Name
	enriched.Surname = input.Surname
	enriched.Patronymic = input.Patronymic

	if err := config.DB.Create(&enriched).Error; err != nil {
		log.Errorf("Ошибка при сохранении в БД: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить"})
		return
	}

	log.Infof("Пользователь успешно создан с ID %d", enriched.ID)
	c.JSON(http.StatusOK, enriched)
}

// listPeople возвращает список людей с фильтрацией и пагинацией
// @Summary List people
// @Description Get list of people, optionally filtered by name and surname, with pagination support
// @Tags people
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.Person
// @Failure 500 {object} map[string]string "Internal error"
// @Router /people/ [get]
func listPeople(c *gin.Context) {
	log := logger.L().WithField("endpoint", "listPeople")
	log.Debug("Вызван endpoint listPeople")

	var people []models.Person
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query := config.DB.Limit(limit).Offset(offset)

	if name := c.Query("name"); name != "" {
		log.WithField("name", name).Debug("Применён фильтр по имени")
		query = query.Where("name = ?", name)
	}
	if surname := c.Query("surname"); surname != "" {
		log.WithField("surname", surname).Debug("Применён фильтр по фамилии")
		query = query.Where("surname = ?", surname)
	}

	if err := query.Find(&people).Error; err != nil {
		log.WithError(err).Error("Не удалось получить список пользователей")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
		return
	}

	log.WithFields(logrus.Fields{
		"count":  len(people),
		"limit":  limit,
		"offset": offset,
	}).Debug("Пользователи успешно получены")

	c.JSON(http.StatusOK, people)
}

// updatePerson обновляет данные существующего человека по ID
// @Summary Update a person
// @Description Update person data by ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body models.Person true "Updated person data"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /people/{id} [put]
func updatePerson(c *gin.Context) {
	log := logger.L().WithField("endpoint", "updatePerson")
	log.Debug("Вызван endpoint updatePerson")

	id := c.Param("id")
	var person models.Person
	if err := config.DB.First(&person, id).Error; err != nil {
		log.WithField("id", id).Warn("Пользователь не найден")
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	var input models.Person
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithError(err).Warn("Невалидный входной JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.WithFields(logrus.Fields{
		"id":   id,
		"data": input,
	}).Debug("Данные для обновления получены")

	if err := config.DB.Model(&person).Updates(input).Error; err != nil {
		log.WithError(err).Error("Не удалось обновить пользователя")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}

	log.WithField("id", id).Info("Пользователь успешно обновлён")
	c.JSON(http.StatusOK, person)
}

// deletePerson удаляет человека по ID
// @Summary Delete a person
// @Description Delete person by ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /people/{id} [delete]
func deletePerson(c *gin.Context) {
	log := logger.L().WithField("endpoint", "deletePerson")
	log.Debug("Вызван endpoint deletePerson")

	id := c.Param("id")

	result := config.DB.Delete(&models.Person{}, id)
	if result.Error != nil {
		log.WithError(result.Error).WithField("id", id).Error("Ошибка при удалении пользователя")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
		return
	}

	if result.RowsAffected == 0 {
		log.WithField("id", id).Warn("Пользователь для удаления не найден")
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	log.WithField("id", id).Info("Пользователь успешно удалён")
	c.Status(http.StatusNoContent)
}
