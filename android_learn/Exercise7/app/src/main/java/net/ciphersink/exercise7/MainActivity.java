package net.ciphersink.exercise7;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;


public class MainActivity extends ActionBarActivity implements View.OnClickListener{

    private Button mStartButton;
    private TextView mDisplay;
    private BroadcastReceiver mMyIntentServiceReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {
            Log.d(Constants.MAD, intent.toString());
            onIntentServiceMsgRecieved(intent);
        }
    };

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        mStartButton = (Button)findViewById(R.id.startButton);
        mDisplay = (TextView)findViewById(R.id.titleTextView);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }

    @Override
    public void onClick(View v) {
        //only called by start button.
        Log.d(Constants.MAD, "Start button clicked");
        mStartButton.setEnabled(false);

        Intent msgIntent = new Intent(this, MyIntentService.class);
        msgIntent.putExtra(Constants.COUNT_REQUEST.KEY_MAXVAL, 10);
        startService(msgIntent);
    }

    private void onIntentServiceMsgRecieved(Intent intent)
    {
        String msg = intent.getStringExtra(Constants.COUNT_REQUEST.KEY_RETMSG);

        if (msg.equals(getString(R.string.taskCompleted)))
        {
            mStartButton.setEnabled(true);
        }

        mDisplay.setText(msg);
    }

    @Override
    public void onStart()
    {
        super.onStart();
        registerReceiver(mMyIntentServiceReceiver, new IntentFilter(Constants.COUNT_REQUEST.INTENT_RESULT));
    }

    @Override
    public void onStop()
    {
        super.onStop();
        unregisterReceiver(mMyIntentServiceReceiver);
    }
}
