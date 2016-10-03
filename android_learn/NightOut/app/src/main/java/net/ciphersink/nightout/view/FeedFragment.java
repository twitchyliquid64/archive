package net.ciphersink.nightout.view;

import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.support.v4.widget.SwipeRefreshLayout;
import android.support.v7.widget.LinearLayoutManager;
import android.support.v7.widget.RecyclerView;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Notification;
import net.ciphersink.nightout.model.Session;

import java.util.ArrayList;


/**
 * Implements the UI responsible for displaying  the feed (Notifications).
 */
public class FeedFragment extends Fragment {

    private SwipeRefreshLayout mSwipeLayout;
    private RecyclerView mRecyclerView;
    private FeedAdapter mAdapter;


    public FeedFragment() {
        // Required empty public constructor
    }

    /**
     * Used internally to refresh the data model from the server and display.
     */
    private void refresh() {
        Log.d(Constants.MAD, "FeedFragment.refresh()");
        new LoadNotificationsTask().execute();
    }

    /**
     * Initialises the UI
     * @param savedInstanceState
     */
    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);
        mSwipeLayout = (SwipeRefreshLayout)getView().findViewById(R.id.feedFragSwipeLayout);
        mRecyclerView = (RecyclerView)getView().findViewById(R.id.feedFragRecyclerView);

        mRecyclerView.setHasFixedSize(true);
        mRecyclerView.setLayoutManager(new LinearLayoutManager(getActivity()));

        mAdapter = new FeedAdapter(new ArrayList<Notification>(), this);
        mRecyclerView.setAdapter(mAdapter);

        mSwipeLayout.setOnRefreshListener(new SwipeRefreshLayout.OnRefreshListener() {
            @Override
            public void onRefresh() {
                refresh();
            }
        });
        refresh();
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(getString(R.string.feed));
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_feed, container, false);
    }


    private class LoadNotificationsTask extends AsyncTask<Void, Void, Void> {

        private ArrayList<Notification> mNotifications;

        @Override
        protected void onPreExecute() {
            if (!mSwipeLayout.isRefreshing())mSwipeLayout.setRefreshing(true);
        }

        @Override
        protected Void doInBackground(Void... dummy) {
            Session session = ((MainActivity)getActivity()).getSession();
            mNotifications = Notification.getNotifications(session);
            Log.d(Constants.MAD, ""+mNotifications.size());
            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            mSwipeLayout.setRefreshing(false);
            mAdapter.setNotifications(mNotifications);
            mAdapter.notifyDataSetChanged();
        }
    }
}
