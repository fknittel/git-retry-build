# Chrome Enterprise Deployment manager AppEngine Service

Chrome Enterprise Deployment manager(a.k.a Poros) is a service which manages lab
configs of all enterprise lab instances.

Services maintained:

- Inventory (Go Service)
- FrontEnd App (React, Redux, TypeScript)

Project Vision & Design Document

- Vision :  go/chrome-enterprise-lab2.0-vision
- Design :  go/enterprise-deployment-manager-design-phase-1

## Backend

Running the server:

```sh
go run main.go
```

This will set up the backend server running on port `8800`

## Frontend

Running the frontend:

```sh
cd frontend
npm install
npm start
```

This will set up the frontend client running on port `3000` with an automatic
proxy to the backend server running on `8800`.  To view the UI, go to
[localhost:3000](http://localhost:3000)

Formatting:

```sh
npm run fix
```
