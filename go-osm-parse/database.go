package main

import "database/sql"
import "time"
import "fmt"

import "runtime"

var tableCreateCmd = `
	CREATE TABLE IF NOT EXISTS Nodes
	(
		NodeID int64,
		Lat float64,
		Lon float64,
		Changeset int64
	);
	CREATE INDEX IF NOT EXISTS NodesID ON Nodes (NodeID);
	CREATE INDEX IF NOT EXISTS NodesLat ON Nodes (Lat);
	CREATE INDEX IF NOT EXISTS NodesLon ON Nodes (Lon);
	
	
	CREATE TABLE IF NOT EXISTS NodeTags
	(
		NodeID int64,
		Key string,
		Value string
	);
	CREATE INDEX IF NOT EXISTS NodeTagsID ON NodeTags (NodeID);
	
	
	
	
	
	
	
	CREATE TABLE IF NOT EXISTS Ways
	(
		WayID int64,
	);
	CREATE INDEX IF NOT EXISTS WaysID ON Ways (WayID);
	
	
	CREATE TABLE IF NOT EXISTS WayTags
	(
		WayID int64,
		Key string,
		Value string
	);
	CREATE INDEX IF NOT EXISTS WayTagsID ON WayTags (WayID);
	
	
	CREATE TABLE IF NOT EXISTS WayRefs
	(
		WayID int64,
		RefID int64
	);
	CREATE INDEX IF NOT EXISTS WayRefsWayID ON WayRefs (WayID);
	CREATE INDEX IF NOT EXISTS WayRefsRefID ON WayRefs (RefID);
	
	
	
	
	
	
	
	CREATE TABLE IF NOT EXISTS Relations
	(
		RelationID int64,
		Changeset int64
	);
	CREATE INDEX IF NOT EXISTS RelationsID ON Relations (RelationID);
	
	CREATE TABLE IF NOT EXISTS RelationTags
	(
		RelationID int64,
		Key string,
		Value string
	);
	CREATE INDEX IF NOT EXISTS RelationsID ON RelationTags (RelationID);
	
	CREATE TABLE IF NOT EXISTS RelationMembers
	(
		RelationID int64,
		Ref int64,
		Type string,
		Role string
	);
	CREATE INDEX IF NOT EXISTS RelationMembersID ON RelationMembers (RelationID);
	CREATE INDEX IF NOT EXISTS RelationMembersRef ON RelationMembers (Ref);
`


func initialise_db(db *sql.DB) {
	
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	
	if _, err = tx.Exec(tableCreateCmd); err != nil {
		panic(err)
	}
	
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}


func db_worker(db *sql.DB, nodeIn *chan *Node, wayIn *chan *Way, relationsIn *chan *Relation){
	workers.Add(1)
	defer workers.Done()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	
	nodes := *nodeIn
	ways  := *wayIn
	relations  := *relationsIn
	count := 0
	for  {
		select {
			case node, ok := <- nodes:
				if !ok{
					nodes = nil
				}else{
					if _, err = tx.Exec("INSERT INTO Nodes (NodeID,Lat,Lon,Changeset) VALUES ($1,$2,$3,$4)", node.Id, node.Lat, node.Lon, node.Changeset); err != nil {
						panic(err)
					}
					
					for _, tag := range node.Tags{
						if _, err = tx.Exec("INSERT INTO NodeTags (NodeID,Key,Value) VALUES ($1,$2,$3)", node.Id, tag.K, tag.V); err != nil {
							panic(err)
						}
					}
				}
				
				
			case way, ok := <- ways:
				if !ok{
					ways = nil
				}else{
					if _, err = tx.Exec("INSERT INTO Ways (WayID) VALUES ($1)", way.Id); err != nil {
						panic(err)
					}
					
					for _, tag := range way.Tags{
						if _, err = tx.Exec("INSERT INTO WayTags (WayID,Key,Value) VALUES ($1,$2,$3)", way.Id, tag.K, tag.V); err != nil {
							panic(err)
						}
					}
					
					for _, ref := range way.Refs{
						if _, err = tx.Exec("INSERT INTO WayRefs (WayID,RefID) VALUES ($1,$2)", way.Id, ref.Ref); err != nil {
							panic(err)
						}
					}
				}
				
				
			case relation, ok := <- relations:
				if !ok{
					relations = nil
				}else{
					if _, err = tx.Exec("INSERT INTO Relations (RelationID, Changeset) VALUES ($1,$2)", relation.Id, relation.Changeset); err != nil {
						panic(err)
					}
					
					for _, tag := range relation.Tags{
						if _, err = tx.Exec("INSERT INTO RelationTags (RelationID,Key,Value) VALUES ($1,$2,$3)", relation.Id, tag.K, tag.V); err != nil {
							panic(err)
						}
					}
					
					for _, ref := range relation.Members{
						if _, err = tx.Exec("INSERT INTO RelationMembers (RelationID,Ref,Type,Role) VALUES ($1,$2,$3,$4)", relation.Id, ref.Ref, ref.Type, ref.Role); err != nil {
							panic(err)
						}
					}
				}
		}

		if ways==nil && nodes==nil && relations==nil{ break }

		count++
		if (count%1000) == 0{
			fmt.Print(count, " entries processed. Intermission commit in progress ... ")
			if err = tx.Commit(); err != nil {
				panic(err)
			}
			
			if (count%10000) == 0{
				fmt.Print("Pause. ")
				runtime.GC()
				time.Sleep(time.Second * 2)
			}
			
			
			tx, err = db.Begin()
			if err != nil {
				panic(err)
			}
			fmt.Println("DONE.")
		}
	}
	
	fmt.Println("Now finalizing commit.")
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	fmt.Println("All entries committed to db.")
}
