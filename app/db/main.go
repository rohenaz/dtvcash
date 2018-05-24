package db

import (
	"errors"
	"fmt"
	"github.com/jchavannes/gorm"
	_ "github.com/jchavannes/gorm/dialects/mysql"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/config"
	"reflect"
	"strings"
	"unicode"
)

const (
	BlockTable          = "Block"
	KeyTable            = "Key"
	TxInTable           = "TxIn"
	TxInKeyTable        = "TxIn.Key"
	TxInTxnOutTable     = "TxIn.TxnOut"
	TxInTxnOutTxnTable  = "TxIn.TxnOut.Transaction"
	TxOutTable          = "TxOut"
	TxOutKeyTable       = "TxOut.Key"
	TxOutTxnInTable     = "TxOut.TxnIn"
	TxOutTxnInTxnTable  = "TxOut.TxnIn.Transaction"
	TransactionTable    = "Transaction"
	TransactionBlockTbl = "Transaction.Block"
)

var conn *gorm.DB

const alreadyExistsErrorMessage = "record already exists"

var alreadyExistsError = errors.New(alreadyExistsErrorMessage)

var dbInterfaces = []interface{}{
	User{},
	Session{},
	CsrfToken{},
	Key{},
	Block{},
	Transaction{},
	TransactionIn{},
	TransactionOut{},
	Peer{},
	MemoTest{},
	MemoPost{},
	MemoSetName{},
	MemoFollow{},
	MemoLike{},
	NodeStatus{},
	MemoSetProfile{},
	UserSettings{},
	Notification{},
	UserAction{},
	MemoPollQuestion{},
	MemoPollOption{},
	MemoPollVote{},
}

func getDb() (*gorm.DB, error) {
	if conn == nil {
		conf := config.GetMysqlConfig()
		var err error
		connectionString := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ")/" + conf.Database + "?parseTime=true&loc=Local"
		conn, err = gorm.Open("mysql", connectionString)
		conn.LogMode(false)
		if err != nil {
			return conn, jerr.Get(fmt.Sprintf("failed to connect to database (host: %s)", conf.Host), err)
		}
		for _, dbInterface := range dbInterfaces {
			result := conn.AutoMigrate(dbInterface)
			if result.Error != nil {
				return result, result.Error
			}
		}
	}
	return conn, nil
}

func IsRecordNotFoundError(e error) bool {
	return hasError(e, "record not found")
}

func IsNoRowsInResultSetError(e error) bool {
	return hasError(e, "sql: no rows in result set")
}

func IsAlreadyExistsError(e error) bool {
	return hasError(e, alreadyExistsErrorMessage)
}

func IsDuplicateEntryError(e error) bool {
	return hasErrorPart(e, "Duplicate entry")
}

func hasError(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(jerr.JError)
	if !ok {
		return e.Error() == s
	}
	for _, errMessage := range err.Messages {
		if errMessage == s {
			return true
		}
	}
	return false
}

func hasErrorPart(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(jerr.JError)
	if !ok {
		return strings.Contains(e.Error(), s)
	}
	for _, errMessage := range err.Messages {
		if strings.Contains(errMessage, s) {
			return true
		}
	}
	return false
}

func create(value interface{}) error {
	db, err := getDb()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Create(value)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}

func Find(out interface{}, where ...interface{}) error {
	return find(out, where...)
}

func find(out interface{}, where ...interface{}) error {
	db, err := getDb()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	result := db.Find(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}

func findPreloadColumns(columns []string, out interface{}, where ...interface{}) error {
	db, err := getDb()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	for _, column := range columns {
		db = db.Preload(column)
	}
	result := db.Find(out, where...)
	if result.Error != nil {
		return jerr.Get("error running query", result.Error)
	}
	return nil
}

func Save(value interface{}) error {
	db, _ := getDb()
	if db.Error != nil {
		return jerr.Get("error getting db", db.Error)
	}
	result := db.Save(value)
	if result.Error != nil {
		return jerr.Get("error saving value", result.Error)
	}
	return nil
}

func save(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db.Error != nil {
		return db
	}
	result := db.Save(value)
	return result
}

func remove(value interface{}) *gorm.DB {
	db, _ := getDb()
	if db.Error != nil {
		return db
	}
	result := db.Delete(value)
	return result
}

func count(where interface{}) (uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	var totalCount uint
	result := db.Model(where).Where(where).Count(&totalCount)
	if result.Error != nil {
		return 0, jerr.Get("error running query", result.Error)
	}
	return totalCount, nil
}

func getColumnName(value interface{}) string {
	return reflect.TypeOf(value).Name()
}

func getIdColumnName(value interface{}) string {
	return strings.ToLower(ToSnake(reflect.TypeOf(value).Name())) + "_id"
}

func getWhereInColumn(value interface{}) string {
	return getIdColumnName(value) + " in (?)"
}

func getArrayColumnName(value interface{}) string {
	name := getColumnName(value)
	if strings.HasSuffix(name, "s") {
		return name + "es"
	} else {
		return name + "s"
	}
}

func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
