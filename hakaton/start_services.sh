#!/bin/bash
workers_num=2
data_folder="data"
start_port=8000

mkdir -p $data_folder

for i in $(eval echo "{1..$workers_num}")
do
    port=$(($start_port + $i))
    echo "$data_folder"/"$i"_data.txt
    go run writer/writer.go -p=$port -f="$data_folder"/"$i"_data.txt &
    pid[$i]=$!
done

sleep 2

ports="$start_port"
if [[ $workers_num -gt 1 ]]
then
    for i in $(eval echo "{1..$workers_num}")
    do
        port=$(($start_port + $i))
        ports="$ports,$port"
    done
fi
echo $ports
go run emiter/emiter.go -w=$ports &
pid[$i]=$!
wait
