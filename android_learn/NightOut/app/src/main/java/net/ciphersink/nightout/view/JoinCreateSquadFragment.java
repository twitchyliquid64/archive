package net.ciphersink.nightout.view;

import android.app.Activity;
import android.net.Uri;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.Toast;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Session;
import net.ciphersink.nightout.model.Squad;
import net.ciphersink.nightout.model.SquadFactory;


public class JoinCreateSquadFragment extends Fragment implements View.OnClickListener{
    private OnFragmentInteractionListener mListener;
    private Button mJoinSquadButton;
    private Button mCreateSquadButton;
    private EditText mSquadKeyEditText;
    private EditText mSquadNameEditText;
    private LinearLayout mControlsContainer;
    private ProgressBar mProgress;

    public JoinCreateSquadFragment() {
        // Required empty public constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(getString(R.string.joinCreateSquadFragTitle));
    }

    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);
        mCreateSquadButton = (Button)getView().findViewById(R.id.joinCreateFragCreateButton);
        mJoinSquadButton = (Button)getView().findViewById(R.id.joinCreateFragJoinButton);
        mSquadKeyEditText = (EditText)getView().findViewById(R.id.joinCreateFragSquadKeyEditText);
        mSquadNameEditText = (EditText)getView().findViewById(R.id.joinCreateFragSquadNameEditText);
        mControlsContainer = (LinearLayout)getView().findViewById(R.id.joinCreateFragControlsContainer);
        mProgress = (ProgressBar)getView().findViewById(R.id.joinCreateFragProgress);

        mCreateSquadButton.setOnClickListener(this);
        mJoinSquadButton.setOnClickListener(this);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_join_create_squad, container, false);
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

    @Override
    public void onClick(View v) {
        if (v.getId() == R.id.joinCreateFragCreateButton) { //create button
            String squadName = mSquadNameEditText.getText().toString();
            if (squadName.equals("")) {
                mSquadNameEditText.setError("Please enter a name");
            } else {
                Log.d(Constants.MAD, "Create squad pressed: " + squadName);
                new NewSquadTask().execute(squadName);
            }
        } else if (v.getId() == R.id.joinCreateFragJoinButton) {//join button
            String squadKey = mSquadKeyEditText.getText().toString();
            Log.d(Constants.MAD, "Join squad pressed: " + squadKey);
            new JoinSquadTask().execute(squadKey);
        }
    }



    public interface OnFragmentInteractionListener {
        // TODO: Update argument type and name
        public void onFragmentInteraction(Uri uri);
        public void newSquadNotify(Squad squad);
    }

    private Session getSession() {
        MainActivity act = (MainActivity)getActivity();
        return act.getSession();
    }

    private class NewSquadTask extends AsyncTask<String, Void, Void> {

        Squad mSquad;

        @Override
        protected void onPreExecute() {
            mControlsContainer.setVisibility(View.GONE);
            mProgress.setVisibility(View.VISIBLE);
        }

        @Override
        protected Void doInBackground(String... squadName) {
            Squad newSquad = SquadFactory.createSquad(squadName[0], getSession().getKey());
            Log.d(Constants.MAD, "Squad access key: " + newSquad.getAccessKey());
            Log.d(Constants.MAD, "Squad name: " + newSquad.getName());
            newSquad.debugPrintMembers();
            mSquad = newSquad;
            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            MainActivity act = (MainActivity)getActivity();
            act.newSquadNotify(mSquad);
            Toast toast = Toast.makeText(getActivity(), getString(R.string.squad_key_good) + " " + mSquad.getName() + "!", Toast.LENGTH_SHORT);
            toast.show();
            mControlsContainer.setVisibility(View.VISIBLE);
            mProgress.setVisibility(View.GONE);
            mSquadNameEditText.setText("");
        }
    }

    private class JoinSquadTask extends AsyncTask<String, Void, Void> {
        Squad mSquad;

        @Override
        protected void onPreExecute() {
            mControlsContainer.setVisibility(View.GONE);
            mProgress.setVisibility(View.VISIBLE);
        }

        @Override
        protected Void doInBackground(String... squadAccessKey) {
            Squad newSquad = SquadFactory.joinSquad(squadAccessKey[0], getSession().getKey());
            if (newSquad != null) {
                Log.d(Constants.MAD, "Squad name: " + newSquad.getName());
                newSquad.debugPrintMembers();
                mSquad = newSquad;
            }
            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            if (mSquad != null) {
                MainActivity act = (MainActivity) getActivity();
                act.newSquadNotify(mSquad);
                Toast toast = Toast.makeText(getActivity(), getString(R.string.squad_key_good) + " " + mSquad.getName() + "!", Toast.LENGTH_SHORT);
                toast.show();
                mSquadKeyEditText.setText("");
            } else {
                Toast toast = Toast.makeText(getActivity(), getString(R.string.squad_key_invalid), Toast.LENGTH_SHORT);
                toast.show();
            }
            mControlsContainer.setVisibility(View.VISIBLE);
            mProgress.setVisibility(View.GONE);
        }
    }


}
