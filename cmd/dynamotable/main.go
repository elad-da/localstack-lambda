package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Input struct {
	Table string `json:"table"`
}

type Response events.APIGatewayProxyResponse

var builtinResponse = Response{
	StatusCode:      200,
	IsBase64Encoded: false,
	Headers: map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Cache-Control":                    "no-cache; no-store",
		"Content-Type":                     "application/json",
		"Content-Security-Policy":          "default-src self",
		"Strict-Transport-Security":        "max-age=31536000; includeSubDomains",
		"X-Content-Type-Options":           "nosniff",
		"X-XSS-Protection":                 "1; mode=block",
		"X-Frame-Options":                  "DENY",
	},
}

type handler struct {
}

func (h handler) Run(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	defaultTable := "my-table"

	if req.Body != "" {
		var input Input
		_ = json.Unmarshal([]byte(req.Body), &input)
		if input.Table != "" {
			defaultTable = input.Table
		}
	}

	fmt.Printf("Table name: %s\n", defaultTable)
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://" + os.Getenv("LOCALSTACK_HOSTNAME") + ":4566",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		builtinResponse.StatusCode = 500
		builtinResponse.Body = err.Error()
		return builtinResponse, nil
	}

	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(defaultTable),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		builtinResponse.StatusCode = 500
		builtinResponse.Body = err.Error()
		return builtinResponse, nil
	}

	builtinResponse.Body = *out.TableDescription.TableName
	return builtinResponse, nil
}

func New() *handler {
	return &handler{}
}

func main() {
	handler := New()
	lambda.Start(handler.Run)
}
