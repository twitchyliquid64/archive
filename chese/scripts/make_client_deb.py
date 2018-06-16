#!/usr/bin/python
#This script should be run from inside the packaging/deb folder.
#./make_deb.py <version> [<path-to-config>]
import sys
from packager import packager

if __name__ == '__main__':
    version = '0.0.1'
    version = sys.argv[1]
    config_path = None
    if len(sys.argv) > 2:
        config_path = sys.argv[2]

    deb_builder = packager.DebPackage('chese-client',
                                      maintainer='Twitchyliquid64 <twitchyliquid64@ciphersink.net>',
                                      description='Chese uses Chrome to make a disposable browsing session for risky security research.',
                                      desktop_file='chese.desktop',
                                      desktop_file_path='deb_resources/chese.desktop',
                                      bin_files={'../chese-client': 'chese-client',
                                                 '../chese-setup': 'chese-setup'},
                                      data_files={'../resources': 'resources'},
                                      config_file_name='client.json',
                                      configuration_dir='etc/chese',
                                      config_data={
                                        "chrome": {
                                          "bin-path": "/usr/bin/google-chrome-stable",
                                          "data-directory": "/tmp/chese"
                                        },
                                        "keyhole": {
                                          "listener": ":8427",
                                          "forwarders": [
                                            {
                                              "type": "tls-pinned",
                                              "address": "localhost:8425",
                                              "ca-cert-path": "/etc/chese/ca-cert.pem",
                                              "client-cert-path": "/etc/chese/client-cert.pem",
                                              "client-key-path": "/etc/chese/client-key.pem"
                                            }
                                          ]
                                        },
                                        "resource-path": "/usr/share/chese-client/resources"
                                      })

    print deb_builder.package(version, config_path)
