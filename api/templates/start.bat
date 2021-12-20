ECHO OFF
ECHO Launching ElasticSearch...
set JAVA_HOME=%~dp0{{.JavaHome}}
START /B %~dp0{{.ESRoot}}bin\elasticsearch.bat
TIMEOUT /t 30
START api.exe serve