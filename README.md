Henk is a very simple reverse HTTP proxy leveraging OpenSSH's built-in reverse
proxy features for authentication and secure tunneling and Let's Encrypt for
adding HTTPS.

WARNING: Henk has not been reviewed for security, use at your own risk.

### Usage

Once henk is set up, run the HTTP service you want to expose locally, for
example on port 8080. Then use SSH to set up a reverse proxy:

    ssh henk@tunnel.host.net -NR /run/henk/foobar:localhost:8080

OpenSSH will create a proxy file on your server at `/run/henk/foobar` and
forward its traffic to your localhost at port 8080. Henk will accept connections
on port 80 and 443 for foobar.tunnel.host.net, obtain a certificate if needed
and proxy traffic to the proxy file. You can then access your service publicly:

    curl https://foobar.tunnel.host.net

You can also tell HTTP backends running on the same host to listen on a socket
in `/run/henk` to instantly make them reachable over HTTPS.

### Setup

- Build henk with `go build`.
- Place the binary at `/usr/bin/henk`.
- Create a henk user and group.
- Place the service file at `/etc/systemd/system/henk.service`.
- Edit it to set the desired base domain.
- Set DNS to forward all subdomains for that base domain to your server.
- Set sshd to disallow login as henk. Reverse proxy options should still work.
- Authorize desired SSH keys for user henk.
- Set sshd to clean up old sockets before creating new ones:

      echo 'StreamLocalBindUnlink yes' >> /etc/ssh/sshd_config

### Security model

Henk relies on the security of OpenSSH to protect the data forwarded to and from
your local HTTP service. It also relies on the authentication of OpenSSH to
specify who is allowed to use henk.

Henk relies on the Go autocert library to obtain Let's Encrypt certificates. It
only obtains certificates for domains whose corresponding socket exists.

Henk has about a 100 lines of Go glue code which is a fairly small attack
surface, but it has not yet been thoroughly reviewed.

Anyone with one of henk's authorized keys can use henk, host anything on any
subdomain and steal subdomain proxies from others. This is a limitation of
OpenSSH's reverse proxy feature. Make sure you only give access to people you
trust and that trust each other.

### FAQ

#### Why do I need to connect as the same henk user?

OpenSSH creates the reverse proxy socket with the connecting user as owner and
no permissions for group or other. If you connect as another user, the henk
daemon won't be able to access the socket.

#### Why does OpenSSH fail to set up the reverse proxy?

Unfortunately, the error messages from OpenSSH aren't very clear. It could be
one of these things:

- You're using SSH connection sharing or `ControlMaster` options: The proxy can
  only be set up by the master connection. Consider just disabling that option
  for this host.
- OpenSSH does not have permission to create the socket: Check the ownership and
  permission bits of `/run/henk` and make sure you connect to the same user that
  runs the henk daemon.
- The path you specified for the socket still exists: Make sure to add
  `StreamLocalBindUnlink yes` to `/etc/ssh/sshd_config` and restart SSHD. Check
  that the henk user has permission to remove the old socket or remove it
  manually.
