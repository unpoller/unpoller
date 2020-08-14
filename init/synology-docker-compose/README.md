# Synology Docker Compose Example

This folder is an example `docker-compose.yml` file and a `.env` file for using docker through a Synology Server.

The `.env` file is a special environment/variable file that `docker-compose` automatically uses to make configuration of your containers easier.

## Directions

First, create the directories where you want the containers to save everything. I use the same base directory for them, like: `/volume1/Docker/unifi-poller`, then inside the `./unifi-poller` directory: `./grafana` and `./influxdb`.

  NOTE: Its best to use the command line over SSH to create these directories, AFTER you have the root shared directory created. (`Control Panel -> Shared Folder -> Create`)

```bash
sudo mkdir -p /volume[#]/[Shared Directory]/unifi-poller/grafana
sudo mkdir -p /volume[#]/[Shared Directory/unifi-poller/influxdbx
```

Where `/volume[#]` is the volume number corresponding to your volumes in Synology;
`[Shared Directory]` is the shared directory from above, and then
create the `grafana and influxdb` directories.

You still have to [do this prep work](https://github.com/unifi-poller/unifi-poller/wiki/Synology-HOWTO#method-2) creating the `unifi-poller` user, which I'll re-iterate here:

#. Create a new user account on the Synology from the Control Panel:
    - Name the user `unifi-poller`
    - Set the password (you don't need to logon as unifipoller and change it)
    - `Disallow Password Change`
    - Assign them to the user group `users`
    - Give them `r/w` permission to the folder you created e.g. `/docker/unifi-poller`
    - Don't assign them **anything** else - the point of this user is for security's sake.
#. SSH into your Synology
#. Run the following command to find the PID of the user you created and set the variable `GRAFANA_LOCAL_USERID` in your `.env` file:
    - `sudo id unifi-poller`
    - `GRAFANA_LOCAL_USERID=1026`

##Variables

For all of the variables used in the docker-compose file, make sure to not only read over but make important edits to the `docker-compose.example.env file. As well as making a copy and naming it `.env`

For the `/local/storage/location/` lines, change those to match your directories.

```bash
#influxdb
INFLUXDB_ADMIN_PASSWORD=changeme
INFLUXDB_LOCAL_VOLUME=/local/storage/location/influxdb

#grafana
GRAFANA_USERNAME=unifi-poller_username
GRAFANA_PASSWORD=changeme
GRAFANA_LOCAL_USERID=1026

UNIFI_PASS=set_this_on_your_controller
UNIFI_URL=https://127.0.0.1:8443
```
