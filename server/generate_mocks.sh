#!/bin/bash

FILES=(
    "lib/jwt.go"
    "lib/pw_hash.go"
    "ws/nats.go"
    "db/app/user_db.go"
    "db/app/conversation_db.go"
)

DEST_PACKAGE="mocks"


for FILE in "${FILES[@]}"; do
    BASENAME=$(basename "$FILE")
    DEST_FILE="${DEST_PACKAGE}/mock_${BASENAME}"

    echo "Generating mock for $FILE -> $DEST_FILE"
    
    mockgen -source="$FILE" -destination="$DEST_FILE" -package=$DEST_PACKAGE
done

echo "Mocks generated successfully."