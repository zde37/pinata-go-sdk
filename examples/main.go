package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zde37/pinata-go-sdk/pinata"
)

func main() {
	auth := pinata.NewAuthWithJWT(os.Getenv("PINATA_JWT"))

	pinataClient := NewPinataClient(auth)
	pinataClient.TestAuthentication()
}
