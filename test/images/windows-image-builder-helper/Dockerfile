# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM mcr.microsoft.com/windows/servercore:ltsc2019 as prep

# busybox
ENV CURL_VERSION=7.57.0 \
    PS_VERSION=6.2.0
WORKDIR /curl
ADD https://skanthak.homepage.t-online.de/download/curl-$CURL_VERSION.cab curl.cab
ADD https://github.com/PowerShell/PowerShell/releases/download/v$PS_VERSION/PowerShell-$PS_VERSION-win-x64.zip /PowerShell/powershell.zip
ADD https://eternallybored.org/misc/netcat/netcat-win32-1.12.zip /netcat/netcat.zip
ADD https://github.com/kubernetes-sigs/windows-testing/raw/master/images/busybox/busybox.exe /bin/busybox.exe

USER ContainerAdministrator

# NOTE(claudiub):we have to create relative path symlinks because docker buildx has an issue when copying
# over symlinks, it prepends "Files\" to the symlink target. "Files\C:\bin\busybox.exe" would be an
# invalid path.
RUN cd C:\bin && FOR /f "tokens=*" %i IN ('.\busybox --list') DO mklink .\%i.exe busybox.exe

RUN expand /R curl.cab /F:* . &\
    del C:\curl\curl.cab &\
    cd C:\PowerShell &\
    tar.exe -xf powershell.zip &\
    del powershell.zip &\
    mklink powershell.exe pwsh.exe &\
    cd C:\netcat &\
    tar.exe -xf netcat.zip &\
    del netcat.zip &\
    setx /M PATH "C:\bin;C:\curl\;%PATH%"

# download bind and prepare a folder just for dig.
RUN powershell -Command "\
      curl.exe 'https://downloads.isc.org/isc/bind9/9.14.10/BIND9.14.10.x64.zip' -o /bind.zip; \
      Expand-Archive -Path C:\bind.zip -DestinationPath C:\bind; \
      $s = [System.Diagnostics.Process]::Start('C:\bind\vcredist_x64.exe', '/quiet'); \
      $s.WaitForExit(); \
      mkdir C:\dig; \
      cp C:\bind\dig.exe C:\dig\; \
      cp C:\bind\nslookup.exe C:\dig\; \
      cp C:\bind\*.dll C:\dig\; \
      cp C:\Windows\System32\vcruntime140.dll C:\dig\; \
      rm C:\bind.zip;"

# httpd
ENV chocolateyUseWindowsCompression false
RUN powershell -Command "\
    iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1')); \
    choco feature disable --name showDownloadProgress;"

RUN choco install apache-httpd --version 2.4.39 -y &\
    choco install php --version 7.2.14 -y --params '"/ThreadSafe "'

# echoserver
RUN choco install openssl.light -y

FROM mcr.microsoft.com/windows/nanoserver:1809

COPY --from=prep /bin /bin
COPY --from=prep /curl/AMD64 /curl/CURL.LIC /curl/
COPY --from=prep /dig /dig
COPY --from=prep /netcat /netcat
COPY --from=prep /PowerShell /PowerShell

# httpd
COPY --from=prep /Windows/System32/vcruntime140.dll /Windows/System32/
COPY --from=prep /Users/ContainerAdministrator/AppData/Roaming/Apache24 /Users/ContainerAdministrator/AppData/Roaming/Apache24
COPY --from=prep /tools/php72 /tools/php72

# echoserver
COPY --from=prep /Windows/System32/vcruntime140.dll /Windows/System32/
COPY --from=prep ["C:/Program Files/OpenSSL", "C:/Program Files/OpenSSL"]

# nginx
#COPY --from=prep /Windows/System32/sync.exe /Windows/System32/SyncHost.exe

# powershell
# set a fixed location for the Module analysis cache
ENV LOCALAPPDATA="C:\Users\ContainerAdministrator\AppData\Local" \
    PSModuleAnalysisCachePath="C:\Users\ContainerAdministrator\AppData\Local\Microsoft\Windows\PowerShell\docker\ModuleAnalysisCache" \
    # Persist %PSCORE% ENV variable for user convenience
    PSCORE="C:\PowerShell\pwsh.exe"

# use downloaded powershell
USER ContainerAdministrator
RUN setx /M PATH "C:\Powershell\;%PATH%" &&\
    echo "" > C:\agnhost.exe &&\
    mklink agnhost-2.exe agnhost.exe &&\
    del agnhost.exe

# intialize powershell module cache
RUN powershell \
        -NoLogo \
        -NoProfile \
        -Command " \
          $stopTime = (get-date).AddMinutes(15); \
          $ErrorActionPreference = 'Stop' ; \
          $ProgressPreference = 'SilentlyContinue' ; \
          while(!(Test-Path -Path $env:PSModuleAnalysisCachePath)) {  \
            Write-Host "'Waiting for $env:PSModuleAnalysisCachePath'" ; \
            if((get-date) -gt $stopTime) { throw 'timout expired'} \
            Start-Sleep -Seconds 6 ; \
          }"
