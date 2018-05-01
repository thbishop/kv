#!/bin/bash -ex

aws_region=$(curl --silent http://169.254.169.254/latest/dynamic/instance-identity/document|grep region|awk -F\" '{print $4}')
aws s3 cp s3://kv-artifacts-${aws_region}/consul.zip /var/tmp/ --region ${aws_region}
cd /var/tmp
unzip consul.zip
chmod +x ./consul
mkdir -p /opt/consul
mv ./consul /opt/consul
