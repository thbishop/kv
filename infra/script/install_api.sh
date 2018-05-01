#!/bin/bash -ex

aws_region=$(curl --silent http://169.254.169.254/latest/dynamic/instance-identity/document|grep region|awk -F\" '{print $4}')
aws s3 cp s3://kv-artifacts-${aws_region}/api.zip /var/tmp/ --region ${aws_region}
cd /var/tmp
unzip /var/tmp/api.zip
mkdir -p /opt/kv/api
mv /var/tmp/api /opt/kv/api/
chmod +x /opt/kv/api/api
useradd kv-api
chown -R kv-api:kv-api /opt/kv
mv /var/tmp/kv-api /etc/init.d/
chmod +x /etc/init.d/kv-api
chkconfig kv-api on
service kv-api start
