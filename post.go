package tistorysdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	CategoryID string `json:"categoryId"`
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
	fmt.Println(ps.apiClient.accessToken)
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

func (ps *PostService) Write(blogName string, post *Post) (*WritePostResponse, error) {
	q := map[string]string{}
	b := map[string]string{
		"access_token":  ps.apiClient.accessToken,
		"output":        "json",
		"blogName":      blogName,
		"title":         post.Title,
		"content":       post.Content,
		"visibility":    post.Visibility.String(),
		"category":      post.CategoryID,
		"published":     fmt.Sprint(time.Time(post.Date).Unix()),
		"slogan":        post.Slogan,
		"tag":           strings.Join(post.Tags.Tag, ","),
		"acceptComment": post.AcceptComment.String(),
		"password":      post.Password,
	}
	raw, err := ps.apiClient.request(http.MethodPost, "post/write", q, b)
	if err != nil {
		return nil, err
	}
	var res WritePostResponse
	if err = json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (ps *PostService) Update(blogName string, post *Post) (*WritePostResponse, error) {
	fmt.Println("----------------------------")
	fmt.Println("Editing post id:", post.ID.String())
	fmt.Println("----------------------------")
	b := map[string]string{
		"access_token":  ps.apiClient.accessToken,
		"output":        "json",
		"blogName":      blogName,
		"postId":        post.ID.String(),
		"title":         post.Title,
		"content":       post.Content,
		"visibility":    post.Visibility.String(),
		"category":      post.CategoryID,
		"published":     fmt.Sprint(time.Time(post.Date).Unix()),
		"slogan":        post.Slogan,
		"tag":           strings.Join(post.Tags.Tag, ","),
		"acceptComment": post.AcceptComment.String(),
		"password":      post.Password,
	}
	raw, err := ps.apiClient.request(http.MethodPost, "post/modify", nil, b)
	if err != nil {
		return nil, err
	}
	var res WritePostResponse
	if err = json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type WritePostResponse struct {
	Status SInt   `json:"status"`
	PostID SInt   `json:"postId"`
	Url    string `json:"url"`
}

func (ps *PostService) UploadFile(blogName string, filename string, filebuf io.Reader) (*UploadFileResponse, error) {
	url := "https://www.tistory.com/apis/post/attach"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("access_token", ps.apiClient.accessToken)
	writer.WriteField("output", "json")
	writer.WriteField("blogName", blogName)
	part, err := writer.CreateFormFile("uploadedfile", filename)
	io.Copy(part, filebuf)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r struct {
		Tistory UploadFileResponse `json:"tistory"`
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error uploading file to Tistory :%s", res.Status)
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r.Tistory, nil
}

type UploadFileResponse struct {
	Url      string `json:"url"`
	Replacer string `json:"replacer"`
}
