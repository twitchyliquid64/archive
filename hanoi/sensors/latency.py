import platform,subprocess,re

def Ping(hostname,timeout=2):
    command="ping -i "+str(timeout)+" -c 3 " + hostname
    proccess = subprocess.Popen(command, stdout=subprocess.PIPE, shell=True)
    d = proccess.stdout.read()
    if '100% packet loss' in d:
        return -1
    if 'unknown host' in d:
        return -2
    ping = d.split("\n")[-2].split(" ")[3].split("/")[1]
    return float(ping)

p = Ping(sensor['dest'], 2)
if p == -1:
    state.ok = True
    state.last_result['val'] = 99999
    state.last_result['reason'] = 'No response'
elif p == -1:
    state.ok = True
    state.last_result['val'] = 99999
    state.last_result['reason'] = 'Unknown host'
else:
    state.ok = True
    state.last_result['val'] = p
write()
