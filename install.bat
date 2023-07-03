@echo off

set SERVICE_NAME=NadiShipper
set SERVICE_DISPLAY_NAME="Nadi Shipper"
set SERVICE_DESCRIPTION="Ship Nadi logs to Collector"
set SERVICE_ARGS="-record"

set SERVICE_PATH="%~dp0shipper.exe"
set SERVICE_START_TYPE=auto
set SERVICE_DEPENDENCIES=

sc create %SERVICE_NAME% binPath= %SERVICE_PATH% start= %SERVICE_START_TYPE% DisplayName= %SERVICE_DISPLAY_NAME% description= %SERVICE_DESCRIPTION% error= normal depend= %SERVICE_DEPENDENCIES%
