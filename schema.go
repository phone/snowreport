package main

import (
	"encoding/xml"
)

type Result struct {
	XMLName xml.Name `xml:"dwml"`
	Datas   []Data   `xml:"data"`
}

type Data struct {
	Type             string        `xml:"type,attr"`
	TimeLayouts      []TimeLayout  `xml:"time-layout"`
	Temperatures     []Temperature `xml:"parameters>temperature"`
	WeatherSummaries []Summary     `xml:"parameters>weather>weather-conditions"`
	WordedForecasts  []string      `xml:"parameters>wordedForecast>text"`
	Icons            []string      `xml:"parameters>conditions-icon>icon-link"`
}

type TimeLayout struct {
	LayoutKey string   `xml:"layout-key,attr"`
	Periods   []Period `xml:"start-valid-time"`
}

type Period struct {
	Name string `xml:"period-name,attr"`
}

type Temperature struct {
	Type       string `xml:"type,attr"`
	Units      string `xml:"units,attr"`
	TimeLayout string `xml:"time-layout,attr"`
	Name       string
	Values     []int `xml:"value"`
}

type Summary struct {
	Condition string `xml:"weather-summary,attr"`
}
