ECHO OFF
ECHO Launching ElasticSearch...
set ES_JAVA_HOME=%~dp0\JSTOR\jdk-11.0.13\
START /B %~dp0\JSTOR\elasticsearch-7.15.2\bin\elasticsearch.bat
TIMEOUT /t 60
START %~dp0\JSTOR\api.exe serve