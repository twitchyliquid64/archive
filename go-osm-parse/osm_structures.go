package main

import "encoding/xml"


type OSM struct {
	XMLName xml.Name `xml:"osm"`
	Version string `xml:"version,attr"`
	Generator string `xml:"generator,attr"`
}

type Bounds struct {
	XMLName xml.Name `xml:"bounds"`
	MinLat string `xml:"minlat,attr"`
	MaxLat string `xml:"maxlat,attr"`
	
	MinLon string `xml:"minlon,attr"`
	MaxLon string `xml:"maxlon,attr"`
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

type Ref struct {
	XMLName xml.Name `xml:"nd"`
	Ref int64 `xml:"ref,attr"`
}

type Member struct {
	XMLName xml.Name `xml:"member"`
	Type string `xml:"type,attr"`
	Ref int64 `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}


type Way struct {
	XMLName xml.Name `xml:"way"`
	Id int64 `xml:"id,attr"`
	Tags []Tag `xml:"tag"`
	Refs []Ref `xml:"nd"`
}


type Node struct {
	XMLName xml.Name `xml:"node"`
	Id int64 `xml:"id,attr"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Changeset int64 `xml:"changeset,attr"`
	Tags []Tag `xml:"tag"`
}

type Relation struct {
	XMLName xml.Name `xml:"relation"`
	Id int64 `xml:"id,attr"`
	Changeset int64 `xml:"changeset,attr"`
	Tags []Tag `xml:"tag"`
	Members []Member `xml:"member"`
}

