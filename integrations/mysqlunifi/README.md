# MYSQL Output Plugin Example

This plugin is not finished and did not get finished for the release of poller v2.
Sorry about that. I'll try to get it working soon! 2/4/20

The code here, and the dynamic plugin provided shows an example of how you can
write your own output for unifi-poller. This plugin records some very basic
data about clients on a unifi network into a mysql database.

You could write outputs that do... anything. An example: They could compare current
connected clients to a previous list (in a db, or stored in memory), and send a
notification if it changes. The possibilities are endless.

You must compile your plugin using the unifi-poller source for the version you're
using. In other words, to build a plugin for version 2.0.1, do this:
```
mkdir -p $GOPATH/src/github.com/unifi-poller
cd $GOPATH/src/github.com/unifi-poller

git clone git@github.com:unifi-poller/unifi-poller.git
cd unifi-poller

git checkout v2.0.1

cp -r <your plugin> plugins/
GOOS=linux make plugins
```
The plugin you copy in *must* have a `main.go` file for `make plugins` to build it.
