#!/bin/bash

PROJECT="original-nomad-125917"
IBASE="instance"
ZONE="us-central1-b"
MACHT="n1-highcpu-8"
NETWORK="default"
MAINTP="MIGRATE"
STODEF="https://www.googleapis.com/auth/devstorage.read_only"
LOGDEF="https://www.googleapis.com/auth/logging.write"
MONDEF="https://www.googleapis.com/auth/monitoring.write"
USRDEF="https://www.googleapis.com/auth/cloud.useraccounts.readonly"
IMAGE="/debian-cloud/debian-8-jessie-v20160301"
DISKSIZE="15"
DISKTYPE="pd-ssd"

for i in {0..0}
do
  INST="$IBASE-$i"
  # create the instance
  echo "create instance $INST"

  gcloud compute --project "$PROJECT" instances create "$INST" --zone "$ZONE" \
  --machine-type "$MACHT" --network "$NETWORK" --maintenance-policy "$MAINTP" \
  --scopes default="$STODEF","$LOGDEF","$MONDEF","$USRDEF" \
  --image "$IMAGE" --boot-disk-size "$DISKSIZE" --boot-disk-type "$DISKTYPE" \
  --boot-disk-device-name "$INST"

  echo "initialize instance $INST"
  echo "  download init script"
  # download initialization script
  gcloud compute ssh "$INST" --project "$PROJECT" --zone "$ZONE" \
  --command "wget https://raw.githubusercontent.com/jmmaloney3/csci-8994/master/initvm.sh"

  # make initialization script executable
  echo "  make init script executable"
  gcloud compute ssh "$INST" --project "$PROJECT" --zone "$ZONE" --command "chmod +x ./initvm.sh"

  # execute initialization script
  echo "  execute init script"
  gcloud compute ssh "$INST" --project "$PROJECT" --zone "$ZONE" --command "./initvm.sh"
  
  echo "instance $INST ready"
done
