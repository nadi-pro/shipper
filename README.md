# Nadi Shipper

Nadi Shipper is a super lightweight log shipper for Nadi app, which transport all the Nadi logs to Nadi API.

## Quick Start

Install Shipper on all the servers you want to monitor.

To download and install Filebeat, use the commands that work with your system:

### Install via bash script (Linux & Mac)

Linux & Mac users can install it directly to `/usr/local/bin/shipper` with:

```bash
sudo bash < <(curl -sL https://raw.githubusercontent.com/nadi-pro/shipper/master/install)
```

### Download static binary (Windows, Linux and Mac)

Run the following command which will download latest version and configure default configuration for Windows.
```batch
powershell -command "(New-Object Net.WebClient).DownloadFile('https://raw.githubusercontent.com/nadi-pro/shipper/master/install.ps1', '%TEMP%\install.ps1') && %TEMP%\install.ps1 && del %TEMP%\install.ps1"
```

### Configuring Nadi

Duplicate [nadi.reference.yaml](nadi.reference.yaml) to `nadi.yaml` and update the respective values:

```yaml
###################### Nadi Configuration ##################################

# This file is a configuration file for Nadi Shipper.
#
# You can find the full configuration reference here:
# https://docs.nadi.pro/nadi-shipper

# ============================== Nadi inputs ===============================

nadi:
  # Nadi API Endpoint
  endpoint: http://nadi.pro/api/

  # Accept Header
  accept: application/vnd.nadi.v1+json

  # Login to Nadi app and create your API Token from
  apiKey:

  # Create an application and copy the Application's token and paste it here.
  token:

  # Set path to Nadi logs. By default the path is /var/log/nadi.
  storage: /var/log/nadi

  # Set the Path for tracker.json
  trackerFile: tracker.json

  # Set this to true if you want to maintain the Nadi log after sending them. Default is false.
  persistent: false

  # Set maximum tries to send over the logs. Default is 3 times.
  maxTries: 3

  # Set maximum time before timeout. Default is 1 minute.
  timeout: 1m
```
