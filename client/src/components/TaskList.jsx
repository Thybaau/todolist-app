import { useEffect, useState } from "react"

export default function TaskList() {
    const [tasks, setTasks] = useState([])

    useEffect(() => {
        fetch('http://localhost:8080/tasks', {
            method: 'GET'
        })
            .then(response => response.json())
            .then(data => setTasks(data))
            .catch(error => console.error('Erreur lors de la récupération des tâches', error));
    }, [])
    return (
        <ul>
            {tasks.map(task => (
                <li key={task.id}>
                    {task.content} - État : {task.state}
                </li>
                ))}
        </ul>
    )
}