import { useState } from 'react'
import TaskList from './components/TaskList'

function App() {
  const [tasks, setTasks] = useState([])

  ])

  return (
    <div className='h-screen bg-slate-800'>
      <div className="max-w-4xl mx-auto pt-20 px-6">
        <h1 className='text-3xl text-slate-400 mb-4'>Todo-List</h1>
        <form className='mb-10'>
          <label htmlFor='todo-item' className='text-slate-50'>Tasks</label>
          <input type='text' className='mt-1 block w-full rounded'></input>
          <button className='mt-4 py-2 px-2 bg-slate-50 rounded min-w-[115px]'>Add a task</button>
        </form>
        <TaskList setTasks={setTasks} tasks={tasks}/>
      </div>
    </div>
  )
}

export default App
