package tistorysdk

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	PostMajor
	Content       string `json:"content"`
	AcceptComment SInt   `json:"acceptComment"`
	Tags          Tags   `json:"tags"`
	Date          SDate  `json:"date"`

	Slogan   string
	Password string
}

type PostMajor struct {
	ID         SInt   `json:"id"`
	Title      string `json:"title"`
	Url        string `json:"postUrl"`
	Visibility SInt   `json:"visibility"`
	CategoryID SInt   `json:"categoryId"`
	Comments   SInt   `json:"comments"`
	Trackbacks SInt   `json:"trackbacks"`
	DateStr    string `json:"date"`
}

type Tags struct {
	Tag []string `json:"tag"`
}

type PostService struct {
	apiClient *Client
}

func (ps *PostService) List(blogName string, pageNum int) (*ListPostResItem, error) {
	q := map[string]string{
		"access_token": ps.apiClient.accessToken,
		"output":       "json",
		"blogName":     blogName,
		"page":         strconv.Itoa(pageNum),
	}
	raw, err := ps.apiClient.request(http.MethodGet, "post/list", q, nil)
	if err != nil {
		return nil, err
	}
	var res ListPostResponse
	json.Unmarshal(*raw, &res)
	return &res.Item, nil
}

type ListPostResponse struct {
	BasicResBody
	Item ListPostResItem `json:"item"`
}

type ListPostResItem struct {
	Url          string      `json:"url"`
	SecondaryUrl string      `json:"secondaryUrl"`
	Page         SInt        `json:"page"`
	Count        SInt        `json:"count"`
	TotalCount   SInt        `json:"totalCount"`
	Posts        []PostMajor `json:"posts"`
}

func (ps *PostService) Read() (*Post, error) {
	return nil, nil
}

type ReadPostResponse struct {
	BasicResBody
	Item ReadPostResItem `json:"item"`
}

type ReadPostResItem struct {
	Url          string `json:"url"`
	SecondaryUrl string `json:"secondaryUrl"`
	Post
}

func (ps *PostService) Write(token string, blogName string, post *Post) (*UpdatePostResponse, error) {
	q := map[string]string{
		"access_token": token,
		"output":       "json",
		"blogName":     blogName,
		"title":        post.Title,
		"content":      post.Content,
		"visibility":   post.Visibility.String(),
		"category":     post.CategoryID.String(),
		// "published":
		"slogan":        post.Slogan,
		"tag":           strings.Join(post.Tags.Tag, ","),
		"acceptComment": post.AcceptComment.String(),
		"password":      post.Password,
	}
	raw, err := ps.apiClient.request(http.MethodPost, "post/write", q, nil)
	if err != nil {
		return nil, err
	}
	var res UpdatePostResponse
	if err = json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}
	return &res, err
}

func (ps *PostService) Update(blogName string, post *Post) (*UpdatePostResponse, error) {
	return nil, nil
}

type UpdatePostResponse struct {
	Status SInt   `json:"status"`
	PostID SInt   `json:"postId"`
	Url    string `json:"url"`
}

func (ps *PostService) UploadFile() (*UploadFileResponse, error) {
	return nil, nil
}

type UploadFileResponse struct {
	Status   SInt   `json:"status"`
	Url      string `json:"url"`
	Replacer string `json:"replacer"`
}
