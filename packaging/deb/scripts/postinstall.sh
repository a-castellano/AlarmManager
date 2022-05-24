#!/bin/sh

mkdir -p /etc/windmaker-alarmmanager

echo "### NOT starting on installation, please execute the following statements to configure windmaker-alarmmanager to start automatically using systemd"
echo " sudo /bin/systemctl daemon-reload"
echo " sudo /bin/systemctl enable windmaker-alarmmanager"
echo "### You can start windmaker-alarmmanager by executing"
echo " sudo /bin/systemctl start windmaker-alarmmanager"
