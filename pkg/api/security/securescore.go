package security

import (
	"encoding/json"
	"fmt"
	"github.com/PPACI/microsoft-defender-ATP-exporter/pkg/azureauth"
	"io/ioutil"
	"net/http"
)

type secureScoreApiAnswer struct {
	Value []SecureScoreValue
}

type SecureScoreValue struct {
	ActiveUserCount   int
	LicensedUserCount int
	CurrentScore      float64
	MaxScore          float64
}

func GetSecureScore(authClient *azureauth.AuthClient) ([]SecureScoreValue, error) {
	accessToken, err := authClient.GetToken()
	if err != nil {
		return []SecureScoreValue{}, err
	}
	req, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/security/secureScores?$top=1", nil)
	if err != nil {
		return []SecureScoreValue{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return []SecureScoreValue{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []SecureScoreValue{}, err
	}
	if resp.StatusCode != 200 {
		return []SecureScoreValue{}, fmt.Errorf(string(body))
	}
	data := secureScoreApiAnswer{}
	if err := json.Unmarshal(body, &data); err != nil {
		return []SecureScoreValue{}, err
	}

	return data.Value, nil
}
