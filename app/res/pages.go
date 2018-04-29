package res

import "github.com/jchavannes/jgo/web"

func SetPageAndOffset(r *web.Response, offset int) {
	var prevOffset int
	if offset > 25 {
		prevOffset = offset - 25
	}
	page := offset / 25 + 1
	r.Helper["PrevOffset"] = prevOffset
	r.Helper["NextOffset"] = offset + 25
	r.Helper["Page"] = page
}
