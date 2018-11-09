package verifier

import "github.com/machinebox/graphql"

type graphqlResponse struct {
	NameResolver struct {
		Responses []Response
	}
}

type Response struct {
	Total int
}

func graphqlRequest() graphq.Request {
	req := graphql.NewRequest(`
		query($names: [name!]!, $sources: [Int!]) {
			nameResolver(names: $names,
									preferredDataSourceIds: $sources) {
				responses {
					total
					preferredResults {
						dataSource {id}
					}
				}
			}
		}
	`)
	return req
}
