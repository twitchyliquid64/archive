# <img src="https://github.com/twitchyliquid64/chese/raw/master/resources/keyhole/static/cheese_pizza_icon_by_typhloser.gif" width="48"> chese

*WORK IN PROGRESS - DO NOT USE*

Chese (not a mis-spell) uses Chrome to make a disposable browsing session for risky security research.

 * All traffic is proxied off-machine via an encrypted tunnel to prevent eavesdropping by an active or malicious attacker on your LAN, or association by remote sites to your current IP address.
   * Proxy tunnel is mutually authenticated with TLS certificates.
 * Chrome runs in a sandbox to minimise the risk of a browser exploit affecting your machine. (TODO)
 * Per-invocation user data directories mean a brand-new chrome every time it is invoked. (TODO)

## Installing the client (.deb, Linux)

The deb will:

1. Install the `chese-client` & `chese-setup` binaries into `/usr/bin`.
2. Install the `resources` folder into `/usr/share/chese-client`.
3. Install a sensible default client config into `/usr/chese/client.json`.

First, copy the CA-cert, your client-cert, and your client-key PEM files onto your local machine. A good location for them is `~/.config/chese` - as the setup utility will install
your client configuration at `~/.config/chese/client.json`.

Now, run the setup utility:

1. Run `chese-setup`.
2. Select 'Edit Client Configuration'.
3. Modify the values to represent your environment. Make sure to use absolute paths everywhere.
4. Press Done.

You should now be able to run `chese-client` and have it connect.

### credits

Thanks to typhloser for the [awesome cheese graphic](https://typhloser.deviantart.com/art/Free-to-Use-Cheese-Pizza-Icon-387660529)!
