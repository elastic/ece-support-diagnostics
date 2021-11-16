package main

import (
	"flag"
	"crypto/tls"
	"fmt"
	"regexp"
	"encoding/json"
)

type Certificate struct {
	hostname string
	port string
	Subject string
	Issuer string
	NotAfter string
	message string
}

// Usage : displaySrvCertExpiration -h "coordinator.domain" -p 12443

func removeLBR(text string) string {
    re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
    return re.ReplaceAllString(text, `\n`)
}

func main() {
    var hname, pname string
 	
    flag.StringVar(&hname, "h", "localhost", "Specify password. Default is localhost")
    flag.StringVar(&pname, "p", "12443", "Specify port. Default is 12443")
    flag.Parse() 

    // fmt.Printf(`{ "hostname" : "` + hname + `", "port" : "` + pname + `"`)
    conf := &tls.Config{
	    InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", hname + ":" + pname, conf)
	if err != nil {
		fmt.Printf(`{ "hostname" : "` + hname + `", "port" : "` + pname + `"` + ", \"message\" : \"Server doesn't support SSL certificate [" + removeLBR(err.Error()) + "]\"}\n")

	} else {
		certificate := Certificate{hostname: hname, port: pname, Issuer: conn.ConnectionState().PeerCertificates[0].Issuer.String(), NotAfter: conn.ConnectionState().PeerCertificates[0].NotAfter.String()}
 
		res, err := json.Marshal(certificate)
		     
		if err != nil {
		    fmt.Println(err)
		}
		     
		fmt.Printf("%s\n",res)
	}
}