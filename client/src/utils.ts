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

export const formatUnixToLocalTime = (unixSeconds: number): string => {
    const date = new Date(unixSeconds * 1000); // Convert to milliseconds
    let hours = date.getHours();
    const minutes = date.getMinutes();

    const ampm = hours >= 12 ? "PM" : "AM";
    hours = hours % 12;
    if (hours === 0) hours = 12; // Midnight or Noon

    const paddedMinutes = minutes.toString().padStart(2, "0");

    return `${hours}:${paddedMinutes} ${ampm}`;
};

export const isEmptyString = (inp: string) => {
    return !/\S/.test(inp);
};
