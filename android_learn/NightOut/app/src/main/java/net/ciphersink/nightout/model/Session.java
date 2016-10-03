package net.ciphersink.nightout.model;

import java.util.ArrayList;

/**
 * Encapsulates all the data associated with a session - ORM style.
 */
public class Session {
    private String mKey;
    private String mUsername;
    private String mName;
    private String mEmail;
    private int mUserId;

    /**
     * Returns all the squadID's which the user (session) is a part of.
     * @return Squad ID's
     */
    public ArrayList<Integer> getSquadIds() {
        return mSquadIds;
    }

    private ArrayList<Integer> mSquadIds;

    /**
     * Creates a session object with the given data.
     * @param key
     * @param userID
     * @param username
     * @param name
     * @param email
     * @param squadIds
     */
    public Session(String key, int userID, String username, String name, String email, ArrayList<Integer> squadIds) {
        mKey = key;
        mUserId = userID;
        mUsername = username;
        mName = name;
        mEmail = email;
        mSquadIds = squadIds;
    }

    public String getKey() {
        return mKey;
    }

    public String getUsername() {
        return mUsername;
    }

    public String getName() {
        return mName;
    }

    public String getEmail() {
        return mEmail;
    }

    public int getUserId() {
        return mUserId;
    }
}
