package swaggerservicehandler

import (
	"fmt"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	opn "telehealers.in/router/restapi/operations"
)

var (
	ImageNotFound = "Image not found"
)

func GetProfilePictures(params opn.GetProfilePicturesNameParams, p *models.Principal) middleware.Responder {
	pngFile, fileReadErr := os.Open("example.png")
	if fileReadErr != nil {
		fmt.Errorf("Error:%v", fileReadErr)
		return opn.NewGetProfilePicturesNameDefault(400).WithPayload(
			&models.Error{Message: &ImageNotFound})
	}
	return opn.NewGetProfilePicturesNameOK().WithPayload(pngFile)
}
