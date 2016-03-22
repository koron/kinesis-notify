@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION

SET CC=com.amazonaws.services.kinesis.multilang.MultiLangDaemon
SET CP=.

FOR %%f IN (lib\*.jar) DO SET CP=!CP!;%%f

java -cp %CP% %CC% kinesis-notify.properties
