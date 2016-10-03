package net.ciphersink.remotecontrol;

import android.app.Service;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.os.BatteryManager;
import android.os.Bundle;
import android.os.IBinder;
import android.util.Log;

import net.ciphersink.remotecontrol.util.StatusUpdater;

import java.util.Timer;
import java.util.TimerTask;

/**
 * Created by xxx on 18/12/15.
 */
public class StatusService extends Service {

    LocationManager locationManager;
    boolean on_charge = false;
    private Timer mHeartbeatTimer = new Timer();

    private final BroadcastReceiver mOnChargeReceiver = new BroadcastReceiver() {
        public void onReceive(Context context, Intent intent) {
            on_charge = true;
            int chargePlug = intent.getIntExtra(BatteryManager.EXTRA_PLUGGED, -1);
            boolean usbCharge = chargePlug == BatteryManager.BATTERY_PLUGGED_USB;
            boolean acCharge = chargePlug == BatteryManager.BATTERY_PLUGGED_AC;
            if(usbCharge)
                StatusUpdater.sendMessage(context, "Device now on USB charge");
            else if(acCharge)
                StatusUpdater.sendMessage(context, "Device now on AC charge");
            else
                StatusUpdater.sendMessage(context, "Device now on charge");
        }
    };

    private final BroadcastReceiver mOffChargeReceiver = new BroadcastReceiver() {
        public void onReceive(Context context, Intent intent) {
            on_charge = false;
            StatusUpdater.sendMessage(context, "Device now running off battery");
        }
    };

    public StatusService() {

    }

    public void onCreate() {
        super.onCreate();
        StatusUpdater.sendMessage(getApplicationContext(), "Android StatusService started");

        this.registerReceiver(this.mOnChargeReceiver, new IntentFilter(Intent.ACTION_POWER_CONNECTED));
        this.registerReceiver(this.mOffChargeReceiver, new IntentFilter(Intent.ACTION_POWER_DISCONNECTED));
        setupPeriodicTimer();

        locationManager = (LocationManager) this.getSystemService(Context.LOCATION_SERVICE);
        locationManager.requestLocationUpdates(LocationManager.PASSIVE_PROVIDER, 0, 0, new MsgLocationListener());

    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        return(START_NOT_STICKY);
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        StatusUpdater.sendMessage(this, "Android StatusService stopped");
        stopPeriodicTimer();
    }


    //region HEARTBEAT_TIMER
    private void setupPeriodicTimer() {
        final Context self = this;
        TimerTask task = new TimerTask() {
            @Override
            public void run() {
                StatusUpdater.sendStatus(self);
            }
        };
        // schedules the task to be run in an interval
        mHeartbeatTimer.scheduleAtFixedRate(task, 5000, 15 * 60 * 1000);
    }

    private void stopPeriodicTimer() {
        //May call destroy before instantiating the mHeartbeatTimer
        mHeartbeatTimer.cancel();
        mHeartbeatTimer.purge();
    }
//endregion HEARTBEAT_TIMER

    @Override
    public IBinder onBind(Intent intent) {
        throw new UnsupportedOperationException("Not yet implemented");
    }





















    public class MsgLocationListener implements LocationListener
    {
        private double lastLat;
        private double lastLon;
        private long lastTransmitted;

        /*
         * Calculate distance between two points in latitude and longitude taking
         * into account height difference. If you are not interested in height
         * difference pass 0.0. Uses Haversine method as its base.
         *
         * lat1, lon1 Start point lat2, lon2 End point el1 Start altitude in meters
         * el2 End altitude in meters
         * @returns Distance in Meters
         */
        public double distance(double lat1, double lat2, double lon1,
                                      double lon2, double el1, double el2) {

            final int R = 6371; // Radius of the earth

            Double latDistance = Math.toRadians(lat2 - lat1);
            Double lonDistance = Math.toRadians(lon2 - lon1);
            Double a = Math.sin(latDistance / 2) * Math.sin(latDistance / 2)
                    + Math.cos(Math.toRadians(lat1)) * Math.cos(Math.toRadians(lat2))
                    * Math.sin(lonDistance / 2) * Math.sin(lonDistance / 2);
            Double c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
            double distance = R * c * 1000; // convert to meters

            double height = el1 - el2;

            distance = Math.pow(distance, 2) + Math.pow(height, 2);

            return Math.sqrt(distance);
        }



        @Override
        public void onLocationChanged(Location location)
        {
            Log.d("CNC", "PASSIVE LOCATION RECEIVED: " + location.getLongitude() + " : " + location.getLatitude() + " -- " + location.getTime());

            if(((lastTransmitted + (60 * 8 * 1000)) < location.getTime()) || (distance(lastLat, location.getLatitude(), lastLon, location.getLongitude(), 10, 10) > 200)) {
                lastLat = location.getLatitude();
                lastLon = location.getLongitude();
                lastTransmitted = location.getTime();
                StatusUpdater.sendLocation(getApplicationContext(), location.getLatitude(),
                        location.getLongitude(),
                        location.getSpeed(),
                        location.getAccuracy(),
                        location.getBearing());
            }
        }

        @Override
        public void onProviderDisabled(String arg0)
        {
            StatusUpdater.sendMessage(getApplicationContext(), "Location provider disabled");
        }

        @Override
        public void onProviderEnabled(String arg0)
        {
            Log.d("CNC", ":)");
        }

        @Override
        public void onStatusChanged(String arg0, int arg1, Bundle arg2)
        {
        }
    }

}
