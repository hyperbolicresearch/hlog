
export type GenObs = {
  channels_count: number[];
  logs_per_channel: { [key: string]: number[] };
  logs_per_sender: { [key: string]: number[] };
  logs_per_level: { [key: string ]: number[] };
  senders_count: number[];
  levels_count: number[];
  total_ingested_logs: number[];
  throughput_per_time: number[];
}

export type LogT = {
  channel: string
  log_id: string
  sender_id: string
  timestamp: number
  level: string
  message:string
  data: object
}

export type LogLevelCount = {
  debug: number;
  info: number;
  warn: number;
  error: number;
  fatal: number;
  [key: string]: number;
};

// panic: connection(localhost:27017[-3970]) incomplete read of message header: read tcp [::1]:57831->[::1]:27017: use of closed network connection