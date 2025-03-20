#!/bin/bash
# Wait for Elasticsearch to be ready
until curl -u elastic:"$ELASTIC_PASSWORD" -s "http://localhost:9200/_cluster/health?wait_for_status=yellow" >/dev/null; do
  echo "Waiting for Elasticsearch..."
  sleep 5
done

# Create the Kibana user using the provided environment variables
curl -u elastic:"$ELASTIC_PASSWORD" -X POST "http://localhost:9200/_security/user/$KIBANA_USER" -H 'Content-Type: application/json' -d "{
  \"password\": \"$KIBANA_PASSWORD\",
  \"roles\": [\"kibana_system\"]
}"
