import React from "react"

export default function App() {
  const [count, setCount] = React.useState(0)
  const increment = () => setCount(count + 1)
  return (
    <div>
      <h1>Count: {count}</h1>
      <button onClick={increment} className="border border-gray-800 rounded p-2 bg-cyan-500 hover:bg-cyan-400">Increment</button>
    </div>
  )
}
