package controllers

import (
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/service/activity"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/gin-gonic/gin"
)

// ActivityListHandler : Show a list of activity
func ActivityListHandler(c *gin.Context) {
	page := c.Query("page")
	pagenum := 1
	offset := 100
	userid := c.Query("userid")
	filter := c.Query("filter")

	var err error
	currentUser := getUser(c)
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions []string
	var values []interface{}
	if userid != "" && userPermission.HasAdmin(currentUser) {
		conditions = append(conditions, "user_id = ?")
		values = append(values, userid)
	}
	if filter != "" {
		conditions = append(conditions, "filter = ?")
		values = append(values, filter)
	}

	activities, nbActivities := activity.GetAllActivities(offset, (pagenum-1)*offset, strings.Join(conditions, " AND "), values...)

	nav := navigation{nbActivities, offset, pagenum, "activities"}
	modelList(c, "site/torrents/activities.jet.html", activities, nav, newSearchForm(c))
}
