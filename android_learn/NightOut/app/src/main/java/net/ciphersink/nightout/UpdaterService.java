package net.ciphersink.nightout;

import android.app.Service;
import android.content.Intent;
import android.content.IntentFilter;
import android.location.Location;
import android.location.LocationListener;
import android.location.LocationManager;
import android.os.AsyncTask;
import android.os.BatteryManager;
import android.os.Bundle;
import android.os.Handler;
import android.os.IBinder;
import android.util.Log;

import net.ciphersink.nightout.model.LocationDispatch;
import net.ciphersink.nightout.model.Session;


/**
 * Transmits the users location to the server periodically (defined by Constants.LOCATION_FIDELITY milliseconds)
 * This service is started when MainActivity is created, and stopped when destroyed.
 */
public class UpdaterService extends Service {
    private String mSessionKey;
    private Handler mPeriodicHandler;
    private LocationManager mLocationManager;
    private Location mLastLocation;

    private final LocationListener mLocationListener = new LocationListener() {
        @Override
        public void onLocationChanged(final Location location) {
            //your code here
            Log.d(Constants.MAD, "UpdaterService.LocationListener.onLocationChanged(): " + location.toString());
            mLastLocation = location;
        }

        @Override
        public void onStatusChanged(String provider, int status, Bundle extras) {
            Log.d(Constants.MAD, "UpdaterService.LocationListener.onStatusChanged(): " + provider);
        }

        /**
         * Implemented for completeness sake, not used in implementation
         * @param provider
         */
        @Override
        public void onProviderEnabled(String provider) {

        }

        /**
         * Implemented for completeness sake, not used in implementation
         * @param provider
         */
        @Override
        public void onProviderDisabled(String provider) {

        }
    };

    public UpdaterService() {
    }

    private void update() {
        Log.d(Constants.MAD, "UpdaterService doing update");

        if((mLastLocation == null) ||
                ((System.currentTimeMillis() - mLastLocation.getTime() > (Constants.LOCATION_FIDELITY * 2)))) {
            mLastLocation = mLocationManager.getLastKnownLocation(LocationManager.NETWORK_PROVIDER);
        }

        if (mLastLocation != null) {
            Log.d(Constants.MAD, mLastLocation.toString());
            new DispatchTask().execute();
        }
    }

    @Override
    public void onCreate() {
        super.onCreate();
        Log.d(Constants.MAD, "UpdaterService onCreate() called");

        mLocationManager = (LocationManager) getSystemService(LOCATION_SERVICE);
        mLocationManager.requestLocationUpdates(LocationManager.GPS_PROVIDER, Constants.LOCATION_FIDELITY, 100, mLocationListener);

        update();
        mPeriodicHandler = new Handler();
        mPeriodicHandler.postDelayed(new Runnable() {
            @Override
            public void run() {
                update();
                mPeriodicHandler.postDelayed(this, Constants.LOCATION_FIDELITY);
            }
        }, Constants.LOCATION_FIDELITY);
    }


    /**
     * Returns the current battery percentage. Its separated from control flow to maintain functional separation.
     * @return Percentage of battery remaining.
     */
    private float getBatteryLevel() {
        Intent batteryIntent = registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        int level = batteryIntent.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        int scale = batteryIntent.getIntExtra(BatteryManager.EXTRA_SCALE, -1);

        // Error checking that probably isn't needed but I added just in case.
        if(level == -1 || scale == -1) {
            return 50.0f;
        }

        return ((float)level / (float)scale) * 100.0f;
    }

    @Override
    public IBinder onBind(Intent intent) {
        throw new UnsupportedOperationException("Not yet implemented");
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId){
        mSessionKey = intent.getExtras().getString(Constants.KEYS.SESSIONKEY);
        return START_REDELIVER_INTENT;
    }

    @Override
    public void onDestroy(){
        super.onDestroy();
        Log.d(Constants.MAD, "UpdaterService onDestroy() called");
    }

    /**
     * Represents the background task that sends to the network.
     */
    private class DispatchTask extends AsyncTask<Void, Void, Void> {

        @Override
        protected Void doInBackground(Void... dummy) {
            float battLevel = getBatteryLevel();
            LocationDispatch.transmit(mSessionKey, mLastLocation, (int)battLevel);
            return null;
        }
    }
}
