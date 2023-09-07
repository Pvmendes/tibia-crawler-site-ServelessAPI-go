package handlers

import (
	character "golang-Serveless-characters/pkg/char"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	name := req.QueryStringParameters["charName"]

	if len(name) > 0 {

		result, err := character.GetCharInfo(name)

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return apiResponse(http.StatusOK, result)

	} else {

		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String("We need name of the char"),
		})

	}
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
