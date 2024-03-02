import { useState } from 'react'
import TaskList from './components/TaskList'

function App() {
  const [tasks, setTasks] = useState([])
  const [taskContent, setTaskContent] = useState('');

  const AddTask = async (event) => {
    event.preventDefault()
    if (taskContent.trim() === '') {
      console.error('Task content cannot be empty');
      return;
    }
    try {
      const response = await fetch('http://localhost:8080/tasks', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: taskContent }),
      });
      if (response.ok) {
        const newTask = await response.json();
        setTasks((prevTasks) => [...prevTasks, newTask]);
        setTaskContent('');
        console.log('Task created !');
      } else {
          const errorData = await response.json();
          const errorMessage = errorData.error || 'Unknown error occured';
          throw new Error(`HTTP error : ${errorMessage}`);
      };
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <div className='h-screen bg-slate-800'>
      <div className="max-w-4xl mx-auto pt-20 px-6">
        <h1 className='text-3xl text-slate-400 mb-4'>Todo-List</h1>
        <form onSubmit={e => AddTask(e)} className='mb-10'>
          <label htmlFor='todo-item' className='text-slate-50'>Tasks</label>
          <input value={taskContent} onChange={e => setTaskContent(e.target.value)} type='text' className='mt-1 block w-full rounded'/>
          <button className='mt-4 py-2 px-2 bg-slate-50 rounded min-w-[115px]'>Add a task</button>
        </form>
        <TaskList setTasks={setTasks} tasks={tasks}/>
      </div>
    </div>
  )
}

export default App
