# 📝 Url Shortening App
 This project is a hands-on implementation of a Markdown Note-Taking App, as suggested in the 
 [Backend Developer Roadmap Project](https://roadmap.sh/projects/url-shortening-service).  



A lightweight web application to create, save, and preview Markdown notes as rendered HTML. Notes are stored in a PostgreSQL database and checked for spelling mistakes.

---

## ✨ Features

- Create, Edit and Delete Short versions of Urls
- See how many times a short url was accessed


---

## 🚀 Requirements

- [Go 1.20+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

---

## ⚙️ How to Run the Project

1. Start the PostgreSQL container:

   ```bash
   docker-compose up -d
2. Run the application:
cd cmd
go run main.go

3. Access http://localhost:8000/ and interact with the UI to perform all the actions

