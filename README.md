## Command Line Quiz 
This is a command-line quiz application that tests your knowledge of Linux commands. The quiz presents a series of Linux command-related problems and evaluates your responses.


# Features
    - Written in go
    - Quiz Format: The problems are presented one by one, and you are required to provide the correct command based on the given description.
    - Time Limit: You have a specified time limit to finish the quiz.
    - Scoring: Your score is calculated based on the number of correct answers.

# Usage
    - Clone the repository or download the source code.
    - Navigate to the project directory in your terminal.
    
```bash
   go run main.go 
```

   or You can use your own question bank

```bash
   go run main.go -json <filename.json> -limit <time_in_seconds>
```

# JSON Format
The quiz problems are defined in a JSON file, where each problem has the following structure:
```json
{
  "Command": "ls -l",
  "Description": "List files and directories. Long format listing.",
  "Examples": ["ls -l"]
}
```

# Dependencies
This application is written in Go (Golang) and does not require external dependencies.