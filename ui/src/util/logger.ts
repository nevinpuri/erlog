import winston from "winston";

class Logger {
  private logger = winston.createLogger({
    transports: [
      new winston.transports.Http({
        host: "localhost",
        port: 8000,
      }),
    ],
  });

  constructor() {}

  public info(data: object) {
    this.logger.info("fkas", data);
  }
}

export function getLogger() {
  return new Logger();
}

let logUrl: string;
export function init(_logUrl: string) {
  if (_logUrl.endsWith("/")) {
    logUrl = _logUrl.slice(0, _logUrl.length - 1);
  } else {
    logUrl = _logUrl;
  }
}

/**
 * You would use it like
 * report("post_created", {userId: 3})
 * @param args
 */
export function report(
  name: string,
  args: object,
  throwError: boolean = false
) {
  // just queue it up or some thing or just send the post request asynchronously
  // return a promise, you can await it if you want
  return new Promise<void>(async (resolve, reject) => {
    if (!logUrl) {
      return reject("Error: please initialize logger");
    }
    try {
      await fetch(logUrl + "/metrics", {
        method: "POST",
        body: JSON.stringify({
          name,
          args,
        }),
      });
    } catch (exception) {
      if (throwError === true) {
        return reject(exception);
      }
      return resolve();
    }
    return resolve();
  });
}

/**
 * Supports all the parent_id shit. Main feature, creating metrics in next js apps in the api route
 */
export function log() {}
