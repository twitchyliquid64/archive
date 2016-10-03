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


geolocator = Nominatim()
for loc in crime.getLocations(dbconn):
	locStr = sanitizeName(loc[0]) + ", NSW, Australia"
	print "Location:", locStr
	location = geolocator.geocode(locStr)

	if location == None:
		location = geolocator.geocode(loc[1] + ", " + locStr)
		print "Location:", loc[1] + ", " + locStr

        if location == None:
                location = geolocator.geocode(sanitizeName(loc[0]).split(" ")[0] + ", NSW, Australia")
                print "Location:", loc[1] + ", " + locStr

	if location == None:
		print "Unknown!"
		continue

	print "Lat:", location.latitude, "Long:", location.longitude
	curs = dbconn.cursor()
	curs.execute("UPDATE crimensw SET pos = Point(%s,%s) WHERE lga = %s", (location.longitude,location.latitude,loc[0]))
	dbconn.commit()
	curs.close()
	time.sleep(0.1)
