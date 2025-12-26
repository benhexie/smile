# SMILE

## Build
This project includes a cross-platform setup script that works on Windows, macOS, and Linux.

### Prerequisites
- [Go](https://go.dev/dl/) (version 1.21 or higher)

### How to Build
Run the following command in your terminal:

```bash
go run setup.go
```

This will:
1. Download dependencies.
2. Compile the project for Windows (generating `smile.exe`).
3. Place the output in the `build/` directory.

### Configuration
Edit `config.prop` to set your `USER_ID` and other preferences before distributing the executable.
