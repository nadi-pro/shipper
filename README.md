# Nadi Shipper

Nadi Shipper is a super lightweight log shipper for Nadi app, which transport all the Nadi logs to Nadi API.

## Quick Start

Install Shipper on all the servers you want to monitor.

To download and install Shipper, use the commands that work with your system:

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

Duplicate [nadi.reference.yaml](nadi.reference.yaml) to `nadi.yaml` and update the following values:

- `apiKey` - Login to Nadi app and create your API Token
- `token` - Create an application and copy the Application's token
- `storage` - Set path to Nadi logs

Then you can test the connection to Nadi:

```bash
$ shipper --test
```

In case of monitoring multiple applications, you need to create custom `nadi.yaml` for each of the application.

By default, Nadi Shipper will run as a service.

But you may use Supervisord to run multiple workers to monitor your applications in a single server.

Sample supervisord setup:

```ini
[program:shipper-app1]
process_name=%(program_name)s
command=/usr/local/bin/shipper --config=/path/to/shipper/config/nadi-app1.yaml --record
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/nadi/nadi-app1.log
stopwaitsecs=3600

[program:shipper-app2]
process_name=%(program_name)s
command=/usr/local/bin/shipper --config=/path/to/shipper/config/nadi-app2.yaml --record
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/nadi/nadi-app2.log
stopwaitsecs=3600
```
