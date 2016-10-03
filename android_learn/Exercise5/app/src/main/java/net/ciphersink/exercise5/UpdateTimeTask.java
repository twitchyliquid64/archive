package net.ciphersink.exercise5;

import android.content.Context;
import android.os.AsyncTask;
import android.view.View;
import android.widget.ListView;
import android.widget.ProgressBar;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.Random;

/**
 * Created by xxx on 25/08/15.
 */
public class UpdateTimeTask extends AsyncTask<Void, Void, Void> {

    private Context mContext;
    private TextView mTextView;
    private ProgressBar mProgressBar;
    private ArrayList<TrainData> mTrainList;
    private TrainAdapter mTrainAdaptor;
    private Random mRand;
    private int mPosition;

    UpdateTimeTask(Context context, TextView textView, ProgressBar progressBar, ArrayList<TrainData> trainList, TrainAdapter trainAdapter, int position)
    {
        mContext = context;
        mTextView = textView;
        mProgressBar = progressBar;
        mTrainList = trainList;
        mTrainAdaptor = trainAdapter;
        mPosition = position;
        mRand = new Random();
    }

    @Override
    protected Void doInBackground(Void... params) {
        TrainData train = mTrainList.get(mPosition);
        train.setArrivalTime(mRand.nextInt(30));
        train.setStatus(mRand.nextInt(2) == 1 ? mContext.getString(R.string.onTimeText): mContext.getString(R.string.lateText));

        try {
            Thread.sleep(1000);
        }
        catch (InterruptedException e)
        {
            //do nothing
        }
        return null;
    }

    protected  void onPostExecute(Void dummy)
    {
        mTrainAdaptor.notifyDataSetChanged();
        mTextView.setVisibility(View.VISIBLE);
        mProgressBar.setVisibility(View.GONE);
    }

    @Override
    protected  void onPreExecute() {
        mTextView.setVisibility(View.GONE);
        mProgressBar.setVisibility(View.VISIBLE);
    }
}
