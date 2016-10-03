package net.ciphersink.nightout;

/**
 * Defines custom interfaces used throughout the application
 */
public class Interfaces {
    /**
     * Fragments that implement this interface may recieve notifications about menu icons pressed
     * in MainActivity
     */
    public interface MenuControllerInterface
    {
        public void menuClicked(int id);
    }
}
