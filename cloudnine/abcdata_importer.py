import pg8000
import geopy
import json

fields = []
data = []
dbconn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")


with open('rawdata/abcphotos.json') as jsonfile:
	data = json.load(jsonfile, encoding="latin-1")

for key in data[0].keys():
	isInt = False
	if key == "score":
		IsInt = True
	fields.append((isInt,key))

createSchema = ""

createSchema += "CREATE TABLE abcphotos ("
for field in fields:
        if field[0]:
                createSchema += "\n\t" + field[1].replace(" ","") + " INT,"
        else:
                createSchema += "\n\t" + field[1].replace(" ","") + " TEXT,"


createSchema += "\n\tpos POINT,"
createSchema = createSchema[0:-1] + "\n)"

print createSchema

curs = dbconn.cursor()
try:
        curs.execute("DROP TABLE abcphotos;")
        dbconn.commit()
except Exception, e:
        dbconn.rollback()
        print e

curs.execute(createSchema)
dbconn.commit()

query = "INSERT INTO abcphotos ("
for field in fields:
        if field[0]:
                query += " " + field[1].replace(" ","") + ","
        else:
                query += " " + field[1].replace(" ","") + ","
queryMain = query[:-1] + ")"

print queryMain

for row in data:
        tempQuery = " VALUES ("
        params = []
	for field in fields:
                params.append(row[field[1]])
                if field[0]:
                        tempQuery += "%s,"
                else:
                        tempQuery += "%s,"
        tempQuery = tempQuery[:-1] + ")"
        q = queryMain + " " + tempQuery
        print q, params, len(params), q.count("%s")
        curs.execute(q, params)
        dbconn.commit()

#all the data but the point is imported at this stage!


curs = dbconn.cursor()
curs.execute("SELECT longitude,latitude FROM abcphotos")
rows = curs.fetchall()
curs.close()

for row in rows:
        curs = dbconn.cursor()
	print row
	if row[0] != '':
	        curs.execute("UPDATE abcphotos SET pos = Point(%s,%s) WHERE longitude = %s AND latitude = %s", (float(row[0]),float(row[1]),row[0],row[1]))
	else:
		curs.execute("DELETE FROM abcphotos WHERE longitude = ''")
        dbconn.commit()
	curs.close()

#now the points are imported!
