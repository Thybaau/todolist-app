# todolist-app
To-do list app, using PostgreSQL as database, Golang for back-end part, and ReactJS/Javascript for front-end. 

Everything is packaged with Docker

### RUNNING DATABASE with dockerfile : 
Go in same repository as dockerfile and write : docker build -t [name_of_your_image] .
After that, do docker run --name postgres-test -p 5432:5432 -d [image_name]
