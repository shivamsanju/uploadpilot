package utils

import (
	"net/http"
	"strconv"
)

func GetSkipLimitSearchParams(r *http.Request) (skip int64, limit int64, search string, err error) {
	query := r.URL.Query()
	s := query.Get("skip")
	l := query.Get("limit")
	search = query.Get("search")

	skip, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	limit, err = strconv.ParseInt(l, 10, 64)
	if err != nil {
		return
	}
	return
}
