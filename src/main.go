package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func genericHandler(specificHandler func(*http.Request, map[string]interface{}) map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData map[string]interface{}
		json.NewDecoder(r.Body).Decode(&requestData)
		fmt.Println("Received data:", requestData)

		response_data := specificHandler(r, requestData)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response_data)
	}
}

func getTableHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}
	fmt.Println("Processing get-table data for user:", username)

	var problems = []map[string]interface{}{}

	for _, problem := range getProblemsFromDatabase(username) {
		problems = append(problems, map[string]interface{}{
			"link":               problem.Link,
			"titleSlug":          problem.TitleSlug,
			"difficulty":         problem.Difficulty,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		})
	}
	fmt.Println("Problems for user", username, ":", problems)
	return map[string]interface{}{
		"message": "Get table data processed",
		"table":   problems,
	}
}

func deleteRowHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	fmt.Println("Processing delete-row data:", data)

	username := r.URL.Query().Get("username")
	problemTitleSlug := r.URL.Query().Get("problemTitleSlug")
	if username == "" || problemTitleSlug == "" {
		fmt.Println("Username or problem title slug not provided")
		return map[string]interface{}{"error": "Username or problem title slug not provided"}
	}

	deleteProblemFromDatabase(username, problemTitleSlug)

	return map[string]interface{}{
		"message": "Delete row data processed",
		"data":    data,
	}
}

func insertRowHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}

	problem := LeetCodeProblem{
		Link:               data["link"].(string),
		TitleSlug:          data["titleSlug"].(string),
		Difficulty:         data["difficulty"].(string),
		RepeatDate:         data["repeatDate"].(string),
		LastCompletionDate: data["lastCompletionDate"].(string),
	}
	upsertProblemIntoDatabase(username, problem)

	return map[string]interface{}{
		"message": "Inserted row data processed",
		"data":    data,
	}
}

func main() {
	godotenv.Load()
	fmt.Println("program running!")

	http.HandleFunc("/get-table", enableCORS(genericHandler(getTableHandler)))
	http.HandleFunc("/delete-row", enableCORS(genericHandler(deleteRowHandler)))
	http.HandleFunc("/insert-row", enableCORS(genericHandler(insertRowHandler)))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
