# ece-support-diagnostics

## ⚠️ Note about built in ECE diagnostics tool

Starting from ECE version 3.3, there is a built-in diagnostics tool shipped with every ECE release.
If you are running ECE 3.3 or later, please follow the steps in the [ECE documentation](https://www.elastic.co/guide/en/cloud-enterprise/current/ece-run-ece-diagnostics.html#ece-run-ece-diagnostics) unless instructed otherwise. 

## Description

The support diagnostic utility is a bash script that you can use to gather ECE logs, metrics and information directly on the the host where the ECE instance is running; the resulting archive file can be then provided to Elastic support for troubleshooting and investigation purposes.

## How to use

* download the [latest release -dist.tar.gz or -dist.zip](https://github.com/elastic/ece-support-diagnostics/releases/latest) - instructions match version `2.x` and higher (do not download source code)
* copy to ECE host and unpack
* run as ECE installation owner.
* using options that make use of REST calls ( `-de`, `-c` ) will require ECE user credentials (`-u admin <-p optional-noprompt-password>`), default APIs will also run. Note `curl` is required when using REST related calls ( -u options )
* repeat for each ECE host relevant to the issue and all hosts with the director role (`-u`, `-c`, `-de` options should only be run once to save space as it queries APIs from coordinator)

Comparing the state of a broken node with the state of the directors is often necessary to pinpoint where the root cause is and fixing the root cause will often allow other problems to self heal.


## Sample execution

```
Usage: ./ece-diagnostics.sh [OPTIONS]

Arguments:
-s|--system #collects elastic logs and system information
-d|--docker #collects docker information
-u|--username <username> - will cause collection of data from ECE APIs (admin user recommended, readonly user will work for v1 APIs and fail for v0 APIs)
-e|--ecehost #Specifies ip/hostname of an ECE coordinator (default:localhost)
-y|--protocol <http/https> #Specifies use of http/https (default:http)
-k|--insecure #Bypass certificate validity checks when using https
-ca|--cacert /path/ca.pem #Specify CA certificate when using https
-x|--port <port> #Specifies ECE port (default:12400)
-zk|--zookeeper <path_to_dest_pgp_public_key> #enables ZK contents dump, requires a public PGP key to cipher the contents
-zk-path|--zookeeper-path <zk_path_to_include> #changes the path of the ZK sub-tree to dump (default: /)
-zk-excluded|--zookeeper-excluded <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle
--zookeeper-excluded-insecure <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle WARNING: This options remove default filters aimed to avoid secrets and sensitive information leaks
--zk-stats|--zookeeper-stats #collects statistics on zookeeper contents and events
-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)

Optional arguments :
-de|--deployment <deploymentID2,deploymentID2> #collects deployment historic plan activity logs (ECE username required), comma separated value allowed. Default to collecting for all unhealthy deployments, pass value "disabled" to not collect any deployment activity logs (requires ECE versions 2.4.0 or higher)
-lh|--log-filter-hours #oldest file to collect in hours (default:72). also applied to docker logs
-p|--password <password> #omiting value or argument will prompt password
-sp|--storage-path Optional - overrides storage path (default:/mnt/data/elastic and auto-detected from runner container inspect if folder does not exist). Works in conjunction with -s|--system
-ds|--disable-sudo #to disable all sudo calls when using option -s|--system

Deprecated argument :
-c|--cluster <clusterID> #collects elasticsearch cluster plan activity logs and restricts docker logs collection - from ECE 2.4.0, please use -de|deployment
-a|--allocator #no action - allocator information is now collected by default

Sample usage:
"./ece-diagnostics.sh -d -s" #collects system and docker level info
"./ece-diagnostics.sh -u admin -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects APIs information
"./ece-diagnostics.sh -e 192.168.1.42 -x 12409 -u admin " #collects API information using custom host and port, prompt for password
"./ece-diagnostics.sh -e 192.168.1.42 -u admin -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" -c e817ac5fbc674aeab132500a263eca71 #collects cluster plan,info for the specified cluster ID
"./ece-diagnostics.sh -de e817ac5fbc674aeab132500a263eca71 -u admin -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj" #collects deployment clusters plan,info for only the specified deployment ID (ECE 2.4+) - when -de is ommited plan activity logs will be collected for all unhealthy deployments
```

## What flags to use?

### Basic
The standard basic set of information (local system logs and docker level, and APIs output from coordinator) can be gathered with:

```
./ece-diagnostics.sh -d -s -u admin -e <IP or Hostname of coordinator>
```

### Including ZooKeeper statistics to help diagnose problems

Most investigations on ZooKeeper stability and availability issues require getting some stats on the ZooKeeper contents distributions (common node name patterns as used in ECE and the number of occurrences of these patterns) and a glimpse on the tail of the transaction log.

It is possible to get this information for the ZooKeeper node running in a director node on which we can execute the ECE diagnostics tool. To do so, adding the following option is enough:

```
--zk-stats|--zookeeper-stats #collects statistics on zookeeper contents and events
```

e.g:

```bash
./ece-diagnostics.sh --zk-stats
```

Which will make the following files to be included in the output bundle:

- The ZooKeeper node path patterns with the number of occurrences listed in CSV format: `elastic/zookeeper_stats/zk_nodetype_stats.csv` e.g:
```
path_type,ephemeral,count,count_ratio,size_min,size_max,size_mean,size_stddev,size_ratio,version_max,version_mean,version_stddev,ctime_min,mtime_max
/blueprint/roles/{hostRole},false,6,0.014150943396226415,241,1004,476.33333333333337,305.38281986167243,0.005014835701840119,0,0.0,0.0,2021-12-13 08:24:55,2021-12-13 08:24:55
/blueprint/roles/{hostRole}/blessings,false,4,0.009433962264150943,46,46,46.0,0.0,3.228585616300146E-4,0,0.0,0.0,2021-12-13 08:24:55,2021-12-13 08:26:09
/blueprint/roles/{hostRole}/pending,true,4,0.009433962264150943,47,47,47.0,0.0,3.298772260132758E-4,2,1.25,0.5,2021-12-13 08:25:09,2021-12-13 08:26:09
/config-store/{configName},false,6,0.014150943396226415,52,98,69.0,15.671630419327785,7.264317636675329E-4,0,0.0,0.0,2021-12-13 08:24:38,2021-12-13 08:29:31
/container_sets/{containerSet},false,14,0.0330188679245283,17,17,17.0,0.0,4.1761053080404065E-4,0,0.0,0.0,2021-12-13 08:24:38,2021-12-13 08:24:39
/container_sets/{containerSet}/containers,false,16,0.03773584905660377,10,10,10.0,0.0,2.8074657533044747E-4,0,0.0,0.0,2021-12-13 08:24:38,2021-12-13 08:24:39
/container_sets/{containerSet}/containers/{container},false,16,0.03773584905660377,1295,3343,2167.6875000000005,606.3443156876023,0.06085708420116194,0,0.0,0.0,2021-12-13 08:24:38,2021-12-13 08:24:39
/container_sets/{containerSet}/containers/{container}/allocation@{host},true,12,0.02830188679245283,14,14,14.0,0.0,2.947839040969699E-4,2,1.1666666666666667,0.38924947208076144,2021-12-13 08:24:56,2021-12-13 08:26:12
/container_sets/{containerSet}/containers/{container}/inspect@{host},true,14,0.0330188679245283,6639,9965,8489.0,996.1174629530395,0.20853504682326476,4,2.0714285714285716,0.9972489631508747,2021-12-13 08:24:56,2021-12-13 08:32:58

```

- A NDJSON representation of the last events in the ZooKeeper translog: `elastic/zookeeper_stats/translog.json`
```json
{"@timestamp":"2021-12-13T08:24:36.442+0000","zxid":4294967298,"type":-10,"type_name":"CREATE_SESSION","path":"","path_type":"null","length":0}
{"@timestamp":"2021-12-13T08:24:36.571+0000","zxid":4294967299,"type":15,"type_name":"CREATE2","path":"/v1","path_type":"","length":0}
{"@timestamp":"2021-12-13T08:24:36.580+0000","zxid":4294967300,"type":15,"type_name":"CREATE2","path":"/v1/secrets","path_type":"/secrets","length":0}
{"@timestamp":"2021-12-13T08:24:36.585+0000","zxid":4294967301,"type":15,"type_name":"CREATE2","path":"/v1/bootstrap","path_type":"/bootstrap","length":0}
{"@timestamp":"2021-12-13T08:24:36.589+0000","zxid":4294967302,"type":15,"type_name":"CREATE2","path":"/v1/bootstrap/client","path_type":"/bootstrap/client","length":0}
{"@timestamp":"2021-12-13T08:24:36.612+0000","zxid":4294967303,"type":-10,"type_name":"CREATE_SESSION","path":"","path_type":"null","length":0}
{"@timestamp":"2021-12-13T08:24:36.652+0000","zxid":4294967304,"type":7,"type_name":"SET_ACL","path":"","path_type":"null","length":0}
{"@timestamp":"2021-12-13T08:24:36.661+0000","zxid":4294967305,"type":7,"type_name":"SET_ACL","path":"","path_type":"null","length":0}
{"@timestamp":"2021-12-13T08:24:36.665+0000","zxid":4294967306,"type":7,"type_name":"SET_ACL","path":"","path_type":"null","length":0}
{"@timestamp":"2021-12-13T08:24:36.668+0000","zxid":4294967307,"type":7,"type_name":"SET_ACL","path":"","path_type":"null","length":0}

```


### Including ZooKeeper contents for deep analysis
Some investigations require to have low level ECE details (stored in Zookeeper) at hand. It is possible to include this information in the bundle using:

```
-zk|--zookeeper <path_to_dest_pgp_public_key>
```

This option will enable the inclusion of a complete dump of the cluster's ZK ensemble contents into the generated bundle.

This behavior can be constrained so the bundle:

- Would include just the contents of a concrete ZK sub-tree, e.g: (just the contents under `/kibanas`) 
```bash
./ece-diagnostics.sh -zk ./support.key.pub -zk-path '/kibanas'
```
- Exclude certain paths:
```bash
./ece-diagnostics.sh -zk ./support.key.pub -zk-excluded '/zookeeper,/locks'
```
- Or both:
```bash
./ece-diagnostics.sh -zk ./support.key.pub -zk-excluded '/container_sets/cloud-uis,/container_sets/zookeeper-servers' -zk-path '/container_sets'
```

**Note**: How the list of excluded trees is a comma separated list of ZK paths.

:warning: Eagle-eyed readers have probably noticed that this functionality **requires a PGP public key**. ECE Zookeeper contents contain potentially secret and/or sensitive information. For this reason, the bundle **will never contain ZK contents in clear text**. The output will contain a PGP encrypted file which can only be read by the owner of the private key counterpart for the provided key. Elastic Support will provide the key, when requesting this mode to be used.

:warning: By default, paths known to contain secrets or sensitive information have been excluded from the bundle. This protection can be deactivated, for those cases on which it is absolutely necessary, using the option `--zookeeper-excluded-insecure`. Please use this option with extreme caution and only when you are willing to show your passwords and certificates to the owner of the public PGP key.

### Using a custom storage path
If you've installed ECE using a STORAGE_PATH different than default (`/mnt/data/elastic`),  you can pass the below flag to the diagnostics script:

```
./ece-diagnostics.sh -d -s -sp /my/custom/storage/path
```
Note : storage path should be corrected automatically if the storage path folder does not exist


## Output
Diagnostic output archive will be written to /tmp folder with file name ece_diag-<ECE_host_IP>-<Timestamp>.tar.gz  
Once you have the file please provide it to your designated support agent, by attaching it to the support case.

