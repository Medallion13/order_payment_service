package main

import "github.com/aws/aws-lambda-go/lambda"

func hello() (string, error) {
	return "hello λ!, im the payment event", nil
}

func main() {
	lambda.Start(hello)
}
