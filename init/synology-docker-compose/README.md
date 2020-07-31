# Synology Docker Compose Example

This folder is an example `docker-compose.yml` file and a `.env` file for using docker through a Synology Server.

## Directions

First, create the directories where you want the containers to save everything. I use the same base directory for them, like: `/volume1/Docker/grafana`, then inside the `./grafana` directory: `./grafana` and `./influxdb`. Its best to use the command line over SSH to create these directories, AFTER you have the primary shared directory already created. (`Control Panel -> Shared Folder -> Create`)

`sudo mkdir -p /volume[#]/[Shared Directory]/Grafana/{grafana,influxdb}`

Where `/volume[#]` is the volume number,
`[Shared Directory]` is the shared directory from above, and then
`{grafana,influxdb}` needs to be copied as-is, curly brackets and all.

You still have to [do this prep work](https://github.com/unifi-poller/unifi-poller/wiki/Synology-HOWTO#method-2) creating the `unifipoller` user, which I'll re-iterate here:

#. Create a new user account on the Synology from the Control Panel:
    - Name the user `grafana`
    - Set the password (you don't need to logon as grafana and change it)
    - `Disallow Password Change`
    - Assign them to the user group `users`
    - Give them `r/w` permission to the folder you created e.g. `/docker/grafana`
    - Don't assign them **anything** else - the point of this user is for security's sake.
#. SSH into your Synology
#. Run the following command to find the PID of the user you created and set the variable `GRAFANA_LOCAL_USERID` in your `.env` file:
    - `sudo id grafana`
    - `GRAFANA_LOCAL_USERID=1026`
