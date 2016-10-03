package net.ciphersink.nightout.view;

import android.os.AsyncTask;
import android.os.Bundle;
import android.support.design.widget.FloatingActionButton;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Notification;

/**
 * Implements the UI and network within which users send messages to other users / squads.
 */
public class NewMessageActivity extends AppCompatActivity implements View.OnClickListener {

    //UI
    private TextView mToText;
    private EditText mContentText;
    private FloatingActionButton mFloatingButton;

    //model
    private String mSessionKey;
    private boolean mSendToSquad = false;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_new_message);

        // store session key in member
        mSessionKey = getIntent().getStringExtra(Constants.KEYS.SESSIONKEY);

        // set mode - are we writing a message for a whole squad or a single user?
        if (getIntent().getIntExtra(Constants.KEYS.MODE, -1) ==
                Constants.KEYS.NOTIFICATION_REQUEST.MODE_SEND_ALL_SQUAD) {
            mSendToSquad = true;
        }

        initUI();
    }

    /**
     * Initialises member variables with UI pointers.
     */
    private void initUI() {
        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);

        mFloatingButton = (FloatingActionButton) findViewById(R.id.fab);
        mFloatingButton.setOnClickListener(this);

        mToText = (TextView)findViewById(R.id.newMessageActToText);
        mContentText = (EditText)findViewById(R.id.newMessageActContentEdit);
        mToText.setText(getIntent().getStringExtra(Constants.KEYS.NAME));
    }

    /**
     * Implements view.OnClickListener - Handles click events, in this case
     * exclusively the send button.
     * @param view
     */
    @Override
    public void onClick(View view) {
        Snackbar.make(view, getString(R.string.sending_online_please_wait), Snackbar.LENGTH_LONG)
                .setAction("Action", null).show();
        mFloatingButton.setVisibility(View.GONE);

        if (!mSendToSquad) {
            new SendMessageNotificationTask().execute();
        } else { //send to squad
            new SendMessageToSquadTask().execute();
        }
    }


    /**
     * Encapsulates the network call of transmitting a notification to a single user.
     */
    private class SendMessageNotificationTask extends AsyncTask<Void, Void, Void> {

        private String mContent;

        @Override
        protected void onPreExecute() {
            mContent = mContentText.getText().toString();
        }


        @Override
        protected Void doInBackground(Void... dummy) {

            Notification message = new Notification(mContent, Constants.NET.NOTIFICATON_TYPE.MESSAGE,
                    getString(R.string.private_message_from) + " " + getIntent().getStringExtra(Constants.KEYS.FROM_NAME));
            message.send(mSessionKey, getIntent().getIntExtra(Constants.KEYS.USER_ID, -1));

            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            finish();
        }
    }


    /**
     * Encasupsulates the network call of transmitting a message to a whole squad.
     */
    private class SendMessageToSquadTask extends AsyncTask<Void, Void, Void> {

        private String mContent;

        @Override
        protected void onPreExecute() {
            mContent = mContentText.getText().toString();
        }


        @Override
        protected Void doInBackground(Void... dummy) {

            Notification message = new Notification(mContent, Constants.NET.NOTIFICATON_TYPE.SQUAD_MESSAGE,
                            getString(R.string.squad_message_to) + " " + getIntent().getStringExtra(Constants.KEYS.NAME) +
                            " " + getString(R.string.from) + " " +
                            getIntent().getStringExtra(Constants.KEYS.FROM_NAME));
            message.sendToSquad(mSessionKey, getIntent().getIntExtra(Constants.KEYS.SQUAD_ID, -1));

            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            finish();
        }
    }
}
