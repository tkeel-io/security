package ciperutil

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var priv_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1Fh8AWKSyZKzam044MxTga7jY45Yc53zKzgVrFM6K8vZje14
A9QevqQGaI4XOLf4xqHVdETD52R0cZ9QqJ1kChpLR2bqPsHUIxavn5DHsV60cYSy
egkQ1Tm7Ty6tH3cdqDwH0o5K5e1XDqhE5fLdr03oPoYme8w+zTvHY/DxvG77LCkQ
oqD5whjwlbNhsuS1/b/xMy6gElyvtAOC3kRx1LJ+3tVFOq20qthU0m33Ch0yl4Wd
zjL353W7+n7ia71aSUkuIyQ8Nm/hUeIjQ/e7ZROHZMExejp72XWyO3lvudvpsUtl
n8YAR2M/8nf4gpIpnnmeBX6zTay5JOIlNlpSvQIDAQABAoIBAHlPE4i3B6Sgal8i
hpvmHs63WrBFYcKrOYA3SipoYZMjoeWzBB0r0WSX0BFrG1kLwNO5IHiL0F8gxGUk
3q91OvGKk/b2lFvF36ssCqjdzTdHI062wD93bmZh1OAnij+vjQHPSajAIYm5TULS
Xon/dSXxG/ycJuASjs4wmHg/7dn3gg4767ZIpP4Z80bDrjqxjabKTj/Jppvnd3jy
cpMOl7mhIAkk5S0qpFu/OxiZ8vwp6DX8kEsCAXur+GcLCIiUUrRQ/RNSGCodPWYN
zNurrY8/jOFqCS4eiBn7cLPl0MepB0EoGbZDWCiiEe7OPbJWV190ulWpDYKtWWH0
HecqhwECgYEA+PdOPtBB6HSVCkrg4ZzbcsBQn3Gf7gm7J8+1XNdkO5D0X9aiUozz
Zv42PuH8z63Ej5UX/IQFZsp7piWM/pncWh56PxGU6GZkZ9i3fq40lBNzVfoxH5IN
QGmfR+VMSwi2dojXe/JdXm4RrQxcm18Lu6rl10HZLiHIy00/MxPp0k8CgYEA2lhQ
k8aitcv9p32lgIjo0pvLSq9zvSea1ZBvi/Oc9zFpfNQXD9/xwnjAj95bMbZLsouV
2NBl7sblKWTk379NVzEg3cDvvb1M2GlWDjPKF5mFvomx6Xtg78g1vZiJ4k1qDVhE
B6ollGuEz5xec3PoYxTArG/SWmyw+Y2WsZ6pgzMCgYAm++J+p3GKiqbDw9HOwhcC
suZs8QfP6CosI7QMY5XIVfxN/XfRYUzDtc6crho+EsGSkg/9dFa8L7yI2ZxoSYNd
gSQ3N/OPKGlUcXTaG0EUZq9KqTCD9wSoL2HHijoWDbk0elzzhZHlNWsDI8pkcc+O
hUYUHLV7KcMdQm1A/D9CtQKBgQCk16QNzRsucVGRT/TRM4vC5Uf3nLqehfYJYkA0
wrkwjqd7TIwUuhfFoHCQrWjgASbpJyT8SWmLebGtLLT2j3EOcNLFWFInR3FquSv1
EPh0FL26ei5NfY5TuS2JdE41AgkdBhRmTPiOMxZTv1Q3ibxagWJtTQbcqc56uYCZ
nJWRrQKBgQCqFINEp75IODP8hM0cyuSWqOp1AaRISVQhSyAp6SNSUUWm/1nJP4EY
tsXi4hUnN6WPWEQYWT4S9n8cLdMHFm5mdF4PKlW4mNSs0SRtqtgWVfuPzGMcfRro
nd2EezVT8aWMzzwnklTct23cxE/QWn3kqSfLBlVyWCDmAC5zcRKpNg==
-----END RSA PRIVATE KEY-----`

var pub_key = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1Fh8AWKSyZKzam044MxT
ga7jY45Yc53zKzgVrFM6K8vZje14A9QevqQGaI4XOLf4xqHVdETD52R0cZ9QqJ1k
ChpLR2bqPsHUIxavn5DHsV60cYSyegkQ1Tm7Ty6tH3cdqDwH0o5K5e1XDqhE5fLd
r03oPoYme8w+zTvHY/DxvG77LCkQoqD5whjwlbNhsuS1/b/xMy6gElyvtAOC3kRx
1LJ+3tVFOq20qthU0m33Ch0yl4WdzjL353W7+n7ia71aSUkuIyQ8Nm/hUeIjQ/e7
ZROHZMExejp72XWyO3lvudvpsUtln8YAR2M/8nf4gpIpnnmeBX6zTay5JOIlNlpS
vQIDAQAB
-----END RSA PUBLIC KEY-----`

func TestGenerateRSA(t *testing.T) {
	priv, pub := GenerateRSA()

	assert.True(t, strings.HasPrefix(priv, "-----BEGIN RSA PRIVATE KEY-----"))
	assert.True(t, strings.HasPrefix(pub, "-----BEGIN RSA PUBLIC KEY-----"))
}

func TestParseRSAPublicKeyFromPEM(t *testing.T) {
	pub :=
		`-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1QpuzBeblZ6MmfgPoph0
9lDnnBlVwUnujyDJPA7LO+ybcGZ/pyfcZ+YmGe/l0Fgf8AJIlGj414OKQ9NmuoS6
3XUt9zEBlkqqqG7YQ7cP7Fm9ZYJaZYor0MigF6ITx/R54BezptSsCgzUmDLakqAo
AW2FBrFCKVoxprdmeHX0GtTn335KOEMOdRh3OkGfYzKlpM0DNkyBxV5sqhtFgso+
nxExiZpG3RtH+SaCYseq9NZ1Py55UzoEEEbueRM8mc2ur3vMiqWATPgjOze6LSWP
7Q27qfKTBVAx9OGzmREk/ucbsxckka/yApgD/lvQ7GHdkpWYx0W+CcMEAh48m5fz
DQIDAQAB
-----END RSA PUBLIC KEY-----`
	pubKey, err := ParseRSAPublicKeyFromPEM(pub)
	assert.Nil(t, err)
	assert.IsType(t, &rsa.PublicKey{}, pubKey)
}

func TestParseRSAPrivateKeyFromPEM(t *testing.T) {
	key, err := ParseRSAPrivateKeyFromPEM(priv_key)
	assert.Nil(t, err)
	assert.IsType(t, &rsa.PrivateKey{}, key)
}

func TestRSAPrivateEncrypt(t *testing.T) {
	key, _ := ParseRSAPrivateKeyFromPEM(priv_key)
	en, err := RSAPrivateEncrypt(key, []byte("Hello"))
	assert.Nil(t, err)
	assert.NotEqual(t, []byte(""), en)
	assert.NotEqual(t, 0, len(en))
}

func TestRSAPublicDecrypt(t *testing.T) {
	priv, _ := ParseRSAPrivateKeyFromPEM(priv_key)
	en, _ := RSAPrivateEncrypt(priv, []byte("Hello"))
	pubkey, _ := ParseRSAPublicKeyFromPEM(pub_key)
	de, _ := base64.StdEncoding.DecodeString(en)
	dedata, err := RSAPublicDecrypt(pubkey, de)
	assert.Nil(t, err)
	assert.Equal(t, "Hello", string(dedata))
}

func TestPublicKeyEncryptAndPrivKeyDecrypt(t *testing.T) {
	priv, _ := ParseRSAPrivateKeyFromPEM(priv_key)
	pubkey, _ := ParseRSAPublicKeyFromPEM(pub_key)
	data := "++++secret_data _____"

	en, err := RSAPublicEncrypt(pubkey, []byte(data))
	assert.Nil(t, err)
	assert.NotEqual(t, data, en)
	in, _ := base64.StdEncoding.DecodeString(en)
	de, err := RSAPrivateDecrypt(priv, in)
	assert.Nil(t, err)
	assert.Equal(t, data, string(de))
}

func TestParseRSAPrivateKeyEnAndPubKeyDecrypt(t *testing.T) {
	priv, _ := ParseRSAPrivateKeyFromPEM(priv_key)
	pubkey, _ := ParseRSAPublicKeyFromPEM(pub_key)
	data := "++++secret_data _____"

	en, err := RSAPrivateEncrypt(priv, []byte(data))
	assert.Nil(t, err)
	assert.NotEqual(t, data, en)
	in, _ := base64.StdEncoding.DecodeString(en)
	de, err := RSAPublicDecrypt(pubkey, in)
	assert.Nil(t, err)
	assert.Equal(t, data, string(de))
}

//func TestFromFrontData(t *testing.T)  {
//	data := ""
//	pubk := `-----BEGIN RSA PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzDhi8IkY+lk8zLmlXrgm
//jxmb94D3HRhPhXO+FfKC/BQhi8iq5Lw+fvcWoQ+uU6V1NJ8a6xaolawjpElarAsn
///RCaTF2X4lIjau7H22oWdePc47rZ2yyV5VCzBzP+meEJZpjerHnXXqQvqOZ51tHr
//p4L/GNSWge9TFjFyszeabe0hZXqbwk3NVSFRyAiXVA5/MHtI3XjmInaSRjN8nBRw
//WE8OkJEN8tB580s9iJwrU1Q6UW0I6NMvEUDpRz5IzwcWJzRYtOqLUkd2IcINUOmp
//WzG5uCUk1rT5VUBJkaXgsvTg9FTISEzTaEs5XVfSIlnGaIari3z3NLeqLAFTzIVm
//TQIDAQAB
//-----END RSA PUBLIC KEY-----
//`
//	k,_ := ParseRSAPublicKeyFromPEM(pubk)
//	de,_ := RSAPublicDecrypt(k, []byte(data))
//
//	fmt.Printf("%#v", de)
//}
