# D7024E Lab Assignment: Creating a Peer-to-Peer Distributed Data Store
This description describes how to set up and run our implementation of a Kademlia-based Distributed Data Store for the D7024E ([Mobile and Distributed Computing Systems](https://www.ltu.se/en/education/syllabuses/course-syllabus?id=D7024E)) course at [Luleå Technical University](https://www.ltu.se/).

## Requirements
- Go 1.25+
- Docker 4.45+

## Setup
1. Clone the repository files to your local machine using the terminal command:
   ```bash
   git clone https://github.com/antonfollinger/kademlia-go
   ```

2. Open the project and run the following command to download dependencies:
   ```bash
   go mod tidy
   ```

## Testing
Run the <i>runTests.sh</i> script located inside the Test folder to run a complete coverage test and generate an accompanying coverage report.
   ```bash
chmod +x runTests.sh
   ```

   ```bash
./runTests.sh
   ```