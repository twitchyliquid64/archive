package net.ciphersink.exercise6;

import android.content.Context;
import android.content.SharedPreferences;
import android.graphics.drawable.Drawable;
import android.preference.PreferenceManager;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;

import org.w3c.dom.Text;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.Locale;

/**
 * Created by xxx on 1/09/15.
 */
public class FriendAdapter extends BaseAdapter{

    private Context mContext;
    private ArrayList<FriendData> mFriends;

    FriendAdapter(Context context, ArrayList<FriendData> friends)
    {
        mContext = context;
        mFriends = friends;
    }

    @Override
    public int getCount() {
        return mFriends.size();
    }

    @Override
    public FriendData getItem(int position) {
        return mFriends.get(position);
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View view;
        if (convertView != null)
        {
            view = convertView;
        }
        else {
            view = LayoutInflater.from(parent.getContext()).inflate(R.layout.friend_item, null);
        }

        FriendData friend = getItem(position);

        TextView name = (TextView)view.findViewById(R.id.friendName);
        name.setText(friend.getName());

        TextView occupation = (TextView)view.findViewById(R.id.friendOccupation);
        occupation.setText(mContext.getString(R.string.prefix_occupation) + " " + friend.getOccupation());

        TextView city = (TextView)view.findViewById(R.id.friendCity);
        city.setText(mContext.getString(R.string.prefix_city) + " " + friend.getCity());

        Date friendsSince = new Date(friend.getFriendSince() * 1000);
        SimpleDateFormat format = new SimpleDateFormat(getDateFormat(), Locale.US);

        TextView friendSince = (TextView)view.findViewById(R.id.friendSince);
        friendSince.setText(mContext.getString(R.string.prefix_since) + " " + format.format(friendsSince));

        ImageButton delButton = (ImageButton)view.findViewById(R.id.deleteFriendButton);
        delButton.setOnClickListener((View.OnClickListener)mContext);

        return view;
    }

    private String getDateFormat()
    {
        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(mContext);
        return prefs.getString(Constants.PREFERENCES.KEY_FRIENDDATEFORMAT, "dd/MM/yyyy");
    }
}
