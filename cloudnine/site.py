from tornado import Server
import TemplateAPI
import pg8000, json, os, autointegrate, random, string
import hashlib, json
from model import crime, population, stories, other
import similarity

#user: cloudnine
#pass: lolsaurus1
#database: cloudnine
#local access: psql -d cloudnine -U cloudnine -W

conn = pg8000.connect(port=5432,host="128.199.113.29", user="cloudnine", password="lolsaurus1", database="cloudnine")


def hash(data):
    return hashlib.sha1(repr(str(data))+",AsdSALTY$vd").hexdigest()

def dummy():
    print "Started"

def gitUpdatedWebHook(response):
    autointegrate.updateRepo()


def locationData(response):
    out = [a[0] for a in crime.getLocations(conn)]
    response.write(json.dumps(out))

def testPage(response):
    loc = json.dumps(crime.getLocations(conn), sort_keys=True, indent=4)
    specific = json.dumps(crime.getScoreByLGA(conn, "Hurstville"), sort_keys=True, indent=4)
    response.write(TemplateAPI.render('example.html', response, {"s": specific, "loc": loc, "exampledata": "LOL"}))

def indexPage(response):
    response.write(TemplateAPI.render('index.html', response, {'header': ''}))

def getStories(response):
    LGA = response.get_field("lga")
    s = stories.getStoriesByLGA(conn, LGA)
    response.write(json.dumps(s))

def getDisplay(response):
    LGA = response.get_field("lga")
    s = similarity.getSimilar(conn, LGA)
    cdata1 = crime.getScoreByLGA(conn, LGA)
    pdata1 = population.getPopulationByLGA(conn, LGA)
    rank1 = other.getRankByLGA(conn, LGA)
    output = [{
        "name": LGA,
        "score": 900,
        "crime": cdata1,
	"population": pdata1,
	"rank": rank1,
    }]

    count = 0
    for possibility in s:
        count += 1
        if count < 4:
            cdata = crime.getScoreByLGA(conn, possibility[0])
            pdata = population.getPopulationByLGA(conn, possibility[0])
	    rank = other.getRankByLGA(conn, LGA)
        else:
            cdata = []
            pdata = []
	    rank = []
        obj = {
            "name": possibility[0],
            "score": possibility[1],
            "crime": cdata,
            "population": pdata,
	    "rank": rank,
        }

        if count < 4:
                output.append(obj)
    response.write(json.dumps(output))

def getSimilarLocations(response):
    LGA = response.get_field("lga")
    s = similarity.getSimilar(conn, LGA)
    output = []
    count = 0
    for possibility in s:
        count += 1
        if count < 5:
            cdata = crime.getScoreByLGA(conn, possibility[0])
            pdata = population.getPopulationByLGA(conn, possibility[0])
	    rank = other.getRankByLGA(conn, LGA)
        else:
            cdata = []
            pdata = []
	    rank = []
        obj = {
            "name": possibility[0],
            "score": possibility[1],
	    "crime": cdata,
            "population": pdata,
	    "rank": rank,
	}    
        
	if count < 5:
	        output.append(obj)
    response.write(json.dumps(output))    



import os

port = 80
if os.name == 'nt':
	port = 8080

server = Server('0.0.0.0', port)
server.register("/", indexPage)
server.register("/remotevent/git/commithook", gitUpdatedWebHook)
server.register("/api/locations", locationData)
server.register("/api/similar", getSimilarLocations)
server.register("/api/display", getDisplay)
server.register("/api/stories", getStories)
server.run(dummy)
