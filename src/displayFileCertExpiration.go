package main

import (
	"flag"
	"crypto/x509"
	"fmt"
	"encoding/pem"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Certificate struct {
	filename string
	Subject string
	Issuer string
	NotAfter string
	message string
}

// Usage : displayFileCertExpiration -f "/location/cert.pem"

func main() {
    var fpath string
 	
    flag.StringVar(&fpath, "f", "notProvided", "Specify PEM full path")
    flag.Parse()

	if fpath != "notProvided" {
		content, err := ioutil.ReadFile(fpath)
		if err != nil {

			certificate := Certificate{filename: fpath, message: "could not read file"}
 
			res, err := json.Marshal(certificate)
		     
		    if err != nil {
		        fmt.Println(err)
		    }
		     
		    fmt.Printf("%s,\n",res)
			os.Exit(0)
		}
	    
	    for block, rest := pem.Decode([]byte(content)); block != nil; block, rest = pem.Decode(rest) {
	        switch block.Type {
	        case "CERTIFICATE":
	        	// Showing certificate information
	            cert, err := x509.ParseCertificate(block.Bytes)
	            if err != nil {
	            	certificate := Certificate{filename: fpath, message: "could not parse certificate"}
 
					res, err := json.Marshal(certificate)
				     
				    if err != nil {
				        fmt.Println(err)
				    }
				     
				    fmt.Printf("%s,\n",res)
	            }
	            
	            certificate := Certificate{filename: fpath, Issuer: cert.Issuer.String(), Subject: cert.Subject.String(), NotAfter: cert.NotAfter.String()}
 
				res, err := json.Marshal(certificate)
				     
			    if err != nil {
			        fmt.Println(err)
			    }
				     
			    fmt.Printf("%s,\n",res)

	        default:
	            //ignoring any other type of blocks
	        }
	    }
	}
}