#!/usr/bin/env bash
# if restart_time_file records 3 times retry attempt. It won't restart service
# after 3 retry times attmpts , delete restart_time_file to restart checking routine
#
host_ip=""
declare -A cluster_query_cmd
cluster_query_cmd[mainnet]="http://$host_ip/mainnet-beta/last6hours"
cluster_query_cmd[testnet]="http://$host_ip/testnet/last6hours"
cluster_query_cmd[devnet]="http://$host_ip/devnet/last6hours"

slack_webhook=""
restart_time_file="$PWD/ping_server_resart_times.out"
restart_times=0
restart_cmd=""

slack_alert(){
	sdata=$(jq --null-input --arg val "$slacktext" '{"text":$val}')
	echo try to sen to slack $sdata
	curl -X POST -H 'Content-type: application/json' --data "$sdata" $slack_webhook
}

write_retry_times() {
	echo $restart_times > $restart_time_file
	echo $(cat $restart_time_file)
}
read_retry_times() {
	restart_times=$(cat $restart_time_file)
	#echo read from file $restart_times
}

# return is_alive, 1: alive 0:not alive (404/408 etc)
alive_check() {
    for retry in 0 1 2
	do
        echo retry=$retry
		if [[ $retry -gt 0 ]];then
			sleep 10
		fi
		alive_status_code=$(curl -o /dev/null -s -w "%{http_code}\n" --connect-timeout 10 ${cluster_query_cmd[$cluster]})
		if [[ $alive_status_code == 204 || $alive_status_code == 200 ]];then
			is_alive=1
			restart_times=0
			write_retry_times 					# reset to zero
            echo $(date --rfc-3339=seconds) => $cluster  $alive_name is alive, status:$alive_status_code
            break
		else
			is_alive=0
            echo $(date --rfc-3339=seconds) => $cluster $alive_name is NOT alive, status:$alive_status_code
		fi
   
	done
}

if [[ ! -f "$restart_time_file"  ]];then
	echo $restart_times > $restart_time_file
else
	read_retry_times
fi
echo $(date --rfc-3339=seconds)  initial restart_times = $restart_times

for cluster in "mainnet"
do
	alive_check
	if [[ $is_alive -eq 0 ]];then
		# restart solana-ping-api.service
		#exec /etc/init.d/solana-ping-api.service restart
		if [[ $restart_times -lt 2 ]];then
			slacktext="{hostname: $HOSTNAME, cluster:$cluster, internal_ip:$host_ip, msg: api last6hours return $alive_status_code and the server is restarting }"
			slack_alert
			restart_times=$(expr $restart_times + 1)
			write_retry_times
			exec /home/sol/ping-api-server/restart-api.sh
			exit 0
		else
			slacktext="{hostname: $HOSTNAME, cluster:$cluster, internal_ip:$host_ip, msg: api last6hours return $alive_status_code and has restarted 2 times."
			slack_alert
		fi	
	fi
done
