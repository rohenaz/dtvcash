package script

import (
	"bytes"
	"github.com/rohenaz/dtvcash/app/bitcoin/memo"
	"github.com/jchavannes/btcd/txscript"
	"strings"
)

func GetScriptString(pkScript []byte) string {
	if len(pkScript) < 5 || ! bytes.Equal(pkScript[0:3], []byte{
		txscript.OP_RETURN,
		txscript.OP_DATA_2,
		memo.CodePrefix,
	}) {
		return ""
	}
	var start = 3
	if pkScript[3] == memo.CodePost {
		start = 4
	}
	if pkScript[3] == memo.CodeSetName {
		start = 4
	}
	data, err := txscript.PushedData(pkScript[start:])
	if err != nil || len(data) == 0 {
		return ""
	}
	var stringArray []string
	for _, bt := range data {
		stringArray = append(stringArray, string(bt))
	}
	return strings.Join(stringArray, " ")
}

func GetMemoType(pkScript []byte) string {
	if len(pkScript) < 5 || ! bytes.Equal(pkScript[0:3], []byte{
		txscript.OP_RETURN,
		txscript.OP_DATA_2,
		memo.CodePrefix,
	}) {
		return ""
	}
	switch pkScript[3] {
	case memo.CodeTest:
		return "Test"
	case memo.CodePost:
		return "Post"
	case memo.CodeSetName:
		return "Set Name"
	}
	return ""
}
