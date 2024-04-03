package pkg

import (
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/response"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/validator"
)

var PackagesInstance *Packages

type Packages struct {
	Validator validator.IValidator
	Response  response.IResponse
}

type Option func(*Packages)

func WithValidator(validator validator.IValidator) Option {
	return func(p *Packages) {
		p.Validator = validator
	}
}

func WithResponse(response response.IResponse) Option {
	return func(p *Packages) {
		p.Response = response
	}

}

func New(opts ...Option) *Packages {
	p := &Packages{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
