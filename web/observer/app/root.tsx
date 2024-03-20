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
import { LinksFunction } from "@remix-run/node";
import {
  HomeIcon,
  ChartPieIcon,
  Bars3BottomLeftIcon,
  QueueListIcon,
  ChartBarSquareIcon,
  BoltIcon,
  LockClosedIcon,
  Cog6ToothIcon,
} from "@heroicons/react/24/outline"

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        <main className="h-screen">
          <Header />
          <div className="flex h-max">
            <SideBar />
            {children}
          </div>
        </main>
        <ScrollRestoration />
        <Scripts />
        <LiveReload />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

export function Header() {
  return (
    <div className="flex justify-between px-6 border-b-[1px] border-[#f2f2f2] h-[8vh] items-center">
      <div className="flex gap-2">
        <img 
          src="hlog_logo.png"
          alt="hlog logo"
          className="w-6 h-6"
          />
        <p className="font-semibold">
          Hlog
          <span className="font-normal"> by </span>
          <span>Hyperbolic Research</span>
        </p>
      </div>
    </div>
  )
}

export function SideBar() {
  const observability_menu = [
    {to: "home", title: "Home", icon: <HomeIcon width={22}/>},
    {to: "dashboard", title: "Dashboard", icon: <ChartPieIcon width={22}/>},
    {to: "live", title: "Live tail", icon: <QueueListIcon width={22}/>},
    {to: "metrics", title: "Metrics", icon: <ChartBarSquareIcon width={22}/>},
    {to: "functions", title: "Functions", icon: <BoltIcon width={22}/>}
  ]

  const administrative_menu = [
    {to: "settings", title: "Settings", icon: <Cog6ToothIcon width={22}/>},
    {to: "admin", title: "Admin", icon: <LockClosedIcon width={22}/>},
  ]

  return (
    <div className="w-64 px-6 py-4 border-r-[1px] h-[92vh]">
      <p className="text-[#808081] text-sm">Observability tools</p>
      <ul className="py-2">
        {
          observability_menu.map((menu, index) => (
            <li className="py-4">
              <NavLink to={menu.to} key={index} className="flex gap-2 pl-3">
                {menu.icon}
                <p className="text-sm">{menu.title}</p>
              </NavLink>
            </li>
          ))
        }
      </ul>
      <p className="text-[#808081] text-sm">Administrative tools</p>
      <ul className="py-2">
        {
          administrative_menu.map((menu, index) => (
            <li className="py-4">
              <NavLink to={menu.to} key={index} className="flex gap-2 pl-4">
                {menu.icon}
                <p className="text-sm">{menu.title}</p>
              </NavLink>
            </li>
          ))
        }
      </ul>
    </div>
  )
}

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: styles },
];
