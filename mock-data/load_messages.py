"""Loads sample message data from 'messages.json'. Change up the variables here with the user and conversation IDs in your database"""

import json
from psycopg2.extras import execute_values
from dotenv import load_dotenv
import psycopg2
import random
import time
import os

load_dotenv()

# Define user IDs
user1_id = 1947592818079846400
user2_id = 1947592805706649600
conversation_id = 1947593581061492736


# Snowflake-like ID generator (64-bit integer with timestamp-based prefix)
def generate_snowflake():
    timestamp = int(time.time() * 1000)  # milliseconds
    random_bits = random.getrandbits(22)  # simulate lower bits
    return (timestamp << 22) | random_bits


# Base timestamp for the conversation
start_time = 1721802000


f = open("./messages.json", "r")
raw = f.read()
conversation = json.loads(
    raw.replace("{user1_id}", str(user1_id)).replace("{user2_id}", str(user2_id))
)
f.close()

# Assign message IDs
for i, message in enumerate(conversation):
    message["id"] = generate_snowflake()
    message["created_at"] = start_time + i * 20


conn = psycopg2.connect(os.environ["MESSAGE_DB_POSTGRES_URL"])

cursor = conn.cursor()

records = [
    (
        msg["id"],
        conversation_id,
        msg["author_id"],
        msg["content"],
        msg["created_at"],
        None,  # edited_at
        False,  # deleted
    )
    for msg in conversation
]

# SQL statement
sql = """
INSERT INTO messages (id, conversation_id, author_id, content, created_at, edited_at, deleted)
VALUES %s
"""

# Execute batch insert
execute_values(cursor, sql, records)

# Commit and close
conn.commit()
cursor.close()
conn.close()
