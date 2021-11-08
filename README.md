# whatbot

Discord bot to randomize a server's icon.

## Usage

Install dependencies: `yarn`

Set the environment variables:

- `DISCORD_TOKEN` = bot token to access Discord's API
- `GUILD_ID` = ID of the guild to set the icon of
- `CHANNELS` = comma-separated list of channel IDs to count messages in
- `START_TIMESTAMP` = what time is considered 0, as UNIX timestamp in milliseconds
- `INTERVAL` = duration in which to group messages, in milliseconds

Run `node whatbot.js`. It will pick a random icon from the `icons` directory in the same directory as the script and set it.
