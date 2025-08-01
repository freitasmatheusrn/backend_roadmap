# 📝 Markdown Note App
 This project is a hands-on implementation of a Markdown Note-Taking App, as suggested in the 
 [Backend Developer Roadmap Project](https://roadmap.sh/projects/markdown-note-taking-app).  



A lightweight web application to create, save, and preview Markdown notes as rendered HTML. Notes are stored in a PostgreSQL database and checked for spelling mistakes.

---

## ✨ Features

- Create and save notes written in Markdown format
- Live HTML rendering of saved notes
- Spell-checking via LanguageTool API
- Notes stored in a PostgreSQL database

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

## 📌 Future Improvements
🔄 Update existing notes (coming soon)

❌ Delete notes (coming soon)

🔍 Search/filter by note name or content (planned)