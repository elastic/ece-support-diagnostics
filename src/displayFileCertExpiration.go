package main

import (
	"flag"
	"crypto/x509"
	"fmt"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// Usage : displayFileCertExpiration -f "/location/cert.pem"

func main() {
    var fpath string
 	
    flag.StringVar(&fpath, "f", "notProvided", "Specify PEM full path")
    flag.Parse()

	if fpath != "notProvided" {
		content, err := ioutil.ReadFile(fpath)
		if err != nil {
			fmt.Printf(`{ "filename" : "` + fpath + "\"message\" : \"could not read file\" }\n")
			os.Exit(0)
		}
	    
	    for block, rest := pem.Decode([]byte(content)); block != nil; block, rest = pem.Decode(rest) {
	        switch block.Type {
	        case "CERTIFICATE":
	        	// Showing certificate information
	            cert, err := x509.ParseCertificate(block.Bytes)
	            if err != nil {
	                fmt.Printf("\n{ \"filename\": \"%s\", \"message\" : \"could not parse certificate\" },")
	            }
	            
	            fmt.Printf("\n{ \"filename\": \"%s\", \"Issuer\" : \"%s\", \"Subject\" : \"%s\", \"NotAfter\" : \"%v\" },", fpath, cert.Issuer , cert.Subject, cert.NotAfter)

	        default:
	            //ignoring any other type of blocks
	        }
	    }
	}
}