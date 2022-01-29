package api

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type Method string

const (
	GET    Method = http.MethodGet
	POST   Method = http.MethodPost
	PUT    Method = http.MethodPut
	DELETE Method = http.MethodDelete
)

type Endpoint struct {
	API     API
	Path    string
	Method  Method
	Handler http.HandlerFunc
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Encoder func(interface{}) ([]byte, error)
type Decoder func([]byte, interface{}) error

type API struct {
	Client
	Encoder Encoder
	Decoder Decoder
}

var ErrMethodNotAllowed = errors.New("method not allowed")

type H[T any] struct {
	Value T
	RW    http.ResponseWriter
}

type Handler interface {
	any
	http.ResponseWriter
}

func Host[T any](ctx context.Context, e *Endpoint) <-chan H[T] {
	out := make(chan H[T])

	e.Handler = func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		v, err := Decode[T](e.API.Decoder, b)
		if err != nil {
			panic(err)
		}

		h := H[T]{
			Value: v,
			RW:    w,
		}

		out <- h
	}

	return out
}

func Query[R any, E error](e Endpoint, body io.Reader) (R, error) {
	req, err := http.NewRequest(string(e.Method), e.Path, body)
	if err != nil {
		return *new(R), err
	}

	resp, err := e.API.Client.Do(req)
	if err != nil {
		return *new(R), err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(R), err
	}

	return Decode[R](e.API.Decoder, b)
}

func Decode[T any](d Decoder, b []byte) (T, error) {
	value := new(T)

	err := d(b, value)
	if err != nil {
		return *new(T), err
	}

	return *value, nil
}

// func Get[U, T any](endpoint Endpoint[U, T]) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != string(endpoint.Method) {
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"message": "Hello World"}`))
// 	}
// }

// type HttpClient interface {
// 	Do(r *http.Request) (*http.Response, error)
// }

// // type End[In, Out any] interface {
// // 	Request(ctx context.Context, in *In) (*Out, error)
// // }

// // TODO: How do I want to handle content-type?

// func Client(
// 	ctx context.Context,
// 	c HttpClient,
// 	url *url.URL,
// ) HttpClient {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	return &client{
// 		c,
// 		ctx,
// 		cancel,
// 	}
// }

// type client struct {
// 	HttpClient
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

// type Encoder[Enc, Dec any] interface {
// 	ContentType() string
// 	Encode(t Enc) (io.Reader, error)
// 	Decode(io.ReadCloser) (*Dec, error)
// }

// func New[In, Out any](
// 	ctx context.Context,
// 	c HttpClient,
// 	e Encoder[In, Out],
// 	method, path string,
// ) Endpoint[In, Out] {
// 	return &endpoint[In, Out]{c, e, ctx, path, method}
// }

// type Endpoint[In, Out any] interface {
// 	Request(ctx context.Context, in *In) (*Out, error)
// }

// type endpoint[In, Out any] struct {
// 	c HttpClient
// 	Encoder[In, Out]
// 	context.Context
// 	Path   string
// 	Method string
// }

// func (e *endpoint[In, Out]) Request(ctx context.Context, data *In) (*Out, error) {
// 	in, err := e.Encoder.Encode(*data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(reflect.TypeOf(data))
// 	fmt.Println(reflect.ValueOf(data))

// 	r, err := http.NewRequestWithContext(ctx, e.Method, e.Path, in)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TODO:? Should this be here?
// 	r.Header.Set("Content-Type", e.ContentType())

// 	resp, err := e.c.Do(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode >= 400 {
// 		return nil, errors.New(resp.Status)
// 	}

// 	// TODO: handle other status codes
// 	out, err := e.Encoder.Decode(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return out, nil
// }

// type Decoder[Out any] interface {
// 	Decode(io.ReadCloser) (*Out, error)
// }

// type Encoder[T any] interface {
// 	Encode(t T) (io.Reader, error)
// }

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
