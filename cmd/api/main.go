package main

import (
	"context"
	"fmt"
  "encoding/json"

	gpt3 "github.com/PullRequestInc/go-gpt3"
  "net/http"
  "strings"
)


type Question struct {
  Question string
}

func main() {
  fmt.Println("listen server!")
  http.HandleFunc("/chatgpt", func(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("token")
    q := new(Question)
    err := json.NewDecoder(r.Body).Decode(q)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	  ctx := context.Background()
	  client := gpt3.NewClient(token)

		questionParam := validatedQuestion(q.Question)
		response(w, client, ctx, questionParam)
  })

  http.ListenAndServe(":5001", nil)
}

func response(w http.ResponseWriter, client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
    fmt.Fprintf(w, resp.Choices[0].Text)
	})
	if err != nil {
    fmt.Fprintf(w, "error: %+v", err)
    return
	}
	fmt.Printf("\n")
}

func validatedQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "clear", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
