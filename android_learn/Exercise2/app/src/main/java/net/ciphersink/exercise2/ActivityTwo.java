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
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.CheckBox;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

/*
 * Allows the user to confirm their contact details.
 */
public class ActivityTwo extends ActionBarActivity implements View.OnClickListener{

    /*
     * Called when the activity is loaded, expects to be called with intent
     * for a result. Asks the user to confirm their contact details, and returns
     * that to the calling activity.
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_two);

        Intent intent = getIntent();

        String name = intent.getStringExtra(Constants.Keys.name);
        ((TextView)findViewById(R.id.nameDisplay)).setText(name);

        String phnum = intent.getStringExtra(Constants.Keys.phone);
        ((TextView)findViewById(R.id.numberDisplay)).setText(phnum);

        String email = intent.getStringExtra(Constants.Keys.email);
        ((TextView)findViewById(R.id.emailDisplay)).setText(email);

        String typ = intent.getStringExtra(Constants.Keys.phoneType);
        ((TextView)findViewById(R.id.typeDisplay)).setText(typ);

        ((Button)findViewById(R.id.subButtonFinal)).setOnClickListener(this);

    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_activity_two, menu);
        return true;
    }

    /*
     * Called only when the submit button is pressed. Packages whether the user
     * agrees or not into a Intent extra, before returning that result to the
     * former activity.
     */
    @Override
    public void onClick(View v)
    {
        Intent intent = new Intent();
        intent.putExtra(Constants.Keys.confirmation_status,
                ((CheckBox)findViewById(R.id.agreeCheckbox)).isChecked() ?
                        getString(R.string.agree_text) : getString(R.string.disagree_text));

        setResult(1, intent);
        finish();
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
