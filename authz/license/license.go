package license

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tkeel-io/security/authz/license/ciperutil"
)

var (
	ErrNoSecretKey           = errors.New("no private key set")
	ErrNoPublicKey           = errors.New("no public key set")
	ErrInvalidLicenseContent = errors.New("invalid license content")
	ErrInvalidParsedInfo     = errors.New("invalid parsed info")
)

type Issuer interface {
	Issue() (string, error)
}

type Parser interface {
	Parse() (map[string]string, error)
}

var _ Parser = &License{}
var _ Issuer = &License{}

type License struct {
	secret  string
	key     string
	Info    map[string]string
	options []Option
	content string
}

func NewLicenseParser(content string, opts ...Option) Parser {
	return &License{content: content, options: opts, Info: map[string]string{}}
}

func NewIssuer(opts ...Option) Issuer {
	return &License{options: opts, Info: map[string]string{}}
}

func (l *License) SetInfo(info map[string]string) {
	if l.Info == nil {
		l.Info = map[string]string{}
	}
	for k, v := range info {
		l.Info[k] = v
	}
}

func (l *License) Set(opts ...Option) {
	l.options = append(l.options, opts...)
}

func (l *License) Issue() (string, error) {
	for i := 0; i < len(l.options); i++ {
		if err := l.options[i](l); err != nil {
			return "", err
		}
	}

	if l.secret == "" {
		return "", ErrNoSecretKey
	}

	aesKey32 := ciperutil.GenerateCipherKey32()
	info, err := json.Marshal(l.Info)
	if err != nil {
		return "", err
	}
	aesencr, err := ciperutil.AESEncrypt(string(info), aesKey32)
	if err != nil {
		return "", err
	}
	privpem, err := ciperutil.ParseRSAPrivateKeyFromPEM(l.secret)
	if err != nil {
		return "", err
	}
	secretSigned, err := ciperutil.RSAPrivateEncrypt(privpem, info)
	if err != nil {
		return "", err
	}

	licenseFormula := "%s%04x%s%s"
	license := fmt.Sprintf(licenseFormula, aesKey32, len(aesencr), aesencr, secretSigned)
	return license, nil
}

func (l *License) Parse() (map[string]string, error) {
	for i := 0; i < len(l.options); i++ {
		if err := l.options[i](l); err != nil {
			return nil, err
		}
	}
	if l.key == "" {
		return nil, ErrNoPublicKey
	}
	return l.parse()
}

func (l *License) parse() (map[string]string, error) {
	keypem, err := ciperutil.ParseRSAPublicKeyFromPEM(l.key)
	if err != nil {
		return nil, err
	}
	if l.content == "" || len(l.content) < 32 {
		return nil, ErrInvalidLicenseContent
	}

	// get aes key and decrypt aes encrypted content.
	aeskey := l.content[:32]
	if len(l.content) < 36 {
		return nil, ErrInvalidLicenseContent
	}
	aesEncryptedLenHex := l.content[32:36]
	aesEncrLen, err := strconv.ParseInt(aesEncryptedLenHex, 16, 64)
	if err != nil {
		return nil, err
	}
	if len(l.content) < 36+int(aesEncrLen) {
		return nil, ErrInvalidLicenseContent
	}
	aesEncr := l.content[36 : 36+int(aesEncrLen)]

	// get rsa encrypted content and decrypt it by public key.
	rsaSigned := l.content[36+int(aesEncrLen):]
	decode, err := base64.StdEncoding.DecodeString(rsaSigned)
	if err != nil {
		return nil, err
	}
	decr, err := ciperutil.AESDecrypt(aesEncr, aeskey)
	if err != nil {
		return nil, err
	}
	rsadecr, err := ciperutil.RSAPublicDecrypt(keypem, decode)
	if err != nil {
		return nil, err
	}

	// compare two decrypted.
	if decr != string(rsadecr) {
		return nil, ErrInvalidParsedInfo
	}

	var info map[string]string
	err = json.Unmarshal(rsadecr, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

type Option func(*License) error

func WithSecret(secret string) Option {
	return func(l *License) error {
		l.secret = secret
		return nil
	}
}

func WithPublicKey(key string) Option {
	return func(license *License) error {
		license.key = key
		return nil
	}
}

func WithInfo(info map[string]string) Option {
	return func(license *License) error {
		license.SetInfo(info)
		return nil
	}
}

func WithLicense(license string) Option {
	return func(l *License) error {
		l.content = license
		return nil
	}

}
