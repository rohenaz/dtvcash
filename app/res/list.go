package res

import (
	"strings"
)

var JsFiles = []string{
	"lib/jquery.min.js",
	"lib/jquery-ui.min.js",
	"js/init.js",
	"js/login.js",
	"js/signup.js",
	"js/key.js",
}

var CssFiles = []string{
	"lib/jquery-ui.min.css",
	"style.css",
}

var MinJsFile = "res/js/min.js"

func GetResJsFiles() []string {
	var fileList []string
	for _, file := range JsFiles {
		if strings.HasPrefix(file, "http") {
			continue
		}
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
		fileList = append(fileList, file)
	}
	fileList = append(fileList, MinJsFile)
	return fileList
}
