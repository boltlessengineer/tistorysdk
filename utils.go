package tistorysdk

import "fmt"

type UtilsService struct {
	apiClient *Client
}

/* func (us *UtilsService) FindCategoryByName(blogName, categoryName string) ([]Category, error) {
	categories, err := us.apiClient.Category.List(blogName)
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.Name == categoryName {
			return []Category{c}, nil
		}
	}
	return []Category{}, nil
} */

func (us *UtilsService) FindCategoryByLabel(blogName, categoryLabel string) ([]Category, error) {
	categories, err := us.apiClient.Category.List(blogName)
	fmt.Println(categories)
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		fmt.Println(c.Label)
		if c.Label == categoryLabel {
			return []Category{c}, nil
		}
	}
	return []Category{}, nil
}
