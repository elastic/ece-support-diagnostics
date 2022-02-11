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
	Filename string
	Subject string
	Issuer string
	NotBefore string
	NotAfter string
	DNSNames []string
}

// Usage : displayFileCertExpiration -f "/location/cert.pem"

func main() {
    var fpath string
 	
    flag.StringVar(&fpath, "f", "notProvided", "Specify PEM full path")
    flag.Parse()

	if fpath != "notProvided" {
		content, err := ioutil.ReadFile(fpath)
		if err != nil {

			certificate := Certificate{Filename: fpath}
 
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
	            	certificate := Certificate{Filename: fpath}
 
					res, err := json.Marshal(certificate)
				     
				    if err != nil {
				        fmt.Println(err)
				    }
				     
				    fmt.Printf("%s,\n",res)
	            }

	            certificate := Certificate{Filename: fpath, Issuer: cert.Issuer.String(), Subject: cert.Subject.String(), NotBefore: cert.NotBefore.String(), NotAfter: cert.NotAfter.String(), DNSNames: cert.DNSNames}
 
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