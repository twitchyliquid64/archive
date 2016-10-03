from tornado import Server
import TemplateAPI
import sqlite3 as lite
import random, string
import json

con = lite.connect('main.db')

def randomword(length):
	return ''.join(random.choice(string.lowercase) for i in range(length))

def createTables():
	cur = con.cursor()
	cur.execute("CREATE TABLE IF NOT EXISTS Users(id INTEGER PRIMARY KEY, Username varchar(255) NOT NULL UNIQUE, Name varchar(255), Password varchar(255), Email varchar(255), ProfilePic BLOB)")
	cur.execute("CREATE TABLE IF NOT EXISTS Sessions(id INTEGER PRIMARY KEY, key varchar(64) NOT NULL UNIQUE, name varchar(255), userID INTEGER)")
	#cur.execute("DROP TABLE SquadMembership")
	cur.execute("CREATE TABLE IF NOT EXISTS Squads(id INTEGER PRIMARY KEY, name varchar(64) NOT NULL, key varchar(64) NOT NULL UNIQUE, creatorID INTEGER)")
	cur.execute("CREATE TABLE IF NOT EXISTS SquadMembership(id INTEGER PRIMARY KEY, squadID INTEGER NOT NULL, userID INTEGER NOT NULL, created TIMESTAMP)")
	cur.execute("CREATE TABLE IF NOT EXISTS Locations(id INTEGER PRIMARY KEY, userID INTEGER NOT NULL, provider varchar(64), lon REAL, lat REAL, precision INT, created TIMESTAMP, batt INTEGER)")
	cur.execute("CREATE TABLE IF NOT EXISTS Notifications(id INTEGER PRIMARY KEY, userID INTEGER NOT NULL, type varchar(64), content varchar(512), created TIMESTAMP, subline varchar(128))")
	con.commit()
	cur.close()

def dummy():
    print "Started"


def checkSession(key): #helper to fetch session data associated with a session key
	cur = con.cursor()
	cur.execute("SELECT id, name, userID FROM Sessions WHERE key = ?", (str(key),))
	res = cur.fetchone()
	cur.close()
	return res

def getUser(userID): #helper to fetch user data by userID
	cur = con.cursor()
	cur.execute("SELECT id, username, name, email FROM Users WHERE id = ?", (str(userID),))
	res = cur.fetchone()
	cur.close()
	return res

def getNotifications(userID):
	cur = con.cursor()
	cur.execute("SELECT * FROM Notifications WHERE userid = ? ORDER BY created desc LIMIT 100", (str(userID),))
	res = cur.fetchall()
	cur.close()
	return res

def getAllMembersAllSquads(userID):
	cur = con.cursor()
	cur.execute("SELECT DISTINCT userID FROM SquadMembership WHERE squadID IN (SELECT DISTINCT squadID FROM SquadMembership WHERE userID = ?)", (str(userID),))
	res = cur.fetchall()
	cur.close()
	return res

def getMemberSquadIDs(userID):
	cur = con.cursor()
	cur.execute("SELECT DISTINCT squadID FROM SquadMembership WHERE userID = ?", (str(userID),))
	res = cur.fetchall()
	cur.close()
	return res

def userIsSquadMember(squadID, userID):
	cur = con.cursor()
	cur.execute("SELECT squadID FROM SquadMembership WHERE userID = ? AND squadID = ?", (str(userID),str(squadID),))
	res = cur.fetchone()
	cur.close()
	if res == None:
		return False
	return True

def squadDetails(squadID):
	cur = con.cursor()
	cur.execute("SELECT id, name, key, creatorID FROM Squads WHERE id = ?", (str(squadID),))
	res = cur.fetchone()
	cur.close()
	return res

def squadDetailsByKey(squadKey):
	cur = con.cursor()
	cur.execute("SELECT id, name, key, creatorID FROM Squads WHERE key = ?", (str(squadKey),))
	res = cur.fetchone()
	cur.close()
	return res


def getAllSquadMembers(squadID):
	cur = con.cursor()
	cur.execute("SELECT id, username, name FROM Users WHERE id IN (SELECT userID FROM SquadMembership WHERE squadID = ?)", (str(squadID),))
	res = cur.fetchall()
	cur.close()
	return res
	

def getSquadDetails(response):	
	key = response.get_field("key")
	squadID = response.get_field("squadid")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		if userIsSquadMember(squadID, u[0]):
			deets = squadDetails(squadID)
			response.write(json.dumps({'details': deets, 'members': getAllSquadMembers(squadID)}))
		else:
			response.write("ERROR")
	else:
		response.write("ERROR")


def getNotificationsHandler(response):
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getNotifications(sess[2])
		response.write(json.dumps(u))
	else:
		response.write("ERROR")


def getSelfDetails(response): #handler to return session and user information
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		response.write(json.dumps({'s': sess, 'u': u, 'squads': getMemberSquadIDs(sess[2])}))
	else:
		response.write("ERROR")


def joinSquad(response):
	key = response.get_field("key")
	squadKey = response.get_field("squadkey")
	squadDetails = squadDetailsByKey(squadKey)
	if not squadDetails:
		response.write("ERROR")
		return
	squadID = squadDetails[0]

	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		cur = con.cursor()
		if userIsSquadMember(squadID, u[0]):
			response.write("ERROR")
		else:
			cur.execute("INSERT INTO SquadMembership(squadID, userID) VALUES(?,?)", (str(squadID), str(u[0]),))
			con.commit()
			cur.close()
			response.write(str(squadID))
	else:
		response.write("ERROR")


def createSquad(response):
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		cur = con.cursor()
		cur.execute("INSERT INTO Squads(name,key,creatorID) VALUES(?,?,?)", (
			response.get_field("name"),
			randomword(7),
			str(u[0]),))
		squadID = cur.lastrowid
		cur.execute("INSERT INTO SquadMembership(squadID, userID) VALUES(?,?)", (str(squadID), str(u[0]),))
		con.commit()
		cur.close()
		response.write(str(squadID))
	else:
		response.write("ERROR")

def insertNotificationHandler(response):
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		cur = con.cursor()
		cur.execute("INSERT INTO Notifications(type, content, userID, created, subline) VALUES(?,?,?,CURRENT_TIMESTAMP,?)",
		(str(response.get_field("type")), 
		str(response.get_field("content")),
		str(response.get_field("userid")),
		str(response.get_field("subline")),))
		con.commit()
		cur.close()
		response.write("GOOD")
	else:
		response.write("ERROR")

def insertAllNotificationsHandler(response): #getAllMembersAllSquads
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		for userD in getAllMembersAllSquads(sess[2]):
			cur = con.cursor()
			cur.execute("INSERT INTO Notifications(type, content, userID, created, subline) VALUES(?,?,?,CURRENT_TIMESTAMP,?)",
			(str(response.get_field("type")), 
			str(response.get_field("content")),
			str(userD[0]),
			str(response.get_field("subline")),))
			con.commit()
			cur.close()
		response.write("GOOD")
	else:
		response.write("ERROR")


def insertSquadNotificationsHandler(response):
	key = response.get_field("key")
	squadID = response.get_field("squadid")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		for userD in getAllSquadMembers(squadID):
			cur = con.cursor()
			cur.execute("INSERT INTO Notifications(type, content, userID, created, subline) VALUES(?,?,?,CURRENT_TIMESTAMP,?)",
			(str(response.get_field("type")), 
			str(response.get_field("content")),
			str(userD[0]),
			str(response.get_field("subline")),))
			con.commit()
			cur.close()
		response.write("GOOD")
	else:
		response.write("ERROR")


def inserLocationsHandler(response):
	key = response.get_field("key")
	sess = checkSession(key)
	if sess:
		u = getUser(sess[2])
		cur = con.cursor()
		cur.execute("INSERT INTO Locations(provider, lat, lon, precision, userID, created, batt) VALUES(?,?,?,?,?,CURRENT_TIMESTAMP, ?)", (
			response.get_field("prov"),
			response.get_field("lat"),
			response.get_field("lon"),
			response.get_field("acc"),
			str(u[0]),
			response.get_field("batt"),))
		con.commit()
		cur.close()
		response.write(str("GOOD"))

def createSessionHandler(response): #handler to create a session given a username and password
	username = response.get_field("username")
	passwd = response.get_field("password")
	print "Creating new session for:", username
	cur = con.cursor()
	cur.execute("SELECT id, username, name FROM Users WHERE username = ? AND password = ?", (username,passwd,))
	res = cur.fetchone()
	cur.close()
	if res:
		key = randomword(20)
		print "Password valid", key
		cur = con.cursor()		
		cur.execute("INSERT INTO Sessions (key, name, userID) VALUES(?,?,?)", (key,res[2], res[0],))
		con.commit()
		cur.close()
		response.write(key)
	else:
		response.write("ERROR")

def getLocationForUserHandler(response):
	key = response.get_field("key")
	sess = checkSession(key)
	uID = response.get_field("userid")
	if sess:
		u = getUser(sess[2])
		cur = con.cursor()
		cur.execute("SELECT id,lat,lon,precision,batt,((julianday() - julianday(created)) * 86400.0), provider FROM Locations WHERE userid = ? ORDER BY created desc", (uID,))
		res = cur.fetchone()
		response.write(json.dumps(res))
	else:
		response.write("ERROR")	

def checkUserExistsHandler(response): #handler which indicates if a specific username is taken or not
	username = response.get_field("username")
	print "Checking:", username
	cur = con.cursor()
	cur.execute("SELECT username FROM Users WHERE username = ?", (username,))
	res = cur.fetchall()
	cur.close()
	response.write("OK" if len(res) == 0 else "EXISTS")

def registerUserHandler(response): #handler to register a specific username
	cur = con.cursor()
	cur.execute("INSERT INTO Users (username,password,name,email) VALUES (?,?,?,?)", (
			response.get_field("username"),
			response.get_field("password"),
			response.get_field("name"),
			response.get_field("email"),))
	con.commit()
	cur.close()
	response.write("OK")



def indexPage(response):
    response.write(TemplateAPI.render('website.html', response, {}))


createTables()
server = Server('0.0.0.0', 80)
server.register("/", indexPage)
server.register("/userexists", checkUserExistsHandler)
server.register("/register", registerUserHandler)
server.register("/session/new", createSessionHandler)
server.register("/session/getdetails", getSelfDetails)
server.register("/squad/new", createSquad)
server.register("/squad/join", joinSquad)
server.register("/squad/getdetails", getSquadDetails)
server.register("/location/push", inserLocationsHandler)
server.register("/location/get", getLocationForUserHandler)
server.register("/session/notifications", getNotificationsHandler)
server.register("/session/notifications/new", insertNotificationHandler)
server.register("/session/notifications/squad/new", insertSquadNotificationsHandler)
server.register("/session/notifications/all/new", insertAllNotificationsHandler)
server.run(dummy)
