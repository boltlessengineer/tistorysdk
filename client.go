package tistorysdk

import (
	"encoding/json"
	"errors"
	"fmt"

	//"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	apiURL   = "https://www.tistory.com/apis/"
	oauthURL = "https://www.tistory.com/oauth"
)

type Client struct {
	httpClient *http.Client
	baseUrl    *url.URL

	ID          string
	SK          string
	RedirectUri string

	accessToken string

	Post     PostService
	Category CategoryService
	Comment  CommentService
}

func NewClient(clientID string, clientSK string, redirect string) *Client {
	u, err := url.Parse(apiURL)
	if err != nil {
		panic(err)
	}
	c := &Client{
		httpClient:  http.DefaultClient,
		baseUrl:     u,
		ID:          clientID,
		SK:          clientSK,
		RedirectUri: redirect,
	}

	c.Post = PostService{apiClient: c}
	c.Category = CategoryService{apiClient: c}
	c.Comment = CommentService{apiClient: c}

	return c
}

// GetAuthCodeURL() returns URL to request Authorization Code
func (c *Client) GetAuthCodeURL(state string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s/%s", oauthURL, "authorize"))
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Add("client_id", c.ID)
	q.Add("redirect_uri", c.RedirectUri)
	q.Add("response_type", "code")
	q.Add("state", state)
	u.RawQuery = q.Encode()
	return u
}

const (
	AccessTokenUrl string = "https://www.tistory.com/oauth/access_token"
)

func GetToken(clientID, clientSK, redirectUri, code string) (string, error) {
	u, err := url.Parse(AccessTokenUrl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSK)
	q.Add("redirect_uri", redirectUri)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	s := string(b)
	const (
		token_format string = "access_token=%s"
		error_format string = "error=%s&error_description=%s"
	)
	var token string
	if _, err := fmt.Sscanf(s, token_format, &token); err != nil {
		var e, ed string
		if _, err := fmt.Sscanf(s, error_format, &e, &ed); err != nil {
			if err != nil {
				return "", err
			}
			return "", errors.New(fmt.Sprintf("[Error - %s] %s", e, ed))
		}
	}

	return token, nil
}

// GetToken retrieves Access Token with Authorization Code
func (c *Client) SetToken(code string) error {
	u, err := url.Parse(AccessTokenUrl)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Add("client_id", c.ID)
	q.Add("client_secret", c.SK)
	q.Add("redirect_uri", c.RedirectUri)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	s := string(b)
	const (
		token_format string = "access_token=%s"
		error_format string = "error=%s&error_description=%s"
	)
	var token string
	if _, err := fmt.Sscanf(s, token_format, &token); err != nil {
		var e, ed string
		if _, err := fmt.Sscanf(s, error_format, &e, &ed); err != nil {
			if err != nil {
				return err
			}
			return errors.New(fmt.Sprintf("[Error - %s] %s", e, ed))
		}
	}
	c.accessToken = token

	return nil
}

func (c *Client) request(method string, urlStr string, queryParams map[string]string, requestBody map[string]string) (*json.RawMessage, error) {
	u, err := c.baseUrl.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.String())

	var buf io.ReadWriter
	// TODO - write request body to buf

	if len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r BasicResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var tmp BasicResBody
		if err := json.Unmarshal(r.Tistory, &tmp); err != nil {
			b, _ := io.ReadAll(res.Body)
			fmt.Println(string(b))
			return nil, err
		}
		return nil, errors.New(tmp.ErrorMessage)
	}

	return &r.Tistory, nil
}

type BasicResponse struct {
	Tistory json.RawMessage `json:"tistory"`
}

type BasicResBody struct {
	Status       SInt   `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type CategoryService struct {
	apiClient *Client
}

type CommentService struct {
	apiClient *Client
}
