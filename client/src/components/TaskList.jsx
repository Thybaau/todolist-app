import { useEffect, useState } from "react"

export default function TaskList() {
    const [tasks, setTasks] = useState([])

    useEffect(() => {
        fetch('http://localhost:8080/tasks', {
            method: 'GET'
        })
            .then(response => response.json())
            .then(data => setTasks(data))
            .catch(error => console.error('Error while getting tasks', error));
    }, [])
    return (
        <ul>
            {tasks.map(task => (
                <li key={task.id} className="p-2 bg-zinc-200 mb-2 rounded flex">
                    {task.content}
                    <button className="ml-auto bg-red-600 w-6 h-6 rounded text-zinc-200">X</button>
                </li>
                ))}
        </ul>
    )
}