package azureauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	loginBaseUrl = "https://login.microsoftonline.com"
	loginScope   = "https://graph.microsoft.com/.default"
)

var tokenCache Token

type AuthClient struct {
	TenantId     string
	clientId     string
	clientSecret string
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	issueTime   time.Time
}

// NewAuthClient Initialize a new AuthClient using azure credentials
func NewAuthClient(tenantId string, clientId string, clientSecret string) *AuthClient {
	return &AuthClient{TenantId: tenantId, clientId: clientId, clientSecret: clientSecret}
}

// getToken give you a access token with still at least 5 minutes of validity.
// Multiple consecutive call will give you the same token.
func (a *AuthClient) GetToken() (string, error) {
	expirationTime := tokenCache.issueTime.Add(time.Duration(tokenCache.ExpiresIn) * time.Second)
	if tokenCache == (Token{}) || expirationTime.Sub(time.Now()) < 5 {
		log.Println("Refreshing Azure token")
		authUrl := fmt.Sprintf("%s/%s/oauth2/v2.0/token", loginBaseUrl, a.TenantId)
		resp, err := http.PostForm(authUrl, url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {a.clientId},
			"scope":         {loginScope},
			"client_secret": {a.clientSecret},
		})
		if err != nil {
			return "", err
		}
		newToken := Token{issueTime: time.Now()}
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return "", errors.New(string(body))
		}
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(body, &newToken); err != nil {
			return "", err
		}
		log.Printf("Received token: %s...\n", newToken.AccessToken[0:20])
		tokenCache = newToken
	}
	return tokenCache.AccessToken, nil
}
