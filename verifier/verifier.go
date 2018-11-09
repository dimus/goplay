package verifier

import (
	"context"
	"fmt"
	"log"

	"github.com/dimus/goplay/util"
	"github.com/machinebox/graphql"
)

type NameInput struct {
	Value string `json:"value"`
}

func Verify(names []string, m *util.Model) {
	client := graphql.NewClient(m.Verifier.URL)
	client.Log = func(s string) { log.Println(s) }
	req := graphqlRequest()

	req.Var("names", jsonNames(names))
	req.Var("sources", []int{1, 2, 3, 4, 5})

	ctx := context.Background()
	var res graphqlResponse
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	for _, v := range res.NameResolver.Responses {
		fmt.Println("total", v.Total)
		fmt.Println("suppliedInput", v.SuppliedInput)
		if len(v.Results) > 0 {
			fmt.Println("quality", v.Results[0].QualitySummary)
			for _, vv := range v.Results[0].MatchedNames {
				fmt.Println("classification", vv.Classification.Path)
				fmt.Println("sourceId", vv.DataSource.ID)
				fmt.Println("sourceTitle", vv.DataSource.Title)
				fmt.Println("name", vv.Name.Value)
				fmt.Println("accepted", vv.AcceptedName.Name.Value)
				fmt.Println("match", vv.MatchType.Kind)
				fmt.Println("editDistance", vv.MatchType.VerbatimEditDistance)
				fmt.Println("stemEditDist", vv.MatchType.StemEditDistance)
			}
		}
		fmt.Println()
	}
}

func jsonNames(names []string) []NameInput {
	res := make([]NameInput, len(names))
	for i, v := range names {
		res[i] = NameInput{Value: v}
	}
	return res
}
