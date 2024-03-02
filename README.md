# Todolist App

![Go version](https://img.shields.io/badge/Go-1.18-blue)

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-project">About Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
  </ol>
</details>

## About Project
Simple To-do list app, to order your tasks.
Using PostgreSQL as database, Golang for back-end part, and ReactJS/Javascript for front-end. Each services are containerized with Docker and run with Docker Compose.

### Built with

* [![Golang][Go]][Go-url]
* [![React][React.js]][React-url]
* [![Postgre][Postgre]][Postgre-url]
* [![Docker][Docker]][Docker-url]

## Getting Started

### Prerequisites

You need to have Docker and Docker Compose installed on your computer to run project services.

### Installation

## Usage

In project root, do the following command :

```shell
docker compose up -d
```

App should now correctly running. Remove the `-d` if you want to see the logs in your terminal.

The project is composed of three services : `database`, `server` and `client`. If you want to see the logs in realtime of a specifical service, write this command and replace `<service_name>` with one of the three services names :

```shell
docker compose logs -f <service_name>
```

If you want to enter in one of the containers :

```shell
docker compose exec <service_name> sh
```
### Functionalities

* Real time tasks list.
* Add new task, settings his content.
* Delete task, by clicking on red button to the right of each tasks.
* Edit task, by clicking on the pencil button to the right of each tasks.
* Check/uncheck task with checkbox button. This will validate and cross task.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/github_username/repo_name.svg?style=for-the-badge
[contributors-url]: https://github.com/github_username/repo_name/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/github_username/repo_name.svg?style=for-the-badge
[forks-url]: https://github.com/github_username/repo_name/network/members
[stars-shield]: https://img.shields.io/github/stars/github_username/repo_name.svg?style=for-the-badge
[stars-url]: https://github.com/github_username/repo_name/stargazers
[issues-shield]: https://img.shields.io/github/issues/github_username/repo_name.svg?style=for-the-badge
[issues-url]: https://github.com/github_username/repo_name/issues
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/linkedin_username
[product-screenshot]: images/screenshot.png

[Go]:     https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/

[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/

[Docker]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/

[Postgre]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[Postgre-url]: https://www.postgresql.org/
