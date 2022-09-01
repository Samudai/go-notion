package notion

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// AuthResponse is the response from the notion authentication endpoint.
type AuthResponse struct {
	AccessToken   string  `json:"access_token"`
	TokenType     string  `json:"token_type"`
	BotID         string  `json:"bot_id"`
	WorkspaceName *string `json:"workspace_name"`
	WorkspaceIcon *string `json:"workspace_icon"`
	WorkspaceID   string  `json:"workspace_id"`
	Owner         Owner   `json:"owner"`
}

// Owner is the owner of a notion workspace.
type Owner struct {
	Type string `json:"type"`
	User User   `json:"user"`
}

// NotionError is the error response from the Notion API.
type NotionError struct {
	Error string `json:"error"`
}

// NotionAuthPayload is the payload for the notion authentication endpoint.
type NotionAuthPayload struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// AuthUser authenticates with notion and returns an user access token.
func AuthUser(code string) (*AuthResponse, error) {
	url := baseURL + "/oauth/token"

	clientID := os.Getenv("NOTION_CLIENT_ID")
	clientSecret := os.Getenv("NOTION_CLIENT_SECRET")
	redirectURI := os.Getenv("NOTION_REDIRECT_URI")

	notionAuthPayload := NotionAuthPayload{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: redirectURI,
	}

	payload := &bytes.Buffer{}
	err := json.NewEncoder(payload).Encode(notionAuthPayload)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to encode filter to JSON: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	auth := base64.URLEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var resp NotionError
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", resp.Error)
	}

	var resp AuthResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
