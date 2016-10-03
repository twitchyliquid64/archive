package net.ciphersink.remotecontrol;

import android.app.Activity;
import android.app.ActivityManager;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.Switch;

import net.ciphersink.remotecontrol.util.StatusUpdater;


public class ControlPanel extends Activity implements View.OnClickListener, CompoundButton.OnCheckedChangeListener{

    private Button mSettingsButton;
    private Button mHeartbeatButton;
    private Switch mStatusServiceSwitch;
    private Switch mLocationServiceSwitch;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_control_panel);
        populateUIVariables();
        setupListeners();

        if (isMyServiceRunning(StatusService.class)){
            mStatusServiceSwitch.setChecked(true);
        }

        if (isMyServiceRunning(LocationService.class)){
            mLocationServiceSwitch.setChecked(true);
        }

    }

    private void populateUIVariables()
    {
        mSettingsButton = (Button)findViewById(R.id.controlPanel_SettingsButton);
        mHeartbeatButton = (Button)findViewById(R.id.controlPanel_sendHeartbeatButton);
        mStatusServiceSwitch = (Switch)findViewById(R.id.controlPanel_passiveTrackerSwitch);
        mLocationServiceSwitch = (Switch)findViewById(R.id.controlPanel_locTrackerSwitch);
    }

    private void setupListeners()
    {
        mSettingsButton.setOnClickListener(this);
        mHeartbeatButton.setOnClickListener(this);
        mStatusServiceSwitch.setOnCheckedChangeListener(this);
        mLocationServiceSwitch.setOnCheckedChangeListener(this);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId())
        {
            case R.id.controlPanel_SettingsButton:
                startActivity(new Intent(this, SettingsActivity.class));
                break;

            case R.id.controlPanel_sendHeartbeatButton:
                StatusUpdater.sendStatus(this);
                break;
        }
    }

    @Override
    public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
        Log.d("CNC", "Changed: " + isChecked);
        if(buttonView.getId() == R.id.controlPanel_passiveTrackerSwitch) {
            if(isChecked){
                Intent serviceIntent = new Intent(this, StatusService.class);
                startService(serviceIntent);
            } else {
                Intent serviceIntent = new Intent(this, StatusService.class);
                stopService(serviceIntent);
            }
        } else if (buttonView.getId() == R.id.controlPanel_locTrackerSwitch) {
            if(isChecked){
                Intent serviceIntent = new Intent(this, LocationService.class);
                startService(serviceIntent);
            } else {
                Intent serviceIntent = new Intent(this, LocationService.class);
                stopService(serviceIntent);
            }
        }
    }



    private boolean isMyServiceRunning(Class<?> serviceClass) {
        ActivityManager manager = (ActivityManager) this.getSystemService(Context.ACTIVITY_SERVICE);
        for (ActivityManager.RunningServiceInfo service : manager.getRunningServices(Integer.MAX_VALUE)) {
            if (serviceClass.getName().equals(service.service.getClassName())) {
                return true;
            }
        }
        return false;
    }
}
