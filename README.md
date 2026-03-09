> SSH based portfolio. Built with Go, Bubbletea and Wish

## Build Instructions

### 1. Prepare Environment & Keys

Create the directory and generate a secure Ed25519 host key

```bash
$ mkdir -p .ssh
$ cd .ssh
$ ssh-keygen -t ed25519 -f .ssh/host_ed25519 -N ""
```

### 2. Build the Application

```bash
$ go build .
```

### 3. Permissions

Grant the binary permission to bind to port 22

```bash
$ sudo setcap 'cap_net_bind_service=+ep' ./ssh-portfolio
```

### 4. Launch

```bash
$ ./ssh-portfolio <PORT>
```
