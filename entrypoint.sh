#!/bin/sh
if [ "$ISBOOTSTRAP" = "TRUE" ]; then
  export PORT=9001
else
  # Extract replica number from hostname (peer.1, peer.2, ...)
  NUM=$(echo $HOSTNAME | awk -F. '{print $2}')
  export PORT=$((9001 + NUM))
fi
exec "$@"