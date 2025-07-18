type LogLevel = "info" | "warn" | "error" | "debug" | "log";

export function log(level: LogLevel, ...args: any[]) {
    const timestamp = new Date().toISOString();
    const prefix = `[${timestamp}] [${level.toUpperCase()}]`;

    switch (level) {
        case "info":
            console.info(prefix, ...args);
            break;
        case "warn":
            console.warn(prefix, ...args);
            break;
        case "error":
            console.error(prefix, ...args);
            break;
        case "debug":
            console.debug(prefix, ...args);
            break;
        default:
            console.log(prefix, ...args);
    }
}
