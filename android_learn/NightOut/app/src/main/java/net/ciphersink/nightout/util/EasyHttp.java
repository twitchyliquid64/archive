package net.ciphersink.nightout.util;

import android.util.Log;

import net.ciphersink.nightout.Constants;

import org.apache.http.HttpResponse;
import org.apache.http.HttpStatus;
import org.apache.http.NameValuePair;
import org.apache.http.StatusLine;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.utils.URLEncodedUtils;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.HttpConnectionParams;

import java.io.ByteArrayOutputStream;
import java.util.ArrayList;

/**
 * Easy way to make a HTTP call to a REST endpoint on the server, and get data back.
 * This acts as a wrapper method for HttpClient.
 */
public class EasyHttp {

    private boolean mDidError;
    private StatusLine mStatusLine;
    private HttpResponse mResponse;
    private String data;

    /**
     * Creates a new HTTP Request object which will hit a given API endpoint.
     * @param endpoint API endpoint (URL)
     * @param params data to be transmitted
     */
    public EasyHttp(String endpoint, ArrayList<NameValuePair> params) {
        String URL = Constants.NET.NET_URI + Constants.NET.ADDRESS + endpoint + "?" + URLEncodedUtils.format(params, "utf-8");

        Log.d(Constants.MAD, "EasyHTTP: Request for " + URL);

        try {
            HttpClient httpclient = new DefaultHttpClient();
            HttpConnectionParams.setTcpNoDelay(httpclient.getParams(), true);

            mResponse = httpclient.execute(new HttpGet(URL));
            mStatusLine = mResponse.getStatusLine();
            Log.d(Constants.MAD, "StatusCode: " + mStatusLine.getStatusCode());

            if (mStatusLine.getStatusCode() == HttpStatus.SC_OK) {
                ByteArrayOutputStream out = new ByteArrayOutputStream();
                mResponse.getEntity().writeTo(out);
                data = out.toString();
                out.close();
            } else {
                //Closes the connection.
                mResponse.getEntity().getContent().close();
                mDidError = true;
                data = "";
            }
        } catch(java.io.IOException e) {
            mDidError = true;
        }
    }

    /**
     * Returns true if the network connection errored (timeout or non-200 error code).
     * @return diderror == true
     */
    public boolean didError() {
        return mDidError;
    }


    public StatusLine getStatusLine() {
        return mStatusLine;
    }


    public HttpResponse getResponse() {
        return mResponse;
    }

    /**
     * Returns as a string the text returned from the remote HTTP call.
     * @return Plaintext response data
     */
    public String getData() {
        return data;
    }

}
