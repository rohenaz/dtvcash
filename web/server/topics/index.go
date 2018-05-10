package topics

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/res"
	"net/http"
	"strings"
)

var indexRoute = web.Route{
	Pattern: res.UrlTopics,
	Handler: func(r *web.Response) {
		preHandler(r)
		offset := r.Request.GetUrlParameterInt("offset")
		searchString := html_parser.EscapeWithEmojis(r.Request.GetUrlParameter("s"))
		topics, err := db.GetUniqueTopics(uint(offset), searchString)
		if err != nil {
			r.Error(jerr.Get("error getting topics from db", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Topics"] = topics
		r.Helper["SearchString"] = searchString
		res.SetPageAndOffset(r, offset)
		if searchString != "" {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?s=%s", strings.TrimLeft(res.UrlTopics, "/"), searchString)
		} else {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?", res.UrlTopics)
		}
		r.Render()
	},
}
