Henk is a very simple reverse HTTP proxy leveraging openSSH's built-in reverse
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

### Other uses

Henk is only about a hundred lines of Go so it should be easy to adapt for other
uses. One thing you can try is to host more permanent websites by having their
backends create a socket in `/run/henk`.
