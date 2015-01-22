package main

import (
	"database/sql"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func GetEtcdVal(ecl *etcd.Client, key string) (string, error) {
	rsp, err := ecl.Get(key, false, false)
	if err != nil {
		return "", err
	}
	return rsp.Node.Value, nil
}

func GetDb() (*sql.DB, error) {
	var (
		err error
		tmp string
		db  *sql.DB
		ks  = []string{
			"/snowreport/mysql/user",
			"/snowreport/mysql/password",
			"/snowreport/mysql/host",
			"/snowreport/mysql/port",
		}
		vals []interface{} = make([]interface{}, len(ks), len(ks))
		ecl                = etcd.NewClient([]string{"http://127.0.0.1:4001/"})
	)

	// this is vaguely flimsy but w/e
	for i, k := range ks {
		tmp, err = GetEtcdVal(ecl, k)
		if err != nil {
			return nil, err
		}
		vals[i] = tmp
	}
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/snow", vals...))
	if err != nil {
		return nil, err
	}

	// let's force the connection and fail now instead of on a query
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetLocations(db *sql.DB) ([]*Location, error) {
	var (
		err  error
		rows *sql.Rows
		ret  = make([]*Location, 0, 10)
	)

	rows, err = db.Query("select * from location")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		l := &Location{}
		err = rows.Scan(&l.Id, &l.Name, &l.Zip, &l.Lat, &l.Lon, &l.Town, &l.State)
		if err != nil {
			log.Println(rows.Err())
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}
