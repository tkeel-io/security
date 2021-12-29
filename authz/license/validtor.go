package license

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	_http      = "http"
	_localhost = "http://127.0.0.1"
	_hostPath  = "license"

	_post = "POST"
	_get  = "GET"
	_put  = "PUT"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnableLicense  = errors.New("license is unable")
)

type Verifier interface {
	Verify() error
}

var _ Verifier = &Validator{}

type Validator struct {
	// Raw is the original license content.
	Raw []byte `json:"Raw"`

	// validate request target remote address.
	remote        string
	httpGetter    http.Client
	httpReq       http.Request
	method        string
	options       []ValidatorOption
	verifyFunc    Verify
	serializeFunc Serialize
}

func NewValidator(remote string, license []byte, vFunc Verify, options ...ValidatorOption) Validator {
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
		Scheme: _http,
		Host:   _localhost,
		Path:   _hostPath,
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
	case _get:
		l.httpReq.URL.RawQuery = url.QueryEscape(fmt.Sprintf("Raw=%s", l.Raw))
	case _post, _put:
		c, err := l.serializeFunc(l)
		if err != nil {
			return err
		}
		_, err = l.httpReq.Body.Read(c)
		if err != nil {
			return err
		}
	default:
		l.method = _get
		l.httpReq.Method = _get
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

	return ErrUnableLicense
}

type ValidatorOption func(*Validator) error
type Verify func(response http.Response) bool
type Serialize func(*Validator) ([]byte, error)

func WithPOST() ValidatorOption {
	return func(l *Validator) error {
		l.method = _post
		return nil
	}
}

func WithGET() ValidatorOption {
	return func(l *Validator) error {
		l.method = _get
		return nil
	}
}

func WithPUT() ValidatorOption {
	return func(l *Validator) error {
		l.method = _put
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
