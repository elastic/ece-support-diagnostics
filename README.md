# ece-support-diagnostics

## Description

The support diagnostic utility is a bash script that you can use to gather ECE logs, metrics and information directly on the the host where the ECE instance is running; the resulting archive file can be then provided to Elastic support for troubleshooting and investigation purposes.

## How to use

* download the [latest version of the script](https://github.com/elastic/ece-support-diagnostics/raw/master/diagnostics.sh) and copye to the ECE host. (Alternatively, you can just run this in the command line of the host: `wget https://github.com/elastic/ece-support-diagnostics/raw/master/diagnostics.sh`)
* Give execution permissions to the file (`chmod +x diagnostics.sh`)
* run as ECE installation owner.
* using options that make use of REST calls ( -a, -c ) will require ECE user credentials (-u readonly -p \<password\>)
* note `curl` is required when using REST related calls ( -a, -c options )
* repeat for each ECE host relevant to the issue and all hosts with the director role 

Comparing the state of a broken node with the state of the directors is often necessary to pinpoint where the root cause is and fixing the root cause will often allow other problems to self heal.


## Sample execution

```
$ ./diagnostics.sh 
ECE Diagnostics
Usage: ./diagnostics.sh [OPTIONS]

Options:
-e|--ecehost #Specifies ip/hostname of the ECE (default:localhost)
-y|--protocol <http/https> #Specifies use of http/https (default:http)
-x|--port <port> #Specifies ECE port (default:12400)
-s|--system #collects elastic logs and system information
-d|--docker #collects docker information
-zk|--zookeeper <path_to_dest_pgp_public_key> #enables ZK contents dump, requires a public PGP key to cipher the contents
-zk-path|--zookeeper-path <zk_path_to_include> #changes the path of the ZK sub-tree to dump (default: /)
-zk-excluded|--zookeeper-excluded <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle
--zookeeper-excluded-insecure <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle WARNING: This options remove default filters aimed to avoid secrets and sensitive information leaks
-sp|--storage-path #overrides storage path (default:/mnt/data/elastic). Works in conjunction with -s|--system
-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)
-c|--cluster <clusterID> #collects cluster plan and info for a given cluster (user/pass required). Also restricts -d|--docker action to a specific cluster
-a|--allocators #gathers allocators information (user/pass required)
-u|--username <username>
-p|--password <password>



Sample usage:
"./diagnostics.sh -d -s" #collects system and docker level info
"./diagnostics.sh -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects allocators information
"./diagnostics.sh -e 192.168.1.42 -x 12409 -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects allocators information using custom host and port
"./diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -d -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects cluster plan,info and docker info only for the specified cluster ID
"./diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects cluster plan,info for the specified cluster ID

Tue Sep  5 13:16:56 CEST 2017 [INFO]:  ECE Diagnostics 
Tue Sep  5 13:16:57 CEST 2017 [INFO]:  Nothing to do.
```

## What flags to use?

### Basic
The standard basic set of information (system and docker level) can be gathered with:

```
./diagnostics.sh -d -s
```

### Including Zookeeper contents for deep analysis
Some investigations require to have low level ECE details (stored in Zookeeper) at hand. It is possible to include this information in the bundle using:

```
-zk|--zookeeper <path_to_dest_pgp_public_key>
```

This option will enable the inclusion of a complete dump of the cluster's ZK ensemble contents into the generated bundle.

This behavior can be constrained so the bundle:

- Would include just the contents of a concrete ZK sub-tree, e.g: (just the contents under `/kibanas`) 
```bash
./diagnostics.sh -zk ./support.key.pub -zk-path '/kibanas'
```
- Exclude certain paths:
```bash
./diagnostics.sh -zk ./support.key.pub -zk-excluded '/zookeeper,/locks'
```
- Or both:
```bash
./diagnostics.sh -zk ./support.key.pub -zk-excluded '/container_sets/cloud-uis,/container_sets/zookeeper-servers' -zk-path '/container_sets'
```

**Note**: How the list of excluded trees is a comma separated list of ZK paths.

:warning: Eagle-eyed readers have probably noticed that this functionality **requires a PGP public key**. ECE Zookeeper contents contain potentially secret and/or sensitive information. For this reason, the bundle **will never contain ZK contents in clear text**. The output will contain a PGP encrypted file which can only be read by the owner of the private key counterpart for the provided key.

:warning: By default, paths known to contain secrets or sensitive information have been excluded from the bundle. This protection can be deactivated, for those cases on which it is absolutely necessary, using the option `--zookeeper-excluded-insecure`. Please use this option with extreme caution and only when you are willing to show your passwords and certificates to the owner of the public PGP key.

### Using a custom storage path
If you've installed ECE using a STORAGE_PATH different than default (`/mnt/data/elastic`), please make sure to pass the below flag to the diagnostics script:

```
./diagnostics.sh -d -s -sp /my/custom/storage/path
```


## Output
Diagnostic output archive will be written to /tmp folder with file name ece_diag-<ECE_host_IP>-<Timestamp>.tar.gz  
Once you have the file please provide it to your designated support agent, by attaching it to the support case.

