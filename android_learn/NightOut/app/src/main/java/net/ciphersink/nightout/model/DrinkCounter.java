package net.ciphersink.nightout.model;

import android.content.Context;
import android.content.SharedPreferences;
import android.os.AsyncTask;

import net.ciphersink.nightout.Constants;

import java.lang.reflect.Array;
import java.util.ArrayList;

/**
 * Interface which encapsulates all communication with the data model for drink counting.
 */
public class DrinkCounter {

    /**
     * Increments the drink counter for the given session by one.
     * @param main Current context (activity or service)
     * @param session Session where the drink needs incrementing
     * @param squads List of squads to notify of the change
     */
    public static void anotherDrink(Context main,Session session, ArrayList<Squad> squads) {
        SharedPreferences sharedPref = main.getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        int newCount = sharedPref.getInt(Constants.KEYS.DRINK_COUNT, 0) + 1;

        SharedPreferences.Editor editor = sharedPref.edit();
        editor.putInt(Constants.KEYS.DRINK_COUNT, newCount);
        editor.commit();
        new PostNotificationTask(session, main, squads).execute();
    }

    /**
     * Returns the current count of the drink counter
     * @param main Current Context
     * @return Number of drinks drunk
     */
    public static int getCount(Context main) {
        SharedPreferences sharedPref = main.getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        return sharedPref.getInt(Constants.KEYS.DRINK_COUNT, 0);
    }

    /**
     * Resets the current drink count.
     * @param main Current Context
     */
    public static void reset(Context main) {
        SharedPreferences sharedPref = main.getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        SharedPreferences.Editor editor = sharedPref.edit();
        editor.putInt(Constants.KEYS.DRINK_COUNT, 0);
        editor.commit();
    }

    /**
     * Encapsulates the background task of transmission to the network
     */
    private static class PostNotificationTask extends AsyncTask<Void, Void, Void> {
        Context mContext;
        ArrayList<Squad> mSquads;
        Session mSession;

        public PostNotificationTask(Session session, Context context, ArrayList<Squad> squads) {
            mContext = context;
            mSquads = squads;
            mSession = session;
        }

        @Override
        protected Void doInBackground(Void... dummy) {

            String dSuffix = " drink.";
            if (getCount(mContext) > 1)dSuffix = " drinks.";

            Notification notification = new Notification(mSession.getName() + " has had another drink!",
                    Constants.NET.NOTIFICATON_TYPE.DRINK_MESSAGE, "The count is now " + getCount(mContext) + dSuffix);

            notification.sendToAll(mSession.getKey());
            return null;
        }
    }
}
