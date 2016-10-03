package net.ciphersink.exercise6;

import android.content.Intent;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;


public class AddFriendActivity extends ActionBarActivity implements View.OnClickListener{

    private Button mAddFriendButton;
    private Button mCancelButton;
    private EditText mName;
    private EditText mCity;
    private EditText mOccupation;

    private FriendDatabaseHelper mFriendDatabase;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_add_friend);

        mAddFriendButton = (Button)findViewById(R.id.addFriendFinishButton);
        mCancelButton = (Button)findViewById(R.id.addFriendCancelButton);
        mName = (EditText)findViewById(R.id.addFriendNameField);
        mCity = (EditText)findViewById(R.id.addFriendCityField);
        mOccupation = (EditText)findViewById(R.id.addFriendOccupationField);

        mAddFriendButton.setOnClickListener(this);
        mCancelButton.setOnClickListener(this);

        mFriendDatabase = new FriendDatabaseHelper(this);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_add_friend, menu);
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
        switch(v.getId())
        {
            case R.id.addFriendFinishButton:
                String name = mName.getText().toString();
                String city = mCity.getText().toString();
                String occupation = mOccupation.getText().toString();

                FriendData friend = new FriendData(name, occupation, city, System.currentTimeMillis() / 1000);

                Log.d(Constants.MAD, "Adding friend:");
                Log.d(Constants.MAD, "\tName: " + name);
                Log.d(Constants.MAD, "\tOccupation: " + occupation);
                Log.d(Constants.MAD, "\tCity: " + city);

                mFriendDatabase.addFriend(friend);
                Log.d(Constants.MAD, "new friend committed.");

                CharSequence text = getString(R.string.friendAddedToast);
                Toast toast = Toast.makeText(this, text, Toast.LENGTH_SHORT);
                toast.show();

                finish();
                break;

            case R.id.addFriendCancelButton:
                finish();
                break;
        }
    }
}
