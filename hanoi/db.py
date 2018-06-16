import sqlite3
import threading
import json
from os.path import expanduser, join

home = expanduser("~")
conn = None



create_sql = '''
CREATE TABLE IF NOT EXISTS sensor_runs(
    id      INTEGER PRIMARY KEY   AUTOINCREMENT,
    name    VARCHAR(64),
    ok      BOOLEAN,
    result  TEXT,
    t TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)


'''

def check():
    curs = conn.cursor()
    curs.execute(create_sql)
    conn.commit()
    curs.close()

def janitorThread():
    threading.Timer(60 * 4, janitorThread).start()
    curs = conn.cursor()
    curs.execute("DELETE FROM sensor_runs WHERE t <= date('now','-2 day');")
    curs.execute("VACUUM;")
    conn.commit()
    curs.close()

def db_start():
    global conn
    conn = sqlite3.connect(":memory:", check_same_thread = False)
    check()
    janitorThread()

def write_sensor_result(name, ok, result):
    curs = conn.cursor()
    curs.execute('INSERT INTO sensor_runs (name, ok, result) VALUES (?, ?, ?)', (name, ok, json.dumps(result),))
    conn.commit()
    curs.close()

def query(name):
    curs = conn.cursor()
    curs.execute("""SELECT ok, result, CAST(strftime('%s', t) as decimal) FROM sensor_runs WHERE name = ? ORDER BY t""", (str(name),))
    res = curs.fetchall()
    curs.close()
    return res
