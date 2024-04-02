import type { MetaFunction } from "@remix-run/node";
import { useEffect, useState } from "react";
import { Bar, Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ChartData,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

type GenObs = {
  channels_count: number[];
  logs_per_channel: { [key: string]: number };
  logs_per_sender: { [key: string]: number };
  logs_per_level: { [key: string ]: number };
  senders_count: number[];
  levels_count: number[];
  total_ingested_logs: number[];
  throughput_per_time: number[];
}

export default function Dashboard() {
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

  useEffect(() => {
    let socket = new WebSocket("ws://localhost:1542/genericobservables");
    socket.onopen = () => {
      socket.send("connection")
    };
    socket.onmessage = (event) => {
      const _data = JSON.parse(event.data);
      setGenObs({
        ...genObs,
        channels_count: [...genObs.channels_count, _data.channels_count].slice(-100),
        senders_count: [...genObs.senders_count, _data.senders_count].slice(-100),
        levels_count: [...genObs.levels_count, _data.levels_count].slice(-100),
        total_ingested_logs: [...genObs?.total_ingested_logs, _data.total_ingested_logs].slice(-100),
        logs_per_level: _data.logs_per_level,
        logs_per_sender: _data.logs_per_sender,
        logs_per_channel: _data.logs_per_channel,
      })
    };

    return () => {
      socket.close();
    }
  })

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

  const bar_options = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: false
      },
    },
  }

  // Those are options for the chartjs chart that is displayed along with
  // the number of total ingested logs.
  const total_ingest_logs_labels = Array.from(Array(genObs.total_ingested_logs.length).keys())
  const total_ingested_logs_data = {
    total_ingest_logs_labels,
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
  const logs_per_level_labels = Object.keys(genObs.logs_per_level)
  const logs_per_level_data: ChartData<"bar", number[], unknown> = {
    labels: logs_per_level_labels,
    datasets: [
      {
        label: "Logs per level",
        data: Object.values(genObs.logs_per_level),
        borderColor: '#1C65F4',
        backgroundColor: '#1C65F4',
        borderWidth: 1,
        borderRadius: 5,
      }
    ],
  }
  const logs_per_sender_labels = Object.keys(genObs.logs_per_sender);
  const logs_per_sender_data: ChartData<"bar", number[], unknown> = {
    labels: logs_per_sender_labels,
    datasets: [
      {
        label: "Logs per level",
        data: Object.values(genObs.logs_per_sender),
        borderColor: '#1C65F4',
        backgroundColor: '#1C65F4',
        borderWidth: 1,
        borderRadius: 5,
      }
    ],
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
                data={total_ingested_logs_data}
                />
            </div>
          </div>
        </article>
      </section>
      <section className="flex gap-2 ">
        {/* logs_per_channel*/}
        <article className="text-white bg-black p-3 rounded-lg h-[20rem] flex-1 flex flex-col gap-4 justify-between">
          <p className="text-[#86898D] text-sm">Logs per levels</p>
          <div>
            <Bar 
              options={bar_options} 
              data={logs_per_level_data}
            />
          </div>
        </article>
        {/* logs_per_sender */}
        <article className="text-white bg-black p-3 rounded-lg h-[20rem] flex-1 flex flex-col gap-4 justify-between">
          <p className="text-[#86898D] text-sm">Logs per senders</p>
          <div>
            <Bar 
              options={bar_options} 
              data={logs_per_sender_data}
            />
          </div>
        </article>
        {/* logs_per_level */}

      </section>
    </div>
  )
}