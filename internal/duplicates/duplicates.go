package duplicates

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"golang.org/x/crypto/ssh"
)

type Result struct {
	IP  string
	Key string
}

func FindDuplicates(subnets []string) ([]Result, map[string][]string) {
	var wg sync.WaitGroup
	results := make(chan Result)

	for _, arg := range subnets {
		str := ipaddr.NewIPAddressString(arg)

		addr, err := str.ToAddress()
		if err != nil {
			log.Println(err)
			continue
		}
		i := addr.Iterator()
		for i.HasNext() {
			wg.Add(1)
			ip := i.Next().WithoutPrefixLen().String()
			go func(ip string) {
				defer wg.Done()
				config := &ssh.ClientConfig{
					User: "user",
					Auth: []ssh.AuthMethod{ssh.Password("pw")},
					HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
						//fmt.Println("Received SSH host key:", ssh.FingerprintSHA256(key))
						results <- Result{
							IP:  ip,
							Key: ssh.FingerprintSHA256(key),
						}
						return nil // Accept any host key
					},
					Timeout: 10 * time.Second,
				}

				client, err := ssh.Dial("tcp", ip+":22", config)
				if err != nil {
					// log.Println(err)
					return
				}

				client.Close()

			}(ip)

		}
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()
	var resultList []Result
	byKey := make(map[string][]string)
	for result := range results {
		resultList = append(resultList, result)
	}
	var duplicates []Result
	for _, result := range resultList {
		for _, compared := range resultList {

			if result.Key == compared.Key && result.IP != compared.IP {
				fmt.Printf("found duplicate: %s and %s\n", result.IP, compared.IP)
				fmt.Printf("%s\n%s\n---\n", result.Key, compared.Key)
				duplicates = append(duplicates, result)
				byKey[result.Key] = append(byKey[result.Key], result.IP)
			}
		}

	}
	return duplicates, byKey

}
