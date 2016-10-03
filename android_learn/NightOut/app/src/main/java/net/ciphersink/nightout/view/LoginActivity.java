package net.ciphersink.nightout.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ProgressBar;
import android.widget.RelativeLayout;
import android.widget.TextView;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Session;
import net.ciphersink.nightout.model.SessionFactory;


/**
 * Implements the UI which users see when they first start the app. It has the login
 * widgets as the name suggests, a register button, and the logo / version number.
 */
public class LoginActivity extends Activity implements View.OnClickListener{

    //ui elements
    private TextView mTitleAndVersion;
    private Button mLoginButton;
    private Button mRegisterButton;
    private EditText mUsername;
    private EditText mPassword;
    private RelativeLayout mUserControlsContainer;
    private ProgressBar mProgress;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_login);
        loadUIMembers();
        setVersionText();

        checkIfCanReuseSession();
    }

    /**
     * Initalises UI member variables with the correct instances.
     */
    private void loadUIMembers()
    {
        mTitleAndVersion = (TextView)findViewById(R.id.loginActTitleAndVersionView);
        mLoginButton = (Button)findViewById(R.id.logincActLoginButton);
        mLoginButton.setOnClickListener(this);
        mRegisterButton = (Button)findViewById(R.id.loginActRegisterButton);
        mRegisterButton.setOnClickListener(this);
        mUsername = (EditText)findViewById(R.id.loginActUsrField);
        mPassword = (EditText)findViewById(R.id.loginActPassField);
        mUserControlsContainer = (RelativeLayout)findViewById(R.id.loginActControlsContainer);
        mProgress = (ProgressBar)findViewById(R.id.loginActProgress);
    }

    /**
     * Sets the TextView with a string representing the application version.
     */
    private void setVersionText()
    {
        String versionString = getString(R.string.versionText);
        versionString += " " + Constants.RELEASE_VERSION;
        versionString += "." + Constants.MAJOR_VERSION;
        versionString += "." + Constants.MINOR_VERSION;
        if (Constants.RELEASE_VERSION < 1){
            versionString += " " + getString(R.string.alphaText);
        }
        else {
            versionString += " " + getString(R.string.betaText);
        }
        mTitleAndVersion.setText(versionString);
    }

    /**
     * Called when the register button is clicked.
     */
    private void onRegisterClicked()
    {
        Intent intent = new Intent(getBaseContext(), RegisterActivity.class);
        startActivityForResult(intent, Constants.KEYS.REGISTER_ACTIVITY.REQUESTCODE_REGISTER);
    }

    /**
     * Called when the login button is clicked.
     */
    private void onLoginClicked()
    {
        hideControlsAndShowProgress();
        mTitleAndVersion.setText(getString(R.string.login_progress_message));
        new LoginManualTask().execute(mUsername.getText().toString(), mPassword.getText().toString());
    }

    /**
     * Implements View.OnClickListener, is attached to all buttons in the activity and routes them
     * appropriately.
     * @param v View
     */
    public void onClick(View v)
    {
        switch (v.getId())
        {
            case R.id.loginActRegisterButton:
                onRegisterClicked();
                break;
            case R.id.logincActLoginButton:
                onLoginClicked();
                break;
        }
    }

    /**
     * Attemps to login using a known sessionkey and username. Failure / success UI actions are defined in
     * LoginSessionkeyTask.
     * @param sessionkey
     * @param username
     */
    private void doAutoLogin(String sessionkey, String username) {
        mUsername.setText(username);
        mTitleAndVersion.setText(getString(R.string.login_progress_message));
        hideControlsAndShowProgress();

        new LoginSessionkeyTask().execute(sessionkey);
    }

    /**
     * Hides user input controls and shows a loading progress.
     */
    private void hideControlsAndShowProgress() {
        mUserControlsContainer.setVisibility(View.GONE);
        mProgress.setVisibility(View.VISIBLE);
    }

    /**
     * Shoes user input controls and shows a loading progress.
     */
    private void showControlsAndHideProgress() {
        mProgress.setVisibility(View.GONE);
        mUserControlsContainer.setVisibility(View.VISIBLE);
    }

    /**
     * Called on initialisation - checks if a valid session key exists that can be used
     * to auto-login.
     */
    private void checkIfCanReuseSession() {
        SharedPreferences sharedPref = getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        String sessionkey = sharedPref.getString(Constants.KEYS.SESSIONKEY, null);
        String username = sharedPref.getString(Constants.KEYS.USERNAME, "");

        if (sessionkey != null) {
            doAutoLogin(sessionkey, username);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {

        //if registration was successful, the necessary details are returned to complete
        //a login.
        if (requestCode == Constants.KEYS.REGISTER_ACTIVITY.REQUESTCODE_REGISTER &&
            resultCode == Constants.KEYS.REGISTER_ACTIVITY.RESPONSECODE_REGISTRATION_SUCCESS) {
            String sessionkey = data.getStringExtra(Constants.KEYS.SESSIONKEY);
            String username = data.getStringExtra(Constants.KEYS.USERNAME);
            Log.d(Constants.MAD, "Login got details from register: " + sessionkey + " " + username);

            doAutoLogin(sessionkey, username);
        }
    }

    /**
     * If a login is successful, this method is called to save the parameters and start mainActivity.
     * @param key
     * @param username
     */
    private void finishLogin(String key, String username) {
        Log.d(Constants.MAD, "Login successful!");

        SharedPreferences sharedPref = getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        SharedPreferences.Editor editor = sharedPref.edit();
        editor.putString(Constants.KEYS.SESSIONKEY, key);
        editor.putString(Constants.KEYS.USERNAME, username);
        editor.commit();

        Intent intent = new Intent(getBaseContext(), MainActivity.class);
        intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
        startActivity(intent);
        showControlsAndHideProgress();
        setVersionText();
        finish();
    }

    /**
     * Encapsulates the network call to validate a session key.
     */
    private class LoginSessionkeyTask extends AsyncTask<String, Void, Void> {
        private Session mSession = null;

        @Override
        protected Void doInBackground(String... key) {
            //check if session is valid
            mSession = SessionFactory.makeFromKey(key[0]);
            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            if (mSession == null) {
                //session invalid - they need to login again
                showControlsAndHideProgress();
                setVersionText();
            } else {
                //session valid - can proceed using sessionkey
                finishLogin(mSession.getKey(), mSession.getUsername());
            }
        }
    }

    /**
     * Encapsulates the network call to login with a username and password.
     */
    private class LoginManualTask extends AsyncTask<String, Void, Void> {
        private Session mSession = null;

        @Override
        protected Void doInBackground(String... authdata) {
            //check if credentials are valid
            mSession = SessionFactory.login(authdata[0], authdata[1]);
            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            if (mSession == null) {
                //credentials incorrect - they need to login again
                showControlsAndHideProgress();
                mTitleAndVersion.setText(getString(R.string.error_authentication_details));
            } else {
                //session valid - can proceed using sessionkey
                finishLogin(mSession.getKey(), mSession.getUsername());
            }
        }
    }
}
