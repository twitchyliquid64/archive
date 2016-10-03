package net.ciphersink.exercise6;

/**
 * Created by xxx on 1/09/15.
 */
public class Constants {
    public static final String MAD = "MAD";

    public class DB {
        public static final String FNAME = "friend.db";
        public static final String FRIENDTABLE = "friends";
        public static final String NAMECOL = "name";
        public static final String OCCUPATIONCOL = "occupation";
        public static final String CITYCOL = "city";
        public static final String FRIENDSINCECOL = "friendsince";
        public static final String IDCOL = "_id";

        public static final String CREATE_FRIENDTABLE_SQL = "CREATE TABLE friends (" +
                "_id integer primary key autoincrement," +
                " name TEXT," +
                " occupation TEXT," +
                " city TEXT," +
                " friendsince LONG);";

    }

    public class ADDFRIEND {
        public static final int RESULT_SUCCESS = 10;
        public static final int RESULT_CANCELLED = 11;
        public class KEYS {
            public static final String NAME = "name";
            public static final String OCCUPATION = "occupation";
            public static final String CITY = "city";
        }
    }

    public class VIEWER {
        public static final int MODE_VIEWALL = 2;
        public static final int MODE_SEARCHRESULTS = 3;
        public static final String MODE_KEY = "mode";
        public static final String FILTERNAME_KEY = "friendname";
    }

    public class PREFERENCES {
        public static final String KEY_FRIENDDATEFORMAT = "preferences_friend_dateformat";
        public static final String KEY_FRIENDCANDELETE = "preferences_friend_canDeleteFriends";
    }
}
