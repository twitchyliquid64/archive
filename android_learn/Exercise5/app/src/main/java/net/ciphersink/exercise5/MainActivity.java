package net.ciphersink.exercise5;

import android.content.Context;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ListView;
import android.widget.ProgressBar;
import android.widget.TextView;

import org.w3c.dom.Text;

import java.util.ArrayList;
import java.util.List;


public class MainActivity extends ActionBarActivity implements View.OnClickListener{

    private ArrayList mTrains;
    private TrainAdapter mTrainAdaptor;
    private ListView mTrainListView;
    private ProgressBar mProgressBar;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        mTrains = new ArrayList();
        for (int i = 0; i < 3; i++) {
            mTrains.add(new TrainData("Platform "+(2*i), 12+i, "On Time", "Wollstonecraft", "15:01"));
            mTrains.add(new TrainData("Platform "+(1+i), 33+i, "Fuckin late", "Kekland", "12:44"));
        }

        mProgressBar = (ProgressBar)findViewById(R.id.activity_main_progressbar);
        mTrainAdaptor = new TrainAdapter(this, mTrains);
        mTrainListView = (ListView)findViewById(R.id.trainListView);
        mTrainListView.setAdapter(mTrainAdaptor);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId())
        {
            case R.id.action_menu_refresh:
                new UpdateTimesTask(this, mTrainListView, mProgressBar, mTrains, mTrainAdaptor).execute();
                return true;
            case R.id.action_menu_add:
                mTrains.add(new TrainData("Platform 24", 1, "Scheduled - unknown", "Nimben", "4:20"));
                mTrainAdaptor.notifyDataSetChanged();
                return true;
            case R.id.action_menu_delete_all:
                mTrains.clear();
                mTrainAdaptor.notifyDataSetChanged();
                return true;
        }

        return super.onOptionsItemSelected(item);
    }

    @Override
    public void onClick(View v) {
        int position = mTrainListView.getPositionForView((View) v.getParent());

        ProgressBar progressBar = (ProgressBar)v.findViewById(R.id.arriveTimeProgress);
        TextView textView = (TextView)v.findViewById(R.id.arrivalTime);

        new UpdateTimeTask(this, textView, progressBar, mTrains, mTrainAdaptor, position).execute();
    }
}
