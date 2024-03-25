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
    <div className="relative overflow-auto w-[80vw]">
      <section className="flex justify-between mx-8 bg-black rounded-[20px] mt-2">
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
            <p className="text-[#808081] text-xl font-light">
              Ingest billions of log units  from an arbitrary number of clients in just 
              a fraction of a second, while enjoying fast analytical queries running !
            </p>
          </div>
        </div>
        <img
          alt="img_for_homepage_banner"
          src="hlog_home_img.png"
          width={350}
          height={350}
          className="mr-2"
        />
      </section>
      <section className="p-8">
        <p className="mb-4">Resources</p>
        <section className="flex gap-2 items-top">
          <article className="bg-[#F3F4F6] rounded-lg w-[15rem] h-[15rem] p-4 flex flex-col gap-2 justify-between">
            <p className="text-4xl text-ellipsis line-clamp-3 h-[60%]">Yet another logging system! But why?</p>
            <p className="text-[#818C99] text-sm line-clamp-3">Get familiar with the core components required to run effectively an instance.</p>
          </article>
          <article className="bg-[#F3F4F6] rounded-lg w-[15rem] h-[15rem] p-4 flex flex-col gap-2 justify-between">
            <p className="text-4xl text-ellipsis line-clamp-3 h-[60%]">Get started in a minute.</p>
            <p className="text-[#818C99] text-sm line-clamp-3">Get familiar with the core components required to run effectively an instance.</p>
          </article>
          <article className="bg-[#F3F4F6] rounded-lg w-[15rem] h-[15rem] p-4 flex flex-col gap-2 justify-between">
            <p className="text-4xl text-ellipsis line-clamp-3 h-[60%]">Define your first metrics.</p>
            <p className="text-[#818C99] text-sm line-clamp-3">Your data have a semi-structured format? Guess what? You can query them. Do you see where I'm going? Let's get started with metrics!</p>
          </article>
        </section>
      </section>
    </div>
  );
}
