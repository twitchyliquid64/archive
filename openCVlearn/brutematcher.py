import numpy as np
import cv2

cap = cv2.VideoCapture(0)


def bruteMatchKnn(img1, img2):
	# Initiate SIFT detector
	orb = cv2.ORB()
	# find the keypoints and descriptors with SIFT
	kp1, des1 = orb.detectAndCompute(img1,None)
	kp2, des2 = orb.detectAndCompute(img2,None)
	# Match descriptors.
	if des1 == None:
		return [], [], []
	if des2 == None:
		return [], [], []
	bf = cv2.BFMatcher()
	matches = bf.knnMatch(des1,des2, k=2)
	# Apply ratio test
	good = []
	try:
		for m,n in matches:
			if m.distance < 0.75*n.distance:
				good.append(m)
	except:
		pass
	return good, kp1, kp2




def bruteMatch(img1, img2):
	# Initiate SIFT detector
	orb = cv2.ORB()
	# find the keypoints and descriptors with SIFT
	kp1, des1 = orb.detectAndCompute(img1,None)
	kp2, des2 = orb.detectAndCompute(img2,None)
	# create BFMatcher object
	bf = cv2.BFMatcher(cv2.NORM_HAMMING, crossCheck=True)
	# Match descriptors.
	if des1 == None:
		return [], [], []
	if des2 == None:
		return [], [], []
	matches = bf.match(des1,des2)
	# Sort them in the order of their distance.
	matches = sorted(matches, key = lambda x:x.distance)
	return matches, kp1, kp2


def captureTemplate():
	while True:
	    ret, template = cap.read()
	    # Display the resulting frame
	    cv2.imshow('set template',template)
	    if cv2.waitKey(1) & 0xFF == ord('q'):
		break
	cv2.destroyAllWindows()
	return template


def loadTemplate(fname):
	return cv2.imread(fname, 1)


template = loadTemplate("acircles_pattern.png")

while True:
	ret, img = cap.read()
	matches, kp1, kp2 = bruteMatch(template, img) #get the matches, and the keypoints for each image
	print len(matches)
	count = 0

	for match in matches:
		x1, y1 = kp2[match.trainIdx].pt #get x and y of match on the CAMERA frame (not the template frame)
		cv2.rectangle(img, (int(x1), int(y1)), (int(x1)+10, int(y1)+10), (0, 255, count*3), 2)


	cv2.imshow('match',img)
	if cv2.waitKey(1) & 0xFF == ord('q'):
		break
	

