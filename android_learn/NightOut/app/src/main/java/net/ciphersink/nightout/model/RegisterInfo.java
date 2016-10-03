package net.ciphersink.nightout.model;


import android.content.Context;
import android.util.Log;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.util.EasyHttp;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;

import java.util.ArrayList;

/**
 * Encapsulates the data associated with registering a user.
 */
public class RegisterInfo {
    private String mUsername;
    private String mPassword;
    private String mEmail;
    private String mName;

    private boolean mUsernameUnique = false;
    private boolean mValidationSuccessful = false;

    //called to validate that all the contents of the model are correct - able
    //to be submitted with success. Specifically, this checks that the username
    //is not already in use, and that the required information is populated.
    //
    //whilst it is resilient against multiple calls, it is not advised. If the caller
    //is sure that validationInformation() has already been called, the cached result
    //can be obtained from validationSuccessful().
    public boolean validateInformation() {
        if(!usernamePopulated())return false;
        if(!namePopulated())return false;
        if(!emailPopulated())return false;
        if(!passwordPopulated())return false;
        if(!validateUsernameUnique())return false;

        mValidationSuccessful = true;
        return true;
    }

    //should only be used if the caller is sure that validateInformation()
    //has already been called. If not, user validateInformation() instead.
    public boolean validationSuccessful() {
        return mValidationSuccessful;
    }

    //should only be used if the caller is sure that validateInformation()
    //has already been called. If not, call validateInformation() first.
    public boolean usernameUnique() {
        return mUsernameUnique;
    }

    //returns false if the username is textually invalid.
    public boolean usernamePopulated() {
        return !mUsername.equals("");
    }

    //returns false if the name is textually invalid.
    public boolean namePopulated() {
        return !mName.equals("");
    }

    //returns false if the email is textually invalid.
    public boolean emailPopulated() {
        return !mEmail.equals("");
    }

    //returns false if the password is textually invalid.
    public boolean passwordPopulated() {
        return !mPassword.equals("");
    }

    private boolean validateUsernameUnique() {
        //get cached value if set-positive
        if(mUsernameUnique)return true;

        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.USERNAME, String.valueOf(mUsername)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.CHECK_USERNAME, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_OK)) {
            mUsernameUnique = true;
            return true;
        }

        return false;
    }

    //Returns true id the user was successfully registered on the server.
    //the user should call this when the registration information is valid and it should be
    //committed to the server.
    public boolean register() {
        ArrayList<NameValuePair> params = new ArrayList<NameValuePair>();
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.USERNAME, String.valueOf(mUsername)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.NAME, String.valueOf(mName)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.PASSWORD, String.valueOf(mPassword)));
        params.add(new BasicNameValuePair(Constants.NET.PARAM_KEYS.EMAIL, String.valueOf(mEmail)));

        EasyHttp request = new EasyHttp(Constants.NET.REST_ENDPOINT.REGISTER, params);

        if(request.didError())return false;
        if(request.getData().equals(Constants.NET.STANDARD_OK)) {
            mUsernameUnique = true;
            return true;
        }

        return false;
    }

    public String getUsername() {
        return mUsername;
    }

    public void setUsername(String mUsername) {
        this.mUsername = mUsername;
    }

    public String getPassword() {
        return mPassword;
    }

    public void setPassword(String mPassword) {
        this.mPassword = mPassword;
    }

    public String getEmail() {
        return mEmail;
    }

    public void setEmail(String mEmail) {
        this.mEmail = mEmail;
    }

    public String getName() {
        return mName;
    }

    public void setName(String mFirstname) {
        this.mName = mFirstname;
    }
}
