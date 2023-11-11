package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/Shubham005x1/MyValidations/validations"
)

// func main() {
// 	router := mux.NewRouter()

//		router.HandleFunc("/employees", PartialUpdateEmployee).Methods("PATCH")
//		log.Println("Server started on :8080")
//		log.Fatal(http.ListenAndServe(":8080", router))
//	}
func init() {
	functions.HTTP("PartialUpdateEmployee", PartialUpdateEmployee)

}

func PartialUpdateEmployee(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	initializeFirestore()
	// Get employee ID from the URL path
	id := r.URL.Query().Get("id")
	logClient, _ = logging.NewClient(ctx, "takeoff-task-3")

	// Ensure the logClient is closed after the function completes.
	defer logClient.Close()

	// Create a logger for this function.
	logger := logClient.Logger("my-log")
	logger.Log(logging.Entry{
		Payload:  " Update method started",
		Severity: logging.Info,
	})
	if id == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		logger.Log(logging.Entry{
			Payload:  "Employee ID is required",
			Severity: logging.Error,
		})
		return
	}

	// Parse request body to get partial employee data
	var partialEmp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&partialEmp)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Log(logging.Entry{
			Payload:  "Invalid request body",
			Severity: logging.Error,
		})
		return
	}
	if partialEmpFirstName, ok := partialEmp["firstname"].(string); ok {
		if validations.ValidNameEntry(partialEmpFirstName) {
			http.Error(w, "Name Cannot contain Numbers please enter valid Name", http.StatusBadRequest)
			return
		}
	}
	if partialEmpLastName, ok := partialEmp["lastname"].(string); ok {
		if validations.ValidNameEntry(partialEmpLastName) {
			http.Error(w, "Last Name Cannot contain Numbers please enter valid Last Name", http.StatusBadRequest)
			return
		}
	}
	if partialEmpEmail, ok := partialEmp["email"].(string); ok {
		err = validations.IsValidEmail(partialEmpEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if partialEmpPhoneNo, ok := partialEmp["phoneNo"].(string); ok {
		err = validations.IsNumberValid(partialEmpPhoneNo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Check if the employee ID is provided and valid

	if partialEmpPassword, ok := partialEmp["Password"].(string); ok {
		if partialEmpPassword == "" {
			http.Error(w, "Password cannot be empty", http.StatusBadRequest)
			return
		}
	}
	if partialEmpID, ok := partialEmp["id"].(string); !ok || partialEmpID != id {
		http.Error(w, "ID in request body must match ID in URL", http.StatusBadRequest)
		logger.Log(logging.Entry{
			Payload:  "ID in request body must match ID in URL",
			Severity: logging.Error,
		})
		return
	}

	// Get a reference to the employee document in Firestore
	collectionRef := client.Collection("employees")

	// Query for the employee with the provided ID
	query := collectionRef.Where("id", "==", id).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying Firestore: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error querying Firestore: %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Check if any documents were found
	if len(docs) == 0 {
		http.Error(w, fmt.Sprintf("Employee with ID %s not found", id), http.StatusNotFound)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Employee with ID %s not found", id),
			Severity: logging.Warning,
		})
		return
	}

	// Get the document ID of the employee
	docID := docs[0].Ref.ID

	// Get the current employee data from Firestore
	doc, err := collectionRef.Doc(docID).Get(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting employee data: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error getting employee data: %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Merge the partial data with the existing employee data
	currentData := doc.Data()

	mergedData := mergeMap(currentData, partialEmp)

	// Update the employee data in Firestore
	_, err = collectionRef.Doc(docID).Set(ctx, mergedData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update employee: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Failed to update employee: %v", err),
			Severity: logging.Error,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Employee updated successfully"))
	logger.Log(logging.Entry{
		Payload:  "Employee updated successfully",
		Severity: logging.Info,
	})
}
