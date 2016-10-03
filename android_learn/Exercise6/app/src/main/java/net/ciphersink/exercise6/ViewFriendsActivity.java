package net.ciphersink.exercise6;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.support.v7.app.ActionBarActivity;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Toast;

import java.util.ArrayList;

/**
 * Created by xxx on 1/09/15.
 */
public class ViewFriendsActivity extends ActionBarActivity  implements View.OnClickListener{

    private ArrayList<FriendData> mFriends;
    private ListView mFriendListView;
    private FriendAdapter mFriendAdapter;
    private FriendDatabaseHelper mDb;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_friendviewer);

        Intent intent = getIntent();
        mDb = new FriendDatabaseHelper(this);

        switch(intent.getIntExtra(Constants.VIEWER.MODE_KEY, Constants.VIEWER.MODE_VIEWALL))
        {
            case Constants.VIEWER.MODE_VIEWALL:
                mFriends = mDb.getAllFriends();
                mDb.close();
                break;

            case Constants.VIEWER.MODE_SEARCHRESULTS:
                mFriends = mDb.getFriendsByNameFilter(intent.getStringExtra(Constants.VIEWER.FILTERNAME_KEY));
                mDb.close();
                break;
        }

        mFriendListView = (ListView)findViewById(R.id.friendDisplayList);
        mFriendAdapter = new FriendAdapter(this, mFriends);
        mFriendListView.setAdapter(mFriendAdapter);
    }

    @Override
    public void onClick(View v) {
        int position = mFriendListView.getPositionForView((View) v.getParent());
        int rowID = mFriends.get(position).getID();
        Log.d(Constants.MAD, "Delete button pressed for row: " + rowID);

        if(canDeleteFriends())
        {
            mDb.deleteFriend(rowID);
            mFriends.remove(position);
            mFriendAdapter.notifyDataSetChanged();
        }
        else {
            CharSequence text = getString(R.string.cannotDeleteFriendsErrorMsg);
            Toast toast = Toast.makeText(this, text, Toast.LENGTH_SHORT);
            toast.show();
        }
    }

    private boolean canDeleteFriends()
    {
        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(this);
        return prefs.getBoolean(Constants.PREFERENCES.KEY_FRIENDCANDELETE, false);
    }
}
