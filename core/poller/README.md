# Unifi

Collect your Unifi Controller Client data and send it to an InfluxDB instance.

![image](https://raw.githubusercontent.com/davidnewhall/unifi/master/grafana-unifi-dashboard.png)

## Deploying


Clone the repository and using `.env.example` create your own `.env` file with your Unifi GUI and InfluxDB credentials.


Set your environment variables before running:

```
source .env ; ./unifi-poller
```

## Copyright & License
Copyright © 2016 Garrett Bjerkhoel. See [MIT-LICENSE](http://github.com/dewski/unifi/blob/master/MIT-LICENSE) for details.
