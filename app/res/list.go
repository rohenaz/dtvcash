package res

import (
	"fmt"
	"regexp"
	"strings"
)

var JsFiles = []string{
	"lib/jquery.min.js",
	"lib/jquery-ui.min.js",
	"js/init.js",
	"js/login.js",
	"js/signup.js",
	"js/key.js",
	"js/memo.js",
}

var CssFiles = []string{
	"lib/jquery-ui.min.css",
	"style.css",
}

var MinJsFile = "res/js/min.js"

const AppendNumber = 6

func GetResJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if strings.HasPrefix(file, "http") {
			continue
		}
		re := regexp.MustCompile(`(.*)(\.js)`)
		file = re.ReplaceAllString(file, fmt.Sprintf("$1-%d$2", AppendNumber))
		fileList = append(fileList, file)
	}
	return fileList
}

func GetMinJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if ! strings.HasPrefix(file, "http") {
			continue
		}
		re := regexp.MustCompile(`(.*)(\.css)`)
		file = re.ReplaceAllString(file, fmt.Sprintf("$1-%d$2", AppendNumber))
		fileList = append(fileList, file)
	}
	fileList = append(fileList, MinJsFile)
	return fileList
}
