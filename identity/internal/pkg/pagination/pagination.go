package pagination

import (
	identity_core_go "github.com/agntcy/identity/internal/pkg/generated/agntcy/identity/core/v1alpha1"
)

// Creates a PagedResponse object
func ConvertToPagedResponse[T any](
	paginationFilter PaginationFilter,
	items *Pageable[T],
) *identity_core_go.PagedResponse {
	var nextPage *int32
	hasNextPage := int64(paginationFilter.GetPage())*int64(paginationFilter.GetLimit()) < items.Total
	if hasNextPage {
		n := paginationFilter.GetPage() + 1
		nextPage = &n
	}

	return &identity_core_go.PagedResponse{
		HasNextPage: &hasNextPage,
		NextPage:    nextPage,
		Total:       items.Total,
		Size:        items.Size,
	}
}
