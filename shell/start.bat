ECHO OFF
ECHO Launching ElasticSearch...
set "PATH=%~dp0"
set ES_JAVA_HOME=%PATH%\JSTOR\jdk-11.0.13\
START /B %PATH%\JSTOR\elasticsearch-7.15.2\bin\elasticsearch.bat
TIMEOUT /t 60
START %PATH%\JSTOR\api-windows.exe serve