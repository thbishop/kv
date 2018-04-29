#!/bin/bash -ex

if [ -f /etc/init.d/functions ]; then
  . /etc/init.d/functions
fi

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

# see if there's an existing cluster with this node marked as failed
if /opt/etcd/etcdctl --discovery-srv kv.dyson-sphere.com --insecure-discovery member list | grep -q ${new_hostname}; then 
  echo "Found this node as an existing member; removing from the cluster"
  etcd_member_id=$(/opt/etcd/etcdctl --discovery-srv kv.dyson-sphere.com --insecure-discovery member list | grep $new_hostname | awk -F ':' '{print $1}')
  echo "ETCD member ID is: ${etcd_member_id}"
  echo "Removing as member...."
  /opt/etcd/etcdctl --discovery-srv kv.dyson-sphere.com --insecure-discovery member remove ${etcd_member_id}
  echo "Cleaning up etcd data dir"
  rm -fr /var/cache/etcd/
  echo "Adding this node to the existing cluster"
  add_member_output=$(/opt/etcd/etcdctl \
    --discovery-srv kv.dyson-sphere.com \
    --insecure-discovery \
    member add ${new_hostname} http://${local_ip}:2380 | grep -v Added)
  etcd_vars=($add_member_output)

  for i in "${etcd_vars[@]}"
  do
    echo "export ${i}" >> /var/tmp/etcd_vars
  done

  source /var/tmp/etcd_vars

  cd /opt/etcd
  daemon ./etcd \
    -initial-advertise-peer-urls http://${local_ip}:2380 \
    -advertise-client-urls http://${local_ip}:2379 \
    -listen-client-urls http://0.0.0.0:2379 \
    -listen-peer-urls http://0.0.0.0:2380 \
    -data-dir /var/cache/etcd/state \
    -name $new_hostname > /var/log/etcd.log 2>&1 &

  exit 0
fi

# otherwise join a new cluster
cd /opt/etcd
daemon ./etcd \
  -discovery-srv kv.dyson-sphere.com \
  -initial-advertise-peer-urls http://${local_ip}:2380 \
  -advertise-client-urls http://${local_ip}:2379 \
  -listen-client-urls http://0.0.0.0:2379 \
  -listen-peer-urls http://0.0.0.0:2380 \
  -data-dir /var/cache/etcd/state \
  -name $new_hostname > /var/log/etcd.log 2>&1 &
