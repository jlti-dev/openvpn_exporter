package main

import(
	"log"
	"os"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)
type service struct{
	Host	string
	Port	string
	client	*Client
}

func readEnv() []service{
	i := 1
	var ret []service
	for{
		s := service{
			Host: os.Getenv(fmt.Sprintf("HOST_%d", i)),
			Port: os.Getenv(fmt.Sprintf("PORT_%d", i)),
		}
		if s.Host != "" && s.Port != "" {
			log.Printf("Loaded service %s", s.Host)
			s.client, _ = NewClient(s.Host, s.Port)
			ret = append(ret, s)
			i ++
		}else{
			break
		}
	}
	log.Printf("Loaded %d services", len(ret))
	return ret
}
func main() {
	log.Println("Starting application")
	services := readEnv()
	log.Println("Starting prometheus")
	NewOvpnCollector("ovpn", services)
	http.Handle("/metrics", promhttp.Handler())
	//Running Prometheus (blocking):
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

