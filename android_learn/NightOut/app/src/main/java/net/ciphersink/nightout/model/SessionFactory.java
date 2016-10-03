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
 * Encapsulates static methods used to create session objects - based on given data.
 */
public class SessionFactory {

    /**
     * Returns a valid and populated session object if the session key was valid.
     * @param key
     * @return valid session or null
     */
    public static Session makeFromKey(String key) {
        Session session;
        ArrayList<Integer> squadIds = new ArrayList<Integer>();

        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(key)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.GET_SESSION, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        try {
            JSONArray userobject = new JSONObject(request.getData()).getJSONArray("u");

            JSONArray squadList = new JSONObject(request.getData()).getJSONArray("squads");
            for (int i = 0; i < squadList.length(); i++) {
                squadIds.add(squadList.getJSONArray(i).getInt(0));
            }

            session = new Session(key, userobject.getInt(0), userobject.getString(1), userobject.getString(2), userobject.getString(3), squadIds);
        }catch (JSONException e) {
            return null;
        }
        return session;
    }

    /**
     * Returns a valid and populated session object if the username and password are correct.
     * @param username
     * @param password
     * @return valid session or null
     */
    public static Session login(String username, String password) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.USERNAME, String.valueOf(username)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.PASSWORD, String.valueOf(password)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.CREATE_SESSION, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        return makeFromKey(request.getData());
    }

}
