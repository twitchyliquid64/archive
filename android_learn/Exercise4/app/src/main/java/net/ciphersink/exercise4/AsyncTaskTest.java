package net.ciphersink.exercise4;

import android.content.Context;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.Spinner;
import android.widget.Toast;


public class AsyncTaskTest extends ActionBarActivity implements AdapterView.OnItemSelectedListener, View.OnClickListener {

    private Spinner mFileSelectorSpinner;
    private Button mDownloadFileButton;
    private Button mDownloadAllButton;
    private Button mCalculateButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_async_task);


        initialiseUIMethodVariables();
        initialiseListeners();
        populateFileSpinner();
    }

    private void initialiseUIMethodVariables()
    {
        mFileSelectorSpinner = (Spinner)findViewById(R.id.file_selector_spinner);
        mDownloadFileButton = (Button)findViewById(R.id.download_button);
        mDownloadAllButton = (Button)findViewById(R.id.download_all_button);
        mCalculateButton = (Button)findViewById(R.id.calculate_millionth_button);
    }

    private void initialiseListeners()
    {
        mFileSelectorSpinner.setOnItemSelectedListener(this);
        mDownloadFileButton.setOnClickListener(this);
        mDownloadAllButton.setOnClickListener(this);
        mCalculateButton.setOnClickListener(this);
    }

    private void populateFileSpinner()
    {
        //setup spinner with populated values
        Spinner spinner = (Spinner) findViewById(R.id.file_selector_spinner);
        ArrayAdapter<CharSequence> adapter = ArrayAdapter.createFromResource(this,
                R.array.downloadable_files, android.R.layout.simple_spinner_item); // associate the array of values, and setup how the spinner items will look
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item); // how the dropdown will look
        spinner.setAdapter(adapter);
    }

    private void onDownloadFilePressed(String filename)
    {
        new DownloadSelectedAsyncTask(this, filename).execute();
    }

    private void onDownloadAllPressed(String[] files)
    {
        new DownloadAllAsyncTask(this, files.length).execute(files);
    }

    private void calculateMillionthTriNumber()
    {
        new CalculateAsyncTask(this).execute();
    }

    @Override
    public void onClick(View v) {
        switch (v.getId())
        {
            case R.id.download_button:
                String fileName = ((Spinner) findViewById(R.id.file_selector_spinner)).getSelectedItem().toString();
                Log.d(Constants.MAD, "Download button pressed: " + fileName);
                onDownloadFilePressed(fileName);
                break;
            case R.id.download_all_button:
                Log.d(Constants.MAD, "Download all button pressed");
                onDownloadAllPressed(getResources().getStringArray(R.array.downloadable_files));
                break;
            case R.id.calculate_millionth_button:
                Log.d(Constants.MAD, "Calculate button pressed");
                calculateMillionthTriNumber();
                break;
        }
    }

    @Override
    public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
        Log.d(Constants.MAD, "onItemSelected()");
        Context context = getApplicationContext();
        CharSequence text = getString(R.string.selected_prefix) + ((Spinner) findViewById(R.id.file_selector_spinner)).getSelectedItem().toString();
        int duration = Toast.LENGTH_SHORT;

        Toast toast = Toast.makeText(context, text, duration);
        toast.show();
    }

    @Override
    public void onNothingSelected(AdapterView<?> parent) {
        // do nothing - nothing is selected
    }



    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_async_task_test, menu);
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
