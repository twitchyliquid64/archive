import subprocess
def updateRepo():
	subprocess.Popen(["git", "pull", "origin", "master"])
