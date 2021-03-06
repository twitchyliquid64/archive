import numpy as np
import cv2.cv as cv
import cv2

cap = cv2.VideoCapture(0)

hue = 0
sat = 0
last = 0

def update_hue(p):
	hue = p
def update_sat(p):
	sat = p
def update_last(p):
	last = p
	

cv.NamedWindow("control",cv.CV_WINDOW_AUTOSIZE)
cv2.createTrackbar("hue", "control", 0, 179, update_hue)
cv2.createTrackbar("sat", "control", 0, 255, update_sat)
cv2.createTrackbar("value", "control", 0, 255, update_last)


while True:
	ret, frameOriginal = cap.read()

	frame = cv2.cvtColor(frameOriginal, cv.CV_BGR2HSV)#convert to HSV

	#open and close to reduce image noise
	#frame = cv2.morphologyEx(frame, cv2.MORPH_OPEN, cv2.getStructuringElement(cv2.MORPH_ELLIPSE,(5,5)))
	#frame = cv2.morphologyEx(frame, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE,(5,5)))

	#filter image based on color
        thresh = cv2.inRange(frame, np.array((hue, sat, last)), np.array((hue+230, sat+180, last+180)))

        #Find contours in the threshold image
        contours,hierarchy = cv2.findContours(thresh,cv2.RETR_LIST,cv2.CHAIN_APPROX_SIMPLE)

        #Finding contour with maximum area and store it as best_cnt
        max_area = 0
	best_cnt = None
        for cnt in contours:
            area = cv2.contourArea(cnt)
            if area > max_area:
                max_area = area
                best_cnt = cnt

        #Finding centroids of best_cnt and draw a circle there
	if best_cnt!= None:
		M = cv2.moments(best_cnt)
		cx,cy = int(M['m10']/M['m00']), int(M['m01']/M['m00'])
		cv2.circle(frameOriginal,(cx,cy),10,255,-1)

	cv2.imshow('control',thresh)
	cv2.imshow('frame2',frameOriginal)
	if cv2.waitKey(1) & 0xFF == ord('q'):
		break

# When everything done, release the capture
cap.release()
cv2.destroyAllWindows()
