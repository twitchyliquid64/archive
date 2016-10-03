package net.ciphersink.nightout;

/**
 * Encapsulates all magic numbers and reused values throughout the application, to
 * reduce duplication and associate values => names.
 */
public class Constants {
    public static int RELEASE_VERSION = 0;
    public static int MAJOR_VERSION = 0;
    public static int MINOR_VERSION = 16;

    public static final String MAD = "MAD";
    public static final int LOCATION_FIDELITY = 45 * 1000;

    public class KEYS {
        public class REMOTE_ACTIVITY_REQUEST {
            public static final int REQUEST_IMAGE_CAPTURE = 1;
        }

        public class NOTIFICATION_REQUEST {
            public static final int MODE_SEND_SINGLE = 0;
            public static final int MODE_SEND_ALL_SQUAD = 1;
        }

        public class REGISTER_ACTIVITY {
            public static final int REQUESTCODE_REGISTER = 0;
            public static final int RESPONSECODE_REGISTRATION_SUCCESS = 1;
        }

        public static final String SESSIONKEY = "sessionkey";
        public static final String USERNAME = "username";
        public static final String NAME = "name";
        public static final String FROM_NAME = "fromname";
        public static final String USER_ID = "userid";
        public static final String MODE = "mode";
        public static final String SQUAD_ID = "squadid";
        public static final String DRINK_COUNT = "drinkcount";
    }

    public class RES {
        public static final String PREFERENCES_FILE = "net.ciphersink.nightout.preferences";
    }

    public class NET {
        public static final String NET_URI = "http://";
        public static final String ADDRESS = "nightout.ciphersink.net";

        public class REST_ENDPOINT {
            public static final String CHECK_USERNAME = "/userexists";
            public static final String REGISTER = "/register";
            public static final String GET_SESSION = "/session/getdetails";
            public static final String CREATE_SESSION = "/session/new";
            public static final String NEW_SQUAD = "/squad/new";
            public static final String JOIN_SQUAD = "/squad/join";
            public static final String GET_SQUAD_DETAILS = "/squad/getdetails";
            public static final String TRANSMIT_LOCATION = "/location/push";
            public static final String GET_USER_LOCATION = "/location/get";
            public static final String GET_NOTIFICATIONS = "/session/notifications";
            public static final String SEND_NOTIFICATION = "/session/notifications/new";
            public static final String SEND_SQUAD_NOTIFICATION = "/session/notifications/squad/new";
            public static final String SEND_ALL_NOTIFICATION = "/session/notifications/all/new";
        }

        public class NOTIFICATON_TYPE {
            public static final String MESSAGE = "message";
            public static final String SQUAD_MESSAGE = "squadmessage";
            public static final String DRINK_MESSAGE = "drinkmessage";
        }

        public class PARAM_KEYS {
            public static final String USERNAME = "username";
            public static final String NAME = "name";
            public static final String EMAIL = "email";
            public static final String PASSWORD = "password";
            public static final String KEY = "key";
            public static final String SQUAD_ID = "squadid";
            public static final String SQUAD_KEY = "squadkey";
            public static final String USER_ID = "userid";
            public static final String TYPE = "type";
            public static final String CONTENT = "content";
            public static final String SUBLINE = "subline";

            public static final String LATITUDE = "lat";
            public static final String LONGITUDE = "lon";
            public static final String ACCURACY = "acc";
            public static final String PROVIDER = "prov";
            public static final String BATTERY = "batt";
        }

        public static final String STANDARD_OK = "OK";
        public static final String STANDARD_ERROR = "ERROR";
    }
}
