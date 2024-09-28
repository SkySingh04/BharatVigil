# store_pcap=$(yq eval '.splitcap.Directory_store_actual_pcap' config.yaml)
# session_pcap=$(yq eval '.splitcap.Directory_session_wise_split_pcap' config.yaml)
# time=$(yq eval '.splitcap.Duration_of_collecting_pcap' config.yaml)
# splitcap_loc=$(yq eval '.splitcap.Splitcap_location' config.yaml)
# echo $store_pcap $session_pcap $time $splitcap_loc
# splitcap_actual_loc= "${splitcap_loc}/SplitCap.exe"
# while true; do
#   timestamp=$(date +%s)
#   sudo tcpdump  -w "$store_pcap/capture_$timestamp.pcap" -G $time -W 1  # Capture for 60 seconds
#   sudo "$splitcap_actual_loc"-r $session_pcap/capture_$timestamp.pcap -s session
# done



#!/bin/bash

store_pcap=$(yq eval '.splitcap.Directory_store_actual_pcap' config.yaml)
session_pcap=$(yq eval '.splitcap.Directory_session_wise_split_pcap' config.yaml)
time=$(yq eval '.splitcap.Duration_of_collecting_pcap' config.yaml)
splitcap_loc=$(yq eval '.splitcap.Splitcap_location' config.yaml)
deletion_flag= $(yq eval '.splitcap.deletion_after_split' config.yaml)
size=$(yq eval '.splitcap.size' config.yaml)
data_thread=$(yq eval '.png_param.threads_alloc' config.yaml)


echo "Store PCAP: $store_pcap"
echo "Session PCAP: $session_pcap"
echo "Duration: $time"
echo "SplitCap Location: $splitcap_loc"

splitcap_actual_loc="${splitcap_loc}/SplitCap.exe"

fifo="/tmp/pcap_queue"
fifo1="/tmp/pcap2mnist_queue"

if [[ ! -p $fifo ]]; then
    mkfifo $fifo
fi


if [[ ! -p $fifo1 ]]; then
    mkfifo $fifo1
fi


# listen_packets() {
#   while true; do
#     timestamp=$(date +%s)
#     pcap_file="$store_pcap/capture_$timestamp.pcap"
    
#     echo "Starting tcpdump to capture packets at $pcap_file"
    
  
#     sudo tcpdump -w "$pcap_file" -G $time -W 1
#     #chmod 755 "$pcap_file"

#     echo "$pcap_file" > $fifo
#     echo "Added $pcap_file to the queue"
#   done
# }

split_packets() {
  while true; do

    
    if read pcap_file < $fifo; then
     
      session_file="$session_pcap/capture_$(date +%s)_session"
      
      echo "Splitting $pcap_file into session files at $session_file"
      
      
      sudo "$splitcap_actual_loc" -r $pcap_file -s session -o "$session_file"
       #chmod 755 "$session_file"
      echo "$session_file" > $fifo1
      if [ "$deletion_flag" = true ]; then 

        rm "$pcap_file"
      fi
    fi
    
  done
}

make_png_mnist(){
  while true;do 
   
     if read pcap_file < $fifo1; then
        echo $pcap_file
        taskset -c 0-$((data_thread-1)) python3 session2png_mnsit.py $pcap_file
        
      fi
  done


}

make_prediction(){
  python3 inferencer.py
}



# listen_packets &  
split_packets &   
make_png_mnist &
make_prediction &

wait
