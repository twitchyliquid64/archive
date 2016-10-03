
import pg8000
import geopy
import json

fields = []
data = []
dbconn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")






curs = dbconn.cursor()
try:
        curs.execute("DROP TABLE abcspatial;")
        dbconn.commit()
except Exception, e:
        dbconn.rollback()
        print e





createSchema = "CREATE TABLE abcspatial (url TEXT, lga varchar(128), score int);"
curs.execute(createSchema)
dbconn.commit()












curs = dbconn.cursor()
curs.execute("SELECT lga,pos FROM crimensw;")
rows = curs.fetchall()
curs.close()

lgaToPos = dict()
for row in rows:
	lgaToPos[row[0]] = row[1]

#SELECT * from abcphotos WHERE pos <-> POINT(147.2009759,-34.8163798) < 1 ORDER BY pos <-> POINT(147.2009759,-34.8163798) LIMIT 5;
print "Internal done."
import ast
for r in lgaToPos:
	location = ast.literal_eval(lgaToPos[r])
	curs = dbconn.cursor()
	curs.execute("SELECT url FROM abcphotos WHERE pos <-> POINT(%s,%s) < 1 ORDER BY pos <-> POINT(%s,%s) LIMIT 5;", (location[0],location[1],location[0],location[1],))
	results = curs.fetchall()
	curs.close()
	count = 0
	for photo in results:
		count = count + 1
		curs = dbconn.cursor()
		curs.execute("INSERT INTO abcspatial(url,lga,score) VALUES(%s,%s,%s);", (photo[0],r,str(count)))
		dbconn.commit()
		curs.close()
		print photo, r
