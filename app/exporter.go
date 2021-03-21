package main
import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type OvpnCollector struct {
	server		[]service
	namespace	string
	//LoadStats
	countBytesIn	*prometheus.Desc
	countBytesOut	*prometheus.Desc
	gaugeNClients	*prometheus.Desc
	//Status
	countReadBytes	*prometheus.Desc
	countWriteBytes	*prometheus.Desc
	gaugeCountCli	*prometheus.Desc
	//OVClient
	countBytesSent	*prometheus.Desc
	countBytesRec	*prometheus.Desc
}
func NewOvpnCollector(ns string, server []service) *OvpnCollector {
	ret := &OvpnCollector{}
	ret.namespace = ns
	ret.server = server

	ret.countBytesIn = prometheus.NewDesc(
		ns + "_total_load_stat_bytes_in",
		"Number of Bytes In (Load-Stat)",
		[]string{"server"}, nil,
	)
	ret.countBytesOut = prometheus.NewDesc(
		ns + "_total_load_stat_bytes_out",
		"Number of Bytes Out (Load-Stat)",
		[]string{"server"}, nil,
	)
	ret.gaugeNClients = prometheus.NewDesc(
		ns + "_load_stat_nclients",
		"Number of Clients (Load-Stat)",
		[]string{"server"}, nil,
	)

	ret.countReadBytes = prometheus.NewDesc(
		ns + "_total_status_bytes_in",
		"Number of Bytes In",
		[]string{"server"}, nil,
	)
	ret.countWriteBytes = prometheus.NewDesc(
		ns + "_total_status_bytes_out",
		"Number of Bytes Out",
		[]string{"server"}, nil,
	)
	ret.gaugeCountCli = prometheus.NewDesc(
		ns + "_number_status_clients",
		"Number of Clients connected",
		[]string{"server"}, nil,
	)

	ret.countBytesRec = prometheus.NewDesc(
		ns + "_total_status_cn_bytes_in",
		"Number of Bytes In (by cn)",
		[]string{"server", "common_name"}, nil,
	)
	ret.countBytesSent = prometheus.NewDesc(
		ns + "_total_status_cn_bytes_out",
		"Number of Bytes Out (by cn)",
		[]string{"server", "common_name"}, nil,
	)
	prometheus.MustRegister(ret)
	return ret
}
func (c *OvpnCollector) Describe (ch chan<- *prometheus.Desc){
	ch <- c.countBytesIn
	ch <- c.countBytesOut
	ch <- c.gaugeNClients
	
	ch <- c.countReadBytes
	ch <- c.countWriteBytes
	ch <- c.gaugeCountCli

	ch <- c.countBytesSent
	ch <- c.countBytesRec

}

func (c *OvpnCollector) Collect (ch chan<- prometheus.Metric) {
	for _, s := range c.server {
		ls, err := s.client.GetStats()
		if err != nil   {
			log.Printf("[%s] load-stats not available\n", s.Host)
		}else{
			ch <- prometheus.MustNewConstMetric(
				c.countBytesIn, //Description
				prometheus.CounterValue, //Type
				float64(ls.BytesIn), //value
				s.Host,
			)
			ch <- prometheus.MustNewConstMetric(
				c.countBytesOut, //Description
				prometheus.CounterValue, //Type
				float64(ls.BytesOut), //value
				s.Host,
			)
			ch <- prometheus.MustNewConstMetric(
				c.gaugeNClients, //Description
				prometheus.GaugeValue, //Type
				float64(ls.NClients), //value
				s.Host,
			)

		}

		stats, err := s.client.GetDetails()
		if err != nil {
			log.Printf("[%s] no detailed stats available\n", s.Host)
		}else{
			ch <- prometheus.MustNewConstMetric(
				c.countReadBytes, //Description
				prometheus.CounterValue, //Type
				float64(stats.ReadBytes), //value
				s.Host,
			)
			ch <- prometheus.MustNewConstMetric(
				c.countWriteBytes, //Description
				prometheus.CounterValue, //Type
				float64(stats.WriteBytes), //value
				s.Host,
			)
			ch <- prometheus.MustNewConstMetric(
				c.gaugeCountCli, //Description
				prometheus.GaugeValue, //Type
				float64(len(stats.ClientList)), //value
				s.Host,
			)
			
			for _, client := range stats.ClientList{
				ch <- prometheus.MustNewConstMetric(
					c.countBytesSent, //Description
					prometheus.CounterValue, //Type
					float64(client.BytesSent), //value
					s.Host, client.CommonName,
				)
				ch <- prometheus.MustNewConstMetric(
					c.countBytesRec, //Description
					prometheus.CounterValue, //Type
					float64(client.BytesReceived), //value
					s.Host, client.CommonName,
				)
			}
		}
	}
}
