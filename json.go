package main

import (
	"bytes"

	simplejson "github.com/bitly/go-simplejson"
)

type json struct {
	*simplejson.Json
}

func newJSON(b []byte) *json {
	jsonBody, err := simplejson.NewJson(b)
	if err != nil {
		panic(err)
	}
	return &json{jsonBody}
}

func (j *json) pretty() string {
	jsonBytes, err := j.EncodePretty()
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func (j *json) buffer() *bytes.Buffer {
	b, err := j.Encode()
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}
