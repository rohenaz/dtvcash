package res

import (
	"fmt"
	"strings"
)

var JsFiles = []string{
	"lib/jquery.min.js",
	"lib/jquery-ui.min.js",
	"lib/pnglib.js",
	"lib/identicon.js",
	"lib/jstz.min.js",
	"lib/bootstrap.min.js",
	"js/init.js",
	"js/login.js",
	"js/signup.js",
	"js/key.js",
	"js/memo.js",
	"js/topics.js",
	"js/profile.js",
	"js/poll.js",
	"js/vote.js",
}

var CssFiles = []string{
	"lib/jquery-ui.min.css",
	"lib/bootstrap.min.css",
	"style.css",
}

var MinJsFile = "res/js/min.js"

var appendNumber = 0

func SetAppendNumber(num int) {
	appendNumber = num
}

func GetResCssFiles() []string {
	var fileList []string
	for _, file := range CssFiles {
		if strings.HasPrefix(file, "http") {
			continue
		}
		fileList = append(fileList, fmt.Sprintf("%s?ver=%d", file, appendNumber))
	}
	return fileList
}

func GetResJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if strings.HasPrefix(file, "http") {
			continue
		}
		fileList = append(fileList, fmt.Sprintf("%s?ver=%d", file, appendNumber))
	}
	return fileList
}

func GetMinJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if ! strings.HasPrefix(file, "http") {
			continue
		}
		fileList = append(fileList, fmt.Sprintf("%s?ver=%d", file, appendNumber))
	}
	fileList = append(fileList, MinJsFile)
	return fileList
}
