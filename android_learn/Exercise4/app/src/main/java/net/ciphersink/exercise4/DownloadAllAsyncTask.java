package net.ciphersink.exercise4;

import android.app.ProgressDialog;
import android.content.Context;
import android.os.AsyncTask;
import android.util.Log;

/**
 * Created by xxx on 18/08/15.
 */
public class DownloadAllAsyncTask extends AsyncTask<String, String, Void> {
    private int mNumFiles;
    private int mFilesDone;
    private ProgressDialog mProgressDialog;
    private Context mContext;

    DownloadAllAsyncTask(Context context, int numFiles)
    {
        mNumFiles = numFiles;
        mContext = context;
        mFilesDone = 0;
    }

    /*
     * Runs in its own thread when .execute is called.
     * Must not interact with UI components.
     */
    @Override
    protected Void doInBackground(String... filenames) {
        Log.d(Constants.MAD, "doInBackground() start");

        for (String file : filenames) {
            publishProgress(mContext.getString(R.string.progress_title_prefix) + file);
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                // don't care if sleep interrupted.
            }
        }
        Log.d(Constants.MAD, "doInBackground() end");
        return null;
    }

    /*
     * Runs on UI thread, called to update UI with progress when
     * publishProgress is called on the background thread.
     */
    protected void onProgressUpdate(String... progress) {
        Log.d(Constants.MAD, "onProgressUpdate() " + progress[0]);
        mProgressDialog.setTitle(progress[0]);
        mProgressDialog.setProgress(100 * mFilesDone / mNumFiles);
        mFilesDone++;
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
        mProgressDialog.setTitle(mContext.getString(R.string.progress_title_prefix));
        mProgressDialog.setProgressStyle(ProgressDialog.STYLE_HORIZONTAL);
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

