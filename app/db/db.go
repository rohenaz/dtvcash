package db

import (
	"fmt"
	"git.jasonc.me/main/memo/app/config"
	"github.com/jchavannes/gorm"
	_ "github.com/jchavannes/gorm/dialects/mysql"
	"github.com/jchavannes/jgo/jerr"
	"reflect"
	"strings"
	"unicode"
)

var conn *gorm.DB

var (
	dbInterfaces = []interface{}{
		User{},
		Session{},
		Settings{},
	}
)

func getDb() (*gorm.DB, error) {
	if conn == nil {
		conf := config.GetMysqlConfig()
		var err error
		connectionString := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ")/" + conf.Database + "?parseTime=true"
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

func isRecordNotFoundError(e error) bool {
	err, ok := e.(jerr.JError)
	if !ok {
		return e.Error() == "record not found"
	}
	for _, errMessage := range err.Messages {
		if errMessage == "record not found" {
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

func count(value interface{}) (uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, jerr.Get("error getting db", err)
	}
	var totalCount uint
	result := db.Model(value).Where(value).Count(&totalCount)
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
