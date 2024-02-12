package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"log"
	"os"
)

func main() {
	headers := make(map[string]string)
	headers["X-OpenAI-Api-Key"] = os.Getenv("OPENAI_APIKEY")

	cfg := weaviate.Config{
		Host:       "edu-demo.weaviate.network",
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: "readonly-demo"},
		Headers:    headers,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	DemoJeopardyQuestionAggregateWithNearTextWhereMultiple(client)
	//DemoJeopardyQuestionAggregateWithNearTextWhere(client)
	//DemoJeopardyQuestionAggregateWithNearTextGrouped(client)
	//DemoJeopardyQuestionAggregate(client)
	//DemoJeopardyQuestionNearText(client)
	//DemoJeopardyQuestionNearObject(client)
	//DemoTweetRaw(client)
	//DemoTweet(client)
	//DemoLondonOlympicsRaw(client)
	//DemoLondonOlympics(client)
	//DemoMajorCities(client)
	//GetMeta(client)
	//GetSchema(client)
}

func DemoJeopardyQuestionAggregateWithNearTextWhereMultiple(client *weaviate.Client) {
	/*
		{
		  Get {
		    JeopardyQuestion (
		      limit: 2
		      nearText: {
		        concepts: ["Intergalactic travel"],
		      }
		      where: {
		        operator: And,
		        operands: [
		        {
		          path: ["question"],
		          operator: Like,
		          valueText: "*rocket*"
		        }
		        {
		          path: ["points"],
		          operator: GreaterThan,
		          valueInt: 400
		        },
		        ]

		      }
		    ) {
		      question
		      answer
		      points
		      _additional {
		        distance
		        id
		      }
		    }
		  }
		}
	*/

	res, err := client.GraphQL().Get().
		WithClassName("JeopardyQuestion").
		WithFields([]graphql.Field{
			{Name: "question"},
			{Name: "answer"},
			{Name: "points"},
			{Name: "_additional", Fields: []graphql.Field{
				{Name: "distance"},
				{Name: "id"},
			}},
		}...).
		WithNearText(client.GraphQL().
			NearTextArgBuilder().
			WithConcepts([]string{"Intergalactic travel"}),
		).
		WithWhere(filters.Where().
			WithOperator(filters.And).
			WithOperands(
				[]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"question"}).
						WithOperator(filters.Like).
						WithValueText("*rocket*"),
					filters.Where().
						WithPath([]string{"points"}).
						WithOperator(filters.GreaterThan).
						WithValueInt(400),
				},
			),
		).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoJeopardyQuestionAggregateWithNearTextWhere(client *weaviate.Client) {
	/*
		{
		  Get {
		    JeopardyQuestion (
		      limit: 2
		      nearText: {
		        concepts: ["Intergalactic travel"],
		      }
		      where: {
		        path: ["question"],
		        operator: Like,
		        valueText: "*rocket*"
		      }
		    ) {
		      question
		      answer
		      _additional {
		        distance
		        id
		      }
		    }
		  }
		}
	*/

	res, err := client.GraphQL().Get().
		WithClassName("JeopardyQuestion").
		WithFields([]graphql.Field{
			{Name: "question"},
			{Name: "answer"},
			{Name: "_additional", Fields: []graphql.Field{
				{Name: "distance"},
				{Name: "id"},
			}},
		}...).
		WithNearText(client.GraphQL().
			NearTextArgBuilder().
			WithConcepts([]string{"Intergalactic travel"}),
		).
		WithWhere(filters.Where().
			WithPath([]string{"question"}).
			WithOperator(filters.Like).
			WithValueText("*rocket*"),
		).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoJeopardyQuestionAggregateWithNearTextGrouped(client *weaviate.Client) {
	/*
		{
		  Aggregate {
		    JeopardyQuestion (
		      nearText: {
		        concepts: ["Intergalactic travel"],
		        distance: 0.2
		      }
		      groupBy: ["round"]
		      ) {
		      groupedBy {
		        path
		        value
		      }
		      meta {
		        count
		      }
		    }
		  }
		}
	*/
	res, err := client.GraphQL().Aggregate().
		WithClassName("JeopardyQuestion").
		WithFields([]graphql.Field{
			{Name: "groupedBy", Fields: []graphql.Field{
				{Name: "path"},
				{Name: "value"},
			}},
			{Name: "meta", Fields: []graphql.Field{
				{Name: "count"},
			}},
		}...).
		WithNearText(client.GraphQL().
			NearTextArgBuilder().
			WithConcepts([]string{"Intergalactic travel"}).
			WithDistance(0.2),
		).
		WithGroupBy("round").
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoJeopardyQuestionAggregate(client *weaviate.Client) {
	/*
		{
		  Aggregate {
		    JeopardyQuestion {
		      answer {
		        count
		        topOccurrences
		        {
		          value
		          occurs
		        }
		      }
		    }
		  }
		}
	*/
	res, err := client.GraphQL().Aggregate().
		WithClassName("JeopardyQuestion").
		WithFields([]graphql.Field{
			{Name: "answer", Fields: []graphql.Field{
				{Name: "count"},
				{Name: "topOccurrences", Fields: []graphql.Field{
					{Name: "value"},
					{Name: "occurs"},
				}},
			}},
		}...).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoJeopardyQuestionNearObject(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "_additional", Fields: []graphql.Field{
			{Name: "distance"},
			{Name: "id"},
		}},
	}
	nearObjectArgumentBuilder := client.GraphQL().NearObjectArgBuilder().
		WithID("c8f8176c-6f9b-5461-8ab3-f3c7ce8c2f5c")

	res, err := client.GraphQL().Get().
		WithClassName("JeopardyQuestion").
		WithFields(fields...).
		WithNearObject(nearObjectArgumentBuilder).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoJeopardyQuestionNearText(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
	}
	nearTextArgumentBuilder := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{"Intergalactic travel"})

	res, err := client.GraphQL().Get().
		WithClassName("JeopardyQuestion").
		WithFields(fields...).
		WithNearText(nearTextArgumentBuilder).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoTweetRaw(client *weaviate.Client) {
	res, err := client.GraphQL().Raw().WithQuery(`
{
	Get {
		WikiCity(
			limit: 3
			nearText: { concepts: "Popular Southeast Asian tourist destination" }
		) {
			city_name
			_additional {
				generate(
					singleResult: {
						prompt: """
						Write a tweet with a potentially surprising fact from {wiki_summary}
						"""
					}
				) {
					singleResult
					error
				}
			}
		}
	}
}
`).Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoTweet(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "city_name"},
		{Name: "wiki_summary"},
	}
	nearTextArgumentBuilder := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{"Popular Southeast Asian tourist destination"})

	searchBuilder := graphql.NewGenerativeSearch().
		SingleResult("Write a tweet with a potentially surprising fact from {wiki_summary}")

	res, err := client.GraphQL().Get().
		WithClassName("WikiCity").
		WithFields(fields...).
		WithNearText(nearTextArgumentBuilder).
		WithLimit(3).
		WithGenerativeSearch(searchBuilder).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoLondonOlympicsRaw(client *weaviate.Client) {
	res, err := client.GraphQL().Raw().WithQuery(`{
	Get {
		WikiCity(
			limit: 1
			ask: {
				question: "When was the London Olympics?"
				properties: ["wiki_summary"]
			}
		) {
			city_name
			country
			lng
			lat
			_additional {
				answer {
					hasAnswer
					property
					result
				}
			}
		}
	}
}
`).Do(context.Background())
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoLondonOlympics(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "city_name"},
		{Name: "country"},
		{Name: "lng"},
		{Name: "lat"},
		{Name: "_additional", Fields: []graphql.Field{
			{Name: "answer", Fields: []graphql.Field{
				{Name: "hasAnswer"},
				{Name: "property"},
				{Name: "result"},
			}},
		}},
	}

	askArgumentBuilder := client.GraphQL().AskArgBuilder().
		WithQuestion("When was the London Olympics?").
		WithProperties([]string{"wiki_summary"})

	res, err := client.GraphQL().Get().
		WithClassName("WikiCity").
		WithFields(fields...).
		WithAsk(askArgumentBuilder).
		WithLimit(1).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func DemoMajorCities(client *weaviate.Client) {
	fields := []graphql.Field{
		{Name: "city_name"},
		{Name: "country"},
		{Name: "lng"},
		{Name: "lat"},
	}
	nearTextArgumentBuilder := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{"Major European city"})

	res, err := client.GraphQL().Get().
		WithClassName("WikiCity").
		WithFields(fields...).
		WithNearText(nearTextArgumentBuilder).
		WithLimit(3).
		Do(context.Background())
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
	fmt.Printf("%v", schema)
}

func GetMeta(client *weaviate.Client) {
	meta, err := client.Misc().MetaGetter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	//fmt.Println(meta.Hostname)
	//fmt.Println(meta.Version)

	//if rec, ok := meta.Modules.(map[string]interface{}); ok {
	//	for key, val := range rec {
	//		log.Printf(" [========>] %s = %s", key, val)
	//	}
	//} else {
	//	fmt.Printf("record not a map[string]interface{}: %v\n", meta.Modules)
	//}

	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
