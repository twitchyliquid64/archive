import pg8000


def getPopulationByLGA(dbconn, LGA):
	try:
		curs = dbconn.cursor()
		curs.execute("SELECT children,yadult,oadult,old,suburb FROM crimensw c, population p WHERE c.lga = %s AND c.pos <-> p.pos < 1 ORDER BY c.pos <-> p.pos LIMIT 1;", (LGA,))
                result = curs.fetchone()
                curs.close()
                return result
	except Exception, e:
		curs.rollback()
		print e
