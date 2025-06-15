# guestbook

```
ssh guestbook@phthallo.com
```

<p style = "text-align: center; font-style: italic">leave me a message! [^1]</p>


Guestbook is exactly what the name suggests - a guestbook for anyone who happens to stumble upon it. There are two parts to this: the SSH app where you can write a small message for me, and an API with a GET endpoint that returns the messages sent.

## Development

1. Make sure you have [Go](https://go.dev/dl/) installed and configured.

2. Clone the repository
    ```
    https://github.com/phthallo/guestbook && cd guestbook
    ```

3. Generate a SSH key in the `guestbook` directory. Skip the passphrase.
    ```
    ssh-keygen -t ed25519 -f ed25519
    ```

3. Install dependencies
    ```
    go get
    ```

4. Start development: 
    ```
    go main.go 
    ```

5. Build it into a handy little `guestbook` executable and run it with `./guestbook`!
    ```
    go build 
    ```

## Environment Variables
| Variable | Explanation |
| -------- | ----------- |
| `API_PORT` | Port to run the API on. |
| `GIN_MODE` | Set it to `release` in production. Set it to anything else for `debug` mode in development.
| `DATABASE_URL` | URL of a Postgres database used for storing messages. |

[^1]: i somehow made all of this without knowing that charmbracelet/wish was a thing. lol, lmao even. 