package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type DSStatus struct {
	Ds_id         string
	Ds_snr        string
	Ds_power      string
	Ds_correct    string
	Ds_uncorrect  string
	Ds_modulation string
	Ds_freq       string
	Ds_width      string
}

type USStatus struct {
	Us_status     string
	Us_type       string
	Us_id         string
	Us_freq       string
	Us_width      string
	Us_power      string
	Us_modulation string
}

type CMInfo struct {
	Cm_conn_ip_prov_mode string
	Cm_conn_wan_mode     string
	Cm_conn_ds_gourpObj  []DSStatus
	Cm_conn_us_gourpObj  []USStatus
}

func main() {
	resp, err := http.Get("http://192.168.178.1/htdocs/cm_info_connection.php")
	//	resp, err := http.Get("https://berthub.eu/tmp/cm_info_connection.php")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.Replace(string(content), "\r\n", "\n", -1), "\n")

	for _, s := range lines {
		if strings.HasPrefix(s, "var cm_conn_json") {
			parts := strings.Split(s, "'")
			lejson := parts[1]

			var f CMInfo
			err := json.Unmarshal([]byte(lejson), &f)
			if err != nil {
				panic(err)
			}

			for _, ds := range f.Cm_conn_ds_gourpObj {
				if flval, err := strconv.ParseFloat(ds.Ds_snr, 32); err == nil {
					fmt.Println("# HELP cable_downstream_snr Signal to noise ratio of this channel")
					fmt.Println("# TYPE cable_downstream_snr gauge")
					fmt.Printf("cable_downstream_snr{id=\"%s\"} %f\n", ds.Ds_id, flval/10)
				}
				if flval, err := strconv.ParseFloat(ds.Ds_power, 32); err == nil {
					fmt.Println("# HELP cable_downstream_power Power of this channel")
					fmt.Println("# TYPE cable_downstream_power gauge")
					fmt.Printf("cable_downstream_power{id=\"%s\"} %f\n", ds.Ds_id, flval/10)
				}
				if flval, err := strconv.ParseFloat(ds.Ds_correct, 32); err == nil {
					fmt.Println("# HELP cable_downstream_correct Correctable errors")
					fmt.Println("# TYPE cable_downstream_correct counter")
					fmt.Printf("cable_downstream_correct{id=\"%s\"} %f\n", ds.Ds_id, flval)
				}
				if flval, err := strconv.ParseFloat(ds.Ds_uncorrect, 32); err == nil {
					fmt.Println("# HELP cable_downstream_uncorrect Uncorrectable errors")
					fmt.Println("# TYPE cable_downstream_uncorrect counter")
					fmt.Printf("cable_downstream_uncorrect{id=\"%s\"} %f\n", ds.Ds_id, flval)
				}
				if flval, err := strconv.ParseFloat(ds.Ds_freq, 32); err == nil {
					fmt.Println("# HELP cable_downstream_freq Frequency of this channel")
					fmt.Println("# TYPE cable_downstream_freq gauge")
					fmt.Printf("cable_downstream_freq{id=\"%s\"} %f\n", ds.Ds_id, flval/10)
				}
				if flval, err := strconv.ParseFloat(ds.Ds_modulation, 32); err == nil {
					fmt.Println("# HELP cable_downstream_modulation Modulation of this channel")
					fmt.Println("# TYPE cable_downstream_modulation gauge")
					fmt.Printf("cable_downstream_modulation{id=\"%s\"} %f\n", ds.Ds_id, flval/10)
				}

			}

			for _, us := range f.Cm_conn_us_gourpObj {
				if flval, err := strconv.ParseFloat(us.Us_power, 32); err == nil {
					fmt.Println("# HELP cable_upstream_power Power on this channel")
					fmt.Println("# TYPE cable_upstream_power gauge")
					fmt.Printf("cable_upstream_power{id=\"%s\"} %f\n", us.Us_id, flval)
				}
				if flval, err := strconv.ParseFloat(us.Us_status, 32); err == nil {
					fmt.Println("# HELP cable_upstream_status Status of this channel")
					fmt.Println("# TYPE cable_upstream_status gauge")
					fmt.Printf("cable_upstream_status{id=\"%s\"} %f\n", us.Us_id, flval)
				}
				if flval, err := strconv.ParseFloat(us.Us_type, 32); err == nil {
					fmt.Println("# HELP cable_upstream_type Type of this channel")
					fmt.Println("# TYPE cable_upstream_type gauge")
					fmt.Printf("cable_upstream_type{id=\"%s\"} %f\n", us.Us_id, flval)
				}
				if flval, err := strconv.ParseFloat(us.Us_freq, 32); err == nil {
					fmt.Println("# HELP cable_upstream_freq Freq of this channel")
					fmt.Println("# TYPE cable_upstream_freq gauge")
					fmt.Printf("cable_upstream_freq{id=\"%s\"} %f\n", us.Us_id, flval)
				}
				if flval, err := strconv.ParseFloat(us.Us_width, 32); err == nil {
					fmt.Println("# HELP cable_upstream_width Width of this channel")
					fmt.Println("# TYPE cable_upstream_width gauge")
					fmt.Printf("cable_upstream_width{id=\"%s\"} %f\n", us.Us_id, flval)
				}
				if flval, err := strconv.ParseFloat(us.Us_modulation, 32); err == nil {
					fmt.Println("# HELP cable_upstream_modulation Modulation of this channel")
					fmt.Println("# TYPE cable_upstream_modulation gauge")
					fmt.Printf("cable_upstream_modulation{id=\"%s\"} %f\n", us.Us_id, flval)
				}

			}

		}
	}
}
