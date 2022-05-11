package license

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _priv_key = `-----BEGIN RSA PRIVATE KEY-----
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

var _pub_key = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1Fh8AWKSyZKzam044MxT
ga7jY45Yc53zKzgVrFM6K8vZje14A9QevqQGaI4XOLf4xqHVdETD52R0cZ9QqJ1k
ChpLR2bqPsHUIxavn5DHsV60cYSyegkQ1Tm7Ty6tH3cdqDwH0o5K5e1XDqhE5fLd
r03oPoYme8w+zTvHY/DxvG77LCkQoqD5whjwlbNhsuS1/b/xMy6gElyvtAOC3kRx
1LJ+3tVFOq20qthU0m33Ch0yl4WdzjL353W7+n7ia71aSUkuIyQ8Nm/hUeIjQ/e7
ZROHZMExejp72XWyO3lvudvpsUtln8YAR2M/8nf4gpIpnnmeBX6zTay5JOIlNlpS
vQIDAQAB
-----END RSA PUBLIC KEY-----`

const (
	username = "username"
	alex     = "alex"
)

func TestIssuer(t *testing.T) {
	issuer := NewIssuer(WithSecret(_priv_key), WithInfo(map[string]string{
		username: alex,
	}))

	license, err := issuer.Issue()
	assert.Nil(t, err)
	fmt.Println(license)
}

func TestParser(t *testing.T) {
	license := "bphBbgzmiegpcryAlvfcpAisFlpywxvy002c1EOR+k7hmyX802ZIaexpXuWxjkZpLm31gSIb6zamRQQ=SkT8fCbsu2qtI7a/sjyeXGBRQSQX4hSBePYX/u626xmBiEW5a6G4sWlfd6L/ESWcnSddaNwShJBwIb4ew6Ub0BA2M4FnZUQ2ZjCl+hpuhaOQLp0Ik7MiU8SDgVhNVQ2guXma54SgKfUJ4facfmRDKVYcxdEiJucs0fFnFmk3if1I+E6OxF1nt+SOp2uFRd56upKz4ylBTB6e6NDSw8H9lhSNHbB1xSECaTV9jeu+CcX7JebrurlCXNHuGw5Lyx+8wQM1BDsi90pJzSrSYyYOzRBCUUvd7GXQsY4gTvyov7C0HbObPoWpot4Kjdj8yosBMGLlY5C39MYYxm6FUth35A=="
	parser := NewParser(license, WithPublicKey(_pub_key))
	info, err := parser.Parse()
	fmt.Printf("%+v\n", info)
	assert.Nil(t, err)
	assert.Equal(t, alex, info[username])
}

func TestNewLicense(t *testing.T) {
	license := &License{}
	time := time.Now()

	// Generate License
	// Set Secret with Option WithSecret()
	// Set Info with Option WithInfo().
	license.Set(WithSecret(_priv_key), WithInfo(map[string]string{
		username: alex,
		"time":   time.String(),
	}))
	l, err := license.Issue()
	assert.Nil(t, err)

	// Parse License
	// Set Public Key with Option WithPublicKey()
	// Set Content with Option WithLicense()
	license.Set(WithLicense(l), WithPublicKey(_pub_key))
	info, err := license.Parse()
	assert.Nil(t, err)
	assert.Equal(t, alex, info[username])
	assert.Equal(t, time.String(), info["time"])
}
