import pg8000
from model import crime

crime_mapping = {
	"Assault": 1.5,
	"Sexual offences": 4,
	"Homicide": 6,
	"Arson": 5,
	"Prohibited and regulated weapons offences": 2,
	"Prostitution offences": 0.9,
	"Disorderly conduct": 1.3,
	"Abduction and kidnapping": 5
}
CRIME_MULTIPLIER = -0.0000006


def getStoriesByLocation(dbconn, LGA):
        try:
                curs = dbconn.cursor()
                curs.execute("SELECT url,score FROM abcspatial WHERE lga = %s ORDER BY score ASC", (str(LGA)))
                result = curs.fetchall()
                curs.close()
                return result
        except Exception, e:
                curs.rollback()
                print e


def getTagsByLocation(dbconn, LGA):
	#SELECT name,score FROM abcspatial NATURAL JOIN abcrelations WHERE lga = 'Blacktown';
        try:
                curs = dbconn.cursor()
                curs.execute("SELECT name FROM abcspatial NATURAL JOIN abcrelations WHERE lga = %s ORDER BY score ASC;", (LGA,))
                result = curs.fetchall()
                curs.close()
                return result
        except Exception, e:
                print e
		dbconn.rollback()

def getLocationsByTags(dbconn,tags, exclusionLGA):
	try:
		curs = dbconn.cursor()
		tags = [tag[0] for tag in tags]
		in_p = ', '.join(map(lambda x: '%s', tags))
		curs.execute("SELECT lga FROM abcrelations NATURAL JOIN abcspatial WHERE name IN ("+in_p+");", tags)
                result = curs.fetchall()
                curs.close()
        except Exception, e:
                dbconn.rollback()
		raise

	outDict = dict()
	for tag in result:
		if tag[0] == exclusionLGA:
			continue
		if tag[0] in outDict:
			outDict[tag[0]] = outDict[tag[0]] + 1
		else:
			outDict[tag[0]] = 1
	return getSListFromDict(outDict)


def AdjustScoresByCrime(conn, sortedList):
	out = sortedList
	for x in xrange(len(out)):
		scoreAdj = ScoreCrimeForLGA(conn, out[x][0])
		out[x] = (out[x][0], out[x][1]+scoreAdj)
	return sorted(out, key=lambda x: x[1], reverse=True)

def ScoreCrimeForLGA(conn, LGA):
	a = crime.getScoreByLGA(conn, LGA)
	scoreAdj = 0
	for b in a:
		category = b[2]
		score = float(b[4])
		if category in crime_mapping:
			score *= float(crime_mapping[category])
		scoreAdj += float(score)
	return scoreAdj * CRIME_MULTIPLIER

def getSListFromDict(l):
	out = []
	for tag in l.keys():
		out.append((tag, l[tag]))
	return sorted(out, key=lambda x: x[1], reverse=True)


def getSimilar(dbconn, LGA):
	tags = getTagsByLocation(dbconn, LGA)
	similarLocationsByTag = getLocationsByTags(dbconn,tags,LGA)
	crimeAdjusted = AdjustScoresByCrime(dbconn, similarLocationsByTag)
	return crimeAdjusted


if __name__ == "__main__":
	conn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")
	print getSimilar(conn, 'Ryde')
