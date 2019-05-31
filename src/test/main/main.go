package main

import (
	"errors"
	"github.com/containous/alice"
	"log"
	"net/http"
)

// const variables
const (
	READONLY         = "readonly"
	MUITIPLEMANIFEST = "manifest"
)

var Middlewares = []string{MUITIPLEMANIFEST, READONLY, MUITIPLEMANIFEST}

type ChainBuilder struct {
	middlewares []string
}

func New(middlewares []string) *ChainBuilder {
	return &ChainBuilder{
		middlewares: middlewares,
	}
}

// CreateChain ...
func (b *ChainBuilder) CreateChain() *alice.Chain {
	chain := alice.New()

	for _, mName := range b.middlewares {
		log.Println(mName)
	}

	for _, mName := range b.middlewares {
		midName := mName
		log.Println(mName)
		chain = chain.Append(func(next http.Handler) (http.Handler, error) {
			constructor, err := b.getMiddleware(midName)
			if err != nil {
				log.Println("err in constructor ...")
				log.Println(err)
				return nil, err
			}
			//log.Println(midName)
			//log.Println(constructor)
			//log.Println(next)
			return constructor(next)
		})
	}
	return &chain
}

func (b *ChainBuilder) getMiddleware(mName string) (alice.Constructor, error) {
	var middleware alice.Constructor

	if mName == READONLY {
		middleware = func(next http.Handler) (http.Handler, error) {
			return NewReadOnly(next)
		}
	}

	if mName == MUITIPLEMANIFEST {
		if middleware != nil {
			return nil, errors.New("middleware is not nil")
		}
		middleware = func(next http.Handler) (http.Handler, error) {
			return NewMultipleManifestHandler(next)
		}
	}

	if middleware == nil {
		return nil, errors.New("middleware does not exist")
	}

	return middleware, nil
}

type readonlyHandler struct {
	next http.Handler
}

func NewReadOnly(next http.Handler) (http.Handler, error) {
	return &readonlyHandler{
		next: next,
	}, nil
}

func (rh readonlyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("Executing readonlyHandler")
	rh.next.ServeHTTP(rw, req)
}

type MultipleManifestHandler struct {
	next http.Handler
}

func NewMultipleManifestHandler(next http.Handler) (http.Handler, error) {
	return &MultipleManifestHandler{
		next: next,
	}, nil
}

// The handler is responsible for blocking request to upload manifest list by docker client, which is not supported so far by Harbor.
func (mh MultipleManifestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("Executing MultipleManifestHandler")
	mh.next.ServeHTTP(rw, req)
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	finalHandler := http.HandlerFunc(final)

	var handlerChain *alice.Chain
	handlerChain = New(Middlewares).CreateChain()
	//log.Println("111111")
	log.Println(handlerChain)

	head, err := handlerChain.Then(finalHandler)
	if err != nil {
		log.Println(err)
	}
	//log.Println(head)
	//log.Println("2222")

	http.Handle("/", head)
	http.ListenAndServe(":3001", nil)
}
