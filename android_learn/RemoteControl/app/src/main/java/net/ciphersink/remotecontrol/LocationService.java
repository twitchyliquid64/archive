package net.ciphersink.remotecontrol;

import android.app.Notification;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.location.Location;
import android.os.Build;
import android.os.Bundle;
import android.os.IBinder;
import android.os.PowerManager;
import android.util.Log;

import com.google.android.gms.common.ConnectionResult;
import com.google.android.gms.common.api.GoogleApiClient;
import com.google.android.gms.location.LocationRequest;
import com.google.android.gms.location.LocationListener;
import com.google.android.gms.location.LocationServices;

import net.ciphersink.remotecontrol.util.StatusUpdater;

/**
 * Created by xxx on 19/12/15.
 */
public class LocationService extends Service implements GoogleApiClient.ConnectionCallbacks,
                                                        GoogleApiClient.OnConnectionFailedListener,
                                                        LocationListener {
    private PowerManager.WakeLock mWakeLock;
    private GoogleApiClient mGoogleApiClient;
    private LocationRequest mLocationRequest;

    public LocationService(){

    }

    @Override
    public void onCreate() {
        super.onCreate();
        StatusUpdater.sendMessage(getApplicationContext(), "Android LocationService started");
        setupToRecieveLocations();
        acquireWakelock();
        setupNotification();
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        return(START_NOT_STICKY);
    }

    @Override
    public void onDestroy() {
        StatusUpdater.sendMessage(getApplicationContext(), "Android LocationService stopped");
        stopRecievingLocations();
        releaseWakeLock();
        cancelNotification();
    }

    @Override
    public IBinder onBind(Intent intent) {
        throw new UnsupportedOperationException("Not yet implemented");
    }



    //region LOCATION_LISTENER_IMPLEMENTATION

    @Override
    public void onConnected(Bundle bundle) {
        LocationServices.FusedLocationApi.requestLocationUpdates(mGoogleApiClient, mLocationRequest, this);
    }

    @Override
    public void onConnectionSuspended(int i) {
        // The connection to Google Play services was lost for some reason. We call connect() to
        // attempt to re-establish the connection.
        mGoogleApiClient.connect();
    }

    @Override
    public void onLocationChanged(Location location) {
        Log.d("CNC", "LOCATION RECEIVED: " + location.getLongitude() + " : " + location.getLatitude() + " -- " + location.getTime());
        StatusUpdater.sendLocation(getApplicationContext(), location.getLatitude(),
                location.getLongitude(),
                location.getSpeed(),
                location.getAccuracy(),
                location.getBearing());
    }

    @Override
    public void onConnectionFailed(ConnectionResult connectionResult) {
        StatusUpdater.sendMessage(getApplicationContext(), "PlayService location connect failed");
    }



    void setupToRecieveLocations() {

        mGoogleApiClient = new GoogleApiClient.Builder(this)
                .addConnectionCallbacks(this)
                .addOnConnectionFailedListener(this)
                .addApi(LocationServices.API)
                .build();
        mGoogleApiClient.connect();
        mLocationRequest = new LocationRequest();

        // Sets the desired interval for active location updates. This interval is
        // inexact. You may not receive updates at all if no location sources are available, or
        // you may receive them slower than requested. You may also receive updates faster than
        // requested if other applications are requesting location at a faster interval.
        mLocationRequest.setInterval(10 * 1000);

        // Sets the fastest rate for active location updates. This interval is exact, and your
        // application will never receive updates faster than this value.
        mLocationRequest.setFastestInterval(5 * 1000);
        mLocationRequest.setPriority(LocationRequest.PRIORITY_HIGH_ACCURACY);
        //mLocationRequest.setSmallestDisplacement(15);
    }

    public void stopRecievingLocations()
    {
        if(mGoogleApiClient != null && mGoogleApiClient.isConnected())
            LocationServices.FusedLocationApi.removeLocationUpdates(mGoogleApiClient, this);
    }

//endregion LOCATION_LISTENER_IMPLEMENTATION


    //region WAKELOCK
    private void acquireWakelock() {
        PowerManager pm = (PowerManager) this.getSystemService(Context.POWER_SERVICE);
        mWakeLock = pm.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "CNCLocationServiceWakeLock");
        mWakeLock.acquire();
    }

    private void releaseWakeLock() {
        if(mWakeLock != null && mWakeLock.isHeld())
            mWakeLock.release();
    }
//endregion WAKELOCK


    //region NOTIFICATION
    private void setupNotification()
    {
        Notification.Builder mNotificationBuild = new Notification.Builder(this)
                .setSmallIcon(R.drawable.ic_media_play)
                .setContentTitle("CNC")
                .setContentText("LocationService running")
                .setContentIntent(PendingIntent.getActivity(this, 0,
                        new Intent(this, ControlPanel.class), 0));
        Notification mNotification;
        if (Build.VERSION.SDK_INT < 16){
            mNotification = mNotificationBuild.getNotification();
        } else {
            mNotification= mNotificationBuild.build();
        }

        startForeground(1337, mNotification);
    }

    private void cancelNotification()
    {
        NotificationManager Nm = (NotificationManager)getSystemService(NOTIFICATION_SERVICE);
        Nm.cancel(1337);
    }
//endregion NOTIFICATION

}


