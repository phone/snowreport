package main

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"database/sql"
	"encoding/xml"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//?lat=xxxxxx& lon=yyyyy&   unit=0&    lg=english&  FcstType=dwml
//http://forecast.weather.gov/MapClick.php?lat=42.4518890380859375&FcstType=dwml&lg=english&long=-71.8797302246093750&unit=0

func FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', 16, 32)
}

var (
	urltmpl       = "http://forecast.weather.gov/MapClick.php"
	qryOperations = map[string]func(*Location) string{
		"lat":      func(l *Location) string { return FloatToString(l.Lat) },
		"lon":      func(l *Location) string { return FloatToString(l.Lon) },
		"unit":     func(l *Location) string { return "0" },
		"lg":       func(l *Location) string { return "english" },
		"FcstType": func(l *Location) string { return "dwml" },
	}
)

type Location struct {
	Id    int
	Name  string
	Zip   string
	Lat   float64
	Lon   float64
	Town  string
	State string
}

func (l *Location) GetUrl() (string, error) {
	u, err := url.Parse(urltmpl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range qryOperations {
		q.Set(k, v(l))
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (l *Location) GetForecasts() ([]*Forecast, error) {
	var (
		err   error
		url   string
		cl    = &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
		rsp   *http.Response
		rslt  *Result
		datum *Data
		ret   = make([]*Forecast, 14)
	)
	url, err = l.GetUrl()
	if err != nil {
		return nil, err
	}

	rsp, err = cl.Get(url)
	if err != nil && rsp == nil {
		return nil, errors.New("Protocol Error")
	}
	if rsp.StatusCode != 200 {
		if rsp != nil {
			io.Copy(ioutil.Discard, rsp.Body)
			rsp.Body.Close()
		}
		return nil, errors.New("Status " + string(rsp.StatusCode))
	}
	rslt = &Result{}
	decoder := xml.NewDecoder(rsp.Body)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(rslt)
	if err != nil {
		return nil, err
	}
	for _, d := range rslt.Datas {
		if d.Type == "forecast" {
			datum = &d
			break
		}
	}
	maxidx := len(datum.WordedForecasts)
	for i := 0; i < maxidx; i++ {
		fcst := &Forecast{
			LocationId: l.Id,
			Index:      i,
		}
		tl14 := TimeLayout{}
		for _, tl := range datum.TimeLayouts {
			tl := tl
			if len(tl.Periods) == maxidx {
				tl14 = tl
				break
			}
		}
		fcst.DateDesc = tl14.Periods[i].Name
		fcst.Summary = datum.WeatherSummaries[i].Condition
		fcst.Forecast = datum.WordedForecasts[i]
		fcst.Icon = datum.Icons[i]
		for _, tmp := range datum.Temperatures {
			if tmp.Type == "minimum" && i < len(tmp.Values) {
				fcst.Low = tmp.Values[i]
			}
			if tmp.Type == "maximum" && i < len(tmp.Values) {
				fcst.High = tmp.Values[i]
			}
		}
		ret[i] = fcst
	}
	return ret, nil
}

type Forecast struct {
	LocationId int
	Index      int
	DateDesc   string
	Summary    string
	Forecast   string
	High       int
	Low        int
	Icon       string
	Date       int64
}

// assumes sql server in ANSI_QUOTES or ANSI modes
const INSSTR string = `INSERT INTO forecast (
	location_id,
	"index",
	datedesc,
	summary,
	forecast,
	high,
	low,
	icon)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	location_id=VALUES(location_id),
	"index"=VALUES("index"),
	datedesc=VALUES(datedesc),
	summary=VALUES(summary),
	forecast=VALUES(forecast),
	high=VALUES(high),
	low=VALUES(low),
	icon=VALUES(icon)
`

func (f *Forecast) Upsert(db *sql.DB) error {
	var (
		err error
	)
	if err = db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(
		INSSTR,
		f.LocationId,
		f.Index,
		f.DateDesc,
		f.Summary,
		f.Forecast,
		f.High,
		f.Low,
		f.Icon,
	)
	if err != nil {
		return err
	}
	return nil
}
