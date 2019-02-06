# ece-support-diagnostics

## Description

The support diagnostic utility is a bash script that you can use to gather ECE logs, metrics and information directly on the the host where the ECE instance is running; the resulting archive file can be then provided to Elastic support for troubleshooting and investigation purposes.

## How to use

* download the [latest release](https://github.com/elastic/ece-support-diagnostics/releases/latest)
* copy to ECE host
* run as ECE installation owner.
* using options that make use of REST calls ( -a, -c ) will require ECE user credentials (-u readonly -p \<password\>)
* note `curl` is required when using REST related calls ( -a, -c options )
* repeat for each ECE host relevant to the issue and all hosts with the director role

Comparing the state of a broken node with the state of the directors is often necessary to pinpoint where the root cause is and fixing the root cause will often allow other problems to self heal.


## Sample execution

```
$ ./ece-diagnostics.sh 
ECE Diagnostics
Usage: ./ece-diagnostics.sh [OPTIONS]

Options:
-e|--ecehost #Specifies ip/hostname of the ECE (default:localhost)
-y|--protocol <http/https> #Specifies use of http/https (default:http)
-x|--port <port> #Specifies ECE port (default:12400)
-s|--system #collects elastic logs and system information
-d|--docker #collects docker information
-sp|--storage-path #overrides storage path (default:/mnt/data/elastic). Works in conjunction with -s|--system
-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)
-c|--cluster <clusterID> #collects cluster plan and info for a given cluster (user/pass required). Also restricts -d|--docker action to a specific cluster
-a|--allocators #gathers allocators information (user/pass required)
-u|--username <username>
-p|--password <password>



Sample usage:
"./ece-diagnostics.sh -d -s" #collects system and docker level info
"./ece-diagnostics.sh -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects allocators information
"./ece-diagnostics.sh -e 192.168.1.42 -x 12409 -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects allocators information using custom host and port
"./ece-diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -d -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects cluster plan,info and docker info only for the specified cluster ID
"./ece-diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects cluster plan,info for the specified cluster ID

Tue Sep  5 13:16:56 CEST 2017 [INFO]:  ECE Diagnostics
Tue Sep  5 13:16:57 CEST 2017 [INFO]:  Nothing to do.
```

## What flags to use?

### Basic
The standard basic set of information (system and docker level) can be gathered with:

```
./ece-diagnostics.sh -d -s
```

### Using a custom storage path
If you've installed ECE using a STORAGE_PATH different than default (`/mnt/data/elastic`), please make sure to pass the below flag to the diagnostics script:

```
./ece-diagnostics.sh -d -s -sp /my/custom/storage/path
```


## Output
Diagnostic output archive will be written to /tmp folder with file name ece_diag-<ECE_host_IP>-<Timestamp>.tar.gz  
Once you have the file please provide it to your designated support agent, by attaching it to the support case.
