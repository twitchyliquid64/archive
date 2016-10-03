package net.ciphersink.exercise6;

import android.content.Context;
import android.content.Intent;
import android.os.Parcelable;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import java.sql.BatchUpdateException;


public class MainActivity extends ActionBarActivity implements View.OnClickListener {

    private Button mAddFriendButton;
    private Button mShowFriendsButton;
    private Button mSearchFriendsButton;
    private EditText mSearchEditText;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        mSearchEditText = (EditText)findViewById(R.id.nameSearchEditText);

        mSearchFriendsButton = (Button)findViewById(R.id.searchFriendsButton);
        mSearchFriendsButton.setOnClickListener(this);

        mAddFriendButton = (Button)findViewById(R.id.addFriendButton);
        mAddFriendButton.setOnClickListener(this);

        mShowFriendsButton = (Button)findViewById(R.id.showFriendButton);
        mShowFriendsButton.setOnClickListener(this);
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
            startActivity(new Intent(getBaseContext(), SettingsActivity.class));
            return true;
        }

        return super.onOptionsItemSelected(item);
    }


    @Override
    public void onClick(View v) {
        switch (v.getId())
        {
            case R.id.addFriendButton:
                Intent intent = new Intent(getBaseContext(), AddFriendActivity.class);
                startActivity(intent);
                break;
            case R.id.showFriendButton:
                Intent intent2 = new Intent(getBaseContext(), ViewFriendsActivity.class);
                intent2.putExtra(Constants.VIEWER.MODE_KEY, Constants.VIEWER.MODE_VIEWALL);
                startActivity(intent2);
                break;
            case R.id.searchFriendsButton:
                Intent intent3 = new Intent(getBaseContext(), ViewFriendsActivity.class);
                intent3.putExtra(Constants.VIEWER.MODE_KEY, Constants.VIEWER.MODE_SEARCHRESULTS);
                intent3.putExtra(Constants.VIEWER.FILTERNAME_KEY, mSearchEditText.getText().toString());
                startActivity(intent3);
                break;
        }
    }
}
