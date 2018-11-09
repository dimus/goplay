package verifier

import (
	"context"
	"fmt"
	"log"

	"github.com/shurcooL/graphql"
)

type Name struct {
	Value string `json:"value"`
}

func Verify(names []string, m *utilModel) VerifyOutput {
	client := graphql.NewClient("https://index.globalnames.org/api/graphql")
	client.Log = func(s string) { log.Println(s) }
	req := graphqlRequest()

	req.Var("names", jsonNames(names)}})
	req.Var("sources", []int{1, 2, 3, 4, 5})

	ctx := context.Background()
	var res response
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	for _, v := range res.NameResolver.Responses {
		fmt.Println("total:", v.Total)
	}
}
