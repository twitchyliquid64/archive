package net.ciphersink.nightout.model;

import android.util.Log;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.util.EasyHttp;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

/**
 * Encapsulates the static methods used to construct squad objects given certain data.
 */
public class SquadFactory {

    /**
     * returns null if accesskey invalid or already joined, otherwise,
     * it returns a squad object represented the squad that was just joined.
     * @param squadAccesskey
     * @param sessionKey
     * @return Squad or null
     */
    public static Squad joinSquad(String squadAccesskey, String sessionKey) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SQUAD_KEY, String.valueOf(squadAccesskey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.JOIN_SQUAD, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        int squadID = Integer.parseInt(request.getData());
        return getByKey(squadID, sessionKey);
    }

    /**
     * Returns a valid and populated squad object.
     * @param squadName
     * @param sessionKey
     * @return Squad
     */
    public static Squad createSquad(String squadName, String sessionKey) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.NAME, String.valueOf(squadName)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.NEW_SQUAD, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        int squadID = Integer.parseInt(request.getData());
        return getByKey(squadID, sessionKey);
    }

    /**
     * Returns a squad object, if the session user is a member of the squad.
     * @param squadID
     * @param sessionKey
     * @return Squad or null
     */
    public static Squad getByKey(int squadID, String sessionKey) {
        Squad squad;

        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SQUAD_ID, String.valueOf(Integer.toString(squadID))));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.GET_SQUAD_DETAILS, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        try {
            JSONArray squadobject = new JSONObject(request.getData()).getJSONArray("details");
            squad = new Squad(squadID, squadobject.getString(1), squadobject.getString(2));

            JSONArray squadMembersObject = new JSONObject(request.getData()).getJSONArray("members");
            for (int i = 0; i < squadMembersObject.length(); i++) {
                JSONArray squadMember = squadMembersObject.getJSONArray(i);
                squad.addMember(new SquadMember(squadMember.getInt(0), squadMember.getString(2), squadMember.getString(1)));
            }
        }catch (JSONException e) {
            //request errorred - return null to signal
            //should never happen - we already checked for error
            return null;
        }
        return squad;
    }

    /**
     * Gets all squads which the user is a part of.
     * @param s
     * @return ArrayList<Squad>
     */
    public static ArrayList<Squad> getSquads(Session s) {
        String sessionKey = s.getKey();
        ArrayList<Integer> squadIds = s.getSquadIds();
        ArrayList<Squad> output = new ArrayList<Squad>();

        for(int i = 0; i < squadIds.size(); i++) {
            Log.d(Constants.MAD, "Now loading squad details: " + i);
            output.add(getByKey(squadIds.get(i).intValue(), sessionKey));
        }

        return output;
    }

}
