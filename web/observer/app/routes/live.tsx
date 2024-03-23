import { useEffect, useState } from "react";
import io from "socket.io-client";

import {
  DocumentIcon,
} from "@heroicons/react/24/outline"

type LogT = {
  channel: string
  log_id: string
  sender_id: string
  timestamp: number
  level: string
  message:string
  data: object
}

export default function Live() {
  const inital_logs : LogT[] = [];
  const [logs, setLogs] = useState<LogT[]>(inital_logs)

  useEffect(() => {
    let socket = new WebSocket("ws://localhost:1337");
    socket.onopen = () => {
      socket.send("Connection");
    };
    socket.onmessage = (event) => {
      console.log(event.data);
      const _data = JSON.parse(event.data);
      setLogs((logs) => [_data, ...logs]);
    };

    return () => {
      socket.close();
    };
  }, [])

  return (
    <section className="p-8 flex flex-col gap-2 overflow-auto">
      {
        logs.map(log => (
          <article key={log.log_id} className="bg-[#FAFAFA] px-4 py-3 rounded-lg flex gap-4 items-center">
            <DocumentIcon width={20} height={20} color="#86898D" />
            <p className="text-sm w-[20%] text-[#86898D]">{new Date(log.timestamp).toISOString()}</p>
            <p className="text-sm w-[10%] line-clamp-1 font-medium">{log.channel}</p>
            <p className="text-sm w-[5%] text-[#1C65F4] font-medium">{log.level}</p>
            <p className="text-sm w-[30%] text-[#1E1E1E] line-clamp-1">{log.message}</p>
            <p className="text-sm w-[30%] text-[#5D5D5D] line-clamp-1">{JSON.stringify(log.data)}</p>
          </article> 
        ))
      }
    </section>
  )
}
