package connection

import (
	"booking/model"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	DB *gorm.DB
}

func NewConnection(dbConfig map[string]string) (*DBConnection, error) {
	var err error
	// Open database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s port=%s TimeZone=Asia/Jakarta", dbConfig["host"], dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"], dbConfig["port"])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Mengaktifkan logger default
	})
	if err != nil {
		return nil, err
	}

	// dbConn.SetMaxOpenConns(10)
	// dbConn.SetConnMaxLifetime(time.Minute * 1)
	// dbConn.SetMaxIdleConns(10)

	// err = db.Ping()
	// if err != nil {
	// 	return nil, err
	// }
	if dbConfig["env"] != "prod" {
		db = db.Debug()
	}
	logrus.Info("Database connection initiated")

	return &DBConnection{
		DB: db,
	}, nil
}

func (c DBConnection) Raw(qry string, args ...interface{}) *gorm.DB {
	logrus.Debugf("Running query: %s", qry)
	return c.DB.Raw(qry, args...)
}

func (c DBConnection) Exec(qry string, args ...interface{}) *gorm.DB {
	logrus.Debugf("Running query: %s", qry)
	return c.DB.Exec(qry, args...)
}

func (c DBConnection) Query(strSQL string, args ...interface{}) (*sql.Rows, error) {
	// if no DBConnection, return
	//
	if c.DB == nil {
		return nil, errors.New("database needs to be initiated first")
	}

	//if strSQL, found = sqlCommandMap[strSQL]; !found {
	// rows, err := c.DB.QueryRow(strSQL, args...)
	// if err != nil {
	// 	return nil, err
	// }
	var rows *sql.Rows
	return rows, nil
}

// SelectQueryByFieldNameSlice parse and gets column value by field name
func (c DBConnection) SelectQueryByFieldNameSlice(strSQL string, args ...interface{}) ([]map[string]string, int, error) {
	var rowret *sql.Rows
	var err error
	if rowret != nil {
		defer rowret.Close()
	}
	if c.DB == nil {
		return nil, 0, errors.New("Please OpenConnection prior Query")
	}

	rowret, err = c.Query(strSQL, args...)
	if err != nil {
		log.Errorf("Error Execute SQL : %s %+v -> %+v", strSQL, args, err)
		return nil, 0, err
	}

	results, rowCount, err := c.GetRowsSlice(rowret)
	if err != nil {
		return nil, 0, err
	}

	return results, rowCount - 1, nil
}

func (c DBConnection) GetRowsSlice(rows *sql.Rows) ([]map[string]string, int, error) {
	var results []map[string]string
	if rows != nil {
		defer rows.Close()
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	counter := 1
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, err
		}

		// initialize the second layer
		single := make(map[string]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			single[columns[i]] = value
		}
		results = append(results, single)
		counter++
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, counter, nil
}

func (c DBConnection) SelectQueryByFieldName(strSQL string, args ...interface{}) (map[int]map[string]string, int, error) {
	var rowret *sql.Rows
	var err error
	if rowret != nil {
		defer rowret.Close()
	}
	if c.DB == nil {
		return nil, 0, errors.New("Please OpenConnection prior Query")
	}

	rowret, err = c.Query(strSQL, args...)
	if err != nil {
		return nil, 0, err
	}

	results, rowCount, err := c.GetRows(rowret)
	if err != nil {
		return nil, 0, err
	}

	return results, rowCount - 1, nil
}

func (c DBConnection) GetRows(rows *sql.Rows) (map[int]map[string]string, int, error) {
	var results map[int]map[string]string
	results = make(map[int]map[string]string)
	if rows != nil {
		defer rows.Close()
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	counter := 1
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, err
		}

		// initialize the second layer
		results[counter] = make(map[string]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[counter][columns[i]] = value
		}
		counter++
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, counter, nil
}

func (c DBConnection) Paginate(filter model.GlobalQueryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// q := r.URL.Query()
		var limit int
		var err error
		if len(filter.Limit) > 0 {
			limit, _ = strconv.Atoi(filter.Limit)
		} else {
			limit = 10
		}

		var page int
		if len(filter.Page) > 0 {
			page, err = strconv.Atoi(filter.Page)
			if err != nil {
				page = 1
			}
		} else {
			page = 1
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func (c DBConnection) Order(filter model.GlobalQueryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// q := r.URL.Query()

		if len(filter.OrderBy) > 0 {
			dir := "ASC"
			if len(filter.OrderDir) > 0 {
				dir = filter.OrderDir
			}
			return db.Order(filter.OrderBy + " " + dir)
		}

		return db
	}
}
