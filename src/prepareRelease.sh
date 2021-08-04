#!/bin/bash

promptECEVersion(){
	echo "Enter ECE Version :"
	read ECE_VERSION
	if [[ -z "$ECE_VERSION" ]]; then
		echo "ERROR: ECE version missing"
		exit
	else
		mytmpdir="$(mktemp -d 2>/dev/null || mktemp -d -t 'mytmpdir')"
		output_folder="${mytmpdir}/ece-support-diagnostics-v${ECE_VERSION}"
		mkdir "$output_folder"
	fi
}

verifyECEVersionMatch(){
	if [[ ! "$(cat ${DIR}/ece-diagnostics.sh | grep -c ECE_DIAG_VERSION=${ECE_VERSION})" -eq 1 ]]; then
		echo "ERROR: Version specified [$ECE_VERSION] does not match ece-diagnostics.sh"
		exit 0
	fi
}

findScriptLocation(){
	DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
}

compileGo(){
	#compile one file per container or this messed up architecure of second file
	docker run --rm -e "GOARCH=amd64" -e "GOOS=linux" -v "$DIR:/usr/src/ecediag" -w /usr/src/ecediag golang:1.17rc1-buster go build displayFileCertExpiration.go
	docker run --rm -e "GOARCH=amd64" -e "GOOS=linux" -v "$DIR:/usr/src/ecediag" -w /usr/src/ecediag golang:1.17rc1-buster go build displaySrvCertExpiration.go
	mv "$DIR"/displayFileCertExpiration "${output_folder}/"
	mv "$DIR"/displaySrvCertExpiration "${output_folder}/"
}

addScript(){
	cp "$DIR"/ece-diagnostics.sh "${output_folder}/"
}

createZip(){
	cd "$mytmpdir" && zip -r "ece-support-diagnostics-v${ECE_VERSION}-dist.zip" "ece-support-diagnostics-v${ECE_VERSION}" 
	tar -czf "ece-support-diagnostics-v${ECE_VERSION}-dist.tar.gz" "ece-support-diagnostics-v${ECE_VERSION}" 
	echo "RELEASE FILES : ${mytmpdir}/ece-support-diagnostics-v${ECE_VERSION}-dist.zip ece-support-diagnostics-v${ECE_VERSION}-dist.tar.gz"
}

cleanup(){
	rm -rf "$output_folder" "ece-support-diagnostics-v${ECE_VERSION}"
}

findScriptLocation
promptECEVersion
verifyECEVersionMatch
compileGo
addScript
createZip
cleanup