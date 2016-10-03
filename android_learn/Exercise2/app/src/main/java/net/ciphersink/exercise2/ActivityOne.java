/*
 * Copyright (C) 2015 Tom D'Netto
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package net.ciphersink.exercise2;

import android.content.Context;
import android.content.Intent;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;
import android.widget.Toast;


/*
 * Builds a layout that allows the user to insert their contact details.
 */
public class ActivityOne extends ActionBarActivity {// implements View.OnClickListener

    /*
     * Called when the app is launched (entrypoint).
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_one);

        //setup spinner with populated values
        Spinner spinner = (Spinner) findViewById(R.id.spinner);
        ArrayAdapter<CharSequence> adapter = ArrayAdapter.createFromResource(this,
                R.array.phone_type_array, android.R.layout.simple_spinner_item);//associate the array of values, and setup how the spinner items will look
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);//how the dropdown will look
        spinner.setAdapter(adapter);

        Button submitButton = (Button)findViewById(R.id.subButton);
        submitButton.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                Context context = getApplicationContext();
                CharSequence text = getString(R.string.submit_pressed_toast);
                int duration = Toast.LENGTH_SHORT;

                Toast toast = Toast.makeText(context, text, duration);
                toast.show();

                Intent intent = new Intent(getBaseContext(), ActivityTwo.class);
                intent.putExtra(Constants.Keys.name, ((EditText)findViewById(R.id.nameInput)).getText().toString());
                intent.putExtra(Constants.Keys.email, ((EditText)findViewById(R.id.emailInput)).getText().toString());
                intent.putExtra(Constants.Keys.phoneType, ((Spinner) findViewById(R.id.spinner)).getSelectedItem().toString());
                intent.putExtra(Constants.Keys.phone, ((EditText) findViewById(R.id.phoneInput)).getText().toString());
                startActivityForResult(intent, 1);
            }
        });

        Button clearButton = (Button)findViewById(R.id.clearButton);
        clearButton.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                ((EditText)findViewById(R.id.emailInput)).setText("");
                ((EditText)findViewById(R.id.phoneInput)).setText("");
                ((EditText)findViewById(R.id.nameInput)).setText("");
            }
        });


        Button exitButton = (Button)findViewById(R.id.exitButton);
        exitButton.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                finish();
            }
        });
    }

    /*
     * Called when returning from the confirmation activity. Checks
     * whether the user agrees and outputs the value via a toast message.
     */
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        Context context = getApplicationContext();
        CharSequence text = data.getStringExtra(Constants.Keys.confirmation_status);
        int duration = Toast.LENGTH_SHORT;

        Toast toast = Toast.makeText(context, text, duration);
        toast.show();
    }


    /*
     * Called when the activity is partially occluded.
     */
    @Override
    protected void onPause()
    {
        Log.d("LIFECYCLE DEBUG", "onPause()");
        super.onPause();
    }

    /*
     * Called when the activity is foreground.
     */
    @Override
    protected void onResume()
    {
        Log.d("LIFECYCLE DEBUG", "onResume()");
        super.onResume();
    }

    /*
     * Called when the user switches to another app
     */
    @Override
    protected void onStop()
    {
        Log.d("LIFECYCLE DEBUG", "onStop()");
        super.onStop();
    }

    /*
     * Called when the user returns to the activity
     */
    @Override
    protected void onRestart()
    {
        Log.d("LIFECYCLE DEBUG", "onRestart()");
        super.onRestart();
    }

    /*
     * Called when the resources of the acitivty are yield
     */
    @Override
    protected void onDestroy()
    {
        Log.d("LIFECYCLE DEBUG", "onDestroy()");
        super.onDestroy();
    }

    @Override
    protected void onStart()
    {
        Log.d("LIFECYCLE DEBUG", "onStart()");
        super.onStart();
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
