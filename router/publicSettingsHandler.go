package router

import (
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/user"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
	"github.com/gin-gonic/gin"
)

// LanguagesJSONResponse : Structure containing all the languages to parse it as a JSON response
type LanguagesJSONResponse struct {
	Current   string            `json:"current"`
	Languages map[string]string `json:"languages"`
}

// SeePublicSettingsHandler : Controller to view the languages and themes
func SeePublicSettingsHandler(c *gin.Context) {
	_, Tlang := publicSettings.GetTfuncAndLanguageFromRequest(c)
	availableLanguages := publicSettings.GetAvailableLanguages()

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, LanguagesJSONResponse{Tlang.Tag, availableLanguages})
	} else {
		changeLanguageTemplate(c, Tlang.Tag, availableLanguages)
	}
}

// ChangePublicSettingsHandler : Controller for changing the current language and theme
func ChangePublicSettingsHandler(c *gin.Context) {
	theme := c.PostForm("theme")
	lang := c.PostForm("language")
	mascot := c.PostForm("mascot")
	mascotURL := c.PostForm("mascot_url")

	messages := msg.GetMessages(c)

	availableLanguages := publicSettings.GetAvailableLanguages()

	if _, exists := availableLanguages[lang]; !exists {
		messages.AddErrorT("errors", "language_not_available")
	}
	// FIXME Are the settings actually sanitized?
	// Limit the mascot URL, so base64-encoded images aren't valid
	if len(mascotURL) > 256 {
		messages.AddErrorT("errors", "mascot_url_too_long")
	}

	_, err := url.Parse(mascotURL)
	if err != nil {
		messages.AddErrorTf("errors", "mascor_url_parse_error", err.Error())
	}

	// If logged in, update user settings.
	user, _ := userService.CurrentUser(c)
	if user.ID > 0 {
		user.Language = lang
		user.Theme = theme
		user.Mascot = mascot
		user.MascotURL = mascotURL
		userService.UpdateRawUser(&user)
	}
	// Set cookie with http and not gin for expires (maxage not supported in <IE8)
	http.SetCookie(c.Writer, &http.Cookie{Name: "lang", Value: lang, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot", Value: mascot, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot_url", Value: mascotURL, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})

	c.Redirect(http.StatusSeeOther, "/")
}
func getDomainName() string {
	domain := config.Conf.Cookies.DomainName
	if config.Conf.Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
