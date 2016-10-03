import csv
import pg8000
import geopy

fields = []
data = []
dbconn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")


with open("rawdata/crimesum.csv") as csvfile:
	reader = csv.DictReader(csvfile)
	data = [row for row in reader]
	row1 = data[0]

	for key in row1.keys():
		isInt = False
		try:
			a = int(key[-4])
			isInt = True
		except:
			pass
		if key == "score":
			IsInt = True
		fields.append((isInt,key))

createSchema = ""

createSchema += "CREATE TABLE crimensw ("
for field in fields:
	if field[0]:
		createSchema += "\n\t" + field[1].replace(" ","") + " INT,"
	else:
		createSchema += "\n\t" + field[1].replace(" ","") + " TEXT,"

createSchema += "\n\t pos POINT,"

createSchema = createSchema[0:-1] + "\n)"


curs = dbconn.cursor()

try:
	curs.execute("DROP TABLE crimensw;")
	dbconn.commit()
except Exception, e:
	dbconn.rollback()
	print e


curs.execute(createSchema)
dbconn.commit()

query = "INSERT INTO crimensw ("
for field in fields:
        if field[0]:
		query += " " + field[1].replace(" ","") + ","
        else:
		query += " " + field[1].replace(" ","") + ","
queryMain = query[:-1] + ")"

for row in data[:-4]:
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
	print params
	curs.execute(q, params)
	dbconn.commit()
