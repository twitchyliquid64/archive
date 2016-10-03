package net.ciphersink.exercise6;

/**
 * Created by xxx on 1/09/15.
 */
public class FriendData {
    private int mID;
    private String mName;
    private String mOccupation;
    private String mCity;
    private long mFriendSince;

    FriendData(String name, String occupation, String city, long friendSince)
    {
        mID = -1;
        mName = name;
        mOccupation = occupation;
        mCity = city;
        mFriendSince = friendSince;
    }

    FriendData(int id, String name, String occupation, String city, long friendSince)
    {
        mID = id;
        mName = name;
        mOccupation = occupation;
        mCity = city;
        mFriendSince = friendSince;
    }

    public int getID() {
        return mID;
    }

    public void setID(int mID) {
        this.mID = mID;
    }

    public String getName() {
        return mName;
    }

    public void setName(String mName) {
        this.mName = mName;
    }

    public String getOccupation() {
        return mOccupation;
    }

    public void setOccupation(String mOccupation) {
        this.mOccupation = mOccupation;
    }

    public String getCity() {
        return mCity;
    }

    public void setCity(String mCity) {
        this.mCity = mCity;
    }

    public long getFriendSince() {
        return mFriendSince;
    }

    public void setFriendSince(long mFriendSince) {
        this.mFriendSince = mFriendSince;
    }

}
