package net.ciphersink.nightout.view;

import android.support.v4.app.Fragment;
import android.support.v7.widget.RecyclerView;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.model.Notification;

import java.util.ArrayList;

/**
 * Maps Notification objects (object model) to Notification Card View (UI).
 */
public class FeedAdapter extends RecyclerView.Adapter<FeedAdapter.ViewHolder>  {

    private ArrayList<Notification> mNotifications;
    private Fragment mFrag;

    /**
     * Encapsulates the instances of different UI elements within the view.
     */
    public static class ViewHolder extends RecyclerView.ViewHolder {
        LinearLayout layout;
        TextView content;
        TextView subLine;
        ImageView icon;
        //put all views here

        public ViewHolder(View layout) {
            super(layout);
            this.layout = (LinearLayout)layout;
            this.content = (TextView)this.layout.findViewById(R.id.notificationCardContent);
            this.subLine = (TextView)this.layout.findViewById(R.id.notificationCardSubLine);
            this.icon = (ImageView)this.layout.findViewById(R.id.notificationCardImage);
        }
    }

    /**
     * Constructs the feed adapter, basing initial data off the given notifications list.
     * @param notifications ArrayList of original notifications
     * @param fragment Fragment which the RecyclerView is contained in
     */
    public FeedAdapter(ArrayList<Notification> notifications, Fragment fragment) {
        mNotifications = notifications;
        mFrag = fragment;
    }

    /**
     * Update the data which underlines the view.
     * @param notifications
     */
    public void setNotifications(ArrayList<Notification> notifications) {
        mNotifications = notifications;
    }

    @Override
    public FeedAdapter.ViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {
        View v = LayoutInflater.from(parent.getContext()).inflate(R.layout.notification_card, parent, false);
        return new ViewHolder(v);
    }


    @Override
    public void onBindViewHolder(ViewHolder holder, int position) {
        Notification notification = mNotifications.get(position);
        holder.content.setText(notification.getContent());
        holder.subLine.setText(notification.getSubLine());

        if (notification.getType().equals(Constants.NET.NOTIFICATON_TYPE.SQUAD_MESSAGE)) {
            holder.icon.setImageDrawable(mFrag.getActivity().getResources().getDrawable(R.drawable.person_many));
        } else if (notification.getType().equals(Constants.NET.NOTIFICATON_TYPE.DRINK_MESSAGE)) {
            holder.icon.setImageDrawable(mFrag.getActivity().getResources().getDrawable(R.drawable.ic_drink));
        } else  {
            holder.icon.setImageDrawable(mFrag.getActivity().getResources().getDrawable(R.drawable.ic_person));
        }
    }

    @Override
    public int getItemCount() {
        return mNotifications.size();
    }
}
