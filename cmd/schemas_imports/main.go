package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"log"
	"os"
	"strings"
)

type JeopardyQuestion struct {
	Round    string `json:"round"`
	Value    int    `json:"value"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

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

	JeopardyQuestionsImport(client)
	//JeopardyQuestionSchemaCreate(client)
	//BatchImport(client)
	//DeleteClass(client, "TestClass")
	//DeleteClassArticle(client)
	//CreateClassArticle(client)
	//GetSchema(client)
	//GetMeta(client)
	//DummyDelete(client)
	//DummyCreate(client)
}

func generateWeaviateId(input string) strfmt.UUID {
	input = strings.ToLower(input)
	hash := md5.Sum([]byte(input))
	u := fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:])
	return strfmt.UUID(u)
}

func JeopardyQuestionsImport(client *weaviate.Client) {
	dat, err := os.ReadFile("./jeopardy_100.json")
	if err != nil {
		panic(err)
	}

	var imports []JeopardyQuestion
	err = json.Unmarshal(dat, &imports)
	if err != nil {
		log.Fatal(err)
	}

	className := "JeopardyQuestion"
	var dataObjs []models.PropertySchema
	for _, q := range imports {
		fmt.Println(q.Question)
		dataObjs = append(dataObjs, map[string]interface{}{
			"round":    q.Round,
			"value":    q.Value,
			"question": q.Question,
			"answer":   q.Answer,
		})
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, dataObj := range dataObjs {
		batcher.WithObjects(&models.Object{
			Class:      className,
			Properties: dataObj,
			ID:         generateWeaviateId((dataObj.(map[string]interface{}))["question"].(string)),
		})
	}

	res, err := batcher.Do(context.Background())
	if err != nil {
		panic(err)
	}

	for i, r := range res {
		fmt.Printf("index %d: %s lastUpdateTimeUnix: %d\n", i, r.ID, r.LastUpdateTimeUnix)
		if r.Result.Errors != nil {
			for _, e := range r.Result.Errors.Error {
				fmt.Printf("error at index %d: %s\n", i, e.Message)
			}
		}
	}
}

func JeopardyQuestionSchemaCreate(client *weaviate.Client) {
	className := "JeopardyQuestion"
	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     "round",
				DataType: []string{"text"},
				ModuleConfig: map[string]any{
					"text2vec-contextionary": map[string]any{
						"skip": true,
					},
				},
			},
			{
				Name:     "value",
				DataType: []string{"int"},
			},
			{
				Name:     "question",
				DataType: []string{"text"},
			},
			{
				Name:     "answer",
				DataType: []string{"text"},
			},
		},
		Vectorizer: "text2vec-contextionary",
		ModuleConfig: map[string]any{
			"text2vec-contextionary": map[string]any{
				"skip":                  false,
				"vectorizePropertyName": false,
			},
		},
	}

	err := client.Schema().ClassCreator().
		WithClass(class).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func BatchImport(client *weaviate.Client) {
	className := "Article"
	var dataObjs []models.PropertySchema
	for i := 0; i < 5; i++ {
		dataObjs = append(dataObjs, map[string]interface{}{
			"title": fmt.Sprintf("Title %v", i),
			"url":   fmt.Sprintf("https://example.com/article/%v", i),
		})
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, dataObj := range dataObjs {
		batcher.WithObjects(&models.Object{
			Class:      className,
			Properties: dataObj,
		})
	}

	res, err := batcher.Do(context.Background())
	if err != nil {
		panic(err)
	}

	//b, err := json.MarshalIndent(res, "", "  ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(b))

	for i, r := range res {
		fmt.Printf("index %d: %s lastUpdateTimeUnix: %d\n", i, r.ID, r.LastUpdateTimeUnix)
		if r.Result.Errors != nil {
			for _, e := range r.Result.Errors.Error {
				fmt.Printf("error at index %d: %s\n", i, e.Message)
			}
		}
	}
}

func DeleteClass(client *weaviate.Client, className string) {
	err := client.Schema().ClassDeleter().
		WithClassName(className).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
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
	idToDelete := "3b2dd386-7700-434f-80bf-d3472a913186" // created by DummyCreate
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
