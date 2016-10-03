package net.ciphersink.remotecontrol.util;

import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.SharedPreferences;
import android.net.Uri;
import android.os.BatteryManager;
import android.preference.PreferenceManager;
import android.util.Log;

import com.koushikdutta.async.http.AsyncHttpClient;
import com.koushikdutta.async.http.AsyncHttpRequest;
import com.koushikdutta.async.http.AsyncHttpResponse;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;

/**
 * Created by xxx on 15/12/15.
 */
public class StatusUpdater {

    private static String getURL(Context context){
        SharedPreferences preferences = PreferenceManager.getDefaultSharedPreferences(context);
        return preferences.getString("server_address", "");
    }
    private static String getAPIKey(Context context){
        SharedPreferences preferences = PreferenceManager.getDefaultSharedPreferences(context);
        return preferences.getString("entity_api_key", "");
    }


    public static void sendMessage(Context context, String message){
        fetch_battery_info(context);
        String url = "";
        try {
            url = "http://" + getURL(context) + "/e/status?key=" + getAPIKey(context) + "&status=" + URLEncoder.encode(message, "UTF-8") + "&style=progress-linear&stylemeta=" + Integer.toString(batPercent);
        }catch (UnsupportedEncodingException e){
            e.printStackTrace();
        }
        Log.d("CNC", "Sending req to " + url);
        AsyncHttpClient.getDefaultInstance().executeString(new AsyncHttpRequest(Uri.parse(url), "GET"), new AsyncHttpClient.StringCallback() {
            // Callback is invoked with any exceptions/errors, and the result, if available.
            @Override
            public void onCompleted(Exception e, AsyncHttpResponse response, String result) {
                if (e != null) {
                    e.printStackTrace();
                    return;
                }
                Log.d("CNC", "Response: " + result);
            }
        });
    }

    public static void sendLocation(Context context, double lat, double lon, float speed, float acc, float course){
        String url = "http://" + getURL(context) + "/e/loc?key=" + getAPIKey(context);
        url += "&lat=" + lat;
        url += "&lon=" + lon;
        url += "&kph=" + speed;
        url += "&course=" + (int)course;
        url += "&acc=" + (int)acc;

        Log.d("CNC", "Sending req to " + url);
        AsyncHttpClient.getDefaultInstance().executeString(new AsyncHttpRequest(Uri.parse(url), "GET"), new AsyncHttpClient.StringCallback() {
            // Callback is invoked with any exceptions/errors, and the result, if available.
            @Override
            public void onCompleted(Exception e, AsyncHttpResponse response, String result) {
                if (e != null) {
                    e.printStackTrace();
                    return;
                }
                Log.d("CNC", "Response: " + result);
            }
        });
    }

    public static void sendStatus(Context context){
        fetch_battery_info(context);
        String msg = "Discharging " + batPercent + "%";
        if (mIsCharging)msg = "Charging " + batPercent + "%";

        String url = "";
        try {
            url = "http://" + getURL(context) + "/e/status?key=" + getAPIKey(context) + "&status=" + URLEncoder.encode(msg, "UTF-8") + "&style=progress-linear&stylemeta=" + Integer.toString(batPercent);
        }catch (UnsupportedEncodingException e){
            e.printStackTrace();
        }
        Log.d("CNC", "Sending req to " + url);
        AsyncHttpClient.getDefaultInstance().executeString(new AsyncHttpRequest(Uri.parse(url), "GET"), new AsyncHttpClient.StringCallback() {
            // Callback is invoked with any exceptions/errors, and the result, if available.
            @Override
            public void onCompleted(Exception e, AsyncHttpResponse response, String result) {
                if (e != null) {
                    e.printStackTrace();
                    return;
                }
                Log.d("CNC", "Response: " + result);
            }
        });
    }


    public static int batPercent;
    public static boolean mIsCharging;
    public static void fetch_battery_info(Context context)
    {
        IntentFilter ifilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        Intent batteryStatus = context.registerReceiver(null, ifilter);
        int status = batteryStatus.getIntExtra(BatteryManager.EXTRA_STATUS, -1);
        boolean isCharging = status == BatteryManager.BATTERY_STATUS_CHARGING ||
                status == BatteryManager.BATTERY_STATUS_FULL;

        if(isCharging)
            mIsCharging = true;
        else
            mIsCharging = false;

        //Determine curernt battery level
        int level = batteryStatus.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        int scale = batteryStatus.getIntExtra(BatteryManager.EXTRA_SCALE, -1);
        float batteryPct = level / (float)scale;
        batPercent = (int)(batteryPct*100);
    }
}
