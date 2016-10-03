package net.ciphersink.nightout.model;

/**
 * Encapsulates the data associated with a squad member.
 */
public class SquadMember {
    private String mName;
    private String mUsername;
    private int mId;

    public SquadMember(int id, String name, String username) {
        mId = id;
        mName = name;
        mUsername = username;
    }

    public String getName() {
        return mName;
    }

    public void setName(String mName) {
        this.mName = mName;
    }

    public String getUsername() {
        return mUsername;
    }

    public void setUsername(String mUsername) {
        this.mUsername = mUsername;
    }

    public int getId() {
        return mId;
    }

    public void setId(int mId) {
        this.mId = mId;
    }
}
