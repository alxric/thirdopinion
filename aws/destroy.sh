#!/bin/bash

##########################
# Destroy EC2 instance
##########################
(cd terraform && terraform destroy -auto-approve)
