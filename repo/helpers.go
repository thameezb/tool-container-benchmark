package repo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Insert is a helper to insert a row in a table and set the new ID.
func insert(db *sqlx.DB, table string, v interface{}) error {
	var (
		idField  *reflect.Value
		elem     = reflect.ValueOf(v).Elem()
		names    = []string{}
		values   = []interface{}{}
		bindvars = []string{}
	)

	for i := 0; i < elem.NumField(); i++ {
		valueField := elem.Field(i)
		typeField := elem.Type().Field(i)
		tag, ok := readTags(typeField)
		if ok {
			if tag == "id" {
				idField = &valueField
			} else if valueField.Interface() != nil {
				names = append(names, tag)
				values = append(values, valueField.Interface())
				bindvars = append(bindvars, "?")
			}
		}
	}

	// Check for SQL injection
	if strings.ContainsAny(table, "[]") {
		logrus.Panicf("invalid table name: %s", table)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(names, ", "),
		strings.Join(bindvars, ", "))

	tx := db.MustBegin()
	if _, err := tx.Exec(query, values...); err != nil {
		_ = tx.Rollback()
		return err
	}

	if idField != nil {
		id := new(int64)
		if err := tx.Get(id, `SELECT @@IDENTITY`); err != nil {
			_ = tx.Rollback()
			return err
		}
		idField.SetInt(*id)
	}

	return tx.Commit()
}

func readTags(field reflect.StructField) (string, bool) {
	tags := strings.Split(field.Tag.Get("db"), ",")
	if len(tags) == 0 {
		return "", false
	}
	if len(tags) > 1 {
		for _, t := range tags[1:] {
			if t == "read-only" {
				return "", false
			}
		}
	}

	return tags[0], true
}
