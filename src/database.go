package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

func create_supabase_client() (*supabase.Client, error) {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("Cannot initalize client", err)
	} else {
		fmt.Println("Initailized supabase client")
	}
	return client, err
}

func upsert_problem_into_database(username string, problem LeetCodeProblem) error {
	client, err := create_supabase_client()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Upsert(map[string]interface{}{
			"username":           username,
			"titleSlug":          problem.titleSlug,
			"link":               problem.link,
			"difficulty":         problem.difficulty,
			"repeatDate":         problem.repeatDate,
			"lastCompletionDate": problem.lastCompletionDate,
		}, "username,titleSlug", "", "").
		Execute()

	if err != nil {
		fmt.Println("Error upserting database:", err)
	}
	fmt.Println("Successfully upserted database entry for user:", username)
	return err
}

func delete_problem_from_database(username string, problem_title_slug string) error {
	client, err := create_supabase_client()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Delete("", "").
		Eq("username", username).
		Eq("titleSlug", problem_title_slug).
		Execute()

	if err != nil {
		fmt.Println("Error deleting database entry:", err)
	}
	fmt.Println("Successfully deleted database entry for user:", username)
	return err
}

func get_problems_from_database(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem

	client, e := create_supabase_client()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	fmt.Println("Getting problems from database for user:", username)
	raw_data, _, err := client.From(table).Select("*", "", false).Eq("username", username).Execute()
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return []LeetCodeProblem{}
	}

	fmt.Println("Raw data:", string(raw_data))

	var rawProblems []map[string]interface{}
	err = json.Unmarshal(raw_data, &rawProblems)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return []LeetCodeProblem{}
	}

	for _, rawProblem := range rawProblems {
		problem := LeetCodeProblem{
			link:               rawProblem["link"].(string),
			titleSlug:          rawProblem["titleSlug"].(string),
			difficulty:         rawProblem["difficulty"].(string),
			repeatDate:         rawProblem["repeatDate"].(string),
			lastCompletionDate: rawProblem["lastCompletionDate"].(string),
		}
		problems = append(problems, problem)
	}

	fmt.Printf("Problems for user %s: %+v\n", username, problems)
	return problems
}
