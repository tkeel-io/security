# tkeel security

```
func init() {
   idp := Provider{
      Issuer:       "http://139.198.108.153:8081/auth/realms/master",
      ClientID:     "tkeelOIDC",
      ClientSecret: "961541d2-c739-4938-9f9a-cbc57a5d037d",
      Scopes:       []string{"openid", "profile", "email"},
      RedirectURL:  "http://localhost:8080/v1/oauth/OIDC/2/callback",
   }
   RegisterOIDCProvider("2", idp)
}
```