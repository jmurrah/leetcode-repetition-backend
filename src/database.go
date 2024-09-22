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

// func upsert_database(username string, problems []LeetCodeProblem) {
// 	user := User{
// 		Username: username,
// 		Problems: problems,
// 	}
// 	client, e := create_supabase_client()
// 	if e != nil {
// 		fmt.Println("Error creating supabase client:", e)
// 		return
// 	}
// 	table := os.Getenv("SUPABASE_TABLE")

//		_, _, err := client.From(table).Insert(user, true, "username", "success", "").Execute()
//		if err != nil {
//			fmt.Println("Error upserting database:", err)
//		}
//		fmt.Println("Successfully upserted database entry for user:", username)
//	}

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

	return err
}

func get_problems_from_database(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem
	var raw_response []map[string]json.RawMessage

	client, e := create_supabase_client()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	fmt.Println("Getting problems from database for user:", username)
	raw_data, _, _ := client.From(table).Select("*", "", false).Eq("username", username).Execute()
	json.Unmarshal(raw_data, &raw_response)

	fmt.Println("Raw response:", raw_response)

	if len(raw_response) == 0 {
		fmt.Println("User not found. Adding user to database:", username)
		return []LeetCodeProblem{}
	}
	json.Unmarshal(raw_response[0]["problems"], &problems)
	return problems
}
