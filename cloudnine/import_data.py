import csv, pg8000, sys

try:
	dbName = sys.argv[1]
	csvFilePath = sys.argv[2]
except:
	print "Usage"
	print 'import_data.py dbName "path\to\csv"' 
	sys.exit(0)

dbconn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")

class Table:
	def __init__(self, name, csvIn):
		self.name = name
		
		reader = csv.DictReader(csvIn)
		columns = reader.fieldnames
		if 'lat' in columns:
			self.lat, self.lon = 'lat', 'lon'
		elif 'longitude' in columns:
			self.lat, self.lon = 'latitude', 'longitude'
		else:
			print "Cant find lat/long params"
			if raw_input('Continue? (y/n): ') != 'y':
				sys.exit(1)
			self.lat, self.lon = None, None
		
		frow = None
		for row in reader:
			frow = row
			break
		
		self.columns = self.determineRowTypes(columns, row)
		csvIn.seek(0)
		
	def determineRowTypes(self, columns, row):
		gcolumns = {}
		if self.lat != None:
			gcolumns['pos'] = 'POINT'
		
		for c in columns:
			if c not in [self.lat, self.lon]:
				ctype = 'TEXT'
				try:
					int(row[c])
					ctype = 'INT'
				except:
					pass
				gcolumns[c] = ctype
		
		return gcolumns
	
	def create(self):
		types = ', '.join(["%s %s" % (k.lower(), self.columns[k]) for k in self.columns.keys()])
		return "create table %s (%s)" % (self.name, types)

	def format_for_insert(self, row, columns):
		values = []
		for c in columns:
			if c == 'pos':
				values.append('POINT(%s,%s)' % (row[self.lon], row[self.lat]))
			elif self.columns[c] == 'INT':
				values.append(str(row[c]))
			else:
				values.append("'%s'" % row[c].replace("'", ""))
		return "(%s)" % ', '.join(values)
	
	def insert(self, rows):
		i = 0
		sortedRows = self.columns.keys()
		values = ', '.join([self.format_for_insert(row, sortedRows) for row in rows])
		statement = "INSERT INTO %s (%s) VALUES %s" % (self.name, ', '.join(sortedRows), values)
		
		return statement

def execute(s):
	print s
	curs = dbconn.cursor()
	curs.execute(s)
	dbconn.commit()
	curs.close()

def drop(tname):
	try:
		execute("DROP TABLE " + tname)
	except:
		dbconn.rollback()

with open(csvFilePath, 'r') as csvIn:
	meta = Table(dbName, csvIn)
	
	drop(meta.name)
	execute(meta.create())
	execute(meta.insert(csv.DictReader(csvIn)))

