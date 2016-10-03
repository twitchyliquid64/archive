package net.ciphersink.exercise4;

import android.app.ProgressDialog;
import android.content.Context;
import android.os.AsyncTask;
import android.util.Log;


public class DownloadSelectedAsyncTask extends AsyncTask<Void, Void, Void> {
    String mFileName;
    ProgressDialog mProgressDialog;
    Context mContext;

    DownloadSelectedAsyncTask(Context context, String fileName)
    {
        mFileName = fileName;
        mContext = context;
    }

    /*
     * Runs in its own thread when .execute is called.
     * Must not interact with UI components.
     */
    @Override
    protected Void doInBackground(Void... dummy) {
        Log.d(Constants.MAD, "doInBackground() start");
        try {
            Thread.sleep(2000);
        }
        catch(InterruptedException e)
        {
            // don't care if sleep interrupted.
        }
        Log.d(Constants.MAD, "doInBackground() end");
        return null;
    }

    /*
     * Runs on UI thread, called to update UI with progress when
     * publishProgress is called on the background thread.
     */
    protected void onProgressUpdate(Integer... progress) {
        // dont do anything
    }

    /*
     * Called prior to running doInBackground, on UI thread.
     * Intended to initialise UI components such as progress dialogs.
     */
    @Override
    protected  void onPreExecute()
    {
        Log.d(Constants.MAD, "onPreExecute()");
        mProgressDialog = new ProgressDialog(mContext);
        mProgressDialog.setTitle(mContext.getString(R.string.progress_title_prefix) + mFileName);
        mProgressDialog.show();
    }

    /*
     * Called following the completion of doInBackground, runs on the UI thread.
     */
    @Override
    protected void onPostExecute(Void dummy) {
        Log.d(Constants.MAD, "onPostExecute()");
        mProgressDialog.dismiss();
    }
}

