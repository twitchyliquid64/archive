package net.ciphersink.nightout.view;

import android.app.Activity;
import android.net.Uri;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import net.ciphersink.nightout.R;

import java.util.Random;

/**
 * Implements a UI which indicates the system is loading (downloading data from remote model)
 */
public class LoadingFragment extends Fragment {

    private OnFragmentInteractionListener mListener;
    private TextView mLoadingMessage;

    public LoadingFragment() {
        // Required empty public constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(getString(R.string.app_name));
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_loading, container, false);
    }

    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);

        mLoadingMessage = (TextView)getView().findViewById(R.id.loadingFragText);
        setRandomMessage();
    }

    /**
     * Used internally to pick a random loading message and display it.
     */
    private void setRandomMessage() {
        String[] candidateMessages = getResources().getStringArray(R.array.loading_messages);
        Random r = new Random();
        int index = r.nextInt(candidateMessages.length-1);
        mLoadingMessage.setText(candidateMessages[index] + " " + getString(R.string.loading_message_suffix));
    }

    @Override
    public void onAttach(Activity activity) {
        super.onAttach(activity);
        try {
            mListener = (OnFragmentInteractionListener) activity;
        } catch (ClassCastException e) {
            throw new ClassCastException(activity.toString()
                    + " must implement OnFragmentInteractionListener");
        }
    }

    @Override
    public void onDetach() {
        super.onDetach();
        mListener = null;
    }

    public interface OnFragmentInteractionListener {
        public void onFragmentInteraction(Uri uri);
    }

}
