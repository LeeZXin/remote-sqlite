package sqlite

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

func openAndCheck(path string) (*sql.DB, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.New("no such file or directory")
	}
	return sql.Open("sqlite3", path)
}

func create(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return errors.New("already exists")
	}
	return os.WriteFile(path, []byte(""), os.ModePerm)
}

func joinPath(dataPath, namespace, dbPath string) string {
	return filepath.Join(dataPath, namespace, dbPath)
}

func DeleteNamespace(dataPath, namespace string) error {
	return os.RemoveAll(filepath.Join(dataPath, namespace))
}

func ShowNamespace(dataPath, namespace string) ([]string, error) {
	ret := make([]string, 0)
	err := filepath.WalkDir(filepath.Join(dataPath, namespace), func(_ string, info os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				// ignore the error because the file maybe deleted during traversing.
				return nil
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".db") {
			ret = append(ret, strings.TrimSuffix(info.Name(), ".db"))
		}
		return nil
	})
	return ret, err
}

func CreateDB(dataPath, namespace, dbPath string) error {
	dir := filepath.Join(dataPath, namespace)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return create(joinPath(dataPath, namespace, dbPath))
}

func ExecuteCommand(dataPath, namespace, dbPath, cmd string) (int64, error) {
	db, err := openAndCheck(joinPath(dataPath, namespace, dbPath))
	if err != nil {
		return 0, err
	}
	defer db.Close()
	res, err := db.Exec(cmd)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func QueryCommand(dataPath, namespace, dbPath, cmd string) ([]map[string]string, error) {
	db, err := openAndCheck(joinPath(dataPath, namespace, dbPath))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	ret := make([]map[string]string, 0)
	values := make([][]byte, len(columns))
	scans := make([]any, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	for rows.Next() {
		row := make(map[string]string, len(columns))
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		for i, v := range values {
			row[columns[i]] = string(v)
		}
		ret = append(ret, row)
	}
	return ret, nil
}

func GetDBSize(dataPath, namespace, dbPath string) (int64, error) {
	stat, err := os.Stat(joinPath(dataPath, namespace, dbPath))
	if err != nil {
		return 0, errors.New("no such file or directory")
	}
	return stat.Size(), nil
}

func DropDB(dataPath, namespace, dbPath string) error {
	return os.Remove(joinPath(dataPath, namespace, dbPath))
}
