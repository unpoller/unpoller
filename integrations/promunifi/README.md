# Unifi

Collect your Unifi client data every 15 seconds and send it to an InfluxDB instance.

![image](https://cloud.githubusercontent.com/assets/79995/19002122/6b81f928-86ff-11e6-8ab4-d67f943588f4.png)

## Deploying

The repository is ready for deployment on Heroku. Steps to deploy:

Clone the repository and using `.env.example` create your own `.env` file with your Unifi GUI and InfluxDB credentials.

Create your heroku application:

```
heroku create [name]
```

Set your environment variables before deploying:

```
heroku config:set $(cat .env | grep -v ^# | xargs)
```

Push to heroku:

```
git push heroku master
```

## Copyright
Copyright © 2016 Garrett Bjerkhoel. See [MIT-LICENSE](http://github.com/dewski/unifi/blob/master/MIT-LICENSE) for details.
