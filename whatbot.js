import { config } from 'dotenv';
config();
import fs from 'fs';
import path from 'path';
import { Client, Intents } from 'discord.js';

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

client.on('ready', async () => {
	console.log(`Logged in as ${client.user.tag}!`);

	try {
		const guild = await client.guilds.fetch(process.env.GUILD_ID);
		const iconsDir = path.join(path.dirname(import.meta.url.substr(7)), 'icons');
		const files = await fs.promises.readdir(iconsDir);
		const newIcon = files[Math.floor(Math.random() * files.length)];
		await guild.edit({ icon: path.join(iconsDir, newIcon) });
	} catch (e) {
		console.error(e);
	} finally {
		client.destroy();
	}
});

client.login(process.env.DISCORD_TOKEN);
