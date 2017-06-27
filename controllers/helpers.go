package controllers

import (
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/gin-gonic/gin"
)

type navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

type searchForm struct {
	common.SearchParam
	Category         string
	ShowItemsPerPage bool
	SizeType         string
	DateType         string
	MinSize          string
	MaxSize          string
	FromDate         string
	ToDate           string
}

// Some Default Values to ease things out
func newNavigation() navigation {
	return navigation{
		MaxItemPerPage: 50,
	}
}

func newSearchForm(c *gin.Context) searchForm {
	sizeType := c.DefaultQuery("sizeType", "m")
	return searchForm{
		Category:         "_",
		ShowItemsPerPage: true,
		SizeType:         sizeType,
		DateType:         c.Query("dateType"),
		MinSize:          c.Query("minSize"),  // We need to overwrite the value here, since size are formatted
		MaxSize:          c.Query("maxSize"),  // We need to overwrite the value here, since size are formatted
		FromDate:         c.Query("fromDate"), // We need to overwrite the value here, since we can have toDate instead and date are formatted
		ToDate:           c.Query("toDate"),   // We need to overwrite the value here, since date are formatted
	}
}
func getUser(c *gin.Context) *model.User {
	user, _, _ := userService.RetrieveCurrentUser(c)
	return &user
}
