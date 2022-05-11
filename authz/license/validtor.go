package license

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	_schemaHTTP         = "http"
	_localhost          = "http://127.0.0.1"
	_requestPathLicense = "license"

	_httpMethodPost = "POST"
	_httpMethodGet  = "GET"
	_httpMethodPut  = "PUT"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidLicense = errors.New("license is invalid")
)

type Verifier interface {
	Verify() error
}

var _ Verifier = &Validator{}

type Validator struct {
	// Raw is the original license content.
	Raw string `json:"Raw"`

	// validate request target remote address.
	remote        string
	httpGetter    http.Client
	httpReq       http.Request
	method        string
	options       []ValidatorOption
	verifyFunc    Verify
	serializeFunc Serialize
}

func NewValidator(remote, license string, vFunc Verify, options ...ValidatorOption) Validator {
	return Validator{
		Raw:        license,
		remote:     remote,
		options:    options,
		verifyFunc: vFunc,
	}
}

func (l *Validator) SetOptions(opts ...ValidatorOption) {
	l.options = append(l.options, opts...)
}

func (l *Validator) Verify() error {
	for i := 0; i < len(l.options); i++ {
		if err := l.options[i](l); err != nil {
			return err
		}
	}

	if l.serializeFunc == nil {
		l.serializeFunc = toJSON
	}

	if l.verifyFunc == nil {
		return ErrInvalidRequest
	}

	u := &url.URL{
		Scheme: _schemaHTTP,
		Host:   _localhost,
		Path:   _requestPathLicense,
	}
	if l.remote != "" {
		var err error
		if u, err = url.Parse(l.remote); err != nil {
			return err
		}
	}

	// Data within a structure as standard.
	l.httpReq.Method = l.method
	l.httpReq.URL = u
	switch l.method {
	case _httpMethodGet:
		l.httpReq.URL.RawQuery = url.QueryEscape(fmt.Sprintf("Raw=%s", l.Raw))
	case _httpMethodPost, _httpMethodPut:
		c, err := l.serializeFunc(l)
		if err != nil {
			return err
		}
		_, err = l.httpReq.Body.Read(c)
		if err != nil {
			return err
		}
	default:
		l.method = _httpMethodGet
		l.httpReq.Method = _httpMethodGet
		l.httpReq.URL.RawQuery = url.QueryEscape(fmt.Sprintf("Raw=%s", l.Raw))
	}

	resp, err := l.httpGetter.Do(&l.httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if l.verifyFunc(*resp) {
		return nil
	}

	return ErrInvalidLicense
}

type ValidatorOption func(*Validator) error
type Verify func(response http.Response) bool
type Serialize func(*Validator) ([]byte, error)

func WithPOST() ValidatorOption {
	return func(l *Validator) error {
		l.method = _httpMethodPost
		return nil
	}
}

func WithGET() ValidatorOption {
	return func(l *Validator) error {
		l.method = _httpMethodGet
		return nil
	}
}

func WithPUT() ValidatorOption {
	return func(l *Validator) error {
		l.method = _httpMethodPut
		return nil
	}
}

func SerializeWithJSON() ValidatorOption {
	return func(l *Validator) error {
		l.serializeFunc = toJSON
		return nil
	}
}

func toJSON(l *Validator) ([]byte, error) {
	return json.Marshal(l)
}
