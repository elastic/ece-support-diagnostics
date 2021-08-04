package main

import (
	"flag"
	"crypto/x509"
	"fmt"
	"regexp"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// Usage : displayFileCertExpiration -f "/location/cert.pem"

func removeLBR(text string) string {
    re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
    return re.ReplaceAllString(text, `\n`)
}

func main() {
    var fpath string
 	
    flag.StringVar(&fpath, "f", "notProvided", "Specify PEM full path")
    flag.Parse()

    fmt.Printf(`{ "filename" : "` + fpath + `", `)
	if fpath != "notProvided" {
		content, err := ioutil.ReadFile(fpath)
		if err != nil {
			fmt.Printf("\"message\" : \"could not read file\" }\n")
			os.Exit(0)
		}
		block, _ := pem.Decode([]byte(content))
	    if block == nil {
	    	fmt.Printf("\"message\" : \"could not parse certificate PEM\" }\n")
	        os.Exit(0)
	    }
	    // cert, err := x509.ParseCertificate(block.Bytes)
	    // if err != nil {
	    //     fmt.Printf("\"message\" : \"failed to parse certificate\" }\n")
	    //     os.Exit(0)
	    // }
	    // fmt.Printf("\"Expiry\": \"%v\" }\n", cert.NotAfter)

	    //above code fails for pem file which contain key first - similar to https://github.com/golang/go/issues/3986
	    cert, err := x509.ParseCertificate(block.Bytes)
	    if err != nil {
	        certs, err := x509.ParseCertificates(block.Bytes)
	        if err != nil {
	        	fmt.Printf("\"message\" : \"failed to parse certificate\" }\n")
			    os.Exit(0)
			}
			fmt.Printf("\"Expiry\": \"")
			for _, certa := range certs {
				fmt.Printf("%v", certa.NotAfter)
			}
			fmt.Printf("\" }\n")
	    }
	    fmt.Printf("\"Expiry\": \"%v\" }\n", cert.NotAfter)
	}
}