import { useEffect, useState } from "react";
import {
  DocumentIcon,
} from "@heroicons/react/24/outline"
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ChartData,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

type LogT = {
  channel: string
  log_id: string
  sender_id: string
  timestamp: number
  level: string
  message:string
  data: object
}

type LogLevelCount = {
  debug: number;
  info: number;
  warn: number;
  error: number;
  fatal: number;
  [key: string]: number;
};

export default function Live() {
  const inital_logs : LogT[] = [];
  const [logs, setLogs] = useState<LogT[]>(inital_logs)
  const level_map: LogLevelCount = {
    "debug": 0,
    "info": 0,
    "warn": 0,
    "error": 0,
    "fatal": 0,
  }
  useEffect(() => {
    let socket = new WebSocket("ws://localhost:1337");
    socket.onopen = () => {
      socket.send("Connection");
    };
    socket.onmessage = (event) => {
      const _data = JSON.parse(event.data);
      setLogs((logs) => [_data, ...logs]);
    };

    return () => {
      socket.close();
    };
  }, [])

  for (let i = 0; i < logs.length; i++) {
    const l: string = logs[i].level
    level_map[l] += 1
  }

  const options = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: false,
      },
    },
  };

  const data: ChartData<"bar", number[], unknown> = {
    labels: Object.keys(level_map),
    datasets: [
      {
        label: "Log level",
        data: Object.values(level_map),
        backgroundColor: '#1C65F4',
      }
    ],
  }

  return (
    <section className="px-8 flex flex-col gap-2 overflow-auto">
      <article className="bg-black p-4 rounded-lg flex justify-between">
        <section className="w-[25%] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Total loaded logs</p>
          <p className="text-white text-5xl">{logs.length}</p>
        </section>
        <section className=" w-[25%] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Channels count</p>
          <p className="text-white text-5xl">{new Set(logs.map(log => log.channel)).size}</p>
        </section>
        <section className="w-[25%] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Timeframe</p>
          <p></p>
          <p className="text-white text-2xl">{
            new Date(logs[logs.length - 1]?.timestamp * 1000).toLocaleTimeString() + 
            " - " + 
            new Date(logs[0]?.timestamp * 1000).toLocaleTimeString()}
          </p>
        </section>
        <section className="w-[22%] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Logs per levels</p>
          <Bar options={options} data={data} />
        </section>
      </article>
      <article className="bg-black px-4 py-3 rounded-lg flex gap-2 items-center sticky top-0">
        <DocumentIcon width={20} height={20} color="white" />
        <p className="font-semibold text-white text-sm w-[20%]">Date and time</p>
        <p className="font-semibold text-white text-sm w-[10%] line-clamp-1">Channel</p>
        <p className="font-semibold text-white text-sm w-[5%]  line-clamp-1">Level</p>
        <p className="font-semibold text-white text-sm w-[30%] line-clamp-1">Message</p>
        <p className="font-semibold text-white text-sm w-[30%] line-clamp-1">Data</p>
      </article>
      {
        logs.map(log => (
          <article 
            key={log.log_id} 
            className="bg-[#FAFAFA] px-4 py-3 rounded-lg flex gap-2 items-center cursor-pointer"
          >
            <DocumentIcon width={20} height={20} color="#86898D" />
            <p className="text-sm w-[20%] text-[#86898D]">{new Date(log.timestamp * 1000).toISOString()}</p>
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
