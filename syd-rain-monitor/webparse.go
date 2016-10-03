package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func getRadarData()(ret map[string]string, err error){
	res, err := http.Get("http://www.bom.gov.au/products/IDR714.loop.shtml")
	if err != nil {
		return map[string]string{}, err
	}
	d, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return map[string]string{}, err
	}

	lines := strings.Split(string(d), "\n")
	r := map[string]string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "theImageNames[0] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[0] = \"")
			r["img1"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "theImageNames[1] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[1] = \"")
			r["img2"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "theImageNames[2] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[2] = \"")
			r["img3"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "theImageNames[3] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[3] = \"")
			r["img4"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "theImageNames[4] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[4] = \"")
			r["img5"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "theImageNames[5] = \"") {
			s1 := strings.TrimPrefix(line, "theImageNames[5] = \"")
			r["img6"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "lat = \"") {
			s1 := strings.TrimPrefix(line, "lat = \"")
			r["lat"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "lon = \"") {
			s1 := strings.TrimPrefix(line, "lon = \"")
			r["lon"] = s1[:len(s1)-2]
		}

		if strings.HasPrefix(line, "Km = ") {
			s1 := strings.TrimPrefix(line, "Km = ")
			r["km"] = s1[:len(s1)-1]
		}



	}

	return r, nil
}
