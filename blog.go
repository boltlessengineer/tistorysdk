package tistorysdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func request(method, reqUrl, token string, queryParams, reqBody map[string]string) (*json.RawMessage, error) {
	baseUrl := "https://www.tistory.com/apis/"
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	u.Parse(reqUrl)
	fmt.Println("tistorysdk request(testing...):", u.String())
	form := url.Values{}
	if reqBody != nil && len(reqBody) > 0 {
		for key, val := range reqBody {
			form.Add(key, val)
		}
	}
	if queryParams != nil && len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		q.Add("access_token", token)
		q.Add("output", "json")
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequest(method, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var bs BasicResponse
	if err := json.NewDecoder(res.Body).Decode(&bs); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		var errbody BasicResBody
		if err := json.Unmarshal(bs.Tistory, &errbody); err != nil {
			b, _ := io.ReadAll(res.Body)
			return nil, errors.New("Can't json unmarshal response.\nRaw response :" + string(b))
		}
		return nil, errors.New(errbody.ErrorMessage)
	}
	return &bs.Tistory, nil
}

func GetBlogsInfo(token string) (*BlogsInfoRes, error) {
	raw, err := request(http.MethodGet, "blog/info", token, nil, nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		BasicResBody
		Item BlogsInfoRes `json:"item"`
	}
	err = json.Unmarshal(*raw, &res)
	return &res.Item, nil
}

type BlogsInfoRes struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Blogs  []Blog `json:"blogs"`
}

type Blog struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	SecondaryUrl string `json:"secondaryUrl"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	// TODO: add more
}
