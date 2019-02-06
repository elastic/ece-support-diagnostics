#!/usr/bin/env bash

VERSION=1.0.0
output_path=/tmp
diag_name=ece_diag_$(hostname)_$(date "+%d_%b_%Y_%H_%M_%S")
diag_folder=$output_path/$diag_name
elastic_folder=$diag_folder/elastic
docker_folder=$diag_folder/docker
docker_logs_folder=$docker_folder/logs

ece_host=localhost
ece_port=12400
protocol=http
user=
password=
cluster_id=
missing_creds=
actions=
storage_path=/mnt/data/elastic

create_folders(){
	while :; do
        	case $1 in
        		system)
			mkdir -p $elastic_folder
			;;
			docker)
			mkdir -p $docker_logs_folder
            		;;
			allocators)
			mkdir -p $elastic_folder/allocators/
			;;
			plan)
			mkdir -p $docker_folder/plan/
			;;
			cluster_info)
			mkdir -p $docker_folder/cluster_info/
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
	rm -rf $diag_folder
}




create_archive(){

	if [ -d $diag_folder ]
		then
			print_msg "Compressing diag file..." "INFO"
                        cd $output_path && tar czf $diag_name.tar.gz $diag_name/* 2>&1
                        print_msg "Diag ready at $output_path/$diag_name.tar.gz" "INFO"
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
	echo "Usage: ./ece-diagnostics.sh [OPTIONS]"
	echo ""
	echo "Options:"
	echo "-e|--ecehost #Specifies ip/hostname of the ECE (default:localhost)"
	echo "-y|--protocol <http/https> #Specifies use of http/https (default:http)"
	echo "-x|--port <port> #Specifies ECE port (default:12400)"
	echo "-s|--system #collects elastic logs and system information"
	echo "-d|--docker #collects docker information"
	echo "-sp|--storage-path #overrides storage path (default:/mnt/data/elastic). Works in conjunction with -s|--system"
	echo "-o|--output-path #Specifies the output directory to dump the diagnostic bundles (default:/tmp)"
	echo "-c|--cluster <clusterID> #collects cluster plan and info for a given cluster (user/pass required). Also restricts -d|--docker action to a specific cluster"
	echo "-a|--allocators #gathers allocators information (user/pass required)"
	echo "-u|--username <username>"
	echo "-p|--password <password>"
	echo ""
	echo "Sample usage:"
	echo "\"./ece-diagnostics.sh -d -s\" #collects system and docker level info"
	echo "\"./ece-diagnostics.sh -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects allocators information"
	echo "\"./ece-diagnostics.sh -e 192.168.1.42 -x 12409 -a -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects allocators information using custom host and port"
	echo "\"./ece-diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -d -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects cluster plan,info and docker info only for the specified cluster ID"
	echo "\"./ece-diagnostics.sh -c e817ac5fbc674aeab132500a263eca71 -u readonly -p oRXdD2tsLrEDelIF4iFAB6RlRzK6Rjxk3E4qTg27Ynj\" #collects cluster plan,info for the specified cluster ID"
	echo ""
}


get_system(){
	#system info
	print_msg "Gathering system info..." "INFO"
	uname -a > $elastic_folder/uname.txt
	cat /etc/*-release > $elastic_folder/linux-release.txt
	top -n1 -b > $elastic_folder/top.txt
	ps -eaf > $elastic_folder/ps.txt
	df -h > $elastic_folder/df.txt

	#network
	sleep 1
	print_msg "Gathering network info..." "INFO"
	sleep 1
	sudo netstat -anp > $elastic_folder/netstat_all.txt 2>&1
	sudo netstat -ntulpn > $elastic_folder/netstat_listening.txt 2>&1
	sudo iptables -L > $elastic_folder/iptables.txt 2>&1
	sudo route -n > $elastic_folder/routes.txt 2>&1

	#mounts
	sudo mount > $elastic_folder/mounts.txt 2>&1
	sudo cat /etc/fstab > $elastic_folder/fstab.txt 2>&1

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
			sar -d -p 1 5 > $elastic_folder/sar_devices.txt 2>&1
			#CPU usage - individual cores - sample 5 times every 1 second
			print_msg "SAR [sampling CPU cores usage]" "INFO"
			sar -P ALL 1 5 > $elastic_folder/sar_cpu_cores.txt 2>&1
			#load average last 1-5-15 minutes - 1 sample
			print_msg "SAR [collect load average]" "INFO"
			sar -q 1 1 > $elastic_folder/sar_load_average_sampled.txt 2>&1
			#memory - sample 5 times every 1 second
			print_msg "SAR [sampling memory usage]" "INFO"
			sar -r 1 5 > $elastic_folder/sar_memory_sampled.txt 2>&1
			#swap - sample once
			print_msg "SAR [collect swap usage]" "INFO"
			sar -S 1 1 > $elastic_folder/sar_swap_sampled.txt 2>&1
			#network
			print_msg "SAR [collect network stats]" "INFO"
			sar -n DEV > $elastic_folder/sar_network.txt 2>&1
		else
			print_msg "'sar' command not found. Please install package 'sysstat' to collect extended system stats" "WARN"
	fi
	print_msg "Grabbing ECE logs" "INFO"
	cd $storage_path && find . -type f -name *.log -exec cp -p --parents \{\} $elastic_folder \;
	print_msg "Checking XFS info" "INFO"
	[[ -x "$(type -P xfs_info)" ]] && xfs_info $storage_path > $elastic_folder/xfs_info.txt 2>&1
}

get_docker(){
	if [ -n $1 ]
		#clusterId is passed as argument - filter on it
		then
			containersId=(`docker ps -a --format "{{.ID}}" --filter="name=$1" `)
        		logNames=(`docker ps -a --format "{{.ID}}__{{.Names}}__{{.Image}}" --filter="name=$1"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g' `)
		#consider all containers
		else
			containersId=(`docker ps -a --format "{{.ID}}"`)
        		logNames=(`docker ps -a --format "{{.ID}}__{{.Names}}__{{.Image}}"  | sed  's/docker\.elastic\.co\///g' | sed 's/[\:\.\/]/_/g' `)
	fi

	print_msg "Grabbing docker logs..." "INFO"
	arrayLength=${#containersId[@]}
	local i=0
	for ((; i<$arrayLength; i++))
	do
		print_msg "Grabbing logs for containerId [${containersId[$i]}]" "INFO"
		docker logs ${containersId[$i]} > $docker_logs_folder/${logNames[$i]}-container.log 2>&1
	done

	print_msg "Grabbing docker ps..." "INFO"
	# output of docker ps -a
	docker ps -a > $docker_folder/ps.txt

	print_msg "Grabbing docker info..." "INFO"
	# output of docker info
	docker info > $docker_folder/info.txt 2>&1

	print_msg "Grabbing docker images..." "INFO"
	# output of docker info
	docker images --all --digests > $docker_folder/images.txt 2>&1

	i=5
	print_msg "Grabbing $i repeated container stats..." "INFO"
	# sample container stats
	while [ $i -ne 0 ] ; do date >> $docker_folder/stats_samples.txt ; print_msg "Grabbing docker stats $i" "INFO"; docker stats --no-stream >> $docker_folder/stats_samples.txt ; i=$((i-1)); done
}


validate_http_creds(){
        if [ -z $user ]
                then missing_creds="$missing_creds user"
        fi
        if [ -z $password ]
                then missing_creds="$missing_creds password"
        fi

}

do_http_request(){

	method=$1
	protocol=$2
	path=$3
	ece_port=$4
	args=$5
	output_file=$6

	#build request
        request="curl -s -X$method -u $user:$password $protocol://$ece_host:$ece_port$path -o $output_file"

	#validation
	validate_http_creds
	if [[ ! -z $missing_creds ]]
		then
			print_msg "Skipping HTTP request [ $path ] because of missing arguments [ $missing_creds ]" "WARN"
                else
					print_msg "Calling [$ece_host:$ece_port$path] with user[$user]" "INFO"
					sleep 1
                                        $request
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
                        get_docker $cluster_id
                        ;;
                        allocators)
                        create_folders allocators
                        do_http_request GET http /api/v1/platform/infrastructure/allocators $ece_port "" $elastic_folder/allocators/allocators.json
                        ;;
			plan)
			validate_http_creds
			if [[ -n $missing_creds ]]
				then print_msg "cannot fetch cluster plan activity without specifying credentials" "WARN"
				else
					if [ -n $cluster_id ]
						then
							create_folders plan
							do_http_request GET http /api/v1/clusters/elasticsearch/$cluster_id/plan/activity $ece_port "" $docker_folder/plan/plan_$cluster_id.json
						else
							print_msg "cannot fetch cluster plan activity without specifying a cluster id. Use option -c|--cluster to specify a cluster ID"	"WARN"
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
                                        		do_http_request GET http "/api/v1/clusters/elasticsearch/$cluster_id" $ece_port "?show_metadata=true&show_plans=true" $docker_folder/cluster_info/cluster_info_$cluster_id.json
                                		else
                                        		print_msg "cannot fetch cluster info without specifying a cluster id. Use option -c|--cluster to specify a cluster ID" "WARN"
                        		fi
			fi
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
	if [ -n $2 ]
		then
			sev="[$2]"
	fi
	echo "`date` $sev:  $1"

}

get_fs_permissions(){
	ls -al $storage_path > $elastic_folder/fs_permissions_storage_path.txt 2>&1
	ls -al /mnt/data > $elastic_folder/fs_permissions_mnt_data.txt 2>&1
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
        -h|--help)
          show_help
          exit
          ;;
        -v|--version)
          echo "-v | --version requested"
          echo $VERSION
          exit
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
		    -u|--user)
                    #user for issuing HTTP requests
		    if [ -z "$2" ]; then
                        die 'ERROR: "-u|--user" requires a username value.'
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
	    	    --)              # End of all options.
            	    shift
                    break
                    ;;
                    -?*)
                    printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
                    ;;
                    *)               # Default case: No more options, so break out of the loop.
                    break
		esac
	    	shift
	done
fi

print_msg "ECE Diagnostics" "INFO"
sleep 1
# go through identified actions and execute
if [ -z "$actions" ]
	then
		: #do nothing
	else
		actions=($actions)
		actionsLength=${#actions[@]}

		for ((i=0; i<$actionsLength; i++))
        		do
        			process_action ${actions[$i]}
		done

fi

create_archive && clean
