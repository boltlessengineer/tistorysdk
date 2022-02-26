package tistorysdk

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

func (us *UtilsService) FindCategory(blogName, categoryLabel string) ([]Category, error) {
	categories, err := us.apiClient.Category.List(blogName)
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.Label == categoryLabel {
			return []Category{c}, nil
		}
	}
	return []Category{}, nil
}
