# Todo-App

This is a simple todo app created using HTMX and golang

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Sumit-Kumar1_Todo-App&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Sumit-Kumar1_Todo-App)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=Sumit-Kumar1_Todo-App&metric=bugs)](https://sonarcloud.io/summary/new_code?id=Sumit-Kumar1_Todo-App)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=Sumit-Kumar1_Todo-App&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=Sumit-Kumar1_Todo-App)[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=Sumit-Kumar1_Todo-App&metric=coverage)](https://sonarcloud.io/summary/new_code?id=Sumit-Kumar1_Todo-App)

## Steps to run the project

- Clone the repo & make sure you install `go v1.23 or newer` and `make`  
- cd `Todo-App` and open a terminal in the directory
- In terminal run `make help` : will list you all available make commands
- or run directly the `make run`
- Open a browser and goto address `localhost:12344`
- Done !! Enjoy adding tasks

## API Specification

- Todo api specification can be found at `openapi/todoApi.yaml` (WIP)

## Requirements

 User Should Be able to do:

- [x] Add new todos to list
- [x] Mark todo as complete
- [x] Delete todos from list
- [ ] Filter by all/activet/complete todos
- [ ] Clear all completed todos
- [ ] Toggle light and dark mode
- [ ] View the optimal layout for the app depending on device screen size
- [ ] See hover states for all interactive elements on the page
- [ ]  **Bonus**: Drag and drop to reorder items on list
- [ ]  **Bonus**: Build this project as a full-stack application