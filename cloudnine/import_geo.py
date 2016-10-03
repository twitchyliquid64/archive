import csv, pg8000, sys, os, time
from geopy import Nominatim

try:
	csvInFileName = sys.argv[1]
	suburbRow = sys.argv[2]
except:
	print "usage"
	print 'import_geo.py "path/to/csv" suburbColumnName'
	sys.exit(0)

def getSim(fname, postfix):
	i = 0
	while True:
		i += 1
		path = "%s_%02d.%s" %(fname, i, postfix)
		if not os.path.exists(path):
			return path

csvOutFileName = getSim(csvInFileName[:-4] + '_geo', 'csv')
geolocator = Nominatim()

def tryHard(callback, n=4):
	o = None
	try:
		o = callback()
	except Exception as ex:
		print ex.message

	if o is not None or n == 0:
		return o

	if n == 2:
		time.sleep(2)

	return tryHard(callback, n-1)

with open(csvInFileName, 'r') as csvIn, open(csvOutFileName, 'w') as csvOut:
	reader = csv.DictReader(csvIn)
	
	columns = reader.fieldnames + ['lat', 'lon']
	print columns

	writer = csv.DictWriter(csvOut, fieldnames=columns)
	writer.writeheader()
	for row in reader:
		location = tryHard(lambda: geolocator.geocode("%s, NSW, Australia" % row[suburbRow]))
		
		if location == None:
			print "Nope [%s]" % row[suburbRow]
			continue
		
		row = dict(row)
		row['lat'] = location.latitude
		row['lon'] = location.longitude
		
		writer.writerow(row)
		
		print row[suburbRow]

print "Done"
