#!/usr/bin/env bash

curl "https://marianoflix.duckdns.org:8181/api/v2?apikey=${TAUTULI_API_KEY}&cmd=get_library_media_info&section_id=3&search=${KEYWORDS}" | jq .
