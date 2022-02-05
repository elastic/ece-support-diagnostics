#!/bin/bash

ECE_DIAG_VERSION=2.0.6

setVariables(){
        #location of scripts
        DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

        output_path=/tmp
        diag_name=ece_diag_$(hostname)_$(date "+%d_%b_%Y_%H_%M_%S")
        diag_folder=$output_path/$diag_name
        initiateLogFile
        elastic_folder=$diag_folder/elastic
        docker_folder=$diag_folder/docker
        docker_logs_folder=$docker_folder/logs
        zookeeper_folder=$elastic_folder/zookeeper_dump
        zookeeper_stats_folder=$elastic_folder/zookeeper_stats
        log_hours=72

        ece_host=localhost
        ece_port=12400
        protocol=http
        user=
        password=
        # cluster_id=
        missing_creds=
        actions=
        deployments=
        storage_path=/mnt/data/elastic
        arr=0 #used to store APIs in 4 arrays
        RED='\033[0;31m'
        NC='\033[0m' # No Color
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
                plan)
                        mkdir -p "$elastic_folder"/plan/
                        ;;
                cluster_info)
                        mkdir -p "$elastic_folder"/cluster_info/
                        ;;
                zookeeper_stats)
                        mkdir -p "$zookeeper_stats_folder"
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
        echo "-k|--insecure #Bypass certificate validity checks when using https"
        echo "-ca|--cacert /path/ca.pem #Specify CA certificate when using https"
        echo "-x|--port <port> #Specifies ECE port (default:12400)"
        echo "-s|--system #collects elastic logs and system information"
        echo "-d|--docker #collects docker information"
        echo "-zk|--zookeeper <path_to_dest_pgp_public_key> #enables ZK contents dump, requires a public PGP key to cipher the contents"
        echo "-zk-path|--zookeeper-path <zk_path_to_include> #selects the ZK sub-tree to dump using the provided path (e.g: /clusters)"
        echo "-zk-excluded|--zookeeper-excluded <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle"
        echo "--zookeeper-excluded-insecure <excluded_paths> #optional, comma separated list of sub-trees to exclude in the bundle WARNING: This options remove default filters aimed to avoid secrets and sensitive information leaks"
        echo "--zk-stats|--zookeeper-stats #collects statistics on zookeeper contents and events"
        echo "-sp|--storage-path #overrides storage path (default:/mnt/data/elastic). Works in conjunction with -s|--system"
        echo "-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)"
        echo "-de|--deployment <deploymentID2,deploymentID2> #collects deployment historic plan activity logs (ECE username required), comma separated value allowed to pass multiple deployments. Default to collecting this for all unhealthy deployments, pass value \"-de disabled\" to not collect any deployment activity logs"
        echo "-u|--username <username>"
        echo "-p|--password <password>"
        echo ""
        echo "Sample usage:"
        echo "\"./diagnostics.sh -d -s\" #collects system and docker level info"
        echo "\"./diagnostics.sh -u admin -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects default ECE APIs information"
        echo "\"./diagnostics.sh -de e817ac5fbc674aeab132500a263eca71 -u admin -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects default APIs information plus deployment plan"
        echo ""
        clean
        exit
}

get_mntr_ZK(){
        if [[ "$(docker ps -q --filter "name=frc-zookeeper-servers-zookeeper" | wc -l)" -eq 1 ]]; then
                mkdir -p "$elastic_folder"
                docker exec frc-zookeeper-servers-zookeeper sh -c 'for i in $(seq 2191 2199); do echo mntr | nc localhost ${i} 2>/dev/null; done' > "$elastic_folder"/zk_mntr.txt
        fi
}

get_certificate_srv(){
        if [[ -f "$DIR"/displaySrvCertExpiration ]]; then
                print_msg "Getting certificate expiration for [${ece_host}:12443]" "INFO"
                if [[ "$protocol" = "http" ]]; then
                        bash -c "${DIR}/displaySrvCertExpiration -h ${ece_host} -p 12443 2>/dev/null" > "$elastic_folder"/certs/coordinator_12443.json
                elif [[ "$protocol" = "https" ]]; then
                        bash -c "${DIR}/displaySrvCertExpiration -h ${ece_host} -p "$ece_port" 2>/dev/null" > "$elastic_folder"/certs/coordinator_12443.json
                fi
        else
                print_msg "Binary missing [${DIR}/displaySrvCertExpiration]" "WARN"
        fi
}

get_certificate_files(){
        if [[ -f "$DIR"/displayFileCertExpiration ]]; then
                print_msg "Getting certificate expiration for PEM files" "INFO"
                echo '[' > "${elastic_folder}/certs/pem_files_expiration.json"
                find "$storage_path" -type f \( -name "*.pem" -o -name "*.crt" \) -exec "${DIR}"/displayFileCertExpiration -f \{\} >> "${elastic_folder}/certs/pem_files_expiration.json" \;
                #remove last character which may be a coma (to obtain valid json array) or newline in case of empty set
                truncate -s-1 "${elastic_folder}/certs/pem_files_expiration.json"
                echo ' ]' >> "${elastic_folder}/certs/pem_files_expiration.json"
        else
                print_msg "Binary missing [${DIR}/displayFileCertExpiration]" "WARN"
        fi
}


get_certificate_expiration(){
        mkdir -p "$elastic_folder"/certs
        if [[ -n "$user" ]]; then
                get_certificate_srv
        fi
        get_certificate_files
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
        timedatectl > "$elastic_folder"/timedatectl.txt

        #network
        sleep 1
        print_msg "Gathering network info..." "INFO"
        sleep 1

        #sudo calls should be located here so they can be disabled
        if [[ -z "$disableSudoCalls" ]]; then
                #system info
                sudo dmesg --ctime > "$elastic_folder"/dmesg-localTZ.txt
                #network
                sudo netstat -anp > "$elastic_folder"/netstat_all.txt 2>&1
                sudo netstat -ntulpn > "$elastic_folder"/netstat_listening.txt 2>&1
                sudo iptables -L -v -n > "$elastic_folder"/iptables.txt 2>&1
                sudo ip rule show > "$elastic_folder"/ip_routing_tables.txt 2>&1
                sudo route -n > "$elastic_folder/"routes.txt 2>&1
                #mounts
                sudo mount > "$elastic_folder"/mounts.txt 2>&1
                sudo cat /etc/fstab > "$elastic_folder"/fstab.txt 2>&1
                sudo cat /etc/sysctl.conf > "$elastic_folder"/sysctl.txt 2>&1
        fi

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
                logNames=($(docker ps -a --format "{{.ID}}__{{.Names}}" --filter="name=$1"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g'))
                #consider all containers
        else
                containersId=($(docker ps -a --format "{{.ID}}"))
                logNames=($(docker ps -a --format "{{.ID}}__{{.Names}}"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g'))
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
        chmod 777 "$zookeeper_cleartext_folder"

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

get_zookeeper_stats(){
        # Check that the current ECE version supports ZK stats
        docker run --rm "$(docker inspect -f '{{ .Config.Image }}' frc-directors-director)"  ls /elastic_cloud_apps/tools/zkstat.jar;

        if [ "$?" -ne "0" ];
                then
                        die "ERROR: ECE Version 3.0 or higher is required"
        fi

        zk_data_path=$(find $storage_path -type d -wholename "*zookeeper/data" 2>/dev/null);

        zkstat() {
                docker run -v "$zookeeper_stats_folder":/target -v "$zk_data_path":/zk_data    \
                        $(docker inspect -f '{{ range .HostConfig.ExtraHosts }} --add-host {{.}} {{ end }}' frc-directors-director) \
                        --rm $(docker inspect -f '{{ .Config.Image }}' frc-directors-director) \
                        java -jar /elastic_cloud_apps/tools/zkstat.jar $@;
        }

        print_msg "Collecting ZooKeeper nodes stats" "INFO"

        stats_file="zk_stats.csv"
        stats_container_path="/target/$stats_file"
        zkstat node-stats /zk_data $stats_container_path
        zkstat nodetype-stats-csv $stats_container_path /target/zk_nodetype_stats.csv
        rm -f "$zookeeper_stats_folder/$stats_file"

        print_msg "Collecting ZooKeeper translog" "INFO"
        zkstat logs-stats /zk_data /target/translog.json

        print_msg "Done collecting ZooKeeper stats" "INFO"
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

        # echo "method[${method}] protocolrequest[${protocolrequest}] path[${path}] ece_port[${ece_port}] args[${args}] output_file[${output_file}]"

        #build request
        request="curl -s -S -X${method} -u ${user}:${password} ${protocolrequest}://${ece_host}:${ece_port}${path} -o ${output_file}"

        if [[ "$insecure" = true ]]; then
                request="${request} -k"
        fi

        if [[ -n "$cacert" ]]; then
                request="${request} --cacert ${cacert}"
        fi

        #validation
        validate_http_creds
        if [[ -n $missing_creds ]]
                then
                        print_msg "Skipping HTTP request [ ${path} ] because of missing arguments [ ${missing_creds} ]" "WARN"
                else
                        print_msg "Calling [${ece_host}:${ece_port}${path}] with user [$user]" "INFO"
                        sleep 1
                        STDERR=$($request 2>&1)
                        if [ ! -s "$output_file" ]; then
                                print_msg "Output from API call is empty - please ensure you are connecting to a coordinator node with -e" "ERROR"
                                print_msg "${STDERR}" "ERROR"
                                if [[ "${ece_host}" = "localhost" ]]; then
                                        if [[ $(docker ps -a | grep -c "frc-admin-consoles-admin-console") -eq 0 ]]; then
                                                print_msg "APIs requested for ${protocolrequest}://localhost:${ece_port} but current host is not a coordinator" "ERROR"
                                        fi
                                fi
                                clean
                                exit
                        elif grep -q "root.unauthenticated" "$output_file"; then
                                print_msg "Diag bundle could not be generated !" "ERROR"
                                print_msg "The supplied authentication is invalid - please use readonly user (or admin user)" "ERROR"
                                print_msg "Please fix credentials or omit APIs call by not specifying any username" "ERROR"
                                clean
                                exit
                        elif grep -q "clusters.cluster_not_found" "$output_file"; then
                                print_msg "Specified Cluster ID is invalid.  The Elasticsearch cluster ID can be found within the endpoint URL" "ERROR"
                        fi
        fi

        #Side note : the following may break json syntax
        #removing any line containing certificate information from output for security deployment when using admin account (instead of the expected readonly)
        if [[ "$(grep -c 'signing' "${output_file}")" -gt 0 ]]; then
                grep -v "signing" "$output_file" > "${output_path}/temp.json" && mv "${output_path}/temp.json" "$output_file"
        fi

        #removing any line containing key for repository (example access_key or secret_key)
        if [[ "$output_file" = "${elastic_folder}/platform/configuration/repositories.json" ]] && [[ "$(grep -c 'key' "${output_file}")" -gt 0 ]]; then
                grep -v "key" "${output_file}" > "${output_path}/temp.json" && mv "${output_path}/temp.json" "${output_file}"
        fi
} 

process_action(){
        while :; do
                case $1 in
                system)
                        create_folders system
                        get_system
                        ;;
                docker)
                        create_folders docker
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
                                                        addApiCall "/api/v1/clusters/elasticsearch/${cluster_id}/plan/activity" "${elastic_folder}/plan/plan_${cluster_id}.json"
                                                else
                                                        print_msg "cannot fetch cluster plan activity without specifying a cluster id. Use option -c|--cluster to specify a cluster ID" "WARN"
                                        fi
                        fi
                        ;;
                cluster_info)
                        validate_http_creds
                        if [[ -n "$missing_creds" ]]
                                then print_msg "cannot fetch cluster info plan without specifying credentials" "WARN"
                                else
                                        if [ -n "$cluster_id" ]
                                                then
                                                        create_folders cluster_info
                                                        addApiCall "/api/v1/clusters/elasticsearch/${cluster_id}?show_security=false&show_metadata=false&show_plans=true&show_plan_logs=true&show_system_alerts=3" "${elastic_folder}/cluster_info/cluster_info_${cluster_id}.json" '2.0.0'
                                                else
                                                        print_msg "cannot fetch cluster info without specifying a cluster id. Use option -c|--cluster to specify a cluster ID" "WARN"
                                        fi
                        fi
                        ;;
                zookeeper)
                        create_folders zookeeper
                        get_zookeeper "$pgp_destination_keypath" "$zk_root" "$zk_excluded"
                        ;;
                zk_stats)
                        create_folders zookeeper_stats
                        get_zookeeper_stats
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
        if [[ "$sev" = "[ERROR]" ]] || [[ "$sev" = "[WARN]" ]]; then
                echo -e "${RED}$(date) ${sev}:  ${1}${NC}" | tee -a "$diag_folder"/ece-diag.log
        else
                echo "$(date) ${sev}:  ${1}" | tee -a "$diag_folder"/ece-diag.log
        fi
}

get_fs_permissions(){
        ls -al "$storage_path" > "$elastic_folder"/fs_permissions_storage_path.txt 2>&1
        ls -al /mnt/data > "$elastic_folder"/fs_permissions_mnt_data.txt 2>&1
}

api_get_platform(){
        do_http_request GET "$protocol" /api/v1/platform "$ece_port" "" "$elastic_folder"/platform/platform.json
}

#compare 2 versions, 0 when lower, 1 when equal, 2 when greater
function vercomp () {
    if [[ $1 == $2 ]]
    then
        return 1
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 2
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 0
        fi
    done
    return 0
}

extractPlatformVersion(){
        ece_version="$(grep version ${elastic_folder}/platform/platform.json | head -1 | cut -d ":" -f2 | cut -d '"' -f2)"
        if [[ ! "$ece_version" =~ [0-9]+\.[0-9]+\.[0-9]+ ]]; then
                if [[ "$(grep -c 'controller_route.controller_not_ready' ${elastic_folder}/platform/platform.json)" -gt 0 ]]; then
                        print_msg "Bypassing APIs; GET platform returned controller_route.controller_not_ready error" "WARN"
                else
                        print_msg "Version could not be found [$ece_version] in platform.json[$(cat ${elastic_folder}/platform/platform.json)]" "WARN"
                        ece_version=$(docker ps -a --filter "name=frc-runners-runner" --format "{{.Image}}" | rev | cut -d ':' -f1 | rev)
                        if [[ ! "$ece_version" =~ [0-9]+\.[0-9]+\.[0-9]+ ]]; then
                                print_msg "Version could not be found [$ece_version] using frc-runners-runner image" "ERROR"
                                print_msg "Cannot run APIs without version information so aborting diagnostics, run without -u to exlude APIs" "ERROR"
                                clean 
                                exit
                                ece_version=
                        fi
                fi
        fi
}

addApiCall(){
        api_url[$arr]="$1"
        api_file[$arr]="$2"
        api_min[$arr]="${3:-2.0.0}"
        api_max[$arr]="$4"
        (( arr=arr+1 ))
}

apis_platform(){
        mkdir -p "${elastic_folder}/platform/license"

        addApiCall '/api/v1/platform/license' "${elastic_folder}/platform/license/license.json" '2.0.0'

        mkdir -p "${elastic_folder}/platform/infrastructure"
        addApiCall '/api/v1/platform/infrastructure/allocators' "${elastic_folder}/platform/infrastructure/allocators.json" '2.0.0'
        addApiCall '/api/v1/platform/infrastructure/blueprinter/roles' "${elastic_folder}/platform/infrastructure/roles.json" '2.3.0'
        addApiCall '/api/v1/platform/infrastructure/constructors' "${elastic_folder}/platform/infrastructure/constructors.json" '2.2.0'
        addApiCall '/api/v1/platform/infrastructure/proxies' "${elastic_folder}/platform/infrastructure/proxies.json" '2.2.0'
        addApiCall '/api/v1/platform/infrastructure/runners' "${elastic_folder}/platform/infrastructure/runners.json" '2.0.0'

        mkdir -p "${elastic_folder}/platform/configuration"
        addApiCall '/api/v1/platform/configuration/instances?show_deleted=false' "${elastic_folder}/platform/configuration/instances.json" '2.0.0'
        addApiCall '/api/v1/platform/configuration/templates/deployments?show_instance_configurations=false' "${elastic_folder}/platform/configuration/deployment_templates.json" '2.0.0'
        addApiCall '/api/v1/platform/configuration/store' "${elastic_folder}/platform/configuration/store.json" '2.2.0'
        addApiCall '/api/v1/platform/configuration/security/realms' "${elastic_folder}/platform/configuration/realms.json" '2.2.0'
        addApiCall '/api/v1/platform/configuration/snapshots/repositories' "${elastic_folder}/platform/configuration/repositories.json"
        # addApiCall '/api/v1/platform/configuration/security/deployment' "${elastic_folder}/platform/configuration/security.json"
}

apis_stacks(){
        mkdir -p "${elastic_folder}/stacks"
        addApiCall '/api/v1/stack/versions' "${elastic_folder}/stacks/versions.json" '2.0.0'
}

apis_users(){
        mkdir -p "${elastic_folder}/users"
        addApiCall '/api/v1/users' "${elastic_folder}/users/users.json" '2.4.0'
}

apis_deployments(){
        mkdir -p "${elastic_folder}/deployments"
        addApiCall '/api/v1/deployments' "${elastic_folder}/deployments/deployments.json" '2.4.0'
        #this call return historical plan activity logs (can be many Mb per deployment)
        if [[ -n "$deployments" ]]; then
                vercomp "$ece_version" '2.4.0'
                if [[ $? -ge 1 ]]; then
                        deployments=($(printf "$deployments" | tr "," " "))
                        deploymentsLength=${#deployments[@]}
                        for ((i=0; i<deploymentsLength; i++))
                        do
                                if [[ ! "${deployments[$i]}" = "disabled" ]]; then
                                        addApiCall "/api/v1/deployments/${deployments[$i]}?show_security=false&show_metadata=false&show_plans=true&show_plan_logs=true&show_plan_history=true&show_plan_defaults=false&convert_legacy_plans=false&show_system_alerts=3&show_settings=true&enrich_with_template=true" "${elastic_folder}/deployments/${deployments[$i]}-detailed.json" '2.4.0'
                                fi
                        done
                else 
                        print_msg "-de|--deployment option has no effect prior to ECE 2.4.0, detected [${ece_version}], use -c|--cluster instead for ES cluster" "WARN"
                fi
        fi
}

apis_clusters_legacy(){
        #only running for ECE versions below 2.4.0 (min version for v1 APIs 2.0.0)
        vercomp "$ece_version" '2.4.0'
        if [[ $? -eq 0 ]]; then
                mkdir -p "${elastic_folder}/clusters"
                addApiCall "/api/v1/clusters/elasticsearch?from=0&size=100&show_security=false&show_metadata=false&show_plans=false&show_plan_defaults=false&convert_legacy_plans=false&enrich_with_template=false&show_system_alerts=3&show_hidden=false&show_settings=false" "${elastic_folder}/clusters/es-clusters.json" '2.0.0' '2.3.2'
                addApiCall "/api/v1/clusters/kibana?from=0&size=100&show_metadata=false&show_plans=false&show_hidden=false&show_plan_defaults=false" "${elastic_folder}/clusters/kibana-clusters.json" '2.0.0' '2.3.2'
        fi
}


apis_v0(){
        mkdir -p "${elastic_folder}/v0containersets"
        addApiCall '/api/v0/regions/ece-region/container-sets/allocators' "${elastic_folder}/v0containersets/allocators.json"
        addApiCall '/api/v0/regions/ece-region/container-sets/proxies' "${elastic_folder}/v0containersets/proxies.json"
        addApiCall '/api/v0/regions/ece-region/container-sets/zookeeper-servers' "${elastic_folder}/v0containersets/zookeeper-servers.json"
}

prepare_apis_arrays(){
        #these just build list of APIs in 4 arrays : api_url, api_file, api_min, api_max
        mkdir -p "${elastic_folder}/platform"
        # api_get_platform
        do_http_request GET "$protocol" /api/v1/platform "$ece_port" "" "$elastic_folder"/platform/platform.json
        extractPlatformVersion

        apis_platform
        apis_stacks
        apis_users
        apis_deployments
        apis_clusters_legacy
        apis_v0
}

run_api(){
        #TODO: if api_file is empty, comput from api_url
        #TODO: add "mkdir -p" here
        do_http_request GET "$protocol" "${api_url[$a]}" "$ece_port" "" "${api_file[$a]}"
}

run_apis(){
        #iterates through arrays and run API when ece_version match min or/and max version
        for ((a=0;a<${#api_url[@]};a++))
        do
                if [[ "${api_min[$a]}" = "" ]]; then
                        if [[ -z "${api_max[$a]}" ]]; then #no min, no max
                                run_api
                        else  #no min, max
                                vercomp "$ece_version" "${api_max[$a]}"
                                if [[ $? -le 1 ]]; then
                                        run_api
                                fi
                        fi
                else
                        if [[ -z "${api_max[$a]}" ]]; then #min, no max
                                vercomp "$ece_version" "${api_min[$a]}"
                                if [[ $? -ge 1 ]]; then
                                        run_api
                                fi
                        else  #min, max
                                vercomp "$ece_version" "${api_max[$a]}"
                                if [[ $? -le 1 ]]; then
                                        vercomp "$ece_version" "${api_min[$a]}"
                                        if [[ $? -ge 1 ]]; then
                                                run_api
                                        fi
                                fi
                        fi

                fi
        done 
}

collect_apis_data(){
        prepare_apis_arrays
        run_apis
}

promptPassword(){
        echo -n "Enter password for ${user} : "
        read -s -r password
}

parseParams(){
      # no arguments -> show help
        if [ "$#" -eq 0 ]; then
                show_help
        # arguments - parse them
        else
                while :; do
                        case $1 in
                        -sp|--storage-path)
                                #changes -s behaviour by
                                #overriding default $storage_path value (/mnd/data/elastic)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-sp|--storage-path" requires a valid full filesystem path to custom storage'
                                else
                                        storage_path=$2
                                        options="${options} -sp ${storage_path}"
                                        shift
                                fi
                                ;;
                        -s|--system)
                                #gather system data
                                actions="$actions system"
                                options="${options} -s"
                                ;;
                        -ds|--disable-sudo)
                                #disable sudo calls
                                disableSudoCalls=true
                                print_msg "Disabling sudo calls is not recommended" "WARN"
                                options="${options} --disable-sudo"
                                ;;
                        -lh|--log-filter-hours)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-sf|--log-filter-hours" requires a valid number of hours'
                                else
                                        log_hours=$2
                                        options="${options} -lh ${log_hours}"
                                        shift
                                fi
                                ;;
                        -o|--output-path)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-o|--output-path" requires a valid full filesystem path'
                                else
                                        output_path=$2
                                        options="${options} -o ${output_path}"
                                        diag_folder=$output_path/$diag_name
                                        elastic_folder=$diag_folder/elastic
                                        docker_folder=$diag_folder/docker
                                        docker_logs_folder=$docker_folder/logs
                                        shift
                                fi
                                ;;
                        -e|--ecehost)
                                if [[ -z "$2" ]] || [[ "$2" = -* ]]; then
                                        die 'ERROR: "-e|--ecehost" requires a hostname/ip value.'
                                else
                                        ece_host=$2
                                        options="${options} -e ${ece_host}"
                                        shift
                                fi
                                ;;
                        -a|--allocator)
                                print_msg '"-a|--allocator" option is deprecated and will be collected' "WARN"
                                options="${options} -a"
                                ;;
                        -u|--username)
                                if [[ -z "$2" ]] || [[ "$2" = -* ]]; then
                                        die 'ERROR: "-u|--user" requires a username (admin or readonly).'
                                else
                                        user=$2
                                        options="${options} -u ${user}"
                                        if [[ "$user" = "readonly" ]]; then
                                                print_msg "Using -u|--username with value [readonly] will cause v0 APIs to fail, prefer [admin] credentials" "WARN"
                                        fi
                                        shift
                                fi
                                ;;
                        -p|--password)
                                #password for issuing HTTP requests
                                if [ -n "$2" ]; then
                                        if [[ -z "$user" ]]; then
                                                print_msg "Password was provided without user - APIs will not run" "WARN"
                                        fi
                                        if [[ ! "$2" = -* ]]; then
                                                password=$2
                                                shift
                                        fi
                                fi
                                ;;
                        -x|--port)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-x|--port" requires a port value.'
                                else
                                        ece_port=$2
                                        options="${options} -x ${ece_port}"
                                        shift
                                fi
                                ;;
                        -d|--docker)
                                #gather docker data
                                actions="$actions docker"
                                options="${options} -d"
                                ;;
                        -de|--deployment)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-de|--deployment" requires a value (comma separated for multiple deployment IDs)'
                                else
                                        deployments=$2
                                        options="${options} -de ${deployments}"
                                        shift
                                fi
                                ;;
                        -y|--protocol)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-y|--protocol" requires a protocol value.'
                                else
                                        protocol=$2
                                        options="${options} -y ${protocol}"
                                        shift
                                fi
                                ;;
                        -k|--insecure)
                                insecure=true
                                options="${options} -k"
                                ;;
                        -ca|--cacert)
                                if [ -z "$2" ]; then
                                        die 'ERROR: "-ca|--cacert" requires a valid ca certificate file.'
                                else
                                        cacert=$2
                                        options="${options} -ca ${cacert}"
                                        if [[ ! -f "$cacert" ]]; then
                                                print_msg "Specified CA file not found [${cacert}]" "WARN"
                                                cacert=
                                        fi
                                        shift
                                fi
                                ;;
                        -c|--cluster)
                        if [ -z "$2" ]; then
                                die 'ERROR: "-c|--cluster" requires a clusterId value.'
                        else
                                cluster_id=$2
                                options="${options} -c ${cluster_id}"
                                actions="$actions plan cluster_info"
                                print_msg '"-c|--cluster" option is deprecated, prefer "-de|--deployment" instead' "WARN"
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
                                        setVariablesZK
                                        pgp_destination_keypath=$2
                                        options="${options} -zk ${pgp_destination_keypath}"
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
                                        options="${options} -zk-path ${zk_root}"
                                        shift
                                fi
                                ;;
                        -zk-excluded|--zookeeper-excluded)
                                # Sets Zookeeper exclusion paths
                                if [ -n "$2" ]; then
                                        zk_excluded="$zk_excluded,$2"
                                        options="${options} -zk-excluded ${zk_excluded}"
                                        shift
                                fi
                                ;;
                        --zookeeper-excluded-insecure)
                                # Sets Zookeeper exclusion paths removing defaults Secret/Sensitive exclusions
                                if [ -n "$2" ]; then
                                        print_msg "This option may lead to the inclusion of secrets and sensitive information within the bundle." "WARN"
                                        zk_excluded="$2"
                                        options="${options} --zookeeper-excluded-insecure ${zk_excluded}"
                                        shift
                                fi
                                ;;
                        --zk-stats|--zookeeper-stats)
                                # Gathers Zookeeper nodes statistics
                                actions="$actions zk_stats"
                                options="${options} --zk-stats"
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
                print_msg "Options used : ${options}" "INFO"
        fi
        if [[ -n "$user" ]] && [[ -z "$password" ]]; then
                promptPassword
        fi
}

findIndentationDeploymentId(){
        #deployments API return pretty printed json with deployment ids, and cluster ids, indentation changed
        #hacky way to parse json
        indentation="$(cat "${elastic_folder}/deployments/deployments.json" | grep " \"id\"" | head -1 | cut -d '"' -f1)"
}

runECEDiag(){
        sleep 1
        # go through identified actions and execute
        if [ -n "$actions" ]
                then
                        actions=($actions) #using word spitting
                        actionsLength=${#actions[@]}

                        for ((i=0; i<actionsLength; i++))
                        do
                                process_action "${actions[$i]}"
                        done

        fi
        if [[ -n "$user" ]]; then
                collect_apis_data
        fi
        #This code iterate deployment ids without plan activity logs (it uses output of another API hence code is run last)
        vercomp "$ece_version" 2.4.0
        if [[ $? -ge 1 ]]; then
                if [[ -f "${elastic_folder}/deployments/deployments.json" ]]; then
                        findIndentationDeploymentId
                        deployment_ids=$(grep -e "^${indentation}\\\"id\\\"" "${elastic_folder}/deployments/deployments.json" | cut -d '"' -f4)
                        deployment_ids=(${deployment_ids})
                        for deployment_id in "${deployment_ids[@]}"
                        do
                                do_http_request GET "$protocol" "/api/v1/deployments/${deployment_id}?show_security=false&show_metadata=false&show_plans=true&show_plan_logs=false&show_plan_history=false&show_plan_defaults=false&convert_legacy_plans=false&show_system_alerts=0&show_settings=true&enrich_with_template=false" "$ece_port" "" "${elastic_folder}/deployments/${deployment_id}.json"
                                #collection of plan activity logs for unhealthy deployments when -de is not specified
                                if [[ -z "$deployments" ]] && [[ "$(grep -ce '\"healthy\" : false,' ${elastic_folder}/deployments/${deployment_id}.json)" -gt 0 ]]; then
                                        do_http_request GET "$protocol" "/api/v1/deployments/${deployment_id}?show_security=false&show_metadata=false&show_plans=true&show_plan_logs=true&show_plan_history=true&show_plan_defaults=false&convert_legacy_plans=false&show_system_alerts=0&show_settings=true&enrich_with_template=false" "$ece_port" "" "${elastic_folder}/deployments/${deployment_id}-detailed.json"
                                fi
                        done
                fi
        fi
        get_mntr_ZK
        get_certificate_expiration
        create_archive && clean
}



verifyStoragePath(){
        #function will attempt to correct storage location - this may deprecate -sp option
        if [[ ! -d "$storage_path" ]]; then
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
        mkdir -p "$diag_folder"
        touch "$diag_folder"/ece-diag.log
        print_msg "ECE Diagnostics ${ECE_DIAG_VERSION}" "INFO"
}

setVariables

parseParams "$@"

verifyStoragePath

runECEDiag