package net.ciphersink.exercise5;

/**
 * Created by xxx on 25/08/15.
 */
public class TrainData {
    private String mPlatform;
    private int mArrivalTime;
    private String mStatus;
    private String mDestination;
    private String mDestinationTime;

    TrainData(String platform, int arrivalTime, String status, String destination, String destinationTime)
    {
        mPlatform = platform;
        mArrivalTime = arrivalTime;
        mStatus = status;
        mDestination = destination;
        mDestinationTime = destinationTime;
    }

    public void setPlatform(String mPlatform) {
        this.mPlatform = mPlatform;
    }

    public void setArrivalTime(int mArrivalTime) {
        this.mArrivalTime = mArrivalTime;
    }

    public void setStatus(String mStatus) {
        this.mStatus = mStatus;
    }

    public void setDestination(String mDestination) {
        this.mDestination = mDestination;
    }

    public void setDestinationTime(String mDestinationTime) {
        this.mDestinationTime = mDestinationTime;
    }


    public String getPlatform() {
        return mPlatform;
    }

    public int getArrivalTime() {
        return mArrivalTime;
    }

    public String getStatus() {
        return mStatus;
    }

    public String getDestination() {
        return mDestination;
    }

    public String getDestinationTime() {
        return mDestinationTime;
    }
}
