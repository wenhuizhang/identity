package pagination

import (
	pyramid_shared_go "github.com/agntcy/pyramid/internal/pkg/generated/agntcy/pyramid/shared/v1alpha1"
)

// Creates a PagedResponse object
func ConvertToPagedResponse[T any](
	paginationFilter PaginationFilter,
	items *Pageable[T],
) *pyramid_shared_go.PagedResponse {
	var nextPage *int32
	hasNextPage := int64(paginationFilter.GetPage())*int64(paginationFilter.GetLimit()) < items.Total
	if hasNextPage {
		n := paginationFilter.GetPage() + 1
		nextPage = &n
	}

	return &pyramid_shared_go.PagedResponse{
		HasNextPage: &hasNextPage,
		NextPage:    nextPage,
		Total:       items.Total,
		Size:        items.Size,
	}
}
