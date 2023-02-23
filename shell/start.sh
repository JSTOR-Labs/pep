#!/bin/bash
NAME=""                                   
while getopts ":u:" options; do 
  case "${options}" in                   
    u) NAME=${OPTARG} ;;
  esac
done

if [ "$NAME" = "" ]; then 
  echo "A username is required. Please use the -u flag to set a username."
  exit 0 
else
 echo "path.data: /home/$NAME/JSTOR/es_data" >> /home/$NAME/JSTOR/es_config/elasticsearch.yml
 cp /home/$NAME/JSTOR/es_config/elasticsearch.yml /etc/elasticsearch
 adduser $NAME elasticsearch
 mkdir /usr/share/elasticsearch/logs
 chmod +x /usr/share/elasticsearch/bin/elasticsearch
 chown -R $NAME:elasticsearch /etc/default/elasticsearch
 chown -R $NAME:elasticsearch /etc/elasticsearch
 chown -R $NAME:elasticsearch /var/log/elasticsearch
 chown -R $NAME:elasticsearch /usr/share/elasticsearch/logs
 chmod +x /home/$NAME/JSTOR/api-chromebook
 
 echo "/usr/share/elasticsearch/bin/elasticsearch & /home/$NAME/JSTOR/api-chromebook serve" >> /home/$NAME/.bashrc
fi
exit 0