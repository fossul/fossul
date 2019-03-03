#/bin/sh
curl -f -k -H 'Content-Type: application/json' -XPOST -d '{"profile": "myprofile","appPlugin": "myapp", "storagePlugin": "mystorage", "backupRetentionDays": 10}' http://localhost:8001/quiesce