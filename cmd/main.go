package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "go.devnw.com/api"
)

type fakeclient struct {
	do func(*http.Request) (*http.Response, error)
}

func (*fakeclient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		Status: "200 OK",
		Body: io.NopCloser(bytes.NewBufferString(`{
			"id": "1",
			"name": "test",
			"description": "test",
			"created": "2019-01-01T00:00:00Z"
		}`)),
	}, nil
}

type TypeIn struct {
	Name string `json:"name"`
}

type TypeOut struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"description"`
	Created string `json:"created"`
}

type FakeEncoder[T any] struct {
	data *T
}

func (fe *FakeEncoder[T]) Encode(data T) (io.Reader, error) {
	output, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(output), nil
}

type FakeDecoder[T any] struct {
	data *T
}

func (fe *FakeDecoder[T]) Decode(body io.ReadCloser) (*T, error) {
	var out = new(T)

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func main() {
	c := &fakeclient{
		do: func(*http.Request) (*http.Response, error) {
			return nil, nil
		},
	}

	e1 := &Endpoint[TypeIn, TypeOut]{
		Client:  c,
		Path:    "/",
		Method:  "GET",
		Encoder: &FakeEncoder[TypeIn]{},
		Decoder: &FakeDecoder[TypeOut]{},
	}

	out, err := e1.Request(&TypeIn{Name: "foo"})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
