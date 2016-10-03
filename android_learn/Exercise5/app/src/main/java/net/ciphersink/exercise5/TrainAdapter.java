package net.ciphersink.exercise5;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.LinearLayout;
import android.widget.TextView;

import org.w3c.dom.Text;

import java.util.ArrayList;


/**
 * Created by xxx on 25/08/15.
 * Context context, ArrayList trainList
 */
public class TrainAdapter extends BaseAdapter {

    private Context mContext;
    private ArrayList mTrainList;

    TrainAdapter(Context context, ArrayList trainList)
    {
        mContext = context;
        mTrainList = trainList;
    }

    @Override
    public int getCount() {
        return mTrainList.size();
    }

    @Override
    public TrainData getItem(int position) {
        return (TrainData)mTrainList.get(position);
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
            view = LayoutInflater.from(parent.getContext()).inflate(R.layout.list_item, null);
        }

        TrainData train = getItem(position);

        TextView arrivalTime = (TextView)view.findViewById(R.id.arrivalTime);
        TextView status = (TextView)view.findViewById(R.id.status);
        TextView destination = (TextView)view.findViewById(R.id.destination);
        TextView destinationTime = (TextView)view.findViewById(R.id.destinationTime);
        TextView platform = (TextView)view.findViewById(R.id.platform);
        LinearLayout aTimeContainer = (LinearLayout)view.findViewById(R.id.leftDisplay);

        arrivalTime.setText(train.getArrivalTime() + " " + mContext.getString(R.string.timePrefix));
        status.setText(train.getStatus());
        destination.setText(train.getDestination());
        destinationTime.setText(train.getDestinationTime());
        platform.setText(train.getPlatform());
        aTimeContainer.setOnClickListener((View.OnClickListener)mContext);


        return view;
    }
}
