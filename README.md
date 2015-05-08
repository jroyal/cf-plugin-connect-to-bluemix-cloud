# cf-plugin-connect-to-bluemix-cloud
Assuming a user is logged into CF and is pushing to an app inside of a Bluemix Org and Space, this plugin will return a string of exports that you can use to enable the use of openstack cli tools.

##Usage
```
$ cf connect-to-bluemix-cloud
```

##Installation

#####Install with binary
- Download the binary [`win64`](https://github.com/jroyal/cf-plugin-connect-to-bluemix-cloud/blob/master/bin/win64/connect-to-bluemix-cloud.exe) [`linux64`](https://github.com/jroyal/cf-plugin-connect-to-bluemix-cloud/blob/master/bin/linux64/connect-to-bluemix-cloud) 
- Install plugin `$ cf install-plugin <binary_name>`



#####Install from Source (need to have [Go](http://golang.org/dl/) installed)
  ```
  $ go get github.com/cloudfoundry/cli
  $ go get github.com/jroyal/cf-plugin-connect-to-bluemix-cloud
  $ cd $GOPATH/src/github.com/jroyal/cf-plugin-connect-to-bluemix-cloud
  $ go build connect-to-bluemix-cloud
  $ cf install-plugin connect-to-bluemix-cloud
  ```
