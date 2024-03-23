import {
  DocumentIcon,
} from "@heroicons/react/24/outline"

export default function Live() {
  const logs = [
    {channel: "transaction", log_id: "1234", sender_id: "client-1", timestamp: Date.now(), level: "info", message: "Some random message that is way longer than necessary", data: {foo: "foo", bar: "bar", baz: "baz"}},
    {channel: "transaction", log_id: "2345", sender_id: "client-2", timestamp: Date.now(), level: "debg", message: "Some random message", data: {foo: "some foo thing", bar: "some bar stuff", baz: "baz"}},
    {channel: "transaction", log_id: "3456", sender_id: "client-1", timestamp: Date.now(), level: "warn", message: "Some random message", data: {foo: "foo", bar: "bar", baz: "baz"}},
    {channel: "transaction", log_id: "4567", sender_id: "client-3", timestamp: Date.now(), level: "info", message: "Some random message", data: {foo: "foo", bar: "bar", baz: "baz"}},
  ]

  return (
    <section className="p-8 flex flex-col gap-2">
      {
        logs.map(log => (
          <article key={log.log_id} className="bg-[#FAFAFA] px-4 py-3 rounded-lg flex gap-4 items-center">
            <DocumentIcon width={20} height={20} color="#86898D" />
            <p className="text-sm w-[20%] text-[#86898D]">{new Date(log.timestamp).toISOString()}</p>
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
