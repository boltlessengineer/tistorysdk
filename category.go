package tistorysdk

import (
	"encoding/json"
	"net/http"
)

type CategoryService struct {
	apiClient *Client
}

type Category struct {
	ID      string `json:"id"`      // Actually, Int
	Name    string `json:"name"`    //
	Parent  string `json:"parent"`  //
	Label   string `json:"label"`   // {Parent}/{Name}
	Entries string `json:"entries"` // Actually, Int
}

func (cs *CategoryService) List(blogName string) ([]Category, error) {
	q := map[string]string{
		"access_token": cs.apiClient.accessToken,
		"output":       "json",
		"blogName":     blogName,
	}
	raw, err := cs.apiClient.request(http.MethodGet, "category/list", q, nil)
	if err != nil {
		return nil, err
	}
	var res ListCategoryResponse
	json.Unmarshal(*raw, &res)
	return res.Item.Categories, nil
}

type ListCategoryResponse struct {
	BasicResBody
	Item struct {
		Categories []Category
	} `json:"item"`
}
