package main

import (
	"bytes"
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

type output struct {
	NameStringID  string              `json:"nameStringId"`
	Parsed        bool                `json:"parsed"`
	Quality       int                 `json:"quality,omitempty"`
	ParserVersion string              `json:"parserVersion"`
	Verbatim      string              `json:"verbatim"`
	Normalized    string              `json:"normalized"`
	CanonicalName canonicalNameOutput `json:"canonicalName"`
	Hybrid        bool                `json:"hybrid"`
	Surrogate     bool                `json:"surrogate"`
	Virus         bool                `json:"virus"`
	Bacteria      bool                `json:"bacteria"`
	Details       namesGroupOutput    `json:"details,omitempty"`
	Positions     positionsOutput     `json:"positions,omitempty"`
}

func newOutput(sn scientificNameNode) *output {
	o := output{}
	return &o
}

// toJSON converts Output to JSON representation.
func (o *output) toJSON(pretty bool) ([]byte, error) {
	if pretty {
		return jsoniter.MarshalIndent(o, "", "  ")
	}
	return jsoniter.Marshal(o)
}

// fromJSON converts JSON representation of Outout to Output object.
func (o *output) fromJSON(data []byte) error {
	r := bytes.NewReader(data)
	return jsoniter.NewDecoder(r).Decode(o)
}

type canonicalNameOutput struct {
	Value       string `json:"value"`
	ValueRanked string `json:"valueRanked"`
}

type namesGroupOutput struct {
	Uninomial uninomialOutput `json:"uninomial,omitempty"`
}

type uninomialOutput struct {
	Value string `json:"value"`
}

type positionsOutput []posOutput

type posOutput struct {
	Type  string
	Start int
	End   int
}

func (p *posOutput) MarshalJSON() ([]byte, error) {
	arr := []interface{}{p.Type, p.Start, p.End}
	return jsoniter.Marshal(arr)
}

func (p *posOutput) UnmarshalJSON(bs []byte) error {
	arr := []interface{}{}
	json.Unmarshal(bs, &arr)
	// TODO: add error handling here.
	p.Type = arr[0].(string)
	p.Start = arr[1].(int)
	p.End = arr[2].(int)
	return nil
}
