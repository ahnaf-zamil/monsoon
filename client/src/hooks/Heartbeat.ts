import type { IWebSocketDispatch, IHeartbeatInit } from "../ws/types";
import { OPCODES } from "../ws/opcodes";
import { dispatchWebSocketEvent } from "../ws/events";
import { useEffect, useRef } from "react";
import { log } from "../utils";

export const useWSHeartbeat = (
    socket: WebSocket | null,
    onDisconnect: () => void,
) => {
    const intervalRef = useRef<NodeJS.Timeout | null>(null);
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);
    const cleanup = () => {
        if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
        }
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
        }
    };

    useEffect(() => {
        if (!socket) return;

        const handleMessage = (e: MessageEvent) => {
            const payload: IWebSocketDispatch<IHeartbeatInit> = JSON.parse(
                e.data,
            );
            if (payload.opcode == OPCODES.HeartbeatInit) {
                log("info", "heartbeat started");

                // cleanup any old loop
                cleanup();

                // start heartbeat loop
                intervalRef.current = setInterval(() => {
                    if (!socket.OPEN) return;

                    dispatchWebSocketEvent(socket, OPCODES.Heartbeat, null);
                    log("debug", "heartbeat sent");

                    // only set timeout if not already set
                    if (!timeoutRef.current) {
                        timeoutRef.current = setTimeout(() => {
                            log(
                                "warn",
                                "no heartbeat ACK received within timeout",
                            );
                            cleanup();
                            onDisconnect();
                        }, payload.data.timeout);
                    }
                }, payload.data.interval);
            }

            if (payload.opcode == OPCODES.HeartbeatAck) {
                log("debug", "heartbeat ack received");

                if (timeoutRef.current) {
                    clearTimeout(timeoutRef.current);
                    timeoutRef.current = null;
                }
            }
        };

        socket.addEventListener("message", handleMessage);
        return () => {
            socket.removeEventListener("message", handleMessage);
            cleanup();
        };
    }, [socket]);
};
