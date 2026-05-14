package opendart_test

import (
	"context"
	"log"

	"github.com/ev3rlit/opendart"
)

func ExampleNew() {
	client, err := opendart.New(opendart.Config{
		APIKey: "test-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = context.Background()
	_ = client
}
