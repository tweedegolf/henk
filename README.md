### Usage

To run the daemon:

    henk tunnel.host.net

It opens port 80 and 443 and mount /run/henk. DNS should be configured to point
`*.tunnel.host.net` to the server it runs on.

To forward https from foobar.tunnel.host.net to port 1111:

    ssh tunnel.host.net -R /run/henk/foobar:localhost:1111

Henk detects the creation of this unix socket, obtain a certificate for
foobar.tunnel.host.net and redirect traffic from https://foobar.tunnel.host.net
to it. You can now use the public endpoint, for example:

    curl https://foobar.tunnel.host.net

Note that the unix socket is not cleaned up by default and this command fails
the second time it is run. Configure the tunnel host to clean up old sockets
before creating new ones and restart sshd with:

    echo 'StreamLocalBindUnlink yes' >> /etc/ssh/sshd_config
    systemctl restart sshd

- [ ] Sockets older than a day are deleted.
- [ ] Henk stops accepting connections for a domain when its socket is deleted.
- [ ] Henk caches certificates for the next time a socket is created.
