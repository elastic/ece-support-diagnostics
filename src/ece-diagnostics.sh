#!/bin/bash

ECE_DIAG_VERSION=2.0.0

setVariables(){
        output_path=/tmp
        diag_name=ece_diag_$(hostname)_$(date "+%d_%b_%Y_%H_%M_%S")
        diag_folder=$output_path/$diag_name
        elastic_folder=$diag_folder/elastic
        docker_folder=$diag_folder/docker
        docker_logs_folder=$docker_folder/logs
        zookeeper_folder=$elastic_folder/zookeeper_dump
        log_hours=72

        ece_host=localhost
        ece_port=12400
        protocol=http
        user=
        password=
        # cluster_id=
        missing_creds=
        actions=
        storage_path=/mnt/data/elastic
} 

setVariablesZK(){
        pgp_destination_keypath=
        zk_root="NONE"

        zk_excluded=""
        # These path patterns are excluded by default for security and privacy reasons
        zk_excluded="$zk_excluded/container_sets/admin-consoles/containers/admin-console(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/admin-consoles/secrets,"
        zk_excluded="$zk_excluded/container_sets/blueprints/containers/blueprint(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/blueprints/secrets,"
        zk_excluded="$zk_excluded/container_sets/client-observers/containers/client-observer,"
        zk_excluded="$zk_excluded/container_sets/client-observers/secrets,"
        zk_excluded="$zk_excluded/container_sets/cloud-uis/containers/cloud-ui(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/cloud-uis/secrets,"
        zk_excluded="$zk_excluded/container_sets/constructors/containers/constructor(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/constructors/secrets,"
        zk_excluded="$zk_excluded/container_sets/curators/containers/curator(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/curators/secrets,"
        zk_excluded="$zk_excluded/container_sets/directors/containers/director(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/container_sets/directors/secrets,"
        zk_excluded="$zk_excluded/container_sets/proxies/containers/proxy(/inspect@[^/]+),"
        zk_excluded="$zk_excluded/container_sets/proxies/secrets,"
        zk_excluded="$zk_excluded/container_sets/proxies/containers/route-server(/inspect@[^/]+)?,"
        zk_excluded="$zk_excluded/services/runners/[^/]+/containers,"
        zk_excluded="$zk_excluded/clusters/[^/]+/secrets,"
        zk_excluded="$zk_excluded/services/allocators/[^/]+/[^/]+/instances,"
        zk_excluded="$zk_excluded/coordinators/secrets,"
        zk_excluded="$zk_excluded/secrets/certificates,"
        zk_excluded="$zk_excluded/services/adminconsole/secrets,"
        zk_excluded="$zk_excluded/services/proxies/secrets,"
        zk_excluded="$zk_excluded/services/cloudui/secrets,"
        zk_excluded="$zk_excluded/services/internaltls/config,"
        zk_excluded="$zk_excluded/clusters/[^/]+/app-auth-secrets,"
        zk_excluded="$zk_excluded/clusters/[^/]+/instances/instance-\d+/certificates,"
        zk_excluded="$zk_excluded/kibanas/[^/]+/instances/instance-\d+/certificates,"
        zk_excluded="$zk_excluded/[a-z]*/[^/]+/plans"
}

create_folders(){
        while :; do
                case $1 in
                system)
                        mkdir -p "$elastic_folder"
                        ;;
                docker)
                        mkdir -p "$docker_logs_folder"
                        ;;
                zookeeper)
                        mkdir -p "$zookeeper_folder"
                        ;;
                --) # End of all options.
                        shift
                        break
                        ;;
                *) # Default case: No more options, so break out of the loop.
                        break
                esac
                shift
        done
}

clean(){
        print_msg "Cleaning temp files..." "INFO"
        rm -rf "$diag_folder"
}

create_archive(){

        if [ -d "$diag_folder" ]
                then
                        print_msg "Compressing diag file..." "INFO"
                        cd "$output_path" && tar czf "$diag_name".tar.gz "$diag_name"/* 2>&1
                        print_msg "Diag ready at ${output_path}/${diag_name}.tar.gz" "INFO"
                else
                        print_msg "Nothing to do." "INFO"
                        exit 1
        fi
}

die() {
        printf '%s\n' "$1" >&2
        exit 1
}

show_help(){
        echo "ECE Diagnostics"
        echo "Usage: ./diagnostics.sh [OPTIONS]"
        echo ""
        echo "Options:"
        echo "-e|--ecehost #Specifies ip/hostname of the ECE Coordinator (default:localhost)"
        echo "-y|--protocol <http/https> #Specifies use of http/https (default:http)"
        echo "-x|--port <port> #Specifies ECE port (default:12400)"
        echo "-s|--system #collects elastic logs and system information"
        echo "-d|--docker #collects docker information"
        echo "-zk|--zookeeper <path_to_dest_pgp_public_key> #enables ZK contents dump, requires a public PGP key to cipher the contents"
        echo "-zk-path|--zookeeper-path <zk_path_to_include> #selects the ZK sub-tree to dump using the provided path (e.g: /clusters)"
        echo "-zk-excluded|--zookeeper-excluded <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle"
        echo "--zookeeper-excluded-insecure <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle WARNING: This options remove default filters aimed to avoid secrets and sensitive information leaks"
        echo "-sp|--storage-path #overrides storage path (default:/mnt/data/elastic). Works in conjunction with -s|--system"
        echo "-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)"
        # echo "-c|--cluster <clusterID> #collects cluster plan and info for a given cluster (ECE user/pass required). Also restricts -d|--docker action to a specific cluster"
        # echo "-a|--allocators #gathers allocators information (ECE user/pass required)"
        echo "-u|--username <username>"
        echo "-p|--password <password>"
        echo ""
        echo "Sample usage:"
        echo "\"./diagnostics.sh -d -s\" #collects system and docker level info"
        echo "\"./diagnostics.sh -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects ECE APIs information"
        # echo "\"./diagnostics.sh -e 192.168.1.42 -x 12409 -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects allocators information using custom host and port"
        # echo "\"./diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -d -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects cluster plan,info and docker info only for the specified cluster ID"
        # echo "\"./diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects cluster plan,info for the specified cluster ID"
        echo ""
}

get_mntr_ZK(){
        if [[ "$(docker ps -q --filter "name=frc-zookeeper-servers-zookeeper" | wc -l)" -eq 1 ]]; then
                mkdir -p "$elastic_folder"
                docker exec frc-zookeeper-servers-zookeeper sh -c 'for i in $(seq 2191 2199); do echo "$(hostname) port is $i" && echo mntr | nc localhost ${i}; done' > "$elastic_folder"/zk_mntr.txt
        fi
}

get_system(){

        #system info
        print_msg "Gathering system info..." "INFO"
        uname -a > "$elastic_folder"/uname.txt
        cat /etc/*-release > "$elastic_folder"/linux-release.txt
        cat /proc/cmdline > "$elastic_folder"/cmdline.txt
        top -n1 -b > "$elastic_folder"/top.txt
        ps -eaf > "$elastic_folder"/ps.txt
        df -h > "$elastic_folder"/df.txt
        sudo dmesg --ctime > "$elastic_folder"/dmesg.txt

        #network
        sleep 1
        print_msg "Gathering network info..." "INFO"
        sleep 1
        sudo netstat -anp > "$elastic_folder"/netstat_all.txt 2>&1
        sudo netstat -ntulpn > "$elastic_folder"/netstat_listening.txt 2>&1
        sudo iptables -L > "$elastic_folder"/iptables.txt 2>&1
        sudo route -n > "$elastic_folder/"routes.txt 2>&1

        #mounts
        sudo mount > "$elastic_folder"/mounts.txt 2>&1
        sudo cat /etc/fstab > "$elastic_folder"/fstab.txt 2>&1

        #fs permissions
        get_fs_permissions

        #SAR
        print_msg "Gathering SAR output..." "INFO"
        sleep 1
        #check sar exists
        if [ -x "$(type -P sar)" ];
                then
                        #sar individual devices - sample 5 times every 1 second
                        print_msg "SAR [sampling individual I/O devices]" "INFO"
                        sar -d -p 1 5 > "$elastic_folder"/sar_devices.txt 2>&1
                        #CPU usage - individual cores - sample 5 times every 1 second
                        print_msg "SAR [sampling CPU cores usage]" "INFO"
                        sar -P ALL 1 5 > "$elastic_folder"/sar_cpu_cores.txt 2>&1
                        #load average last 1-5-15 minutes - 1 sample
                        print_msg "SAR [collect load average]" "INFO"
                        sar -q 1 1 > "$elastic_folder"/sar_load_average_sampled.txt 2>&1
                        #memory - sample 5 times every 1 second
                        print_msg "SAR [sampling memory usage]" "INFO"
                        sar -r 1 5 > "$elastic_folder"/sar_memory_sampled.txt 2>&1
                        #swap - sample once
                        print_msg "SAR [collect swap usage]" "INFO"
                        sar -S 1 1 > "$elastic_folder"/sar_swap_sampled.txt 2>&1
                        #network
                        print_msg "SAR [collect network stats]" "INFO"
                        sar -n DEV > "$elastic_folder"/sar_network.txt 2>&1
                else
                        print_msg "'sar' command not found. Please install package 'sysstat' to collect extended system stats" "WARN"
        fi
        print_msg "Grabbing ECE logs" "INFO"
        cd "$storage_path" && find . -type f \( -name "*.log" -o -name "*.ndjson" \) -mmin -$((log_hours*60)) -exec cp --preserve=timestamps --parents \{\} "$elastic_folder" \;
        print_msg "Checking XFS info" "INFO"
        [[ -x "$(type -P xfs_info)" ]] && xfs_info "$storage_path" > "$elastic_folder"/xfs_info.txt 2>&1
}

get_docker(){
        if [[ -n "$1" ]]; then
                #clusterId is passed as argument - filter on it
                containersId=($(docker ps -a --format "{{.ID}}" --filter="name=$1"))
                logNames=($(docker ps -a --format "{{.ID}}__{{.Names}}__{{.Image}}" --filter="name=$1"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g'))
                #consider all containers
        else
                containersId=($(docker ps -a --format "{{.ID}}"))
                logNames=($(docker ps -a --format "{{.ID}}__{{.Names}}__{{.Image}}"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g'))
        fi

        print_msg "Grabbing docker logs..." "INFO"
        arrayLength=${#containersId[@]}
        local i=0
        for ((; i<arrayLength; i++))
        do
                print_msg "Grabbing logs for containerId [${containersId[$i]}]" "INFO"
                docker logs --since "${log_hours}h" "${containersId[$i]}" > "${docker_logs_folder}/${logNames[$i]}"-container.log 2>&1
        done

        print_msg "Grabbing docker ps..." "INFO"
        # output of docker ps -a
        docker ps -a > "$docker_folder"/ps.txt

        print_msg "Grabbing docker info..." "INFO"
        # output of docker info
        docker info > "$docker_folder"/info.txt 2>&1

        print_msg "Grabbing docker images..." "INFO"
        # output of docker info
        docker images --all --digests > "$docker_folder"/images.txt 2>&1

        i=5
        print_msg "Grabbing $i repeated container stats..." "INFO"
        # sample container stats
        while [ $i -ne 0 ] ; do date >> "$docker_folder"/stats_samples.txt ; print_msg "Grabbing docker stats $i" "INFO"; docker stats --no-stream >> "$docker_folder"/stats_samples.txt ; i=$((i-1)); done
}

encrypt_file(){

        # This function imitates the behaviour of `gpg2 --recipient-file PUBLIC_KEY_FILE -e FILE` which is only available from gpg 2.1.14
    
        public_key_file=$1
        target_file=$2
        
        temp_keyring=$(mktemp -d)

        gpg2 --homedir "$temp_keyring" --import "$public_key_file"

        recipient=$(gpg2 --homedir "$temp_keyring" -k | grep uid | grep -o '<.\+\@.\+>' | sed 's/[<>]//g' | head -n1)

        gpg2 --homedir "$temp_keyring" --trust-model always --batch --recipient "$recipient" -e "$target_file"
        gpg_result=$?
        
        rm -r "$temp_keyring"

        return "$gpg_result"
}

get_zookeeper(){
        public_key_path=$1
        root_node=""
        if [[ -z "$2" || "$2" == "NONE" ]]
                then
                        die 'ERROR: "-zk|--zookeeper" requires a path for the sub-tree to dump. WARNING: Using "/" might include secret/sensitive information in the bundle.'
                #Path for sub-tree root has been passed
                else
                        root_node="$2"
        fi

        if [ -n "$3" ]
                #List of sub-trees to exclude from the bundle has been passed
                then
                        excluded_nodes="$3"
                #No excluded sub-trees
                else
                        excluded_nodes=","
        fi

        # Check that the current ECE version supports ZK dumps
        docker run --rm "$(docker inspect -f '{{ .Config.Image }}' frc-directors-director)"  ls /elastic_cloud_apps/shell/scripts/dumpZkContents.sc;

        if [ "$?" -ne "0" ];
                then
                        die "ERROR: ECE Version 2.5 or higher is required"
        fi

        # Note that this is the directory (sibling to $elastic_folder) which will contain the clear temporary
        # ZK bundle in clear text prior to encryption. It will be deleted automatically.
        zookeeper_cleartext_folder=$(mktemp -d $elastic_folder/../zookeeper_dump_temporary.XXXX)

        #Collect result at $zookeeper_cleartext_folder/zkdump.zip
        #This is done outside the bundle directory to avoid accidental inclusion of
        #ZK contents in clear text within the bundle.

        docker run --env SHELL_JAVA_OPTIONS="-Dfound.shell.exec=/elastic_cloud_apps/shell/scripts/dumpZkContents.sc -Dfound.shell.exec-params=pathsToSkip=${excluded_nodes};rootPath=${root_node};outputPath=/target/zkdump.zip" \
               -v "$zookeeper_cleartext_folder":/target -v ~/.found-shell:/elastic_cloud_apps/shell/.found-shell \
               --env SHELL_ZK_AUTH=$(docker exec -it frc-directors-director bash -c 'echo -n $FOUND_ZK_READWRITE') \
               $(docker inspect -f '{{ range .HostConfig.ExtraHosts }} --add-host {{.}} {{ end }}' frc-directors-director) \
               --rm $(docker inspect -f '{{ .Config.Image }}' frc-directors-director) \
               /elastic_cloud_apps/shell/run-shell.sh;
        
        #Cipher dump file and remove the one in clear text

        # gpg2 --recipient-file $public_key_path -e $zookeeper_folder/zkdump.zip #Ideally we'd use this but it requires a version not so ubiquitous.
        encrypt_file "$public_key_path" "$zookeeper_cleartext_folder"/zkdump.zip;
        encryption_result=$?
        
                # Collect the encrypted version of the ZK contents bundle and include it in the general bundle.
                mv "$zookeeper_cleartext_folder"/zkdump.zip.gpg "$zookeeper_folder "
        rm -r "$zookeeper_cleartext_folder" # Then, delete the temporary directory.

        if [ "$encryption_result" -ne "0" ];
                then
                        die "ERROR: Failed to encrypt ZK dump bundle"
        fi
}

validate_http_creds(){
        if [ -z "$user" ]
                then missing_creds="$missing_creds user"
        fi
        if [ -z "$password" ]
                then missing_creds="$missing_creds password"
        fi
}

do_http_request(){

        method=$1
        protocolrequest=$2
        path=$3
        ece_port=$4
        args=$5
        output_file=$6

        #build request
        request="curl -s -S -X$method -u $user:$password $protocolrequest://$ece_host:$ece_port$path -o $output_file"

        #validation
        validate_http_creds
        if [[ -n $missing_creds ]]
                then
                        print_msg "Skipping HTTP request [ $path ] because of missing arguments [ $missing_creds ]" "WARN"
                else
                        print_msg "Calling [$ece_host:$ece_port$path] with user [$user]" "INFO"
                        sleep 1
                        STDERR=$($request 2>&1)
                        if [ ! -s "$output_file" ]; then
                                print_msg "Output from API call is empty - please ensure you are connecting to a coordinator node with -e" "ERROR"
                                print_msg "${STDERR}" "ERROR"
                        elif grep -q "root.unauthenticated" "$output_file"; then
                                print_msg "The supplied authentication is invalid - please use ECE admin user/pass" "ERROR"
                                clean
                                exit
                        elif grep -q "clusters.cluster_not_found" "$output_file"; then
                                print_msg "Specified Cluster ID is invalid.  The Elasticsearch cluster ID can be found within the endpoint URL" "ERROR"
                        fi
        fi
}

process_action(){
        while :; do
                case $1 in
                system)
                        verifyStoragePath
                        create_folders system
                        get_system
                        ;;
                docker)
                        create_folders docker
                        ;;
                allocators)
                        create_folders allocators
                        do_http_request GET $protocol /api/v1/platform/infrastructure/allocators $ece_port "" $elastic_folder/allocators/allocators.json
                        do_http_request GET $protocol /api/v1/platform $ece_port "" $elastic_folder/allocators/platform.json
                        do_http_request GET $protocol /api/v1/clusters/elasticsearch $ece_port "" $elastic_folder/allocators/elasticsearch-clusters.json
                        get_docker "$cluster_id"
                        ;;
                plan)
                        validate_http_creds
                        if [[ -n "$missing_creds" ]]
                                then print_msg "cannot fetch cluster plan activity without specifying credentials" "WARN"
                                else
                                        if [ -n "$cluster_id" ]
                                                then
                                                        create_folders plan
                                                        do_http_request GET "$protocol" /api/v1/clusters/elasticsearch/"$cluster_id"/plan/activity "$ece_port" "" "${docker_folder}/plan/plan_${cluster_id}.json"
                                                else
                                                        print_msg "cannot fetch cluster plan activity without specifying a cluster id. Use option -c|--cluster to specify a cluster ID"        "WARN"
                                        fi
                        fi
                        ;;
                cluster_info)
                        validate_http_creds
                        if [[ -n $missing_creds ]]
                                then print_msg "cannot fetch cluster info plan without specifying credentials" "WARN"
                                else
                                        if [ -n $cluster_id ]
                                                then
                                                        create_folders cluster_info
                                                        do_http_request GET $protocol "/api/v1/clusters/elasticsearch/$cluster_id" $ece_port "?show_metadata=true&show_plans=true" $docker_folder/cluster_info/cluster_info_$cluster_id.json
                                                else
                                                        print_msg "cannot fetch cluster info without specifying a cluster id. Use option -c|--cluster to specify a cluster ID" "WARN"
                                        fi
                        fi
                        ;;
                zookeeper)
                        create_folders zookeeper
                        get_zookeeper "$pgp_destination_keypath" "$zk_root" "$zk_excluded"
                        ;;
                --)              # End of all options.
                        shift
                        break
                        ;;
                *)               # Default case: No more options, so break out of the loop.
                        break
                esac
                shift
        done
}

print_msg(){
        #$1 msg
        #$2 sev
        local sev=
        if [ -n "$2" ]
                then
                        sev="[${2}]"
        fi
        echo "$(date) ${sev}:  ${1}" | tee -a "$diag_folder"/ece-diag.log

}

promptPassword(){
        echo -n "Enter password for ${user} : "
        read -s password
get_fs_permissions(){
        ls -al "$storage_path" > "$elastic_folder"/fs_permissions_storage_path.txt 2>&1
        ls -al /mnt/data > "$elastic_folder"/fs_permissions_mnt_data.txt 2>&1
}

}

#BEGIN

# no arguments -> show help
if [ "$#" -eq 0 ]; then
        show_help
# arguments - parse them
else
        while :; do
                case $1 in
                -s|--system)
                        #gather system data
                        actions="$actions system"
                        ;;
                -lh|--log-filter-hours)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-sf|--log-filter-hours" requires a valid number of hours'
                                else
                                        log_hours=$2
                                        shift
                                fi
                                ;;
                -sp|--storage-path)
                        #changes -s behaviour by
                        #overriding default $storage_path value (/mnd/data/elastic)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-sp|--storage-path" requires a valid full filesystem path to custom storage'
                        else
                                storage_path=$2
                                shift
                        fi
                        ;;
                -o|--output-path)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-o|--output-path" requires a valid full filesystem path'
                        else
                                output_path=$2
                                diag_folder=$output_path/$diag_name
                                elastic_folder=$diag_folder/elastic
                                docker_folder=$diag_folder/docker
                                docker_logs_folder=$docker_folder/logs
                                shift
                        fi
                        ;;
                -e|--ecehost)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-e|--ecehost" requires a hostname/ip value.'
                        else
                                ece_host=$2
                                shift
                        fi
                        ;;
                -a|--allocators)
                        #gather allocators data
                        actions="$actions allocators"
                        ;;
                        -u|--username)
                        #user for issuing HTTP requests
                        if [ -z "$2" ]; then
                                die 'ERROR: "-u|--username" requires a username value.'
                        else
                                user=$2
                                shift
                        fi
                        ;;
                -p|--password)
                        #password for issuing HTTP requests
                        if [ -z "$2" ]; then
                                die 'ERROR: "-p|--password" requires a password value.'
                        else
                                password=$2
                                shift
                        fi
                        ;;
                -x|--port)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-x|--port" requires a port value.'
                        else
                                ece_port=$2
                                shift
                        fi
                        ;;
                -d|--docker)
                        #gather docker data
                        actions="$actions docker"
                        ;;
                -y|--protocol)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-y|--protocol" requires a protocol value.'
                        else
                                protocol=$2
                                shift
                        fi
                        ;;
                -c|--cluster)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-c|--cluster" requires a clusterId value.'
                        else
                                cluster_id=$2
                                actions="$actions plan cluster_info"
                                shift
                        fi
                        ;;
                -zk|--zookeeper)
                        # First check PGP tools are available
                        gpg2 --help 2>/dev/null > /dev/null;
                        if [ "$?" -ne "0" ]; then
                                die 'ERROR: "-zk|--zookeeper" requires `gnupg2` to be installed in the system.'
                        fi
                        
                        if [ -z "$2" ]; then
                                die 'ERROR: "-zk|--zookeeper" requires a PGP destination public key.'
                        else
                                pgp_destination_keypath=$2
                                actions="$actions zookeeper"
                                shift
                        fi
                        ;;
                -zk-path|--zookeeper-path)
                        # Sets Zookeeper target sub-tree
                        if [ -z "$2" ]; then
                                die 'ERROR: This options requires a path string'
                        else
                                zk_root=$2
                                shift
                        fi
                        ;;
                -zk-excluded|--zookeeper-excluded)
                        # Sets Zookeeper exclusion paths
                        if [ -n "$2" ]; then
                                zk_excluded="$zk_excluded,$2"
                                shift
                        fi
                        ;;
                --zookeeper-excluded-insecure)
                        # Sets Zookeeper exclusion paths removing defaults Secret/Sensitive exclusions
                        if [ -n "$2" ]; then
                                print_msg "WARNING!! This option may lead to the inclusion of secrets and sensitive information within the bundle."
                                zk_excluded="$2"
                                shift
                        fi
                        ;;
                --)             # End of all options.
                        shift
                        break
                        ;;
                -?*)
                        printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
                        ;;
                *)              # Default case: No more options, so break out of the loop.
                        break
                esac
                shift
        done
        if [[ -n "$user" ]] && [[ -z "$password" ]]; then
                promptPassword
        fi
fi



runECEDiag(){
        sleep 1
        # go through identified actions and execute
        if [ -z "$actions" ]
                then
                        : #do nothing
                else
                        actions=("$actions")
                        actionsLength=${#actions[@]}

                        for ((i=0; i<actionsLength; i++))
                                do
                                        process_action "${actions[$i]}"
                        done

        fi
        create_archive && clean
}



verifyStoragePath(){
        if [[ ! -d "$storage_path" ]]; then
                #docker inspect frc-runners-runner | grep logs:/app/logs returns :
                #"/test/julien/elastic/10.0.2.15/services/runner/logs:/app/log",
                local sto_path
                sto_path="$(docker inspect frc-runners-runner 2>/dev/null | grep logs:/app/logs | cut -d ':' -f1 | cut -d '"' -f2)"
                sto_path="$(dirname ${sto_path/\/services\/runner\/logs/})"
                if [[ -d "${sto_path}" ]]; then
                        print_msg "Storage path [${storage_path}] is not accessible, correcting to [${sto_path}]" "WARN"
                        storage_path="$sto_path"
                else 
                        print_msg "Storage path [${storage_path}] is not accessible, found [${sto_path}] but folder not valid" "ERROR"
                        print_msg "-sp|--storage-path #overrides storage path (default:/mnt/data/elastic)." "INFO"
                        clean
                        exit 0
                fi
        fi
}

initiateLogFile(){
        touch "$diag_folder"/ece-diag.log
        print_msg "ECE Diagnostics ${ECE_DIAG_VERSION}" "INFO"
        print_msg "Arguments used [$*]" "INFO"
}

setVariables

parseParams "$@"

initiateLogFile "$@"

runECEDiag

get_mntr_ZK