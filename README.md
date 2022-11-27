# Discord CMD Bot
Simple discord Bot to run system command

## How to run
1. clone this repo
2. compile it
    ```shell
    go build
    ```
3. run it with your token and config (check example.yaml for available config)
    ```shell
   ./discord-cmd-bot -t <TOKEN> -c <config path>
    ```
4. call the bot in discord with this command ```%run <your command>```