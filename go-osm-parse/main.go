package main

import (
	"database/sql"
	_ "github.com/cznic/ql/driver"
	"encoding/xml"
	"runtime"
	"sync"
	"fmt"
	"os"
)

var workers sync.WaitGroup

//result of the parse operation
var bounds Bounds

func process(decoder *xml.Decoder, nodes *chan *Node, ways *chan *Way, relations *chan *Relation){
	var inElement string
	
	for {
		t, _ := decoder.Token()// Read tokens from the XML document in a stream.
		if t == nil {
			break
		}

		switch se := t.(type) {
			case xml.StartElement:
				inElement = se.Name.Local

				if inElement == "bounds" {
					decoder.DecodeElement(&bounds, &se)
					fmt.Println( "Latitude:", bounds.MinLat, "to", bounds.MaxLat, "- Longitude:", bounds.MinLon, "to", bounds.MaxLon)
				}else if inElement == "node" {
					var node Node
					decoder.DecodeElement(&node, &se)
					*nodes <- &node
				}else if inElement == "way" {
					var way Way
					decoder.DecodeElement(&way, &se)
					*ways <- &way
					//fmt.Println(way)
				}else if inElement == "relation" {
					var relation Relation
					decoder.DecodeElement(&relation, &se)
					*relations <- &relation
					//fmt.Println(relation)
				}else {
					//fmt.Println(inElement)
				}
		}

	}
}



func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(" ~===== OSM parser version 0.9 Alpha =====~ ")
	
	xmlFile, err := os.Open("input.osm")//DISABLE YOUR AV CUZ ITS A FUCKIN BIG FILE AND AVs LIKE TO SCAN FILES BEFORE IT LETS THEM OPEN!!!
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	
	db, err := sql.Open("ql", "osmdump.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	initialise_db(db)
	defer db.Close()
	
	
	fmt.Print("Initializing ... ")
	decoder := xml.NewDecoder(xmlFile)
	waysChan := make(chan *Way, 5000)
	nodesChan := make(chan *Node, 5000)
	relationsChan := make(chan *Relation, 5000)
	fmt.Println("DONE.")
	
	fmt.Print("Launching DB workers ... ")
	go db_worker(db, &nodesChan, &waysChan, &relationsChan)
	fmt.Println("DONE.")

	fmt.Println("Processing started.")
	process(decoder, &nodesChan, &waysChan, &relationsChan)
	close(waysChan)
	close(nodesChan)
	close(relationsChan)
	fmt.Println("Waiting for all db workers to complete.")
	workers.Wait()
	fmt.Println("Processing completed, flushing disk caches.")
}
