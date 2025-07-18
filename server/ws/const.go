package ws

import "time"

const HEARTBEAT_TIMEOUT = 10 * time.Second
const HEARTBEAT_CHECK_INTERVAL = 5 * time.Second // Interval for checking for dead clients
const HEARTBEAT_INTERVAL = 5 * time.Second       // Heartbeat interval sent to client

const (
	OpHeartbeat     EventOpCode = "heartbeat"      // Sent by client to heartbeat
	OpHeartbeatInit EventOpCode = "heartbeat_init" // Tells client to start sending heartbeat pings
	OpHeartbeatAck  EventOpCode = "heartbeat_ack"  // Heartbeat acknowledgement sent by server upon receiving heartbeat

	OpRoomSync EventOpCode = "room_sync" // Provides client with its rooms upon connection
)
