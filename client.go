package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// type End[In, Out any] interface {
// 	Request(ctx context.Context, in *In) (*Out, error)
// }

// TODO: How do I want to handle content-type?

type Endpoint[In, Out any] struct {
	Client
	ctx     context.Context
	Path    string
	Method  string
	Encoder Encoder[In]
	Decoder Decoder[Out]
}

func (e *Endpoint[In, Out]) Test(in any) {
	switch in.(type) {
	case *In, In:
		fmt.Println("Switch In")
	case *Out, Out:
		fmt.Println("Switch Out")
	default:
		panic(errors.New("invalid input type"))
	}
}

func (e *Endpoint[In, Out]) Request(ctx context.Context, data *In) (*Out, error) {
	in, err := e.Encoder.Encode(*data)
	if err != nil {
		return nil, err
	}

	fmt.Println(reflect.TypeOf(data))
	fmt.Println(reflect.ValueOf(data))

	r, err := http.NewRequest(e.Method, e.Path, in)
	if err != nil {
		return nil, err
	}

	// TODO: set these up with the proper headers...
	r.Header.Set("Content-Type", "application/json")

	resp, err := e.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	// TODO: handle other status codes

	out, err := e.Decoder.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type Decoder[Out any] interface {
	Decode(io.ReadCloser) (*Out, error)
}

type Encoder[T any] interface {
	Encode(t T) (io.Reader, error)
}

// type Encoder[In] interface {
// 	Encode(data *In) io.Reader
// }

// type Decoder[Out] interface {
// 	Decode(r io.Reader, out *Out) error
// }

// type Comparable interface {
// 	Compare(other comparable) bool
// }

// type Pageable interface {
// 	TotalPages() int
// }

// type PageOf[Out] struct {
// 	Page
// 	Data []Out
// }

// type Page struct {
// 	PageNumber int
// 	PageSize   int
// 	Total      int
// 	TotalPages int
// }

// type PageableEndpoint[In, Out comparable] Endpoint[In, Out]

// func (e *PageableEndpoint[In, Out]) TotalPages() int {
// 	return e.Page.TotalPages
// }

// func (e *PageableEndpoint[In, Out]) Resources(ctx context.Context, out chan<- interface{}) {
// 	defer close(out)

// 	var data <-chan interface{}
// 	var err error

// 	go func(data chan<- interface{}) {
// 		defer handleRoutinePanic(e.lstream)

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case d, ok := <-data:
// 				if ok {
// 					var page PageOf[Out]
// 					if page, ok = d.(PageOf[Out]); ok {
// 						for _, v := range page.Data {
// 							out <- &v
// 						}
// 					} else {
// 						e.lstream.Send(log.Error("unable to cast paged return as page", err))
// 					}

// 	return out, nil
// }
