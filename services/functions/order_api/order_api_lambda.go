package main

import "github.com/aws/aws-lambda-go/lambda"

func hello() (string, error) {
	return "hello λ!, I'm the order API", nil
}

func main() {
	lambda.Start(hello)
}
