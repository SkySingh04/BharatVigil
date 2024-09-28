# listen_packets(){
# time=10
# while true; do

# pcap_loc=$(date +file_%Y-%m-%d_%H-%M-%S.pcap)
# touch $pcap_loc
# chmod 777 $pcap_loc
# sudo tshark -P -a duration:$time -w $pcap_loc -F pcap
# done
# }
# listen_packets


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

store_pcap=$(yq eval '.splitcap.Directory_store_actual_pcap' ../config.yaml)
session_pcap=$(yq eval '.splitcap.Directory_session_wise_split_pcap' ../config.yaml)
time=$(yq eval '.splitcap.Duration_of_collecting_pcap' ../config.yaml)
splitcap_loc=$(yq eval '.splitcap.Splitcap_location' ../config.yaml)
deletion_flag= $(yq eval '.splitcap.deletion_after_split' ../config.yaml)
size=$(yq eval '.splitcap.size' ../config.yaml)
data_thread=$(yq eval '.png_param.threads_alloc' ../config.yaml)
mnist_file_loc=$(yq eval '.png_param.python_file_loc' ../config.yaml)
model_python_file_loc=$(yq eval '.model.model_python_file_loc' ../config.yaml)
png_store_dir=$(yq eval '.png_param.png_saving_dir' ../config.yaml)
# echo "Store PCAP: $store_pcap"
# echo "Session PCAP: $session_pcap"
# echo "Duration: $time"
# echo "SplitCap Location: $splitcap_loc"
# echo "inferencer_file_loc: $mnist_file_loc" 
splitcap_actual_loc="${splitcap_loc}/SplitCap.exe"

fifo="/tmp/pcap_queue"
fifo1="/tmp/pcap2mnist_queue"
fifo2="/tmp/mnist2model_queue"
if [[ ! -p $fifo ]]; then
    mkfifo $fifo
fi


if [[ ! -p $fifo1 ]]; then
    mkfifo $fifo1
fi

if [[ ! -p $fifo2 ]]; then
    mkfifo $fifo2
fi

listen_packets() {
  while true; do
    timestamp=$(date +%s)
    pcap_file="$store_pcap/capture_$timestamp.pcap"
    
    #echo "Starting tcpdump to capture packets at $pcap_file"
    
  
    #sudo tcpdump -w "$pcap_file" -G $time -W 1
    #chmod 755 "$pcap_file"

touch $pcap_file
chmod 777 $pcap_file
#echo "Listening to $pcap_file"
my_ip=$( hostname -I )
sudo tshark -P -a duration:$time -w $pcap_file -f "dst host ${my_ip}" -F pcap
#sleep $time
    echo "$pcap_file" > $fifo
    #echo "Added $pcap_file to the queue"
  done
}

split_packets() {
  while true; do

    
    if read pcap_file < $fifo; then
      filename=$(basename "$pcap_file" .pcap )
      session_file="${session_pcap}/${filename}_session/"
      
      #echo "Splitting $pcap_file into session files at $session_file"
      
      
      sudo "$splitcap_actual_loc" -r $pcap_file -s session -o "$session_file"  > /dev/null 2>&1

       #chmod 755 "$session_file"
      echo "$session_file" > $fifo1
      #echo  " Added $session_file  to fifo1 for png "
      if [ "$deletion_flag" = true ]; then 

        rm "$pcap_file"
      fi
    fi
    
  done
}

make_png_mnist(){
  while true; do 
     
     if read pcap_file < $fifo1; then
        #echo $pcap_file
        taskset -c 0-$((data_thread-1)) python3 $mnist_file_loc $pcap_file  > /dev/null 2>&1
        fol_name=$(basename "$pcap_file")
        png_loc="${png_store_dir}/${fol_name}/"
        echo "$png_loc" > $fifo2
        #echo Added $png_loc to $fifo2
        #echo "$session_file" > $fifo1
      fi
  done


}

make_prediction(){
  # while true; do 
  # if read pcap_file < $fifo2; then 
  #   echo "Lauda $pcap_file"
  # fi 
  # done
  python3 $model_python_file_loc  > /dev/null 2>&1
}



listen_packets &  
split_packets &   
make_png_mnist &
make_prediction &

wait
