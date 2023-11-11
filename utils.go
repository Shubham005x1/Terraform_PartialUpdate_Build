package content

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
)

type Employee struct {
	ID        string `firestore:"id" json:"id"`
	FirstName string `firestore:"firstname" json:"firstname"`
	LastName  string `firestore:"lastname" json:"lastname"`
	Email     string `firestore:"email" json:"email"`
	Password  string `firestore:"password" json:"password"`
	PhoneNo   string `firestore:"phoneNo" json:"phoneNo"`
	Role      string `firestore:"role" json:"role"`
}

var (
	client     *firestore.Client
	logClient  *logging.Client
	onceClient sync.Once
)

func initializeFirestore() {
	onceClient.Do(func() {
		ctx := context.Background()

		// Initialize Firestore with the service account key
		var err error
		client, err = firestore.NewClient(ctx, "takeoff-task-3")
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	})
}

// Function to merge two maps
func mergeMap(destination, source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		if _, ok := destination[key]; !ok {
			// Key doesn't exist in destination map, you can handle it here
			fmt.Printf("Key %s doesn't exist in destination map\n", key)
			continue // Skip this key
		}
		destination[key] = value
	}
	return destination
}
