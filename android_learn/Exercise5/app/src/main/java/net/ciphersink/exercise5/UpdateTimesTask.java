package net.ciphersink.exercise5;

import android.content.Context;
import android.os.AsyncTask;
import android.view.View;
import android.widget.ListView;
import android.widget.ProgressBar;

import java.util.ArrayList;
import java.util.Random;

/**
 * Created by xxx on 25/08/15.
 */
public class UpdateTimesTask extends AsyncTask<Void, Void, Void> {

    private Context mContext;
    private ListView mListView;
    private ProgressBar mProgressBar;
    private ArrayList<TrainData> mTrainList;
    private TrainAdapter mTrainAdaptor;
    private Random mRand;

    UpdateTimesTask(Context context, ListView listView, ProgressBar progressBar, ArrayList trainList, TrainAdapter trainAdapter)
    {
        mContext = context;
        mListView = listView;
        mProgressBar = progressBar;
        mTrainList = trainList;
        mTrainAdaptor = trainAdapter;
        mRand = new Random();
    }

    @Override
    protected Void doInBackground(Void... params) {

        for(int i = 0; i < mTrainList.size(); i++)
        {
            TrainData train = mTrainList.get(i);
            train.setArrivalTime(mRand.nextInt(45));
        }

        try {
            Thread.sleep(3000);
        }
        catch (InterruptedException e)
        {
            //do nothing
        }
        return null;
    }

    protected  void onPostExecute(Void dummy)
    {
        mListView.setVisibility(View.VISIBLE);
        mProgressBar.setVisibility(View.INVISIBLE);
        mTrainAdaptor.notifyDataSetChanged();
    }

    @Override
    protected  void onPreExecute() {
        mListView.setVisibility(View.INVISIBLE);
        mProgressBar.setVisibility(View.VISIBLE);
    }
}
