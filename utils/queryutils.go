package utils

import (
	"strconv"
	"strings"
)

func AppendWhere(baseQuery *string, baseParam *[]interface{}, appendedQuery string, appendedParam string) {
	if len(appendedParam) > 0 {
		if len(*baseQuery) > 0 {
			*baseQuery += " AND "
		}
		*baseQuery += appendedQuery
		*baseParam = append(*baseParam, appendedParam)
	}
}

func AppendWhereRaw(baseQuery *string, appendedQuery string) {
	if len(appendedQuery) > 0 {
		if len(*baseQuery) > 0 {
			*baseQuery += " AND "
		}
		*baseQuery += appendedQuery
	}
}

func AppendWhereLike(baseQuery *string, baseParam *[]interface{}, appendedQuery string, appendedParam string) {
	if len(appendedParam) > 0 {
		if len(*baseQuery) > 0 {
			*baseQuery += " AND "
		}
		*baseQuery += appendedQuery
		*baseParam = append(*baseParam, "%"+appendedParam+"%")
	}
}

func AppendOrderBy(baseQuery *string, orderBy string, orderDirection string) {
	if len(orderBy) > 0 {
		*baseQuery += " ORDER BY " + orderBy
		if len(orderDirection) > 0 {
			if strings.ToLower(orderDirection) == "desc" {
				*baseQuery += " DESC "
			}
		}
	}
}

func AppendComma(baseQuery *string, baseParam *[]interface{}, appendedQuery string, value string) {
	if len(*baseQuery) > 0 {
		*baseQuery += " , "
	}

	if len(value) > 0 {
		*baseQuery += appendedQuery
		*baseParam = append(*baseParam, value)
	} else {
		*baseQuery += strings.ReplaceAll(appendedQuery, "?", "NULL")
	}
}

func AppendCommaNotNull(baseQuery *string, baseParam *[]interface{}, appendedQuery string, value string) {
	if len(*baseQuery) > 0 {
		*baseQuery += " , "
	}

	*baseQuery += appendedQuery
	*baseParam = append(*baseParam, value)
}

func AppendCommaRaw(baseQuery *string, appendedQuery string) {
	if len(appendedQuery) > 0 {
		if len(*baseQuery) > 0 {
			*baseQuery += " , "
		}
		*baseQuery += appendedQuery
	}
}

func AppendLimit(baseQuery *string, pageParam string, perPageParam string) {
	page, _ := strconv.Atoi(pageParam)
	page = GetPageValue(page)
	perPage, _ := strconv.Atoi(perPageParam)
	perPage = GetPerPageValue(perPage)
	offset := (page - 1) * perPage
	*baseQuery += " LIMIT " + strconv.Itoa(perPage) + " OFFSET " + strconv.Itoa(offset)
}

func GetPerPageValue(perPage int) int {
	if perPage == 0 {
		perPage = 10
	}
	return perPage
}

func GetPageValue(page int) int {
	if page == 0 {
		page = 1
	}
	return page
}
