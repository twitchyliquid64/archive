package net.ciphersink.nightout.model;


import android.content.Intent;
import android.location.Location;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.util.EasyHttp;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;

import java.util.ArrayList;

/**
 * Encapsulates the logic for the transmission of location records to the server.
 */
public class LocationDispatch {


    /**
     * Transmits a location Records to the server
     * @param sessionKey session which the location is associated with
     * @param location location object to transmit
     * @param battLevel battery percentage
     * @return success == true
     */
    public static boolean transmit(String sessionKey, Location location, int battLevel) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.LATITUDE, String.valueOf(location.getLatitude())));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.LONGITUDE, String.valueOf(location.getLongitude())));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.ACCURACY, String.valueOf(location.getAccuracy())));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.PROVIDER, String.valueOf(location.getProvider())));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.BATTERY, String.valueOf(battLevel)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.TRANSMIT_LOCATION, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return false;
        }
        return true;
    }
}
