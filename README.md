# guestbook

```
ssh guestbook.phthallo.com
```

<p align = "center"><i>leave me a message!</i></p>


Guestbook[^1] is exactly what the name suggests - a guestbook for anyone who happens to stumble upon it. There are two parts to this: the SSH app where you can write a small message for me, and an API with a GET endpoint (`/entries?limit=10`) that returns the messages sent.

You can also customise the name attributed to your message by specifying a username when you SSH.

## Development

1. Make sure you have [Go](https://go.dev/dl/) installed and configured.

2. Clone the repository
    ```
    https://github.com/phthallo/guestbook && cd guestbook
    ```

3. Generate a SSH key in the `guestbook` directory. Skip the passphrase.
    ```
    ssh-keygen -t ed25519 -f id_ed25519
    ```

3. Install dependencies
    ```
    go get
    ```

4. Start development: 
    ```
    go main.go 
    ```

5. Build it into a handy little `guestbook` executable and run it with `./guestbook`! The API will be accessible at `localhost:{API_PORT}` and the SSH server will listen on the specified `SSH_PORT`. 
    ```
    go build 
    ```

## Docker

Generate the `ed25519` key outside of the container to avoid it being recreated every time you rebuild it, which will cause scary looking `REMOTE HOST IDENTIFICATION CHANGED` errors for people who might be revisiting.

Then, use the provided `docker-compose.yml`.

```
docker compose up -d
```


## Environment Variables
| Variable | Explanation |
| -------- | ----------- |
| `API_PORT` | Port to run the API on. |
| `SSH_PORT` | Port to listen on for SSH. |
| `GIN_MODE` | Set it to `release` in production. Set it to anything else for `debug` mode in development.
| `DATABASE_URL` | URL of a Postgres database used for storing messages. |
| `TERM` | Set it to `xterm-256color`. For colours. |
|  `HOSTNAME` | Set it to `127.0.0.1` for development. |

[^1]: i somehow made all of this without knowing that charmbracelet/wish was a thing. lol, lmao even. 