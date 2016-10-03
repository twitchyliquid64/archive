package net.ciphersink.nightout.view;

import android.app.Activity;
import android.net.Uri;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;

import net.ciphersink.nightout.R;

/**
 * Implements the UI which you see when you first start using the app (helpful messages).
 */
public class FirstStartFragment extends Fragment {
    private OnFragmentInteractionListener mListener;
    private Animation mArrowAnimation;
    private ImageView mArrowImage;

    public FirstStartFragment() {
        // Required empty public constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        MainActivity act = (MainActivity)getActivity();
        act.setBarTitle(getString(R.string.app_name));
    }

    @Override
    public void onActivityCreated(Bundle savedInstanceState) {
        super.onActivityCreated(savedInstanceState);

        mArrowImage = (ImageView)getView().findViewById(R.id.firstStartFragArrowImg);
        mArrowAnimation = AnimationUtils.loadAnimation(getActivity(), R.anim.firststart_rotate_arrow);
        mArrowImage.startAnimation(mArrowAnimation);
        mArrowAnimation.setFillAfter(true);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_first_start, container, false);
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
        // TODO: Update argument type and name
        public void onFragmentInteraction(Uri uri);
    }

}
