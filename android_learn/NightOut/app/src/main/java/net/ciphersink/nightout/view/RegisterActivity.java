package net.ciphersink.nightout.view;

import android.app.Activity;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Bitmap;
import android.os.AsyncTask;
import android.provider.MediaStore;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ImageView;
import android.widget.ProgressBar;
import android.widget.RelativeLayout;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.RegisterInfo;
import net.ciphersink.nightout.model.Session;
import net.ciphersink.nightout.model.SessionFactory;

import java.io.File;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;


/**
 * Implements the UI and controller for users who wish to register
 * a new account on the system.
 */
public class RegisterActivity extends ActionBarActivity implements View.OnClickListener {

    // UI
    private ImageView mProfilePic;
    private ImageButton mTakeProfilePic;
    private Button mRegisterButton;
    private ProgressBar mProgress;
    private RelativeLayout mLayout;
    private EditText mName;
    private EditText mUsername;
    private EditText mPassword;
    private EditText mEmail;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register);

        loadUIMembers();
        registerListeners();
    }

    /**
     * Fetches and stores pointers to each of the UI widgets.
     */
    private void loadUIMembers() {
        mProfilePic = (ImageView) findViewById(R.id.registerActAvatarDisplay);
        mTakeProfilePic = (ImageButton) findViewById(R.id.registerActAvatarTakePicture);
        mRegisterButton = (Button) findViewById(R.id.registerActRegister);
        mProgress = (ProgressBar) findViewById(R.id.registerActProgressBar);
        mLayout = (RelativeLayout) findViewById(R.id.registerActControlsContainer);
        mName = (EditText) findViewById(R.id.registerActNameField);
        mUsername = (EditText) findViewById(R.id.registerActUsernameField);
        mPassword = (EditText) findViewById(R.id.registerActPasswordField);
        mEmail = (EditText) findViewById(R.id.registerActEmailField);

    }

    /**
     * Initialises buttons to fire onClick() when pressed.
     */
    private void registerListeners() {
        mTakeProfilePic.setOnClickListener(this);
        mRegisterButton.setOnClickListener(this);
    }

    /**
     * Finishes the activity and returns to the login activity with information needed to complete a login.
     * Called by the registration task if registration is successful,
     * @param sessionkey
     * @param username
     */
    public void finishActivityRegistrationSuccessful(String sessionkey, String username) {
        Intent intent = new Intent();
        intent.putExtra(Constants.KEYS.SESSIONKEY, sessionkey);
        intent.putExtra(Constants.KEYS.USERNAME, username);
        setResult(Constants.KEYS.REGISTER_ACTIVITY.RESPONSECODE_REGISTRATION_SUCCESS, intent);
        finish();
    }

    /**
     * Called when the camera activity finishes with a photograph.
     * @param requestCode
     * @param resultCode
     * @param data
     */
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {

        switch (requestCode) {
            case Constants.KEYS.REMOTE_ACTIVITY_REQUEST.REQUEST_IMAGE_CAPTURE:
                if (resultCode == RESULT_OK) {
                    Bundle bundle = data.getExtras();
                    Bitmap bitmap = (Bitmap) bundle.get("data");
                    mProfilePic.setImageBitmap(bitmap);
                    //saveProfilePic(bundle);
                }
                break;
            default:
                super.onActivityResult(requestCode, resultCode, data);
        }
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // No action bar icons
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        return super.onOptionsItemSelected(item);
    }

    /**
     * Initialises a camera activity to take a profile picture.
     * @return success == true
     */
    private boolean takePicture() {
        Intent intent = new Intent(MediaStore.ACTION_IMAGE_CAPTURE);
        if (intent.resolveActivity(getPackageManager()) != null) {
            startActivityForResult(intent, Constants.KEYS.REMOTE_ACTIVITY_REQUEST.REQUEST_IMAGE_CAPTURE);
            return true;
        }
        return false;
    }


    /**
     * Implements View.OnClickListener - routes clicks to buttons in the UI.
     * @param v
     */
    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.registerActAvatarTakePicture:
                takePicture();
                break;
            case R.id.registerActRegister:
                new RegisterTask(this).execute();
        }
    }


    /**
     * Encapsulates the data model initialisation and the network call used in registering a new user.
     */
    private class RegisterTask extends AsyncTask<Void, Void, Void> {
        private Context mContext;
        private RegisterInfo data;

        private Session mSession;

        RegisterTask(Context context) {
            mContext = context;
            data = new RegisterInfo();
        }

        @Override
        protected Void doInBackground(Void... dummy) {
            if (data.validateInformation()) { //make sure information is complete.
                Log.d(Constants.MAD, "Validation successful - now registering");

                boolean success = data.register();

                if (success) { //if the network call reported success
                    Log.d(Constants.MAD, "Registration was successful!");
                    mSession = SessionFactory.login(data.getUsername(), data.getPassword()); //auto login with this information
                    Log.d(Constants.MAD, mSession != null ? "Session not null" : "Session is null");

                    if (mSession != null) {
                        Log.d(Constants.MAD, "Name: " + mSession.getName());
                        Log.d(Constants.MAD, "Username: " + mSession.getUsername());
                        Log.d(Constants.MAD, "Email: " + mSession.getEmail());
                        Log.d(Constants.MAD, "UserID: " + mSession.getUserId());
                        Log.d(Constants.MAD, "Key: " + mSession.getKey());
                    }
                }
            }
            return null;
        }

        @Override
        protected void onPreExecute() {
            Log.d(Constants.MAD, "RegisterTask - onPreExecute()");

            // hide UI elements and show a progress bar
            mLayout.setVisibility(View.GONE);
            mProgress.setVisibility(View.VISIBLE);

            // setup our data model with what the user entered
            data.setName(mName.getText().toString());
            data.setEmail(mEmail.getText().toString());
            data.setUsername(mUsername.getText().toString());
            data.setPassword(mPassword.getText().toString());
        }


        @Override
        protected void onPostExecute(Void dummy) {
            Log.d(Constants.MAD, "RegisterTask - onPostExecute()");

            if (!data.validationSuccessful()) { // must have not put in valid input - show an error

                if (!data.namePopulated())
                    mName.setError(mContext.getString(R.string.error_missing_name));
                if (!data.emailPopulated())
                    mEmail.setError(mContext.getString(R.string.error_missing_email));
                if (!data.passwordPopulated())
                    mPassword.setError(mContext.getString(R.string.error_missing_password));

                if (!data.usernamePopulated()) {
                    mUsername.setError(mContext.getString(R.string.error_missing_username));
                } else {
                    if (!data.usernameUnique())
                        mUsername.setError(mContext.getString(R.string.error_existing_username));
                }
            } else { //registration was done because validation was successful
                if (mSession != null) {
                    //return to login activity with information necessary to auto-login
                    finishActivityRegistrationSuccessful(mSession.getKey(), mSession.getUsername());
                } else {
                    Log.e(Constants.MAD, "Session does not exist!");
                }
            }

            mProgress.setVisibility(View.GONE);
            mLayout.setVisibility(View.VISIBLE);
        }
    }
}
