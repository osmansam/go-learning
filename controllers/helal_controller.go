package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/models"
	"github.com/osmansam/react/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var helalCollection *mongo.Collection = configs.GetCollection(configs.DB, "helals")
var validate = validator.New()
// create a helal
func CreateHelal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var helal models.Helal
	defer cancel()
	//parse the body into the helal struct
	if err := c.BodyParser(&helal); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//use the validator library to validate required fields
	if validationErr := validate.Struct(&helal); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newHelal := models.Helal{
		Name: helal.Name,
		Surname: helal.Surname,
		Age: helal.Age,
	}

	result, err := helalCollection.InsertOne(ctx, newHelal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})

	}

	return c.Status(http.StatusCreated).JSON(responses.HelalResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
//get all helal
func GetAllHelal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var helals []models.Helal
	defer cancel()

	results, err := helalCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleHelal models.Helal
		if err = results.Decode(&singleHelal); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		helals = append(helals, singleHelal)
	}

	return c.Status(http.StatusOK).JSON(
		responses.HelalResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": helals}},
	)
}

// get a helal
func GetHelal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var helal models.Helal
	helalIdStr := c.Params("id")

	// Convert helalIdStr to an ObjectID
	helalId, err := primitive.ObjectIDFromHex(helalIdStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "Invalid ID format", Data: &fiber.Map{"data": err.Error()}})
	}
	if err := helalCollection.FindOne(ctx, bson.M{"_id": helalId}).Decode(&helal); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.HelalResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": helal}})
}
//delete helal
func DeleteHelal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helalIdStr := c.Params("id")
	helalId,err:=primitive.ObjectIDFromHex(helalIdStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "Invalid ID format", Data: &fiber.Map{"data": err.Error()}})
	}
	result,err:=helalCollection.DeleteOne(ctx,bson.M{"_id":helalId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.HelalResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})
}

//update helal
func UpdateHelal(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	helalIdStr := c.Params("id")
	helalId, err := primitive.ObjectIDFromHex(helalIdStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "Invalid ID format", Data: &fiber.Map{"data": err.Error()}})
	}

	var updates bson.M
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.HelalResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// (Optional) Perform validation on updates here...

	update := bson.M{
		"$set": updates,
	}
	result, err := helalCollection.UpdateOne(ctx, bson.M{"_id": helalId}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// If no document was updated
	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.HelalResponse{Status: http.StatusNotFound, Message: "Helal not found", Data: nil})
	}

	var updatedHelal models.Helal
	if err := helalCollection.FindOne(ctx, bson.M{"_id": helalId}).Decode(&updatedHelal); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.HelalResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedHelal}})
}
//handle search
func SearchHelal(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var helals []models.Helal
    searchKey := c.Query("searchKey")

    // Define the regex pattern to match anywhere in the string
    pattern := ".*" + searchKey + ".*"
    regex := primitive.Regex{Pattern: pattern, Options: "i"} // "i" is for case-insensitive

    // Build the query filter
    filter := bson.M{"$or": []bson.M{
        {"name": regex},
        {"surname": regex},
    }}

    results, err := helalCollection.Find(ctx, filter)

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    // Reading from the db
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleHelal models.Helal
        if err = results.Decode(&singleHelal); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.HelalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }

        helals = append(helals, singleHelal)
    }

    return c.Status(http.StatusOK).JSON(responses.HelalResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": helals}})
}

