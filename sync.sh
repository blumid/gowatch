#!/bin/bash


# hackerone
wget -O HackerOne.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/master/data/hackerone_data.json
jq 'map(.targets |= with_entries(if .key == "in_scope" or .key == "out_of_scope" then .value |= map(if .asset_identifier then .asset = .asset_identifier | del(.asset_identifier) else . end | if .asset_type then .type = .asset_type | del(.asset_type) else . end) else . end))' HackerOne.json > temp.json && mv temp.json HackerOne.json


# intigriti
wget -O Intigriti.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json
jq 'map(.targets |= with_entries(if .key == "in_scope" or .key == "out_of_scope" then .value |= map(if .endpoint then .asset = .endpoint | del(.endpoint) else . end) else . end))' Intigriti.json > temp.json && mv temp.json Intigriti.json


# BugCrowd.json
wget -O BugCrowd.json -A json https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json
jq 'map(.targets |= with_entries(if .key == "in_scope" or .key == "out_of_scope" then .value |= map(if .target then .asset = .target | del(.target) else . end | if .type == "website" then .type = "url" else . end) else . end))' BugCrowd.json > temp.json && mv temp.json BugCrowd.json


exit