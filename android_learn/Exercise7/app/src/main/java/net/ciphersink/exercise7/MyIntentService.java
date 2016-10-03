package net.ciphersink.exercise7;

import android.app.IntentService;
import android.content.Intent;
import android.util.Log;

/**
 * Created by xxx on 8/09/15.
 */
public class MyIntentService extends IntentService {

    public MyIntentService() {
        super("MyIntentService");//required by IntentService
    }

    @Override
    protected void onHandleIntent(Intent intent)
    {
        Log.d(Constants.MAD, "MyIntentService.onHandleIntent(): " + intent.toString());
        int maxValue = intent.getIntExtra(Constants.COUNT_REQUEST.KEY_MAXVAL, 0);

        for(int i = 0; i < maxValue; i++)
        {
            sendMessage("" + (i+1) + " " + getString(R.string.combinationCountStr) + " " + maxValue);
            try {
                Thread.sleep(1000, 0);
            }
            catch (InterruptedException e)
            {
                //dont handle interrupt exceptions
            }
        }
        sendMessage(getString(R.string.taskCompleted));
    }

    private void sendMessage(String inMsg)
    {
        Log.d(Constants.MAD, "Intent message: " + inMsg);
        Intent i = new Intent(Constants.COUNT_REQUEST.INTENT_RESULT);
        i.putExtra(Constants.COUNT_REQUEST.KEY_RETMSG, inMsg);
        sendBroadcast(i);
    }

    @Override
    public void onDestroy()
    {
        sendMessage(getString(R.string.serviceStopMsg));
    }
}
