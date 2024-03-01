import { useEffect, useState } from "react"
import DeleteForeverRoundedIcon from '@mui/icons-material/DeleteForeverRounded';
import ModeEditOutlineRoundedIcon from '@mui/icons-material/ModeEditOutlineRounded';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';

export default function TaskList({tasks, setTasks}) {
    const [editableTaskId, setEditableTaskId] = useState(null);
    const [editedContent, setEditedContent] = useState('');

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

    async function saveEditTask(id) {
        try {
            const response = await fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ content: editedContent })
            });
            if (response.ok) {
                fetch('http://localhost:8080/tasks', {
                    method: 'GET'
                })
                    .then(response => response.json())
                    .then(data => setTasks(data))
                    .catch(error => console.error('Error while getting tasks', error));
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.error || 'Unknown error occured';
                throw new Error(`HTTP error : ${errorMessage}`);
            }
        } catch (error) {
            console.error(error);
        }
        setEditableTaskId(null);
    };

    const handleEditClick = (id, content) => {
        setEditableTaskId(id);
        setEditedContent(content);
    };

    const handleInputKeyDown = (e, id) => {
        if (e.key === 'Enter') {
          handleSaveEdit(id);
        }
    };

    async function changeTaskState(task) {
        try {
            const response = await fetch(`http://localhost:8080/tasks/state/${task.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },})
            if (response.ok) {
                fetch('http://localhost:8080/tasks', {
                    method: 'GET'
                })
                    .then(response => response.json())
                    .then(data => setTasks(data))
                    .catch(error => console.error('Error while getting tasks', error))};
        } catch (error) {
            console.error(error);
        }
    }

    return (
        <ul>
            {tasks.length === 0 && (<li className="text-slate-50 text-md"> No task yet</li>)}
            {tasks.length > 0 &&
            tasks.map(task => (
                <li key={task.id} className={`p-2 bg-zinc-200 mb-2 rounded flex justify-between ${task.state ? "line-through" : ""}`}>
                    <input type="checkbox" className="rounded-full h-6 w-6 appearance-none border border-gray-700 checked:bg-gray-400 checked:border-transparent ml-2"
                    checked={task.state || false} onChange={() => changeTaskState(task)}/>
                    {editableTaskId === task.id ? (
                        <div className="flex">
                        <input className="bg-gray-100 border border-gray-600 text-gray-700 rounded-lg" type="text" value={editedContent}
                            onChange={(e) => setEditedContent(e.target.value)}
                            onKeyDown={(e) => handleInputKeyDown(e, task.id)}
                            />
                        <button className="ml-2 text-gray-800" onClick={() => saveEditTask(task.id)}><CheckCircleIcon fontSize="large"/></button>
                        <button className="ml-2 text-red-800" onClick={() => setEditableTaskId(null)}><CancelIcon fontSize="large"/></button>
                        </div>
                    ) : <div className="">{task.content}</div>}
                    <div className="flex mr-5">
                        <button className="w-6 h-6 rounded text-gray-800 flex items-center justify-center" onClick={() => handleEditClick(task.id, task.content)}>
                            <ModeEditOutlineRoundedIcon fontSize="large"/>
                        </button>
                        <button onClick={() => deleteTask(task.id)} className="ml-7 w-6 h-6 rounded text-red-800 flex items-center justify-center">
                            <DeleteForeverRoundedIcon fontSize="large"/>
                        </button>
                    </div>
                </li>
                ))}
        </ul>
    )
}