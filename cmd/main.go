package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	. "go.atomizer.io/stream"
	"go.devnw.com/api"
)

type Test struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// var handlers []http.HandlerFunc

	// handlers = append(handlers, Get(Endpoint[string, string]{
	// 	Method: "GET",
	// 	Path:   "/",
	// }))

	a := api.API{
		Client: http.DefaultClient,
		Encoder: func(interface{}) ([]byte, error) {
			return nil, errors.New("not implemented")
		},
		Decoder: func(b []byte, v interface{}) error {
			str, ok := v.(*string)
			if !ok {
				return fmt.Errorf("expected *string, got %T", v)
			}

			*str = string(b)

			return nil //errors.New("not implemented")
		},
	}

	endpoints := []*api.Endpoint{
		{
			API:    a,
			Method: "GET",
			// Path:   "https://benjiv.com",
			Path: "/",
		},
		{
			API:    a,
			Method: "GET",
			// Path:   "https://benjiv.com",
			Path: "/test",
		},
	}

	// body, err := Query[string, error](e, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// println(body)

	incomingRequests := make(chan any)

	go Pipe(
		ctx,
		Intercept(
			ctx,
			api.Host[string](ctx, endpoints[0]),
			func(ctx context.Context, value api.H[string]) (any, bool) {
				fmt.Println("intercepted")
				return value, true
			}), incomingRequests)

	go Pipe(
		ctx,
		Intercept(
			ctx,
			api.Host[Test](ctx, endpoints[1]),
			func(ctx context.Context, value api.H[Test]) (any, bool) {
				fmt.Println("intercepted test endpoint")
				return value, true
			}), incomingRequests)

	for _, e := range endpoints {
		spew.Dump(e)

		http.Handle(e.Path, e.Handler)
	}

	go http.ListenAndServe("127.0.0.1:5000", http.DefaultServeMux)

	for o := range incomingRequests {
		spew.Dump(o)
	}
}

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

// type C[U Endpoint[In, Out], In, Out any] map[string]U

// type fakeclient struct {
// 	do func(*http.Request) (*http.Response, error)
// }

// func (*fakeclient) Do(*http.Request) (*http.Response, error) {
// 	return &http.Response{
// 		Status: "200 OK",
// 		Body: io.NopCloser(bytes.NewBufferString(`{
// 			"id": "1",
// 			"name": "test",
// 			"description": "test",
// 			"created": "2019-01-01T00:00:00Z"
// 		}`)),
// 	}, nil
// }

// type TypeIn struct {
// 	Name string `json:"name"`
// }

// type TypeOut struct {
// 	ID      string `json:"id"`
// 	Name    string `json:"name"`
// 	Desc    string `json:"description"`
// 	Created string `json:"created"`
// }

// type CClient struct {
// 	client Client
// }

// func Request[T any](
// 	e Encoder[T],
// 	method, url string,
// 	body *T,
// ) (*http.Request, error) {
// 	in, err := e.Encode(*body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return http.NewRequest(method, url, in)
// }

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

// type FakeEncoder[Enc, Dec any] struct{}

// func (fe *FakeEncoder[Enc, Dec]) Encode(data Enc) (io.Reader, error) {
// 	output, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bytes.NewBuffer(output), nil
// }

// func (fe *FakeEncoder[Enc, Dec]) Decode(body io.ReadCloser) (*Dec, error) {
// 	var out = new(Dec)

// 	data, err := io.ReadAll(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(data, out)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return out, nil
// }

// func (fe *FakeEncoder[Enc, Dec]) ContentType() string {
// 	return "application/json"
// }

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

// type Query struct {
// }

// type Vulnerability struct {
// }

// type Asset struct {
// }

// type httpClient interface {
// 	Do(r *http.Request) (*http.Response, error)
// }

// type Requester struct {
// 	client api.HttpClient
// 	Vulns  api.Endpoint[Query, Vulnerability]
// 	Assets api.Endpoint[Query, Asset]
// }

// func Client(ctx context.Context, h httpClient) *Requester {
// 	url, err := url.Parse("http://localhost:8080")
// 	if err != nil {
// 		panic(err)
// 	}

// 	r := &Requester{
// 		client: api.Client(ctx, h, url),
// 	}

// 	r.Vulns = api.New[Query, Vulnerability](
// 		ctx,
// 		r.client,
// 		&FakeEncoder[Query, Vulnerability]{},
// 		http.MethodGet,
// 		"/api/v1/vulnerabilities",
// 	)

// 	r.Assets = api.New[Query, Asset](
// 		ctx,
// 		r.client,
// 		&FakeEncoder[Query, Asset]{},
// 		http.MethodGet,
// 		"/api/v1/assets",
// 	)

// 	return r
// }

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	c := Client(ctx, &fakeclient{})

// 	out, err := c.Vulns.Request(ctx, &Query{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(out)

// 	// out, err := Vulns[TypeIn, TypeOut](ctx, &TypeIn{})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println(out)

// 	// out, err := Root.Request(ctx, &TypeIn{Name: "test"})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// ends := C[Endpoint[TypeIn, TypeOut], TypeIn, TypeOut]{
// 	// 	"endpoint": Endpoint[TypeIn, TypeOut]{
// 	// 		Client: &fakeclient{
// 	// 			do: func(*http.Request) (*http.Response, error) {
// 	// 				return nil, nil
// 	// 			},
// 	// 		},
// 	// 		Path:    "/",
// 	// 		Method:  "GET",
// 	// 		Encoder: &FakeEncoder[TypeIn]{},
// 	// 		Decoder: &FakeDecoder[TypeOut]{},
// 	// 	}}

// 	// fmt.Println(out)

// 	// c := &fakeclient{
// 	// 	do: func(*http.Request) (*http.Response, error) {
// 	// 		return nil, nil
// 	// 	},
// 	// }

// 	// cc := &CClient{
// 	// 	client: c,
// 	// }

// 	// out, err := cc.Request(ctx, TYPES, &TypeIn{})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// v := As[*TypeOut](out)

// 	// e1 := &Endpoint[TypeIn, TypeOut]{
// 	// 	Client:  c,
// 	// 	Path:    "/",
// 	// 	Method:  "GET",
// 	// 	Encoder: &FakeEncoder[TypeIn]{},
// 	// 	Decoder: &FakeDecoder[TypeOut]{},
// 	// }

// 	// // m["/api/v1/scanners"] = e1

// 	// out, err := e1.Request(ctx, &TypeIn{Name: "foo"})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// vals, err := url.ParseQuery("foo=bar&baz=qux")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// e1.Test(&TypeIn{Name: "foo"})
// 	// e1.Test(TypeIn{Name: "foo"})

// 	// e1.Test(TypeOut{Name: "foo"})

// 	// e1.Test(&TypeOut{Name: "foo"})

// 	// fmt.Println(vals)

// 	// fmt.Println(out)
// }
