package main

import (
	"flag"
	"crypto/tls"
	"fmt"
	"regexp"
)

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

    fmt.Printf(`{ "hostname" : "` + hname + `", "port" : "` + pname + `"`)
	conn, err := tls.Dial("tcp", hname + ":" + pname, nil)
	if err != nil {
		fmt.Printf(", \"message\" : \"Server doesn't support SSL certificate err:\\n" + removeLBR(err.Error()) + `"`)
	} else {
		expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
		fmt.Printf(", \"Issuer\": \"%s\", \"Expiry\": \"%v\"", conn.ConnectionState().PeerCertificates[0].Issuer, expiry)
	}
	fmt.Printf(" }\n")
}