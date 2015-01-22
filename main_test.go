package main

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestUrlTemplate(t *testing.T) {
	l := &Location{
		Id:    0,
		Name:  "Wachusett",
		Zip:   "01541",
		Lat:   42.451887494213565,
		Lon:   -71.87972699870119,
		Town:  "Princeton",
		State: "MA",
	}
	u, err := l.GetUrl()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if u != "http://forecast.weather.gov/MapClick.php?FcstType=dwml&lat=42.4518890380859375&lg=english&lon=-71.8797302246093750&unit=0" {
		t.Fail()
	}
}

func TestGetForecasts(t *testing.T) {
	l := &Location{
		Id:    9,
		Name:  "Wachusett",
		Zip:   "01541",
		Lat:   42.451887494213565,
		Lon:   -71.87972699870119,
		Town:  "Princeton",
		State: "MA",
	}
	forecasts, err := l.GetForecasts()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if len(forecasts) != 14 {
		log.Printf("%d forecasts instead of 14\n", len(forecasts))
		t.FailNow()
	}
	for _, fc := range forecasts {
		t.Log(*fc)
	}
}

func XTestPullWeather(t *testing.T) {
	var (
		url = "http://forecast.weather.gov/MapClick.php?lat=42.45189&lon=-71.87972699870119&unit=0&lg=english&FcstType=dwml"
		cl  = &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 4}}
	)
	rsp, err := cl.Get(url)
	if err != nil {
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		log.Print(err)
	} else {
		if rsp.StatusCode != 200 {
			log.Println("status", rsp.StatusCode)
		} else {
			log.Println(rsp.StatusCode)
			if err != nil {
				log.Println(err)
			} else {
				rslt := &Result{}
				decoder := xml.NewDecoder(rsp.Body)
				decoder.CharsetReader = charset.NewReader
				err := decoder.Decode(rslt)
				if err != nil {
					log.Print(err)
				}
				log.Printf("%v\n", rslt)
				for _, d := range rslt.Datas {
					if d.Type == "forecast" {
						log.Println("in forecast")
						log.Printf("we have %d tmps", len(d.Temperatures))
						for _, temp := range d.Temperatures {
							log.Println("in temperature")
							log.Printf("we have %d tmps", len(temp.Values))
							for _, tval := range temp.Values {
								log.Printf("tempval: %d", tval)
							}
						}
						log.Printf("we have %d icons", len(d.Icons))
						log.Printf("we have %d timelayouts", len(d.TimeLayouts))
						log.Printf("we have %d forecasts", len(d.WordedForecasts))
						log.Printf("we have %d forecasts", len(d.WeatherSummaries))
					}
				}
			}
		}
	}
}

func TestEtcd(t *testing.T) {
	ecl := etcd.NewClient([]string{"http://127.0.0.1:4001/"})
	rsp, err := ecl.Get("/snowreport/mysql/user", false, false)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(rsp.Node.Value)
}

func TestGetDb(t *testing.T) {
	_, err := GetDb()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestMysql(t *testing.T) {
	var (
		user     string
		password string
		host     string
		port     string
		err      error
		db       *sql.DB
		ecl      *etcd.Client
	)
	ecl = etcd.NewClient([]string{"http://127.0.0.1:4001/"})
	user, err = GetEtcdVal(ecl, "/snowreport/mysql/user")
	if err != nil {
		t.FailNow()
	}
	password, err = GetEtcdVal(ecl, "/snowreport/mysql/password")
	if err != nil {
		t.FailNow()
	}
	host, err = GetEtcdVal(ecl, "/snowreport/mysql/host")
	if err != nil {
		t.FailNow()
	}
	port, err = GetEtcdVal(ecl, "/snowreport/mysql/port")
	if err != nil {
		t.FailNow()
	}
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/snow", user, password, host, port))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = db.Ping()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func XTestDataCollection(t *testing.T) {
	l := &Location{
		Id:    1,
		Name:  "Wachusett",
		Zip:   "01541",
		Lat:   42.451887494213565,
		Lon:   -71.87972699870119,
		Town:  "Princeton",
		State: "MA",
	}
	forecasts, err := l.GetForecasts()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	db, err := GetDb()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for _, fc := range forecasts {
		fc.Upsert(db)
	}
}

func TestGetLocations(t *testing.T) {
	ls, err := GetLocations()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for _, l := range ls {
		t.Log(*l)
	}
}
