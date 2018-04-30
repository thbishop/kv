#!/bin/bash -ex

aws s3 cp s3://kv-artifacts-us-west-2/etcd-v3.3.4-linux-amd64.tar.gz /var/tmp/ --region us-west-2
cd /var/tmp
tar zxf etcd-v3.3.4-linux-amd64.tar.gz
mkdir /opt/etcd
cp etcd-v3.3.4-linux-amd64/etcd* /opt/etcd/
local_ip=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
az=$(curl http://169.254.169.254/latest/meta-data/placement/availability-zone/)
new_hostname="etcd-${az}.kv.dyson-sphere.com"

cat << EOF > /var/tmp/dns-record-update.json
{
    "HostedZoneId": "",
    "ChangeBatch": {
        "Comment": "updating etcd record",
        "Changes": [
            {
                "Action": "UPSERT",
                "ResourceRecordSet": {
                    "Name": "${new_hostname}",
                    "Type": "A",
                    "TTL": 60,
                    "ResourceRecords": [
                        {
                            "Value": "${local_ip}"
                        }
                    ]
                }
            }
        ]
    }
}
EOF

aws route53 change-resource-record-sets --hosted-zone-id Z1M9OS85ZKADDR --cli-input-json file:///var/tmp/dns-record-update.json

sleep 90

hostname $new_hostname

cd /opt/etcd
nohup ./etcd \
  -discovery-srv kv.dyson-sphere.com \
  -initial-advertise-peer-urls http://${local_ip}:2380 \
  -advertise-client-urls http://${local_ip}:2379 \
  -listen-client-urls http://0.0.0.0:2379 \
  -listen-peer-urls http://0.0.0.0:2380 \
  -data-dir /var/cache/etcd/state \
  -name $new_hostname > /var/log/etcd.log 2>&1 &
