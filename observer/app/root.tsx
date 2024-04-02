import { useEffect, useState } from "react";
import styles from "./tailwind.css?url";
import {
  Links,
  LiveReload,
  Meta,
  NavLink,
  Outlet,
  Scripts,
  ScrollRestoration,
} from "@remix-run/react";
import { useLoaderData } from "@remix-run/react";
import { json, LinksFunction } from "@remix-run/node";
import {
  HomeIcon,
  ChartPieIcon,
  QueueListIcon,
  ChartBarSquareIcon,
  BoltIcon,
  LockClosedIcon,
  Cog6ToothIcon,
  RectangleStackIcon,
  BellIcon,
  BookOpenIcon,
} from "@heroicons/react/24/outline"
import { GenObs, LogT } from "~/utils/types";

export const loader = async () => {
  const url = "http://localhost:1542/liveinit"
  let response = await fetch(url)
  response = await response.json()
  return json(response)
}

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body className="flex justify-center">
        <main className="h-screen w-full flex self-center">
          <SideBar />
          <div className="flex flex-col w-full items-center">
            <Header />
            {children}
          </div>
        </main>
        <ScrollRestoration />
        <Scripts />
        {/* <LiveReload /> */}
      </body>
    </html>
  );
}

export default function App() {
  // logs are the preloaded logs that we fetched before
  // rendering the page for displaying purposes, before the live tailing
  // begins.
  const initial_loaded_logs = useLoaderData<typeof loader>();
  const inital_logs : LogT[] = initial_loaded_logs as unknown as LogT[];
  const [logs, setLogs] = useState<LogT[]>(inital_logs);

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
        channels_count: [...genObs.channels_count, _data.channels_count].slice(-50),
        senders_count: [...genObs.senders_count, _data.senders_count].slice(-50),
        levels_count: [...genObs.levels_count, _data.levels_count].slice(-50),
        total_ingested_logs: [...genObs?.total_ingested_logs, _data.total_ingested_logs].slice(-50),
        logs_per_level: _data.logs_per_level,
        logs_per_sender: _data.logs_per_sender,
        logs_per_channel: _data.logs_per_channel,
      })
    };

    return () => {
      socket.close();
    }
  })

  return <Outlet context={[ logs, genObs,]}/>;
}

export function Header() {
  return (
    <div className="px-8 h-[5rem] flex flex-shrink-0 justify-between items-center w-full">
      <div className="flex items-center gap-2">
        <img 
          alt="project-name"
          src="randname.jpg"
          className="w-10 h-10 rounded-lg"
        />
        <div>
          <p className="text-sm">Hypercluster</p>
          <p className="text-xs text-[#808081] font-light">dev-781227</p>
        </div>
      </div>
      <div className="flex gap-4 items-center">
        <BellIcon width={30} height={30} />
        <div className="bg-[#1C65F4] px-3 py-2 rounded-lg flex items-center gap-2">
          <RectangleStackIcon width={24} color="white" />
          <p className="text-sm text-white">Create channel</p>
        </div>
        <div className="flex items-center gap-2">
          <img 
            alt="profile-img"
            src="profile_img.jpg"
            className="w-10 h-10 rounded-full object-cover"
          />
          <div className="flex flex-col justify-center leading-tight">
            <p className="text-sm">Nahum Maurice</p>
            <p className="text-xs text-[#808081] font-light">hyperbolic@research.com</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export function SideBar() {
  const observability_menu = [
    {to: "/", title: "Home", icon: <HomeIcon width={22}/>},
    {to: "/dashboard", title: "Dashboard", icon: <ChartPieIcon width={22}/>},
    {to: "/live", title: "Live tail", icon: <QueueListIcon width={22}/>},
    {to: "/metrics", title: "Metrics", icon: <ChartBarSquareIcon width={22}/>},
    {to: "/functions", title: "Functions", icon: <BoltIcon width={22}/>},
    {to: "/resources", title: "Resources", icon: <BookOpenIcon width={22}/>},
  ]

  const administrative_menu = [
    {to: "/settings", title: "Settings", icon: <Cog6ToothIcon width={22}/>},
    {to: "/admin", title: "Admin", icon: <LockClosedIcon width={22}/>},
  ]

  return (
    <div className="w-[18rem] min-w-[18rem] px-6 py-4 h-full bg-[#F3F5F6]">
      <div className="flex gap-2 items-center mb-12">
        <img 
          src="hlog_logo.png"
          alt="hlog logo"
          className="w-8 h-8
          "
        />
        <p className="font-semibold">
          hlog
        </p>
      </div>
      <p className="text-[#808081] text-sm mb-4 font-extralight">Observability tools</p>
      {
        observability_menu.map((menu, index) => (
          <NavLink 
            to={menu.to} 
            key={index} 
            className={({isActive}) => 
              isActive 
                ? "flex gap-2 pl-3 py-3 bg-[#E8EAEF] rounded-lg font-light" 
                : "flex gap-2 pl-3 py-3 font-light"
            }
          >
            {menu.icon}
            <p className="text-sm">{menu.title}</p>
          </NavLink>
        ))
      }
      <p className="text-[#808081] text-sm my-4 font-extralight">Administrative tools</p>
      {
        administrative_menu.map((menu, index) => (
          <NavLink 
            to={menu.to} 
            key={index} 
            className={({isActive}) => 
              isActive 
                ? "flex gap-2 pl-3 py-3 bg-[#E8EAEF] rounded-lg font-light" 
                : "flex gap-2 pl-3 py-3 font-light"
            }
          >
            {menu.icon}
            <p className="text-sm">{menu.title}</p>
          </NavLink>
        ))
      }
    </div>
  )
}

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: styles },
];
