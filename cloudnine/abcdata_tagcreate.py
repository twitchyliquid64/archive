import csv
import pg8000
from geopy import Nominatim
from model import crime
import time

fields = []
data = []
dbconn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")

def sanitizeName(inp):
	a = inp.replace("Northern ", "")
	a = a.replace("Upper ", "")
	a = a.replace("Lower ", "")
	return a





#delete existing database
curs = dbconn.cursor()
try:
        curs.execute("DROP TABLE abcrelations;")
        dbconn.commit()
except Exception, e:
        dbconn.rollback()
        print e



createSchema = "CREATE TABLE abcrelations (url TEXT, name varchar(128));"
curs.execute(createSchema)
dbconn.commit()



curs = dbconn.cursor()
curs.execute("SELECT url,keywords,subjects FROM abcphotos")
rows = curs.fetchall()
curs.close()

for row in rows:
	relationSet = []
	for tag in row[1].split(","):
		tag = tag.strip().replace(":","").upper()
		if tag == "":
			continue
		tag = (tag[:125] + '..') if len(tag) > 125 else tag
		relationSet.append(tag)
        for tag in row[2].split(","):
                tag = tag.strip().replace(":","").upper()
                if tag == "":
                        continue
		tag = (tag[:125] + '..') if len(tag) > 125 else tag
                relationSet.append(tag)

	for relation in relationSet:
		print row[0], relation
		curs = dbconn.cursor()
		curs.execute("INSERT INTO abcrelations (url,name) VALUES (%s,%s)", (row[0],relation))
		dbconn.commit()
		curs.close()
