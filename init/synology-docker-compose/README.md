# Synology Docker Compose Example

This folder is an example `docker-compose.yml` file and a `.env` file for using docker through a Synology Server.

The `.env` file is a special environment/variable file that `docker-compose` automatically uses to make configuration of your containers easier.

## Directions

### Directories

First, create the directories where you want the containers to save everything. I use the same base directory for them, like: `/volume1/Docker/unifi-poller`, then inside the `./unifi-poller` directory: `./grafana` and `./influxdb`.

  NOTE: Its best to use the command line over SSH to create these directories, AFTER you have the root shared directory created. (`Control Panel -> Shared Folder -> Create`)

```bash
sudo mkdir -p /volume[#]/[Shared Directory]/unifi-poller/grafana
sudo mkdir -p /volume[#]/[Shared Directory]/unifi-poller/influxdbx
```

Where `/volume[#]` is the volume number corresponding to your volumes in Synology;
`[Shared Directory]` is the shared directory from above, and then
creating the `grafana and influxdb` directories.

### User Accounts

You still have to [do this prep work](https://github.com/unifi-poller/unifi-poller/wiki/Synology-HOWTO#method-2), creating the `unifi-poller` user, which I'll re-iterate here:

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

### Spin Up the Containers

At this point, you are able to run `sudo docker-compose up -d` from within the directory that you have the `docker-compose.yml` file and the `.env` file saved on your Synology.

And now, using your Synology's Web GUI, we have to create the Influx Database.

### Create Influx Database

#. Click `Containers` and then double click the running `influxdb1` container
#. Switch to the `terminal` tab
#. Click the drop down next to `Create` and select `launch with command`
#. Enter `bash` and click `ok`
#. Select `bash` from the left hand side. You should now see an active `command prompt`
#. In the command prompt, enter these commands: (note: pasting IS possible! You have to right click in the terminal window and select `paste`)

`influx`

After a couple of seconds you should be in the InfluxDB shell.
Run the following commands in the InfluxDB shell, then close the window. The `unifipoller` username is the read-only user you created in the Unifi Web Management page.

```
CREATE DATABASE unifi
USE unifi
CREATE USER unifipoller WITH PASSWORD 'unifipoller' WITH ALL PRIVILEGES
GRANT ALL ON unifi TO unifipoller
exit
```

### Grafana Login and Final Setup

From here, your three containers should show as running, no alerts of auto-restarting or other issues coming from your synology's web GUI.

Make sure to double check the log output of all three containers. If there are any issues, they **should** appear here.

Then, in your browser, go to `http://{ip address of your synology}:3000`. The default login is `admin:admin`. But, if you used the `INFLUXDB_ADMIN_USER` and `INFLUXDB_ADMIN_PASSWORD` in the variable file, use that login here. 

#. Click `Add Your First Data Source` on the home page
#. Select the `influxdb` source option
#. Set the following fields:
  - Name = `UniFi InfluxDB` (or whatever name you want) and set to default
    URL = http://influxdb1:8086
    Database = unifi
    Username = unifipoller
    Password = unifipoller

No other fields need to be changed or set on this page.
Click Save & Test

    You should get green banner above the save and test that says 'Data Source is Working'
    To return to the homepage click the icon with 4 squares on the left nav-bar and select home



## Variables

For all of the variables used in the docker-compose file, you'll find them in the `docker-compose.example.env` file. Please, `cp docker-compose.example.env .env` and open it in your favorite text editor. `nano .env`

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
