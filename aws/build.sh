#!/bin/bash

STEP_COUNTER=1

print_step() {
    echo ""
    echo "***************************************"
    echo ""
    echo "Step ${STEP_COUNTER}: $1"
    echo ""
    echo "***************************************"
    echo ""
}

##########################
# Build EC2 instance
##########################
print_step "Building EC2 instance"
(cd terraform && terraform apply -auto-approve)
sleep 5
let STEP_COUNTER++

##########################
# Get EC2 IP Address
##########################
print_step "Fetching IP address for created instance"
EC2_IP=$(aws ec2 describe-instances | jq . | grep PublicIpAddress | awk -F ":" '{print $2}' | sed 's/[", ]//g')
let STEP_COUNTER++

##########################
# Create ansible host file
##########################
print_step "Creating Ansible hosts file for ip $EC2_IP"
cat >ansible/hosts <<EOL
[all]
$EC2_IP
[all:vars]
ansible_connection=ssh
ansible_ssh_user=ec2-user
EOL
let STEP_COUNTER++

##########################
# Run Ansible play
##########################
print_step "Running ansible play"
(cd ansible && ansible-playbook -i hosts ec2.yml)
let STEP_COUNTER++

print_step "Thirdopinion now up and running at http://${EC2_IP}:8080"
