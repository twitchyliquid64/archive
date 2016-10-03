package net.ciphersink.nightout.view;

import android.support.v4.app.Fragment;
import android.support.v7.widget.RecyclerView;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageButton;
import android.widget.LinearLayout;
import android.widget.TextView;

import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Squad;
import net.ciphersink.nightout.model.SquadMember;

/**
 * Implements the UI <-> Data model mapping in SquadDisplayFragments' recycler view.
 */
public class SquadMemberAdapter extends RecyclerView.Adapter<SquadMemberAdapter.ViewHolder> {

    private Squad mSquad;
    private Fragment mFrag;

    /**
     * Encapsulates the pointers to UI elements
     */
    public static class ViewHolder extends RecyclerView.ViewHolder {
        LinearLayout layout;
        TextView nameView;
        TextView lowLineView;
        ImageButton trackButton;
        ImageButton messageButton;

        public ViewHolder(View layout) {
            super(layout);
            this.layout = (LinearLayout)layout;
            nameView = (TextView) this.layout.findViewById(R.id.squadMemberViewMemberName);
            lowLineView = (TextView) this.layout.findViewById(R.id.squadMemberViewLowLine);
            trackButton = (ImageButton)this.layout.findViewById(R.id.squadMemberViewFindButton);
            messageButton = (ImageButton)this.layout.findViewById(R.id.squadMemberViewMessageButton);
        }
    }

    public SquadMemberAdapter(Squad squad, Fragment fragment) {
        mSquad = squad;
        mFrag = fragment;
    }

    @Override
    public SquadMemberAdapter.ViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {

        View v = LayoutInflater.from(parent.getContext()).inflate(R.layout.view_squad_member_card, parent, false);
        return new ViewHolder(v);
    }

    /**
     * Sets up the data for a given card with a given data model position.
     * @param holder Card
     * @param position DM index
     */
    @Override
    public void onBindViewHolder(ViewHolder holder, int position) {
        SquadMember member = mSquad.getMembers().get(position);
        holder.nameView.setText(member.getName());
        holder.lowLineView.setText(member.getUsername());
        holder.trackButton.setOnClickListener((View.OnClickListener) mFrag);
        holder.messageButton.setOnClickListener((View.OnClickListener) mFrag);
    }

    @Override
    public int getItemCount() {
        return mSquad.getMembers().size();
    }
}
