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
 * Encapsulates the Notification Data model, and all factory methods associated with it.
 */
public class Notification {
    private String mContent;
    private String mType;
    private String mSubLine;

    /**
     * Construct a Notification with the given information
     * @param content
     * @param type
     * @param subLine
     */
    public Notification(String content, String type, String subLine) {
        this.mContent = content;
        this.mType = type;
        this.mSubLine = subLine;
    }

    public String getContent() {
        return mContent;
    }

    public String getType() {
        return mType;
    }

    public String getSubLine() {
        return mSubLine;
    }

    /**
     * Transmits the notification to the server, such that it will appear for a given user.
     * @param sessionKey
     * @param userID
     * @return Success == true
     */
    public boolean send(String sessionKey, int userID) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.USER_ID, String.valueOf(userID)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.CONTENT, String.valueOf(mContent)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.TYPE, String.valueOf(mType)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SUBLINE, String.valueOf(mSubLine)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.SEND_NOTIFICATION, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return false;
        }
        return true;
    }

    /**
     * Transmits a notification to the server, such that an entire squad can read it.
     * @param sessionKey
     * @param squadID
     * @return success == true
     */
    public boolean sendToSquad(String sessionKey, int squadID) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SQUAD_ID, String.valueOf(squadID)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.CONTENT, String.valueOf(mContent)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.TYPE, String.valueOf(mType)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SUBLINE, String.valueOf(mSubLine)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.SEND_SQUAD_NOTIFICATION, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return false;
        }
        return true;
    }


    /**
     * Transmits a notifications to all individuals in all participating squads.
     * @param sessionKey
     * @return success == true
     */
    public boolean sendToAll(String sessionKey) {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(sessionKey)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.CONTENT, String.valueOf(mContent)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.TYPE, String.valueOf(mType)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.SUBLINE, String.valueOf(mSubLine)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.SEND_ALL_NOTIFICATION, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return false;
        }
        return true;
    }

    /**
     * Gets all notifications associated with the user in the given session
     * @param session
     * @return Arraylist of notifications
     */
    public static final ArrayList<Notification> getNotifications(Session session) {
        ArrayList<Notification> notifications = new ArrayList<Notification>();

        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.KEY, String.valueOf(session.getKey())));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.GET_NOTIFICATIONS, params);

        if(request.didError())return null;
        if(request.getData().equals(Constants.NET.STANDARD_ERROR)) {
            return null;
        }

        try {
            JSONArray obj = new JSONArray(request.getData());

            for (int i = 0; i < obj.length(); i++) {
                JSONArray data = obj.getJSONArray(i);
                notifications.add(new Notification(data.getString(3), data.getString(2), data.getString(5)));
            }
        }catch (JSONException e) {
            //request errorred - return null to signal
            //should never happen - we already checked for error
            Log.e(Constants.MAD, e.toString());
            return null;
        }

        return notifications;
    }
}
