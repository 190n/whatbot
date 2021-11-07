# whatbot

Discord bot to randomize a server's icon.

## Usage

Install dependencies: `yarn`

Set the environment variables:

- `DISCORD_TOKEN` = bot token to access Discord's API
- `GUILD_ID` = ID of the guild to set the icon of

Run `node whatbot.js`. It will pick a random icon from the `icons` directory in the same directory as the script and set it.
