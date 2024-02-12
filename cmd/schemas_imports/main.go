package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
	"log"
	"os"
)

func main() {
	headers := make(map[string]string)
	headers["X-OpenAI-Api-Key"] = os.Getenv("OPENAI_APIKEY")

	cfg := weaviate.Config{
		Host:       os.Getenv("WEAVIATE_HOST"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WEAVIATE_API_KEY")},
		//Headers:    headers,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//DeleteClassArticle(client)
	CreateClassArticle(client)
	//GetSchema(client)
	//GetMeta(client)
	//DummyDelete(client)
	//DummyCreate(client)
}

func DeleteClassArticle(client *weaviate.Client) {
	className := "Article"
	err := client.Schema().ClassDeleter().
		WithClassName(className).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func CreateClassArticle(client *weaviate.Client) {
	className := "Article"

	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
			{
				Name:     "body",
				DataType: []string{"text"},
			},
			{
				Name:     "url",
				DataType: []string{"text"},
				ModuleConfig: map[string]any{
					"text2vec-openai": map[string]any{
						"skip": true,
					},
				},
			},
		},
		Vectorizer: "text2vec-openai",
	}

	err := client.Schema().ClassCreator().
		WithClass(class).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func DummyDelete(client *weaviate.Client) {
	idToDelete := "05dfb8ce-85f0-4c5a-aa08-6f25301bffca" // created by DummyCreate
	err := client.Data().Deleter().
		WithClassName("TestClass").
		WithID(idToDelete).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func DummyCreate(client *weaviate.Client) {
	res, err := client.Data().Creator().
		WithClassName("TestClass").
		WithProperties(map[string]interface{}{
			"name": "dummy",
		}).Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func GetSchema(client *weaviate.Client) {
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func GetMeta(client *weaviate.Client) {
	meta, err := client.Misc().MetaGetter().Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
