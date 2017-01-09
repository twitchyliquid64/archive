@echo off

@set GOPATH=%GOPATH%;C:\coding\projects\go-osm-parse;C:\coding\lib\go
echo !!!! RUN USING ::: go run -tags purego main.go osm_structures.go database.go !!!!
pause
go run -tags purego main.go osm_structures.go database.go
rem ql.exe -db="osmdump.db" "SELECT * FROM NodeTags,Nodes WHERE Nodes.NodeID == NodeTags.NodeID"
pause


rem %windir%\system32\cmd.exe
