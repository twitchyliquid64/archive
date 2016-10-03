package net.ciphersink.nightout.view;


import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.support.v7.widget.LinearLayoutManager;
import android.support.v7.widget.RecyclerView;
import android.text.ClipboardManager;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Toast;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.Interfaces;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Squad;
import net.ciphersink.nightout.model.SquadMember;

/**
 * Implements the UI where users can see a list of all the members of a squad they are in,
 * and perform actions on the members / squad.
 */
public class SquadDisplayFragment extends Fragment implements Interfaces.MenuControllerInterface,
    View.OnClickListener {

    // ui
    private RecyclerView mRecyclerView;
    private RecyclerView.Adapter mAdapter;
    private RecyclerView.LayoutManager mLayoutManager;

    //model
    private Squad mSquad;
    private int mSquadListIndex;

    public SquadDisplayFragment(Squad squad, int squadListIndex) {
        mSquad = squad;
        mSquadListIndex = squadListIndex;
    }

    public SquadDisplayFragment() {
        //required default constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(mSquad.getName());
        act.initialiseMenu(R.id.mainMenuShare, this);
        act.initialiseMenu(R.id.mainMenuMessage, this);
    }

    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);

        //init UI
        mRecyclerView = (RecyclerView)getView().findViewById(R.id.squadDisplayFragRecyclerView);
        mRecyclerView.setHasFixedSize(true);

        mLayoutManager = new LinearLayoutManager(getActivity());
        mRecyclerView.setLayoutManager(mLayoutManager);

        mAdapter = new SquadMemberAdapter(mSquad, this);
        mRecyclerView.setAdapter(mAdapter);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_squad_display, container, false);
    }

    /**
     * Implements MenuControllerInterface - called from MainActivity when a Action bar menu
     * is pressed. In this case, that could be the share or message icons.
     * @param resId
     */
    public void menuClicked(int resId) {
        Log.d(Constants.MAD, "Frag got menu press: " + resId);

        switch (resId) {
            case R.id.mainMenuShare:
                // set squadkey to clipboard
                ClipboardManager _clipboard = (ClipboardManager) getActivity().getSystemService(Context.CLIPBOARD_SERVICE);
                _clipboard.setText(mSquad.getAccessKey());

                // display a message
                Toast toast = Toast.makeText(getActivity(), getString(R.string.squad_key_copied) + " (" + mSquad.getAccessKey() + ")", Toast.LENGTH_SHORT);
                toast.show();
                break;

            case R.id.mainMenuMessage:
                // start the NewMessageActivity in send-to-the-whole-squad mode.
                Intent messageActivityIntent = new Intent(getActivity(), NewMessageActivity.class);
                messageActivityIntent.putExtra(Constants.KEYS.NAME, mSquad.getName());
                messageActivityIntent.putExtra(Constants.KEYS.FROM_NAME, ((MainActivity)getActivity()).getSession().getName());
                messageActivityIntent.putExtra(Constants.KEYS.SESSIONKEY, ((MainActivity)getActivity()).getSession().getKey());
                messageActivityIntent.putExtra(Constants.KEYS.MODE, Constants.KEYS.NOTIFICATION_REQUEST.MODE_SEND_ALL_SQUAD);
                messageActivityIntent.putExtra(Constants.KEYS.SQUAD_ID, mSquad.getId());
                startActivity(messageActivityIntent);
                break;
        }
    }

    /**
     * Called from the adapter whenever a cardView button is pressed.
     * @param v
     */
    public void onClick(View v) {
        // get index of cardview.
        int position = mRecyclerView.getChildAdapterPosition((View) v.getParent().getParent().getParent());
        Log.d(Constants.MAD, "SquadDisplayFragment onClick(): " + position);
        // get the corresponding SquadMember
        MainActivity act = (MainActivity) getActivity();
        SquadMember member = mSquad.getMembers().get(position);

        // do different actions depending on which button was pressed
        switch (v.getId()) {
            case R.id.squadMemberViewFindButton:
                act.loadPaneTrackerPage(member); //track / find my mate button pressed - swap out fragment
                break;

            case R.id.squadMemberViewMessageButton:
                Log.d(Constants.MAD, "Message button pressed");
                Intent messageActivityIntent = new Intent(getActivity(), NewMessageActivity.class);
                messageActivityIntent.putExtra(Constants.KEYS.NAME, member.getName());
                messageActivityIntent.putExtra(Constants.KEYS.SESSIONKEY, ((MainActivity)getActivity()).getSession().getKey());
                messageActivityIntent.putExtra(Constants.KEYS.FROM_NAME, ((MainActivity)getActivity()).getSession().getName());
                messageActivityIntent.putExtra(Constants.KEYS.USER_ID, member.getId());
                messageActivityIntent.putExtra(Constants.KEYS.MODE, Constants.KEYS.NOTIFICATION_REQUEST.MODE_SEND_SINGLE);
                startActivity(messageActivityIntent);
                break;
        }
    }
}
