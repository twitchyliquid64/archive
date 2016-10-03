import pg8000

def getLocations(dbconn):
	try:
		curs = dbconn.cursor()
		curs.execute("SELECT DISTINCT lga,'' FROM crimensw")
		result = curs.fetchall()
		curs.close()
		return result
	except Exception, e:
		curs.rollback()
		print e


def getScoreByLGA(dbconn, LGA):
	try:
		curs = dbconn.cursor()
		curs.execute("SELECT lga, statisticaldivisionorsubdivision, offencecategory, subcategory, score FROM crimensw WHERE lga=%s", (LGA,))
                result = curs.fetchall()
                curs.close()
                return result
	except Exception, e:
		curs.rollback()
		print e
