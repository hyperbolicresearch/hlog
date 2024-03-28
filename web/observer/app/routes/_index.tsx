import { useEffect, useState } from "react";
import type { MetaFunction } from "@remix-run/node";
import { Line } from "react-chartjs-2";
// import 'chart.js/auto';
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

export const loader = async () => {
  return null
}

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

export default function Index() {
  const data = useLoaderData<typeof loader>();
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

  const line_options = {
    responsive: true,
    aspectRatio: 7,
    maintainAspectRatio: true,
    plugins: {
      legend: { display: false, },
      title: { display: false, }
    },
    elements: {
      point:{ radius: 1, }
    },
    scales: {
      x: { display: false },
      y: { 
        // display: false, 
        // beginAtZero: true,
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
      <section className="flex gap-2 overflow-auto h-[7rem]">
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
        {/* logs_per_channel*/}
        {/* total_ingested_logs */}
        <article className="text-white bg-black p-3 rounded-lg h-full flex-1 flex flex-col justify-between">
          <p className="text-[#86898D] text-sm">Total ingested logs</p>
          <div className="flex items-end justify-between gap-16 h-auto">
              <p className="text-white text-5xl w-[15%]">
                {genObs.total_ingested_logs[genObs.total_ingested_logs.length - 1] || 0}
              </p>
            <div className="flex-1">
              <Line 
                options={line_options} 
                data={log_ingested_logs_data}
              />
            </div>
          </div>
        </article>

        {/* logs_per_sender */}
        {/* logs_per_level */}
      </section>
    </div>
  )
}
