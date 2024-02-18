import { useEffect, useState } from "react"

export default function TaskList({tasks, setTasks}) {

    useEffect(() => {
        fetch('http://localhost:8080/tasks', {
            method: 'GET'
        })
            .then(response => response.json())
            .then(data => setTasks(data))
            .catch(error => console.error('Error while getting tasks', error));
    }, [])

    function deleteTask(id){
        fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    const errorMessage = errorData.detail || errorData.error || 'Unknown error occurred';
                    throw new Error(`HTTP error status: ${response.status}, Message: ${errorMessage}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log(data);
            setTasks(tasks.filter(task => task.id !== id));
        })
        .catch(error => {
            console.error('Error while deleting task : ', error.message);
        });
    }

    return (
        <ul>
            {tasks.length === 0 && (<li className="text-slate-50 text-md"> No task yet</li>)}
            {tasks.length > 0 &&
            tasks.map(task => (
                <li key={task.id} className="p-2 bg-zinc-200 mb-2 rounded flex">
                    {task.content}
                    <button onClick={() => deleteTask(task.id)} className="ml-auto bg-red-600 w-6 h-6 rounded text-zinc-200">X</button>
                </li>
                ))}
        </ul>
    )
}