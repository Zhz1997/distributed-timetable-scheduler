# distributed-timetable-scheduler

This is a distribute timetable schedular for solving course scheduling problem.

Data is stored in data/data1 folder, for detail about the data please see the report under report folder.

To run the code:
    - cd src/main
    - go run tsmaster.go ../data/data1 (start the master, and pass the data to the master)
    - go run tsworker.go (start worker, you can start as many workers as you want. The more worker, the less time needed to generate 3000 generations)