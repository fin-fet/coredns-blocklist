# coredns-blocklist
A plugin to block domains based on a given list. Also logs a warning.

## Usage
This plugin can load blockfiles from disk by passing a path. 
It can also load from an HTTP URL. Multiple blocklists are allowed per server.
If the domain is not in the blocklist, it will go to the next plugin.

```
. {
    log
    blocklist blocklist.txt
    blocklist https://raw.githubusercontent.com/blocklistproject/Lists/master/scam.txt
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

**match_subdomains**  
Whether or not to match subdomains. For example if the record is `mango.wah`, if set to true
`sub.mango.wah` would match, but if set to false it would not. Must be true or false, defaults to true.