<div align="center">
<img src = "/assets/landing.png">
</div>

# Brolang

Brolang is a fun programming language written in Golang built for fun.

- Beta version is live at [brolang.ankushsingh.tech](https://brolang.ankushsingh.tech)

## 🚀 Getting Started

You can get started with BROLANG by visiting [brolang.ankushsingh.tech](https://brolang.ankushsingh.tech)

- To set up locally, you can clone the repository and run the following commands:

```bash
git clone https://github.com/ankush-web-eng/brolang
go mod download
go build main.go
go run main.go
```

### Using Docker

You can also run the project using Docker. Follow the steps below:

1. Pull the pre-built Docker image from Docker Hub:

```bash
docker pull deshwalankush23/brolang
```

2. Run the Docker container:

```bash
docker run -d -p 8080:8080 deshwalankush23/brolang
```

This command will run the container and map port 8080 of the container to port 8080 on your local machine.

### Architecture

- BROLANG was primarily built on an **Event-Based Architecture**
- The code has been commented in main.go and api/handler/compiler_handler.go which belonfgs to this architecture.
- On production, this application is working on Client-Server architecture because of this independent student-developer's tight budget, moreover the code has been commented in context/CodeContext.tsx.

### Open Source

- The code itself is self-explanatory and has been well documented in Go-Doc format.
- Contributions, issues and feature requests are welcome. Feel free to check the [issues page](/issues) if you want to contribute.
- I am open for all suggestions related to syntax, errors and tech-stack, your opinion is valuable.
- Either you can contribute directly or leave a quick DM [here](https://x.com/whyankush07)

## 📝 License

BROLANG is licensed under the MIT License. See [LICENSE](LICENSE) for more information.