package main

import (
	// "fmt"
	// character "golang-Serveless-characters/pkg/char"
	"golang-Serveless-characters/pkg/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// fmt.Println(character.GetCharInfo("lucy+soul"))
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req)
	default:
		return handlers.UnhandledMethod()
	}
}
