package net.ciphersink.nightout.view;

import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.support.v4.app.FragmentManager;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.google.android.gms.maps.CameraUpdateFactory;
import com.google.android.gms.maps.GoogleMap;
import com.google.android.gms.maps.OnMapReadyCallback;
import com.google.android.gms.maps.SupportMapFragment;
import com.google.android.gms.maps.model.Circle;
import com.google.android.gms.maps.model.CircleOptions;
import com.google.android.gms.maps.model.LatLng;
import com.google.android.gms.maps.model.Marker;
import com.google.android.gms.maps.model.MarkerOptions;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.Interfaces;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.LocationRecord;
import net.ciphersink.nightout.model.SquadMember;

/**
 * Implements the Fragment users see when finding a mate. It encapsulates UI elements
 * and network code for displaying their loction on a map, battery percentage, lock type etc
 */
public class TrackerFragment extends Fragment  implements OnMapReadyCallback, Interfaces.MenuControllerInterface {
    // data model
    private SquadMember mTrackingIndividual;

    // ui
    private LinearLayout mMapContainer;
    private SupportMapFragment mMap;
    private TextView mLastSeenText;
    private TextView mProviderText;
    private TextView mBattText;
    private ImageView mLoadingImage;
    private Animation mLoadingAnimation;
    private GoogleMap mGoogleMap;

    // result from network call.
    private LocationRecord lastRecord;

    public TrackerFragment(SquadMember individual) {
        mTrackingIndividual = individual;
    }

    public TrackerFragment() {}

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(getString(R.string.tracker) + " " + mTrackingIndividual.getName());
        act.initialiseMenu(R.id.mainMenuRefresh, this); // tell MainActivity to show the refesh menu icon.
    }

    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);

        // init UI elements
        mMap = (SupportMapFragment) getActivity().getSupportFragmentManager()
                .findFragmentById(R.id.trackerFragMapContainer);

        mLastSeenText = (TextView) getView().findViewById(R.id.trackerFragLastSeenText);
        mProviderText = (TextView)getView().findViewById(R.id.trackerFragProviderText);
        mMapContainer = (LinearLayout)getView().findViewById(R.id.trackerFragMapContainer);
        mBattText = (TextView)getView().findViewById(R.id.trackerFragBattText);

        mLastSeenText.setText(getString(R.string.last_seen) + " " + getString(R.string.loading_default_message));

        //animate the tracking icon - such that it looks like a loading animation
        mLoadingImage = (ImageView)getView().findViewById(R.id.trackerFragLoadingImage);
        new UpdateLocationTask().execute();
    }


    /**
     * Places a google map fragment in its corresponding holder.
     */
    private void loadMap() {
        mLoadingAnimation.cancel();
        mLoadingImage.clearAnimation();
        mLoadingImage.setVisibility(View.INVISIBLE);
        mMapContainer.setVisibility(View.VISIBLE);
        //create a map fragment
        FragmentManager fm = getChildFragmentManager();
        mMap =  SupportMapFragment.newInstance();
        fm.beginTransaction().replace(R.id.trackerFragMap, mMap).commit();

        mMap.getMapAsync(this);
    }

    /**
     * Called by play services API when the google map is ready. We initialise
     * the map to mark the location of our mate and plot the accuracy of the lock.
     * @param map
     */
    @Override
    public void onMapReady(GoogleMap map) {
        // Add a marker in Sydney, Australia, and move the camera.
        LatLng loc = new LatLng(lastRecord.getLat(), lastRecord.getLon());
        map.moveCamera(CameraUpdateFactory.newLatLng(loc));
        map.moveCamera(CameraUpdateFactory.zoomTo(12));
        mGoogleMap = map;
        drawMarkerWithCircle(loc, lastRecord.getPrecision());
    }

    /**
     * Draws a circle around a position in google maps
     * @param position
     * @param radiusInMeters
     */
    private void drawMarkerWithCircle(LatLng position, double radiusInMeters){
        int strokeColor = 0xffff0000; //red outline
        int shadeColor = 0x44ff0000; //opaque red fill

        CircleOptions circleOptions = new CircleOptions().center(position).radius(radiusInMeters).fillColor(shadeColor).strokeColor(strokeColor).strokeWidth(8);
        Circle mCircle = mGoogleMap.addCircle(circleOptions);

        MarkerOptions markerOptions = new MarkerOptions().position(position);
        Marker mMarker = mGoogleMap.addMarker(markerOptions);
    }

    /**
     * Called from MainActivity when a menu icon is pressed - in this case exclusively the refresh button.
     * @param id
     */
    @Override
    public void menuClicked(int id) {
        if (id == R.id.mainMenuRefresh) {
            new UpdateLocationTask().execute();
        }
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_tracker, container, false);
    }


    /**
     * Encapsulates network logic for downloading the current location of the squad member.
     */
    private class UpdateLocationTask extends AsyncTask<Void, Void, Void> {

        @Override
        protected void onPreExecute() {
            // hide map and show loading animation
            mLoadingImage.setVisibility(View.VISIBLE);
            mMapContainer.setVisibility(View.INVISIBLE);
            mLoadingAnimation = AnimationUtils.loadAnimation(getActivity(), R.anim.tracker_loadingrotate);
            mLoadingImage.startAnimation(mLoadingAnimation);

            MenuItem refreshButton = ((MainActivity)getActivity()).getRefreshMenuButton();
            refreshButton.setEnabled(false);
        }

        @Override
        protected Void doInBackground(Void... dummy) {
            String sessKey = ((MainActivity)getActivity()).getSession().getKey();
            lastRecord = LocationRecord.getLatestById(sessKey, mTrackingIndividual.getId());

            Log.d(Constants.MAD, "Got record data of age: " + lastRecord.getAge());

            try {
                Thread.sleep(1200, 0);
            } catch (InterruptedException e) {

            }
            return null;
        }


        @Override
        protected void onPostExecute(Void dummy) {
            MenuItem refreshButton = ((MainActivity)getActivity()).getRefreshMenuButton();
            refreshButton.setEnabled(true);
            // update map with data
            loadMap();
            // update UI elements with data
            mLastSeenText.setText(getString(R.string.last_seen) + " " + lastRecord.getAge() + " " + getString(R.string.seconds_ago));
            mProviderText.setText(lastRecord.getProvider() + " " + getString(R.string.lock));
            mBattText.setText(getString(R.string.battery) + " " + lastRecord.getBattery() + "%");
        }
    }
}