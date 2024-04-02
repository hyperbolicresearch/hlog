import { useEffect, useState } from "react";
import type { MetaFunction } from "@remix-run/node";
import { createPortal } from "react-dom";
import { Line } from "react-chartjs-2";
import { json } from "@remix-run/node";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { useLoaderData } from "@remix-run/react";
import moment from 'moment';
import {
  DocumentIcon,
} from "@heroicons/react/24/outline"


export const meta: MetaFunction = () => {
  return [
    { title: "hlog | Home" },
    { name: "description", content: `
        Welcome to Hlog, an open source project developed by Hyperbolic Research
        aiming to provide a performant implementation of a log aggregation and
        management mechanism, leveraging existing technologies without 
        compromising the user experience.
      `
    },
  ];
};

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

type GenObs = {
  channels_count: number[];
  logs_per_channel: { [key: string]: number[] };
  logs_per_sender: { [key: string]: number[] };
  logs_per_level: { [key: string ]: number[] };
  senders_count: number[];
  levels_count: number[];
  total_ingested_logs: number[];
  throughput_per_time: number[];
}

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

export const loader = async () => {
  const url = "http://localhost:1542/liveinit"
  let response = await fetch(url)
  response = await response.json()
  return json(response)
}

export default function Index() {
  // GenObs (general observables) are the data needed to display the
  // high-level statistics of the logging system, including the number
  // of channels, senders, levels used so far and total ingested logs.
  const initial_gen_obs: GenObs = {
    logs_per_channel: {},
    logs_per_sender: {},
    logs_per_level: {},
    channels_count: [],
    senders_count: [],
    levels_count: [],
    total_ingested_logs: [],
    throughput_per_time: [],
  };
  const [genObs, setGenObs] = useState<GenObs>(initial_gen_obs);

  // logs are the preloaded logs that we fetched before
  // rendering the page for displaying purposes, before the live tailing
  // begins.
  const initial_loaded_logs = useLoaderData<typeof loader>();
  const inital_logs : LogT[] = initial_loaded_logs as unknown as LogT[];
  const [logs, setLogs] = useState<LogT[]>(inital_logs)

  // The following two states (displayModal and modalItem) are used to
  // set the display of the modal that renders the clicked log.
  const [displayModal, setDisplayModal] = useState<boolean>(false)
  const [modalItem, setModalItem] = useState<LogT>()
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

  useEffect(() => {
    let socket = new WebSocket("ws://localhost:1542/genericobservables");
    socket.onopen = () => {
      socket.send("connection")
    };
    socket.onmessage = (event) => {
      const _data = JSON.parse(event.data);
      setGenObs({
        ...genObs,
        channels_count: [...genObs.channels_count, _data.channels_count],
        senders_count: [...genObs.senders_count, _data.senders_count],
        levels_count: [...genObs.levels_count, _data.levels_count],
        total_ingested_logs: [...genObs?.total_ingested_logs, _data.total_ingested_logs],
      })
    };

    return () => {
      socket.close();
    }
  })

  // TODO here, upon opening, we connect to the WS server
  // to receive the newest logs. But there are several things we need
  // to work with:
  // Sometimes, we are not connected, even if we are on this page.
  // take care of reconnecting everytime this is open.
  useEffect(() => {
    let socket = new WebSocket("ws://localhost:1542/live");
    socket.onopen = () => {
      socket.send("Connection");
    };
    socket.onmessage = (event) => {
      const _data = JSON.parse(event.data);
      if (Object.keys(_data).length > 0 ) {
        setLogs((logs) => [_data, ...logs].slice(0, 100));
      }
    };

    return () => {
      socket.close();
    };
  }, [])

  // Those are options for the chartjs chart that is displayed along with
  // the number of total ingested logs.
  const line_options = {
    responsive: true,
    aspectRatio: 6,
    maintainAspectRatio: true,
    plugins: {
      legend: { display: false, },
      title: { display: false, }
    },
    elements: {
      point:{ radius: 0, }
    },
    scales: {
      x: { display: false },
      y: { 
        display: false,
      },
    }
  }
  const labels = Array.from(Array(genObs.total_ingested_logs.length).keys())
  const log_ingested_logs_data = {
    labels,
    datasets: [
      {
        label: "Log ingested logs",
        data: genObs.total_ingested_logs,
        borderColor: '#1C65F4',
        backgroundColor: '#1C65F4',
        borderWidth: 1,
      }
    ]
  }

  return (
    <div className="px-8 flex flex-col gap-2 overflow-auto w-full max-w-screen-xl">
      <section className="flex gap-2 overflow-auto h-[7rem] min-h-[7rem]">
        {/* channels_count */}
        <article className="text-white bg-black p-3 rounded-lg h-full w-[10rem] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Channels count</p>
          <p className="text-white text-5xl">
            {genObs.channels_count[genObs.channels_count.length - 1] || 0}
          </p>
        </article>
        {/* senders_count */}
        <article className="text-white bg-black p-3 rounded-lg h-full w-[10rem] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Senders count</p>
          <p className="text-white text-5xl">
            {genObs.senders_count[genObs.senders_count.length - 1] || 0}
          </p>
        </article>
        {/* levels_count */}
        <article className="text-white bg-black p-3 rounded-lg h-full w-[10rem] flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Levels count</p>
          <p className="text-white text-5xl">
            {genObs.levels_count[genObs.levels_count.length - 1] || 0}
          </p>
        </article>
        {/* total_ingested_logs */}
        <article className="text-white bg-black p-3 rounded-lg h-full flex-1 flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Total ingested logs</p>
          <div className="flex items-end justify-between h-auto">
              <p className="text-white text-5xl">
                {genObs.total_ingested_logs[genObs.total_ingested_logs.length - 1] || 0}
              </p>
            <div className="w-[60%]">
              <Line 
                options={line_options} 
                data={log_ingested_logs_data}
                />
            </div>
          </div>
        </article>
        {/* logs_per_channel*/}
        {/* logs_per_sender */}
        {/* logs_per_level */}
      </section>

      {displayModal && createPortal(
        <Modal onClick={onClickDisplayModal} log={modalItem}/>, 
        document.body
      )}

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

    </div>
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

