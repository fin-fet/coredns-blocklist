# coredns-blocklist
A plugin to block domains based on a given list. Also logs a warning.

## Usage
This plugin can load blockfiles from disk by passing a path.

```
. {
    log
    blocklist blocklist.txt
    forward . 1.1.1.1
}
```

### Options
**reload**  
Reload specifies how often the plugin will load new blockfile data. A value of 0 will disable
automatic reloading. Must be a valid Go duration (ex: "5s" or "10h3m").

ex:
```
. {
    blocklist blocklist.txt {
        reload 30m
    }
}
```

**response**  
Response dictates how the plugin will respond to blocked queries. Default is `nxdomain`.

Options are:
* Standard responses
    * nxdomain
    * refused
* Extended (ENDS0 EDE) responses
    * other
    * blocked
    * censored
    * filtered
    * prohibited

EDE options can also have optional text, for example:
```
. {
    blocklist malware-domains.txt {
        repsonse blocked "Known malware domain"
    }
}
```
