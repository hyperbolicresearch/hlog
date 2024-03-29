import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "hlog | Functions" },
    { name: "description", content: `
        Welcome to Hlog, an open source project developed by Hyperbolic Research
        aiming to provide a performant implementation of a log aggregation and
        management mechanism, leveraging existing technologies without compromising
        the user experience.
      `
    },
  ];
};

export default function Functions() {
    return <></>
}