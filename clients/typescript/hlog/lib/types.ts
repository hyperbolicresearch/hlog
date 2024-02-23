/** Types definition for hlog */

export enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARN = 'WARN',
  ERROR = 'ERROR',
  FATAL = 'FATAL'
}

export type Config = {
  client_id: string;
  kafka_server: string;
  kafka_username: string;
  kafka_password: string;
  channel_id: string;
  default_level: LogLevel;
}

export type Log = {
  log_id: string;
  sender_id: string;
  timestamp: number;
  level: string;
  message: string;
  data: any;
}
