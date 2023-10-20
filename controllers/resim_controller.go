package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/models"
	"github.com/osmansam/react/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var resimCollection *mongo.Collection = configs.GetCollection(configs.DB, "resims")

func CreateResim(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Handle file upload
	file, err := c.FormFile("resim")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OsmanResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Save the temporary file to the disk
	tempFilePath := "./temp/" + file.Filename
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error temp file path ", Data: &fiber.Map{"data": err.Error()}})
	}

	// Upload to Cloudinary
	imageURL, err := UploadToCloudinary(tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error uploading to cloudinary", Data: &fiber.Map{"data": err.Error()}})
	}

	// Clean up temp file
	os.Remove(tempFilePath)

	// Create the Resim model with the image URL
	resim := models.Resim{
		Adi:   c.FormValue("adi"),
		Resim: imageURL,
	}

	// Insert into MongoDB
	result, err := resimCollection.InsertOne(ctx, resim)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OsmanResponse{Status: http.StatusInternalServerError, Message: "error inserting to database", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.OsmanResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

// get all resim
func GetAllResim(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var resims []models.Resim
	defer cancel()

	results, err := resimCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ResimResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleResim models.Resim
		if err = results.Decode(&singleResim); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ResimResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		resims = append(resims, singleResim)
	}

	return c.Status(http.StatusOK).JSON(
		responses.ResimResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": resims}},
	)
}