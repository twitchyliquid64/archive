package net.ciphersink.nightout.model;

import android.util.Log;

import net.ciphersink.nightout.Constants;

import java.lang.reflect.Member;
import java.util.ArrayList;

/**
 * Encapsulates the data contained within a row of the squad table (remote) - ORM style.
 */
public class Squad {
    private String mName;
    private String mAccessKey;
    private int mId;
    private ArrayList<SquadMember> mMembers;

    public String getAccessKey() {
        return mAccessKey;
    }

    public void setAccessKey(String mAccessKey) {
        this.mAccessKey = mAccessKey;
    }

    public Squad(int squadID, String name, String accesskey) {
        mName = name;
        mAccessKey = accesskey;
        mId = squadID;
        mMembers = new ArrayList<SquadMember>();
    }

    public void debugPrintMembers() {
        for (SquadMember m : mMembers) {
            Log.d(Constants.MAD, "Member name: " + m.getName() + " (" + m.getUsername() + ")" + m.getId());
        }
    }

    public void addMember(SquadMember member) {
        mMembers.add(member);
    }

    public ArrayList<SquadMember> getMembers() {
        return mMembers;
    }

    public String getName() {
        return mName;
    }

    public void setName(String mName) {
        this.mName = mName;
    }

    public int getId() {
        return mId;
    }

    public void setId(int mId) {
        this.mId = mId;
    }
}
