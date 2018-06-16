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

    deb_builder = packager.DebPackage('chese-server',
                                      maintainer='Twitchyliquid64 <twitchyliquid64@ciphersink.net>',
                                      description='Chese uses Chrome to make a disposable browsing session for risky security research.',
                                      bin_files={'../chese-server': 'chese-server',
                                                 '../chese-setup': 'chese-setup'},
                                      data_files={'../resources': 'resources'},
                                      config_file_name='server.json',
                                      configuration_dir='etc/chese',
                                      postinst='deb_resources/postinst',
                                      config_data={
                                            "listener": ":8425",
                                            "tls": {
                                                "server-cert-path": "/etc/chese/serv-cert.pem",
                                                "server-key-path": "/etc/chese/serv-key.pem",
                                                "ca-cert-path": "/etc/chese/ca-cert.pem"
                                            },
                                            "resource-path": "/usr/share/chese-server/resources"
                                        })

    print deb_builder.package(version, config_path)
