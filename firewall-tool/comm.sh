listen_packets(){
time=10
while true; do

pcap_loc=$(date +file_%Y-%m-%d_%H-%M-%S.pcap)
touch $pcap_loc
chmod 777 $pcap_loc
sudo tshark -P -a duration:$time -w $pcap_loc -F pcap
done
}
listen_packets