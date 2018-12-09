package main

import "strconv"

type simpleOutput struct {
	ID              string
	Verbatim        string
	Canonical       string
	CanonicalRanked string
	Authorship      string
	Year            int
	Quality         int
}

func newSimpleOutput(sn scientificNameNode) *simpleOutput {
	so := simpleOutput{
		ID:              sn.VerbatimID,
		Verbatim:        sn.Verbatim,
		Canonical:       sn.CanonicalNode.Value,
		CanonicalRanked: sn.CanonicalNode.ValueRanked,
		Quality:         1,
	}
	return &so
}

func (so *simpleOutput) toSlice() []string {
	yr := strconv.Itoa(so.Year)
	if yr == "0" {
		yr = ""
	}

	qual := strconv.Itoa(so.Quality)
	if qual == "0" {
		qual = ""
	}
	res := []string{
		so.ID,
		so.Verbatim,
		so.Canonical,
		so.CanonicalRanked,
		so.Authorship,
		yr,
		qual,
	}
	return res
}
