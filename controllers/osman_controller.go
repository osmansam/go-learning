package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/models"
	"github.com/osmansam/react/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var osmanCollection *mongo.Collection = configs.GetCollection(configs.DB, "osmans")
func CreateOsman(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var osman models.Osman
	defer cancel()
	//parse the body into the osman struct
	if err := c.BodyParser(&osman); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OsmanResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//use the validator library to validate required fields
	if validationErr := validate.Struct(&osman); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OsmanResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newOsman := models.Osman{
		Name: osman.Name,
		Helal: osman.Helal,
	}

	result, err := osmanCollection.InsertOne(ctx, newOsman)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})

	}

	return c.Status(http.StatusCreated).JSON(responses.OsmanResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

//get all osman
func GetAllOsman(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var osmans []models.Osman
	defer cancel()

	results, err := osmanCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleOsman models.Osman
		if err = results.Decode(&singleOsman); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
    // Fetch the Helal document for this Osman
    var helal models.Helal
    helalErr := helalCollection.FindOne(ctx, bson.M{"_id": singleOsman.Helal}).Decode(&helal)

    if helalErr != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": helalErr.Error()}})
    }

    // Assign the Helal details to the Osman struct (you can adjust this as needed)
    singleOsman.HelalDetails = helal

		osmans = append(osmans, singleOsman)
	}

	return c.Status(http.StatusOK).JSON(
		responses.OsmanResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": osmans}},
	)
}