package net.ciphersink.nightout.model;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.util.EasyHttp;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

/**
 * Created by xxx on 27/10/15.
 */
public class LocationRecord {
    private Double mLat;
    private Double mLon;
    private int mPrecision;
    private String mProvider;
    private int mAge;
    private int mBattery;

    public Double getLat() {
        return mLat;
    }

    public void setLat(Double mLat) {
        this.mLat = mLat;
    }

    public Double getLon() {
        return mLon;
    }

    public void setLon(Double mLon) {
        this.mLon = mLon;
    }

    public int getPrecision() {
        return mPrecision;
    }

    public void setPrecision(int mPrecision) {
        this.mPrecision = mPrecision;
    }

    public String getProvider() {
        return mProvider;
    }

    public void setProvider(String mProvider) {
        this.mProvider = mProvider;
    }

    public int getAge() {
        return mAge;
    }

    public void setAge(int mAge) {
        this.mAge = mAge;
    }

    public int getBattery() {
        return mBattery;
    }

    public void setBattery(int mBattery) {
        this.mBattery = mBattery;
    }

    public static LocationRecord getLatestById(String sessionKey, int userID) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.USER_ID, String.valueOf(Integer.toString(userID))));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.GET_USER_LOCATION, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        LocationRecord location;
        try {
            JSONArray obj = new JSONArray(request.getData());
            location = new LocationRecord();
            location.setLat(obj.getDouble(1));
            location.setLon(obj.getDouble(2));
            location.setPrecision(obj.getInt(3));
            location.setBattery(obj.getInt(4));
            location.setAge((int)obj.getDouble(5));
            location.setProvider(obj.getString(6));

            return location;

        }catch (JSONException e) {
            //request errorred - return null to signal
            //should never happen - we already checked for error
            return null;
        }
    }
}
