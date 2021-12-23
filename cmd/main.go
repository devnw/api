package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "go.devnw.com/api"
)

// type Clt struct {
// 	c    Client
// 	ends map[string]End
// 	e    Encoder
// }

// type End struct {
// 	Path   string
// 	Method string
// }

// func (e End) Test[T any]() {
// 	switch t.(type) {
// 	case *T, T:
// 		fmt.Println("Switch T")
// 	default:
// 		panic(errors.New("invalid input type"))
// 	}
// }

// func Req[In, Out any](client Clt, in *, as type) (*Out, error) {
// 	in, err := e.Encode(*in)
// 	if err != nil {
// 		return nil, err
// 	}

// 	r, err := http.NewRequest(e.Method, e.Path, in)
// 	if err != nil {
// 		return nil, err
// 	}
// }

// func RE[In, Out any](e End)

type C[U Endpoint[In, Out], In, Out any] map[string]U

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

type CClient struct {
	client Client
}

func Request[T any](
	e Encoder[T],
	method, url string,
	body *T,
) (*http.Request, error) {
	in, err := e.Encode(*body)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, in)
}

// func (c *CClient) Request(ctx context.Context, endpoint EP, data any) (any, error) {

// 	e1 := &Endpoint[TypeIn, TypeOut]{
// 		Client:  c.client,
// 		Path:    "/",
// 		Method:  "GET",
// 		Encoder: &FakeEncoder[TypeIn]{},
// 		Decoder: &FakeDecoder[TypeOut]{},
// 	}

// 	switch endpoint {
// 	case TYPES:
// 		d, ok := data.(*TypeIn)
// 		if !ok {
// 			return nil, errors.New("invalid data type")
// 		}

// 		return e1.Request(ctx, d)
// 	default:
// 		return nil, errors.New("invalid endpoint")
// 	}
// }

// type endconstraint interface {
// 	Endpoint[TypeIn, TypeOut] | Endpoint[TypeOut, TypeIn]
// }

// type Client[T endconstraint] struct {
// 	endpoints map[string]endconstraint

// func Request[In, Out any](ctx context.Context, ep string, in *In) (*Out, error) {
// 	switch ep {
// 	case "/api/v1/scanners":
// 		return
// 	case "/api/v1/notifications":
// 	}
// }

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

// type vuln interface {
// 	struct {
// 		ID      string `json:"id"`
// 		Name    string `json:"name"`
// 		Desc    string `json:"description"`
// 		Created string `json:"created"`
// 	} | struct {
// 		ID   string `json:"id"`
// 		Name string `json:"name"`
// 		Desc string `json:"description"`
// 	}
// }

// func Vulns[In any, Out vuln](ctx context.Context, data *In) (*Out, error) {
// 	return (&Endpoint[In, Out]{
// 		Client: &fakeclient{
// 			do: func(*http.Request) (*http.Response, error) {
// 				return nil, nil
// 			},
// 		},
// 		Path:    "/",
// 		Method:  "GET",
// 		Encoder: &FakeEncoder[In]{},
// 		Decoder: &FakeDecoder[Out]{},
// 	}).Request(ctx, data)
// }

type Query struct {
}

type Vulnerability struct {
}

type Asset struct {
}

func NewClient(c Client) *CC {
	return &CC{
		Vulns: Endpoint[Query, Vulnerability]{
			Client:  c,
			Path:    "/api/v1/vulnerabilities",
			Method:  "GET",
			Encoder: &FakeEncoder[Query]{},
			Decoder: &FakeDecoder[Vulnerability]{},
		},
		Assets: Endpoint[Query, Asset]{
			Client:  c,
			Path:    "/api/v1/assets",
			Method:  "GET",
			Encoder: &FakeEncoder[Query]{},
			Decoder: &FakeDecoder[Asset]{},
		},
	}
}

type CC struct {
	Client
	Vulns  Endpoint[Query, Vulnerability]
	Assets Endpoint[Query, Asset]
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewClient(&fakeclient{})

	out, err := c.Vulns.Request(ctx, &Query{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)

	// out, err := Vulns[TypeIn, TypeOut](ctx, &TypeIn{})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(out)

	// out, err := Root.Request(ctx, &TypeIn{Name: "test"})
	// if err != nil {
	// 	panic(err)
	// }

	// ends := C[Endpoint[TypeIn, TypeOut], TypeIn, TypeOut]{
	// 	"endpoint": Endpoint[TypeIn, TypeOut]{
	// 		Client: &fakeclient{
	// 			do: func(*http.Request) (*http.Response, error) {
	// 				return nil, nil
	// 			},
	// 		},
	// 		Path:    "/",
	// 		Method:  "GET",
	// 		Encoder: &FakeEncoder[TypeIn]{},
	// 		Decoder: &FakeDecoder[TypeOut]{},
	// 	}}

	// fmt.Println(out)

	// c := &fakeclient{
	// 	do: func(*http.Request) (*http.Response, error) {
	// 		return nil, nil
	// 	},
	// }

	// cc := &CClient{
	// 	client: c,
	// }

	// out, err := cc.Request(ctx, TYPES, &TypeIn{})
	// if err != nil {
	// 	panic(err)
	// }

	// v := As[*TypeOut](out)

	// e1 := &Endpoint[TypeIn, TypeOut]{
	// 	Client:  c,
	// 	Path:    "/",
	// 	Method:  "GET",
	// 	Encoder: &FakeEncoder[TypeIn]{},
	// 	Decoder: &FakeDecoder[TypeOut]{},
	// }

	// // m["/api/v1/scanners"] = e1

	// out, err := e1.Request(ctx, &TypeIn{Name: "foo"})
	// if err != nil {
	// 	panic(err)
	// }

	// vals, err := url.ParseQuery("foo=bar&baz=qux")
	// if err != nil {
	// 	panic(err)
	// }

	// e1.Test(&TypeIn{Name: "foo"})
	// e1.Test(TypeIn{Name: "foo"})

	// e1.Test(TypeOut{Name: "foo"})

	// e1.Test(&TypeOut{Name: "foo"})

	// fmt.Println(vals)

	// fmt.Println(out)
}
