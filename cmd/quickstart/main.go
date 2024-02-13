package main

// https://weaviate.io/developers/weaviate/quickstart

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"net/http"
	"os"
)

func main() {
	headers := make(map[string]string)
	headers["X-OpenAI-Api-Key"] = os.Getenv("OPENAI_APIKEY")

	cfg := weaviate.Config{
		Host:   os.Getenv("WEAVIATE_HOST"),
		Scheme: "http",
		//Scheme: "https",
		//AuthConfig: auth.ApiKey{Value: os.Getenv("WEAVIATE_API_KEY")},
		Headers: headers,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	QuestionsGenerativeGrouped(client)
	//QuestionsGenerativeSingle(client)
	//QuestionsWhere(client)
	//QuestionsNearText(client)
	//QuestionsImport(client)
	//QuestionSchemaCreate(client)
}

func QuestionsGenerativeGrouped(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	generativeSearch := graphql.NewGenerativeSearch().
		GroupedResult("Write a tweet with emojis about these facts.")

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(2).
		WithGenerativeSearch(generativeSearch).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonOutput))
}

func QuestionsGenerativeSingle(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	generativeSearch := graphql.NewGenerativeSearch().
		SingleResult("Explain {answer} as you might to a five-year-old.")

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(2).
		WithGenerativeSearch(generativeSearch).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonOutput))
}

func QuestionsWhere(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	where := filters.Where().
		WithPath([]string{"category"}).
		WithOperator(filters.Equal).
		WithValueText("ANIMALS")

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithWhere(where).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", result)
}

func QuestionsNearText(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", result)
}

func QuestionsImport(client *weaviate.Client) {
	// Retrieve the data
	data, err := http.DefaultClient.Get("https://raw.githubusercontent.com/weaviate-tutorials/quickstart/main/data/jeopardy_tiny.json")
	if err != nil {
		panic(err)
	}
	defer data.Body.Close()

	// Decode the data
	var items []map[string]string
	if err := json.NewDecoder(data.Body).Decode(&items); err != nil {
		panic(err)
	}

	// convert items into a slice of models.Object
	objects := make([]*models.Object, len(items))
	for i := range items {
		objects[i] = &models.Object{
			Class: "Question",
			Properties: map[string]any{
				"category": items[i]["Category"],
				"question": items[i]["Question"],
				"answer":   items[i]["Answer"],
			},
		}
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			panic(res.Result.Errors.Error)
		}
	}
}

func QuestionSchemaCreate(client *weaviate.Client) {
	className := "Question"
	class := &models.Class{
		Class:      className,
		Vectorizer: "text2vec-contextionary",
		ModuleConfig: map[string]any{
			"text2vec-contextionary": map[string]any{
				"skip":                  false,
				"vectorizePropertyName": false,
			},
			"generative-openai": map[string]interface{}{},
		},
	}

	err := client.Schema().ClassCreator().
		WithClass(class).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}
