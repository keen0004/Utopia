package tests

import (
	"os"
	"testing"
	"utopia/internal/database"
)

var (
	dbfile = "./test.db"
)

func TestDBOpen(t *testing.T) {
	os.Remove(dbfile)

	db := database.NewDatabase(dbfile)
	if db == nil {
		t.Errorf("Database is nil")
		return
	}
	defer os.Remove(dbfile)

	err := db.Open()
	if err != nil {
		t.Errorf("Open database failed with error: %v", err)
		return
	}

	db.Close()
}

func TestSQL(t *testing.T) {
	db := database.NewDatabase(dbfile)
	if db == nil {
		t.Errorf("Database is nil")
		return
	}

	err := db.Open()
	if err != nil {
		t.Errorf("Open database failed with error: %v", err)
		return
	}
	defer db.Close()

	_, err = db.ExecSql("create table test(id number, name char(32), salary float);")
	if err != nil {
		t.Errorf("Execute sql failed with error: %v", err)
		return
	}

	rows, err := db.ExecSql("insert into test(id, name, salary) values(?,?,?);", 1, "zhangshan", 32.5)
	if err != nil {
		t.Errorf("Execute sql failed with error: %v", err)
		return
	} else if rows != 1 {
		t.Errorf("Expect 1 but %d", rows)
		return
	}

	rows, err = db.ExecSql("insert into test(id, name, salary) values(?,?,?);", 2, "lisi", 18.9)
	if err != nil {
		t.Errorf("Execute sql failed with error: %v", err)
		return
	} else if rows != 1 {
		t.Errorf("Expect 1 but %d", rows)
		return
	}
}

func TestReopen(t *testing.T) {
	db := database.NewDatabase(dbfile)
	if db == nil {
		t.Errorf("Database is nil")
		return
	}
	defer os.Remove(dbfile)

	err := db.Open()
	if err != nil {
		t.Errorf("Open database failed with error: %v", err)
		return
	}
	defer db.Close()

	err = db.Open()
	if err != nil {
		t.Errorf("Open database failed with error: %v", err)
		return
	}

	values, err := db.Query("select id, name, salary from test where id = ?;", 2)
	if err != nil {
		t.Errorf("Execute sql failed with error: %v", err)
		return
	}

	if len(values) != 1 || values[0][0].(int) != 2 || values[0][1].(string) != "lisi" || values[0][2].(float64) != 18.9 {
		t.Errorf("Expect values failed")
		return
	}
}
