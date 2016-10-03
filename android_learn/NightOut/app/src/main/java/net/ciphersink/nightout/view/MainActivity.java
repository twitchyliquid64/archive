package net.ciphersink.nightout.view;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.net.Uri;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.design.widget.FloatingActionButton;
import android.util.Log;
import android.view.SubMenu;
import android.view.View;
import android.support.design.widget.NavigationView;
import android.support.v4.view.GravityCompat;
import android.support.v4.widget.DrawerLayout;
import android.support.v7.app.ActionBarDrawerToggle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.ImageButton;
import android.widget.TextView;

import net.ciphersink.nightout.Constants;
import net.ciphersink.nightout.Interfaces;
import net.ciphersink.nightout.R;
import net.ciphersink.nightout.UpdaterService;
import net.ciphersink.nightout.model.DrinkCounter;
import net.ciphersink.nightout.model.Session;
import net.ciphersink.nightout.model.SessionFactory;
import net.ciphersink.nightout.model.Squad;
import net.ciphersink.nightout.model.SquadFactory;
import net.ciphersink.nightout.model.SquadMember;

import java.util.ArrayList;

/**
 * Encapsulates the navigation drawer and a fragment container,
 * within which the majority of functionality is implemented.
 */
public class MainActivity extends AppCompatActivity
        implements NavigationView.OnNavigationItemSelectedListener,
        View.OnClickListener,
        FirstStartFragment.OnFragmentInteractionListener,
        LoadingFragment.OnFragmentInteractionListener,
        JoinCreateSquadFragment.OnFragmentInteractionListener {

    //ui elements
    private Toolbar mToolbar;
    private FloatingActionButton mFloatingActionButton;
    private DrawerLayout mDrawer;
    private ActionBarDrawerToggle mDrawerToggle;
    private NavigationView mNavigationView;
    private TextView mNavNameDisplay;
    private TextView mSublineDisplay;
    private SubMenu mSquadSubMenu;
    private boolean mHasCreatedSquadSubmenu;
    private MenuItem mShareButton;
    private MenuItem mMessageButton;
    private MenuItem mRefreshButton;
    private ImageButton mAnotherDrinkButton;
    private TextView mDrinkText;

    //data model
    private Session mSession;
    private ArrayList<Squad> mSquads;
    private Interfaces.MenuControllerInterface mMenuController;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        initialiseUI(); // Load member variables
        loadPaneLoadingPage(); //Swap in the fragment with a loading progress / message

        //download user data for this session - will call onDataInitDone() when complete
        new InitialDownloadTask().execute();
    }

    /**
     * Called to set up UI state and save UI components to member variables.
     */
    private void initialiseUI() {
        mToolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(mToolbar);

        mFloatingActionButton = (FloatingActionButton) findViewById(R.id.fab);

        mDrawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        mDrawerToggle = new ActionBarDrawerToggle(
                this, mDrawer, mToolbar, R.string.navigation_drawer_open, R.string.navigation_drawer_close);

        mDrawer.setDrawerListener(mDrawerToggle);
        mDrawerToggle.syncState();

        mNavigationView = (NavigationView) findViewById(R.id.nav_view);
        mNavigationView.setNavigationItemSelectedListener(this);

        mNavNameDisplay = (TextView) findViewById(R.id.mainActNameText);
        mSublineDisplay = (TextView) findViewById(R.id.mainActSubLine);

        mAnotherDrinkButton = (ImageButton) findViewById(R.id.mainActDrinkButton);
        mAnotherDrinkButton.setOnClickListener(this);
        mDrinkText = (TextView) findViewById(R.id.mainActDrinkCountText);

        updateDrinkCount();
    }

    /**
     * Sets up the UI with user data. Called at the end of the session data
     * download. Also starts the updater service.
     */
    private void onDataInitDone() {
        mNavNameDisplay.setText(mSession.getName());
        mSublineDisplay.setText(mSession.getUsername());

        //put the menu items for the squad on the navigationView
        for(int i = 0; i < mSquads.size(); i++) {
            addSquadMenu( mSquads.get(i).getName(), i);
        }

        if (mSquads.size() == 0) {
            loadPaneFirstStart();
        } else {
            loadPaneFeedPage();
        }

        Intent serviceIntent = new Intent(MainActivity.this, UpdaterService.class);
        serviceIntent.putExtra(Constants.KEYS.SESSIONKEY, mSession.getKey());
        startService(serviceIntent);
    }

    //called to hide all action bar menus - done before a new fragment is loaded.
    private void disableMenus() {
        mMenuController = null;
        if(mShareButton != null)mShareButton.setVisible(false);
        if(mMessageButton != null)mMessageButton.setVisible(false);
        if(mRefreshButton != null)mRefreshButton.setVisible(false);
    }

    //swaps in a fragment to the main view which shows a loading message / indeterminate progress.
    private void loadPaneLoadingPage() {
        disableMenus();
        LoadingFragment pane = new LoadingFragment();
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }

    //swaps in the fragment which shows the first-start help.
    private void loadPaneFirstStart() {
        disableMenus();
        FirstStartFragment pane = new FirstStartFragment();
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }


    //swaps in a fragment to the main view which shows a menu for creating / joining squads.
    private void loadPaneJoinCreateSquadPage() {
        disableMenus();
        JoinCreateSquadFragment pane = new JoinCreateSquadFragment();
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }

    //swaps in a fragment to the main view which shows the controls / members for a specific squad.
    private void loadPaneSquadPage(int squadIndex) {
        disableMenus();
        SquadDisplayFragment pane = new SquadDisplayFragment(mSquads.get(squadIndex), squadIndex);
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }

    //given an individuals data, loads and initialises the tracker fragment.
    public void loadPaneTrackerPage(SquadMember individual) {
        disableMenus();
        TrackerFragment pane = new TrackerFragment(individual);
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }

    //loads and initialises the Feed fragment.
    public void loadPaneFeedPage() {
        disableMenus();
        FeedFragment pane = new FeedFragment();
        getSupportFragmentManager().beginTransaction()
                .replace(R.id.pane_container, pane).commit();
    }

    //deletes the current sessionkey and finishes the activity.
    private void logout() {
        SharedPreferences sharedPref = getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
        SharedPreferences.Editor editor = sharedPref.edit();
        editor.putString(Constants.KEYS.SESSIONKEY, "");
        editor.commit();

        finish();
    }

    /**
     * Called by a loaded Fragment to set the visibility of a ActionBar Menu, and set the callback
     * associated with their clicks.
     * @param id ID of the menu item.
     * @param callback Object which implements MenuControllerInterface
     */
    public void initialiseMenu(int id, Interfaces.MenuControllerInterface callback) {
        mMenuController = callback;
        switch (id) {
            case R.id.mainMenuMessage:
                mMessageButton.setVisible(true);
                break;
            case R.id.mainMenuShare:
                mShareButton.setVisible(true);
                break;
            case R.id.mainMenuRefresh:
                mRefreshButton.setVisible(true);
        }
    }

    /**
     * Called to get the current session (cached data model)
     * @return Session
     */
    public Session getSession() {
        return mSession;
    }

    /**
     * Called to set the text of the title of the app bar.
     * @param title
     */
    public void setBarTitle(String title) {
        setTitle(title);
    }

    @Override
    public void onBackPressed() {
        DrawerLayout drawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        if (drawer.isDrawerOpen(GravityCompat.START)) {
            drawer.closeDrawer(GravityCompat.START);
        } else {
            super.onBackPressed();
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.main, menu);

        // menu items now exist - find their references and store them. (UI model)
        mShareButton = menu.findItem(R.id.mainMenuShare);
        mMessageButton = menu.findItem(R.id.mainMenuMessage);
        mRefreshButton = menu.findItem(R.id.mainMenuRefresh);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        if (mMenuController != null)mMenuController.menuClicked(id);
        return super.onOptionsItemSelected(item);
    }

    /**
     * Called on the selection of a navigation item. First processes known ID's (ie: non
     * dynamic) before resolving the dynamic ones (squad menu items).
     * @param item
     * @return
     */
    @SuppressWarnings("StatementWithEmptyBody")
    @Override
    public boolean onNavigationItemSelected(MenuItem item) {
        // Handle navigation view item clicks here.
        int id = item.getItemId();

        if (id == R.id.nav_logout) {
            logout();
        } else if (id == R.id.nav_newsquad) {
            loadPaneJoinCreateSquadPage();
        } else if (id == R.id.nav_summary) {
            loadPaneFeedPage();
        } else if (id == R.id.nav_reset_drinks) {
            DrinkCounter.reset(this);
            updateDrinkCount();
        } else { //squad menuitem selected
            int index = item.getItemId();

            if( index < mSquads.size()) { // sanity check
                Log.d(Constants.MAD, "Squad menuItem selected: " + index + " (" + mSquads.get(index).getName() + ")");
                loadPaneSquadPage(index);
            }
        }

        // close the drawer
        DrawerLayout drawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        drawer.closeDrawer(GravityCompat.START);
        return true;
    }

    //Called with messages from fragments.
    public void onFragmentInteraction(Uri uri) {
        //not used
    }

    /**
     * Called from JoinCreateSquadFragment fragment when a new squad is joined/created.
     * Addes a new squad to the cached data model.
     * @param squad
     */
    public void newSquadNotify(Squad squad) {
        Log.d(Constants.MAD, "New squad notified: " + squad.getName());
        mSquads.add(squad);
        addSquadMenu(squad.getName(), mSquads.size()-1);
    }

    public MenuItem getRefreshMenuButton() {
        return mRefreshButton;
    }

    /**
     * Called to create a new squad menuItem in the NavigationView
     * @param squadName
     * @param index
     */
    private void addSquadMenu(String squadName, int index) {
        Menu squadNavMenu = (Menu)mNavigationView.getMenu();

        if (!mHasCreatedSquadSubmenu) {
            mSquadSubMenu = squadNavMenu.addSubMenu("Squads");
            mHasCreatedSquadSubmenu = true;
        }

        MenuItem newSquadRow = mSquadSubMenu.add(1, index, 1, squadName);
        newSquadRow.setIcon(R.drawable.ic_person);

        // bug in support library  ticket #176300
        // means hacky solution required to force update
        Menu m = mNavigationView.getMenu();
        MenuItem mi = m.getItem(m.size() - 1);
        mi.setTitle(mi.getTitle());
    }

    // updates the UI with the current drink count
    private void updateDrinkCount(){
        mDrinkText.setText(DrinkCounter.getCount(this) + " drinks");
    }

    // Implements View.OnClickListener, recieves all button clicks and routes them
    @Override
    public void onClick(View v) {
        if (v.getId() == R.id.mainActDrinkButton) {
            DrinkCounter.anotherDrink(this, mSession, mSquads);
            updateDrinkCount();
        }
    }

    @Override
    protected void onDestroy() {
        Intent myIntent = new Intent(MainActivity.this, UpdaterService.class);
        stopService(myIntent);
        super.onDestroy();
    }


    /**
     * Encapsulates the network calls for downloading data on activity startup.
     * Calls onDataInitDone() to finish UI init with new data.
     */
    private class InitialDownloadTask extends AsyncTask<Void, Void, Void> {

        @Override
        protected Void doInBackground(Void... dummy) {
            SharedPreferences sharedPref = getSharedPreferences(Constants.RES.PREFERENCES_FILE, Context.MODE_PRIVATE);
            mSession = SessionFactory.makeFromKey(sharedPref.getString(Constants.KEYS.SESSIONKEY, null));
            mSquads = SquadFactory.getSquads(mSession);

            return null;
        }

        @Override
        protected void onPostExecute(Void dummy) {
            onDataInitDone(); //inits the UI with model
        }
    }
}
