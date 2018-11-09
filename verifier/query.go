package verifier

import "github.com/machinebox/graphql"

type graphqlResponse struct {
	NameResolver struct {
		Responses []Response
	}
}

type Response struct {
	Total         int
	SuppliedInput string
	Results       []struct {
		QualitySummary string
		MatchedNames   []MatchedName
	}
	PreferredResults []PreferredResult
}

type MatchedName struct {
	Classification
	DataSource
	Name
	AcceptedName
	MatchType
}

type PreferredResult struct {
	DataSource
	Name
	TaxonID string
}

type DataSource struct {
	ID    int
	Title string
}

type Name struct {
	Value string
}

type Classification struct {
	Path string
}

type AcceptedName struct {
	Name struct {
		Value string
	}
}

type MatchType struct {
	Kind                 string
	VerbatimEditDistance int
	StemEditDistance     int
}

func graphqlRequest() *graphql.Request {
	req := graphql.NewRequest(`
		query($names: [name!]!, $sources: [Int!]) {
			nameResolver(names: $names,
									preferredDataSourceIds: $sources,
								  bestMatchOnly: true) {
				responses {
					total
					suppliedInput
					results {
						qualitySummary
						matchedNames {
							synonym
							classification { path }
							dataSource { id title }
							name { value }
							acceptedName { name { value } }
							matchType {
								kind
								verbatimEditDistance
								stemEditDistance
							}
						}
					}
					preferredResults {
						dataSource {id title}
						name { value } 
						taxonId
						acceptedName { name { value } }
					}
				}
			}
		}
	`)
	return req
}
