ubee-prometheus
---------------
If you have a Ubee modem of the kind used by Ziggo in The Netherlands, this small tool may help you get its data to prometheus.

It works at least on the Ubee UBC1318.

To run, compile ('go build') and launch the binary.

This will periodically poll the modem, which is assumed to be reachable on http://192.168.178.1/ and then keep the results ready for Prometheus.

The URL we poll is: `/htdocs/cm_info_connection.php`

ubee-prometheus is currently hardcoded to listen on port 10000. This will improve once I learn how to do configuration and argument parsing in Go.

No matter how often prometheus polls, ubee-prometheus will only poll your Ubee modem once a minute. This is because a poll takes **20 seconds**, at least on my modem. The result of a Ubee poll is stored in memory, making sure that prometheus gets an anwer quickly.

This also means that at startup, the prometheus URL will serve empty data, until the modem has delivered statistics.

The statistics extracted are:
    * SNR for downstream
    * Power for up and downstream
    * Frequency for up and downstream
    * Width for up and downstream
    * Modulation for up and downstream
    * "Type" for upstream
    * Correctable errors for downstream
    * Uncorrectable errors for downstream



