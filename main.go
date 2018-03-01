package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type route struct {
	Login    string
	Key      string
	Labels   map[string]string
	Balances []string
}

type conf []*route

var (
	addr       = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	configFile = flag.String("config", "/etc/octopush_exporter.yml", "octopush conf file.")
	gauge      *prometheus.GaugeVec
)

func (c *conf) Describe(ch chan<- *prometheus.Desc) {
	gauge.Describe(ch)
}

func (c *conf) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	fmt.Println("Collecting...")
	for _, r := range *c {
		wg.Add(1)
		go func(r *route) {
			defer wg.Done()
			labels := r.Labels
			balances := getBalances(*r)
			for _, balance := range r.Balances {
				if balances == nil {
					labels["balance"] = balance
					deleted := gauge.Delete(labels)
					if deleted {
						fmt.Printf("deleting errored %v\n", labels)
					}
				} else {
					found := false
					for foundbalance, _ := range balances {
						if foundbalance == balance {
							found = true
						}
					}
					if !found {
						labels["balance"] = balance
						deleted := gauge.Delete(labels)
						if deleted {
							fmt.Printf("deleting inexisting%v\n", labels)
						}
					}
				}
			}
			newbalances := []string{}
			for balance, value := range balances {
				if balance == "" {
					continue
				}
				newbalances = append(newbalances, balance)
				labels["balance"] = balance
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					gauge.Delete(labels)
				} else {
					gauge.With(labels).Set(v)
				}
			}
			if len(newbalances) > 0 {
				r.Balances = append(r.Balances, newbalances...)
			}
		}(r)
	}
	wg.Wait()
	fmt.Println("Collected...")
	gauge.Collect(ch)
}

func (c *conf) readConf() error {

	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("yamlFile.Get error: %v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	var c conf
	err := c.readConf()
	if err != nil {
		panic("Can not read file")
	}
	log.Printf("Loaded %d accounts", len(c))

	r := c[0]
	keys := make([]string, 0)
	for k, _ := range r.Labels {
		keys = append(keys, k)
	}
	keys = append(keys, "balance")
	gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "octopush_balance",
			Help: "Balance",
		},
		keys,
	)
	prometheus.Register(&c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
