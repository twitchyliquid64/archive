import pg8000


def getStoriesByLGA(dbconn, LGA):
	try:
		curs = dbconn.cursor()
		curs.execute("SELECT url,title,date FROM abcspatial NATURAL JOIN abcphotos WHERE lga = %s;", (LGA,))
                result = curs.fetchall()
                curs.close()
                return sorted(result, key=lambda x: str(x[2]).split("/")[2], reverse=True)
	except Exception, e:
		curs.rollback()
		print e



if __name__ == "__main__":
	conn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")
	print getStoriesByLGA(conn, 'Murray')
