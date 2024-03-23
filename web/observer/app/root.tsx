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
  QueueListIcon,
  ChartBarSquareIcon,
  BoltIcon,
  LockClosedIcon,
  Cog6ToothIcon,
  RectangleStackIcon,
  BellIcon,
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
        <main className="h-screen w-screen flex">
          <SideBar />
          <div className="flex flex-col">
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
  return <Outlet />;
}

export function Header() {
  return (
    <div className="px-8 h-[5rem] flex flex-shrink-0 justify-between items-center w-[80vw]">
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
    {to: "/functions", title: "Functions", icon: <BoltIcon width={22}/>}
  ]

  const administrative_menu = [
    {to: "/settings", title: "Settings", icon: <Cog6ToothIcon width={22}/>},
    {to: "/admin", title: "Admin", icon: <LockClosedIcon width={22}/>},
  ]

  return (
    <div className="w-[20vw] px-6 py-4 h-full bg-[#F3F5F6]">
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
