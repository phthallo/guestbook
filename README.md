# guestbook

```ssh guestbook@phthallo.com```

leave me a message![^1]


There are two parts to this: the SSH app where you can write a small message for me, and a simple API with a GET endpoint that can be called to retrieve messages.

# Development

1. Make sure you have [Go](https://go.dev/dl/) installed.

2. Clone the repository
    ```
    https://github.com/phthallo/guestbook
    ```

3. Install dependencies
    ```
    go get
    ```

4. Start development: 
    ```
    go main.go 
    ```

5. Build `guestbook`
    ```
    go build 
    ```

[^1]: i somehow made all of this without knowing that charmbracelet/wish was a thing. lol, lmao even. 