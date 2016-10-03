#!/usr/bin/python

from pygame import mixer # Load the required library
import time
import librosa
from sendUDP import send

def wheel(wheelPos):
	wheelPos = 255 - wheelPos
	if wheelPos < 85:
		return (255 - (wheelPos * 3), 0, wheelPos*3)
	if wheelPos < 170:
		wheelPos -= 85
		return (0, wheelPos * 3, 255 - (wheelPos*3))
	wheelPos -= 170
	return (wheelPos * 3, 255 - (wheelPos*3), 0)

NUM_PIX = 140
class Chaser(object):
	def __init__(self):
		self.out = [(0,0,0) for z in range(NUM_PIX)]
		self.q = 0
		self.inteval = 5
	def step(self, color=[200,0,0]):
		self.q = (self.q + 1) % self.inteval
		for x in xrange(self.q, NUM_PIX, self.inteval):
			self.out[x] = wheel(x)
		send(self.out)
		for x in xrange(self.q, NUM_PIX, 3):
			self.out[x] = (0,0,0)

print "-=== NOW ANALYSING MP3 ===-"

y, sr = librosa.load('bigger_than_love.mp3')
tempo, beat_frames = librosa.beat.beat_track(y=y, sr=sr)
beat_times = librosa.frames_to_time(beat_frames, sr=sr)

#print beat_times

mixer.init()
mixer.music.load('bigger_than_love.mp3')
mixer.music.play()
c = Chaser()

while True:
	p = mixer.music.get_pos() - 170
	#print p, beat_times[0]*1000
	diff = abs(p - (beat_times[0]*1000))
	if diff < 1 or ((p+7)>beat_times[0]*1000):
		print "BEAT", diff
		beat_times = beat_times[1:]
		c.step()
	if p > (1000*115):
		break


