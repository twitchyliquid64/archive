import numpy as np
import cv2.cv as cv
import cv2

cap = cv2.VideoCapture(0)
template = cv2.imread('template.jpg',0)

while True:
	# Load the images
	ret, img = cap.read()

	# Convert them to grayscale
	imgg =cv2.cvtColor(img,cv2.COLOR_BGR2GRAY)

	# SURF extraction
	surf = cv2.SURF(5000)
	surf.upright = True
	kp, des = surf.detectAndCompute(imgg,None)

	img2 = cv2.drawKeypoints(img,kp,None,(255,0,0),4)

	# Display the resulting frame
	cv2.imshow('frame',img2)
	if cv2.waitKey(1) & 0xFF == ord('q'):
		break

# When everything done, release the capture
cap.release()
cv2.destroyAllWindows()
