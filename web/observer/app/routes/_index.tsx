import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "Hlog by Hyperbolic Research" },
    { name: "description", content: `
        Welcome to Hlog, an open source project developed by Hyperbolic Research
        aiming to provide a performant implementation of a log aggregation and
        management mechanism, leveraging existing technologies without compromising
        the user experience.
      `
    },
  ];
};

export default function Index() {
  return (
    <div className="flex justify-between mx-8 bg-black rounded-[20px] mt-6">
      <div className="p-8 flex flex-col justify-between">
        <div className="flex items-center">
          <img 
            src="hlog_logo.png"
            alt="hlog logo"
            className="w-12 h-12"
          />
          <p className="font-semibold text-white">
            hlog
          </p>
        </div>
        <div className="flex flex-col gap-5">
          <h1 className="text-white text-6xl">
            Highly performant, scalable logging system.
          </h1>
          <p className="text-[#808081] text-xl">
            Ingest billions of log units  from an arbitrary number of clients in just a fraction of a second, while enjoying fast analytical queries running !</p>
        </div>
      </div>
      <img
        alt="img_for_homepage_banner"
        src="hlog_home_img.png"
        width={400}
        height={400}
        className="mr-2"
      />
    </div>
  );
}
