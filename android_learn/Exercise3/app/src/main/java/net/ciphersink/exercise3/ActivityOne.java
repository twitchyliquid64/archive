package net.ciphersink.exercise3;

import android.content.pm.ActivityInfo;
import android.content.res.Configuration;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;


public class ActivityOne extends ActionBarActivity implements View.OnFocusChangeListener, View.OnClickListener {

    private EditText mBinNr;
    private EditText mBinType;
    private EditText mBinSize;
    private Button mClearButton;
    private Button mResetQtyButton;
    private Button mRotateButton;
    private boolean mScreenIsPortrait = true;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Log.d(Constants.MAD, "onCreate()");
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_one);
        initialiseControls();
    }

    /**
     * fdgdfgdfg
     */
    private void initialiseControls()
    {
        mBinType = (EditText)findViewById(R.id.activity_one_type_edittext);
        mBinSize = (EditText)findViewById(R.id.binSize);
        mBinNr = (EditText)findViewById(R.id.nr);
        mClearButton = (Button)findViewById(R.id.clearButton);
        mResetQtyButton = (Button)findViewById(R.id.resetButton);
        mRotateButton = (Button)findViewById(R.id.rotateButton);

        mBinNr.setOnFocusChangeListener(this);
        mRotateButton.setOnClickListener(this);

        //these buttons are not present in a landscape view, so we get a nullref
        try {
            mClearButton.setOnClickListener(this);
            mResetQtyButton.setOnClickListener(this);
        }
        catch (NullPointerException e)
        {
            // these buttons do not exist in landscape mode - ignore nullref
        }

    }

    @Override
    public void onClick(View v)
    {
        Log.d(Constants.MAD, "onClick(): " + v.getId());
        if (v.getId() == R.id.clearButton)
        {
            onClearAllClicked();
        }
        if (v.getId() == R.id.resetButton)
        {
            mBinNr.setText("");
        }
        if (v.getId() == R.id.rotateButton)
        {
            rotateScreen();
        }
    }

    private void onClearAllClicked()
    {
        //mBinType = null;
        try {
            mBinSize.setText("");
            mBinNr.setText("");
            mBinType.setText("");
        }
        catch (Exception e)
        {
            Log.e(Constants.MAD, e.toString());
        }
    }

    private void rotateScreen()
    {
        Log.d(Constants.MAD, "rotateScreen() currently portrait: " + mScreenIsPortrait);
        mScreenIsPortrait = !mScreenIsPortrait;
        if (mScreenIsPortrait)
        {
            Log.d(Constants.MAD, "locking to portrait");
            this.setRequestedOrientation(ActivityInfo.SCREEN_ORIENTATION_PORTRAIT);
        }else
        {
            Log.d(Constants.MAD, "locking to landscape");
            this.setRequestedOrientation(ActivityInfo.SCREEN_ORIENTATION_REVERSE_LANDSCAPE);
        }
    }

    @Override
    public void onConfigurationChanged(Configuration newConfig)
    {
        int ot = getResources().getConfiguration().orientation;
        switch (ot) {
            case Configuration.ORIENTATION_LANDSCAPE:
                Log.d(Constants.MAD, "Reporting landscape configuration");
                break;
            case Configuration.ORIENTATION_PORTRAIT:
                Log.d(Constants.MAD, "Reporting portrait configuration");
                break;
        }
        setContentView(R.layout.activity_one);
        initialiseControls();
        super.onConfigurationChanged(newConfig);
    }

    @Override
    public void onFocusChange(View v, boolean getsFocus)
    {
        Log.d(Constants.MAD, "onFocusChange()");
        if (v.getId() == R.id.nr && getsFocus)
        {
            CharSequence text = getString(R.string.binQtyHasFocus);
            int duration = Toast.LENGTH_SHORT;

            Toast toast = Toast.makeText(getApplicationContext(), text, duration);
            toast.show();
        }
    }

    @Override
    public void onSaveInstanceState(Bundle savedInstanceState) {
        Log.d(Constants.MAD, "onSaveInstanceState()");
        savedInstanceState.putString(Constants.KEYS.BINTYPE, mBinType.getText().toString());
        savedInstanceState.putBoolean(Constants.KEYS.PORTRAITSTATE, mScreenIsPortrait);
        super.onSaveInstanceState(savedInstanceState);
    }

    @Override
    public void onRestoreInstanceState(Bundle savedInstanceState) {
        Log.d(Constants.MAD, "onRestoreInstanceState()");
        mBinType.setText(savedInstanceState.getString(Constants.KEYS.BINTYPE));
        mScreenIsPortrait = savedInstanceState.getBoolean(Constants.KEYS.PORTRAITSTATE);
        super.onRestoreInstanceState(savedInstanceState);
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_activity_one, menu);
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
}
