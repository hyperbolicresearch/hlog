import { createPortal } from "react-dom";
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
import moment from 'moment';

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
  const [displayModal, setDisplayModal] = useState<boolean>(false)
  const [modalItem, setModalItem] = useState<LogT>()

  // TODO here, upon opening, we connect to the WS server
  // to receive the newest logs. But there are several things we need
  // to work with:
  // Sometimes, we are not connected, even if we are on this page.
  // take care of reconnecting everytime this is open.
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

  const onClickDisplayModal = (
    item: LogT | undefined
  ) => {
    if (displayModal === false) {
      if (item) {
        setModalItem(item)
      }
    }
    setDisplayModal(!displayModal);
  }

  const timeAgo = (interval: number) : string => {
    let str = "("
    const days = Math.floor(interval / 84600)
    let remaining = Math.floor(interval % 84600)
    const hours = Math.floor(remaining / 3600)
    remaining = Math.floor(remaining % 3600)
    const minutes = Math.floor(remaining / 60)
    remaining = Math.floor(remaining % 60)
    const seconds = Math.floor(remaining)

    if (days != 0) { str += `${days} days `}
    if (hours != 0) { str += `${hours} hours `}
    if (minutes != 0) { str += `${minutes} minutes `}
    if (seconds != 0) { str += `${seconds} seconds`} else {
      str += "just now"
    }
    return str + ")"
  }

  return (
    <section className="px-8 flex flex-col gap-2 overflow-auto w-full max-w-screen-xl">
      {displayModal && createPortal(
        <Modal onClick={onClickDisplayModal} log={modalItem}/>, 
        document.body
      )}
      <article className="bg-black p-4 rounded-[20px] flex justify-between">
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
          <div>
            <p className="text-[#86898D] text-sm ">
              {timeAgo(logs[0]?.timestamp - logs[logs.length - 1]?.timestamp)}
            </p>
            <p className="text-white text-2xl">{
              new Date(logs[logs.length - 1]?.timestamp * 1000).toLocaleTimeString() + 
              " - " + 
              new Date(logs[0]?.timestamp * 1000).toLocaleTimeString()}
            </p>
          </div>
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
            onClick={() => onClickDisplayModal(log)}
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

type ModalProps = {
  log: LogT | undefined
  onClick: (
    item: LogT | undefined
  ) => void
}

const Modal: React.FC<ModalProps> = (props) : JSX.Element => {
  return (
    <div 
      className="z-50 bg-black bg-opacity-50 h-[100%] w-[100%] absolute top-0 left-0 flex justify-end"
      onClick={(e) => props.onClick(props.log)}
    >
      <div 
        className="h-[100%] w-[50%] bg-white p-4 overflow-auto"
        onClick={e => { e.stopPropagation(); }}
      >
        <div className="flex justify-between gap-4">
          <div className="flex gap-8">
            <div className="gap-2 items-center">
              <p className="text-sm text-[#5D5D5D]">Log ID</p>
              <p className="text-sm">{props.log?.log_id}</p>
            </div>
            <div className="gap-2 items-center">
              <p className="text-sm text-[#5D5D5D]">Channel</p>
              <p className="text-sm">{props.log?.channel}</p>
            </div>
            <div className="gap-2 items-center">
              <p className="text-sm text-[#5D5D5D]">From</p>
              <p className="text-sm">
                {props.log && moment(props.log?.timestamp * 1000).fromNow()}
              </p>
            </div>
          </div>
          <div className="bg-black py-2 px-6 rounded-md">
            <p className="text-white">{props.log?.level}</p>
          </div>
        </div>
        <div className="pt-4">
          <p className="text-sm text-[#5D5D5D]">Message</p>
          <p className="text-4xl">{props.log?.message}</p>
        </div>
        <div className="pt-4">
          <p className="text-sm text-[#5D5D5D]">Sender ID</p>
          <p>{props.log?.sender_id}</p>
        </div>
        <div className="pt-4">
          <p className="text-sm text-[#5D5D5D]">Data</p>
          <pre className="p-4 text-sm bg-[#FAFAFA] mt-2 rounded-md overflow-auto">
            {JSON.stringify(props.log?.data, null, 2)}
          </pre>
        </div>
      </div>
    </div>
  )
}
