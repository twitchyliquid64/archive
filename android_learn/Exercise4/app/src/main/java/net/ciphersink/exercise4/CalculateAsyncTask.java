package net.ciphersink.exercise4;

import android.app.ProgressDialog;
import android.content.Context;
import android.os.AsyncTask;
import android.util.Log;
import android.widget.Toast;

/**
 * Created by xxx on 18/08/15.
 */
public class CalculateAsyncTask extends AsyncTask<Void, Long, Long> {
    private ProgressDialog mProgressDialog;
    private Context mContext;

    CalculateAsyncTask(Context context)
    {
        mContext = context;
    }

    /*
     * Runs in its own thread when .execute is called.
     * Must not interact with UI components.
     */
    @Override
    protected Long doInBackground(Void... filenames) {
        Log.d(Constants.MAD, "doInBackground() start");

        long result = 0;
        for (long i = 1; i <= 10000000; i++) {
            if ((i % 329) == 0) {
                publishProgress(i);
            }
            result += i;

            if (isCancelled())
            {
                Log.d(Constants.MAD, "doInBackground() cancelled");
                break;
            }
        }

        Log.d(Constants.MAD, "doInBackground() end");
        return result;
    }

    /*
     * Runs on UI thread, called to update UI with progress when
     * publishProgress is called on the background thread.
     */
    protected void onProgressUpdate(Long... progress) {
        mProgressDialog.setProgress((int)(progress[0] + 0));
    }

    /*
     * Called prior to running doInBackground, on UI thread.
     * Intended to initialise UI components such as progress dialogs.
     */
    @Override
    protected  void onPreExecute()
    {
        Log.d(Constants.MAD, "onPreExecute()");

        String fname = "";
        mProgressDialog = new ProgressDialog(mContext);
        mProgressDialog.setTitle(mContext.getString(R.string.calculating_message));
        mProgressDialog.setProgressStyle(ProgressDialog.STYLE_HORIZONTAL);
        mProgressDialog.setMax(10000000);
        mProgressDialog.show();
    }

    /*
     * Called following the completion of doInBackground, runs on the UI thread.
     */
    @Override
    protected void onPostExecute(Long result) {
        Log.d(Constants.MAD, "onPostExecute()");
        mProgressDialog.dismiss();

        Toast toast = Toast.makeText(mContext, result.toString(), Toast.LENGTH_SHORT);
        toast.show();
    }
}

