package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"log"
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

//! Create a prometheus line for DS
func doDSField(value string, proname string, id int, description string, protype string, factor float64) string {
	ret := ""

	name := "cable_downstream_"+proname;
	ret += "# HELP "+name+" "+description+"\n"
	ret += "# TYPE "+name+" "+protype+"\n"
	
	if flval, err := strconv.ParseFloat(value, 32) ; err == nil {
		ret += fmt.Sprintf("%s{id=\"%d\"} %f\n", name, id, flval*factor)
	} else {
		log.Fatalln("Unable to convert value: "+err.Error())
		return ""
	}
	return ret
}


//! Create a prometheus line for US
func doUSField(value string, proname string, id int, description string, protype string) string {
	ret := ""

	name := "cable_upstream_"+proname;
	ret += "# HELP "+name+" "+description+"\n"
	ret += "# TYPE "+name+" "+protype+"\n"
	
	if flval, err := strconv.ParseFloat(value, 32) ; err == nil {
		ret += fmt.Sprintf("%s{id=\"%d\"} %f\n", name, id, flval)
	} else {
		log.Fatalln("Unable to convert value: "+err.Error())
		return ""
	}
	return ret
}


func getPrometheus() string {
	var ret string
	
	//	resp, err := http.Get("http://192.168.178.1/htdocs/cm_info_connection.php")
	resp, err := http.Get("https://berthub.eu/tmp/cm_info_connection.php")
	if err != nil {
		log.Fatalln("Error reading response from modem " + err.Error())
		return ""
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//		panic(err)
		log.Fatalln("Error reading response from modem " + err.Error())
		return ""
	}

	lines := strings.Split(strings.Replace(string(content), "\r\n", "\n", -1), "\n")

	for _, s := range lines {
		if strings.HasPrefix(s, "var cm_conn_json") {
			parts := strings.Split(s, "'")
			lejson := parts[1]

			var f CMInfo
			err := json.Unmarshal([]byte(lejson), &f)
			if err != nil {
				log.Fatalln("Error parsing JSON from modem " + err.Error())
				return ""
			}

			for _, ds := range f.Cm_conn_ds_gourpObj {
				id, err := strconv.Atoi(ds.Ds_id)
				if(err != nil) {
					log.Fatalln("Error parsing downstream id "+ err.Error())
					return ""
				}
				ret += doDSField(ds.Ds_snr, "snr", id,  "Signal to noise ratio of this channel", "gauge", 0.1)
				ret += doDSField(ds.Ds_power, "power", id,  "Power of this channel", "gauge", 1.0)
				ret += doDSField(ds.Ds_correct, "correct", id,  "Correctable errors", "counter", 1.0)
				ret += doDSField(ds.Ds_uncorrect, "uncorrect", id,  "Uncorrectable errors", "counter", 1.0)
				ret += doDSField(ds.Ds_freq, "freq", id,  "Frequency of this channel", "gauge", 1.0)
				ret += doDSField(ds.Ds_width, "width", id,  "Width of this channel", "gauge", 1.0)
				ret += doDSField(ds.Ds_modulation, "modulation", id,  "Modulation of this channel", "gauge", 1.0)
			}

			for _, us := range f.Cm_conn_us_gourpObj {
				id, err := strconv.Atoi(us.Us_id)
				if(err != nil) {
					log.Fatalln("Error parsing upstream id "+ err.Error())
					return ""
				}

				ret += doUSField(us.Us_power, "power", id,  "Power of this channel", "gauge")
				ret += doUSField(us.Us_status, "status", id,  "Status of this channel", "gauge")
				ret += doUSField(us.Us_type, "type", id,  "Type of this channel", "gauge")
				ret += doUSField(us.Us_freq, "freq", id,  "Frequency of this channel", "gauge")
				ret += doUSField(us.Us_width, "width", id,  "Width of this channel", "gauge")
				ret += doUSField(us.Us_modulation, "modulation", id,  "Modulation of this channel", "gauge")
			}
		}
	}
	return ret
}


func main() {
	fmt.Printf("%s", getPrometheus())
}
