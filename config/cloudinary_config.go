package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cld *cloudinary.Cloudinary

func SetupCloudinary() {
	var err error
	Cld, err = cloudinary.NewFromURL(
		"cloudinary://" +
			os.Getenv("CLOUDINARY_API_KEY") + ":" +
			os.Getenv("CLOUDINARY_API_SECRET") + "@" +
			os.Getenv("CLOUDINARY_CLOUD_NAME"),
	)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
}
