package internal

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Updates pagination buttons
func handlePagination(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pages := []Pagination{}
	startPage := index - 2
	endPage := index + 2

	if startPage < 1 {
		startPage = 1
	}

	if endPage > queue.Len()/PageSize {
		endPage = queue.Len() / PageSize
	}

	if (index <= 2 || index >= endPage-2 || index == endPage || index == endPage-1) && endPage-startPage+1 < 5 && endPage < queue.Len()/PageSize {
		endPage = startPage + 4
	}

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, Pagination{
			Number: i,
			Active: i == index,
		})
	}

	ctx.HTML(200, "pagination.html", pages)
}
