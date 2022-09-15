# Invoke a go lambda function to create a dynamodb table

This example shows how to invoke a lambda function through apigateway to create a dynamodb table

## Requirements

You should have installed:

- Docker.
- Docker-compose.
- Npm.

## Prepare environment

Install Dependencies
```bash
make deps
```

## Run

In order to run the example, simply run:
```bash
make run
```
## Usage

The Serverless engine will provide an endpoint for making an HTTP request, in the following format:
```bash
http://localhost:4566/restapis/<id>/local/_user_request
```
Add the path configured in the `serverless.yml` file to the URL and make an HTTP request using cURL, Postman, etc.

For Example:

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"table": "my-table-name"}' \
    http://localhost:4566/restapis/a9dd16qva6/local/_user_request/table
```
