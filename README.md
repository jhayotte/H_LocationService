# Microservice H_LocationService

This project returns all coordinates of a specific driver during the last N minute.

##Dependency
Location service fetchs its data from a message queue system (**NSQ**) and store them in **REDIS**. So both of them are required to have this project up and running.

##Endpoint
An endpoint allow to retrieve all coordinates of a driver during the last N minutes through: **_/drivers/:id/coordinates?minutes=5_**

##How to

### NSQ
```
$ nsqlookupd & 
$ nsqd --lookupd-tcp-address=127.0.0.1:4160 &
$ nsqadmin --lookupd-http-address=127.0.0.1:4161 &
```
