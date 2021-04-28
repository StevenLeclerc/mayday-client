### MayDay-Client v1.2.5

... is the client used by system ops to send log to mayday-core

-----

Mayday client will send log every 2 seconds, or every 1000 messages logged.

The client will try to open a config.json under is path. (not recursive)

The config.json must be formatted like that :

```json
{
  "apiKey": "12345abcdss",
  "defaultHostname": "Tester",
  "debug": true,
  "serverUrl": "localhost:4545",
  "logConfigs": [
    {
      "logFilePath": "/var/log/bhrick3.log",
      "channels": [
        "5f82148bc21f9ac0bb6ae822"
      ],
      "logAllFile": true
    },
    {
      "logFilePath": "/tmp/nick-fury.log",
      "channels": [
        "5f82148bc21f97644b6ae98s"
      ],
      "logAllFile": true
    }
  ]
}

```

The `apiKey` is provided within mayday frontend, under the logFetcher creation process. Mayday-Client will use the ops
hostname, but ou can force one by setting the field: `defaultHostname`
The `serverUrl` is by default to mayday official servers, but you can set your own mayday backend url. The `debug` mode
set to false by default.
`logConfigs` is an array of config, they should indicate the absolute path to the log file (`logFilePath`) and within
which
`channels` the log should be inserted. The `channels` are created within the mayday interface. At last, you can decide
to send all the content of the log file by setting true to `logAllFile`. Otherwise, only new inserted line will be sent.
