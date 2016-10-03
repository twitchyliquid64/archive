package net.ciphersink.exercise6;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;
import android.util.Log;

import java.util.ArrayList;

/**
 * Created by xxx on 1/09/15.
 */
public class FriendDatabaseHelper extends SQLiteOpenHelper {

    public FriendDatabaseHelper(Context context)
    {
        super(context, Constants.DB.FNAME, null, 1);
    }

    @Override
    public void onCreate(SQLiteDatabase db)
    {
        db.execSQL(Constants.DB.CREATE_FRIENDTABLE_SQL);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion)
    {
        //Dont want to loose data while testing - TODO uncomment in production
        //Log.d(Constants.MAD, "Deleting data - Upgrading!");
        //db.execSQL("DROP TABLE IF EXISTS " + Constants.DB.FRIENDTABLE);
        //onCreate(db);
    }

    public void deleteFriend(int rowID)
    {
        SQLiteDatabase db = getWritableDatabase();
        db.execSQL("DELETE FROM " + Constants.DB.FRIENDTABLE + " WHERE " + Constants.DB.IDCOL + " = ?", new Object[]{rowID});
    }

    public ArrayList<FriendData> getFriendsByNameFilter(String nameByFilter)
    {
        SQLiteDatabase db = getWritableDatabase();
        ArrayList<FriendData> retData = new ArrayList<FriendData>();

        String[] cols = {Constants.DB.IDCOL,
                         Constants.DB.NAMECOL,
                         Constants.DB.OCCUPATIONCOL,
                         Constants.DB.CITYCOL,
                         Constants.DB.FRIENDSINCECOL};

        Cursor getAllQry = db.query(Constants.DB.FRIENDTABLE, cols,
                                    Constants.DB.NAMECOL + " LIKE ?",
                                    new String[]{"%" + nameByFilter + "%"}, null, null, null);

        if (getAllQry.moveToFirst()) {
            do {
                int _id = getAllQry.getInt(0);
                String name = getAllQry.getString(1);
                String occupation = getAllQry.getString(2);
                String city = getAllQry.getString(3);
                Long friendSince = getAllQry.getLong(4);

                if (name != null)
                {
                    retData.add(new FriendData(_id, name, occupation, city, friendSince));
                }

            } while(getAllQry.moveToNext());
        }

        return retData;
    }

    public ArrayList<FriendData> getAllFriends()
    {
        SQLiteDatabase db = getWritableDatabase();
        ArrayList<FriendData> retData = new ArrayList<FriendData>();

        String[] cols = {Constants.DB.IDCOL,
                Constants.DB.NAMECOL,
                Constants.DB.OCCUPATIONCOL,
                Constants.DB.CITYCOL,
                Constants.DB.FRIENDSINCECOL};

        Cursor getAllQry = db.query(Constants.DB.FRIENDTABLE, cols, null, null, null, null, null);

        if (getAllQry.moveToFirst()) {
            do {
                int _id = getAllQry.getInt(0);
                String name = getAllQry.getString(1);
                String occupation = getAllQry.getString(2);
                String city = getAllQry.getString(3);
                Long friendSince = getAllQry.getLong(4);

                if (name != null)
                {
                    retData.add(new FriendData(_id, name, occupation, city, friendSince));
                }

            } while(getAllQry.moveToNext());
        }

        return retData;
    }

    public void addFriend(FriendData friend)
    {
        SQLiteDatabase db = getWritableDatabase();

        ContentValues friendValues = new ContentValues();
        friendValues.put(Constants.DB.NAMECOL, friend.getName());
        friendValues.put(Constants.DB.CITYCOL, friend.getCity());
        friendValues.put(Constants.DB.OCCUPATIONCOL, friend.getOccupation());
        friendValues.put(Constants.DB.FRIENDSINCECOL, friend.getFriendSince());

        db.insert(Constants.DB.FRIENDTABLE, null, friendValues); //TODO: Use insertOrThrow
    }
}
