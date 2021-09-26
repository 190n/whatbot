# whatbot

Discord bot to randomize a server's icon.

## Usage

```
Usage of ./whatbot:
  -g string
        ID of the guild to modify
  -i string
        directory to choose an icon from (default "icons")
  -t string
        discord token
```

When run, whatbot will connect, set the server icon to a random file from the specified directory, and then immediately disconnect. To change the icon regularly, you will need another tool set up to run whatbot at regular intervals. Keep in mind that Discord limits you to 2 icon changes every 5 minutes. If this limit has been met, whatbot will hang.
