import { config } from 'dotenv';
config();
import fs from 'fs';
import path from 'path';
import { Client, Intents } from 'discord.js';
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const data = require('./data');

const {
	DISCORD_TOKEN:   token,
	GUILD_ID:        guildID,
	START_TIMESTAMP: startTimestamp,
	INTERVAL:        interval,
} = process.env;
const channels = process.env.CHANNELS.split(',');

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

async function fetchChannel(channel, mostRecentCounted, ignoreAfter) {
	let oldestMessageSeen = { createdTimestamp: Infinity }, done = false, total = 0;
	const counts = {};
	let newMostRecent = -1;
	while (oldestMessageSeen.createdTimestamp > mostRecentCounted && !done) {
		const messages = await channel.messages.fetch({
			limit: 100,
			before: oldestMessageSeen?.id,
		});

		if (newMostRecent < 0) {
			newMostRecent = Date.now();
		}

		total += messages.size;
		console.log(`fetched ${messages.size} messages, total ${total}`);

		if (messages.size == 0) {
			done = true;
			break;
		}

		for (const m of messages.values()) {
			if (m.createdTimestamp <= mostRecentCounted || m.createdTimestamp > ignoreAfter) {
				if (m.createdTimestamp <= mostRecentCounted) {
					done = true;
				}
				continue;
			}
			if (m.createdTimestamp < oldestMessageSeen.createdTimestamp) {
				oldestMessageSeen = m;
			}
			const block = Math.floor((m.createdTimestamp - startTimestamp) / interval);
			counts[block] ??= 0;
			counts[block] += 1;
		}
	}

	console.log(counts);
	return [counts, newMostRecent];
}

function merge(...countsObjects) {
	console.log('\n\n\n', countsObjects, '\n\n\n');
	const counts = {};
	for (const countsObj of countsObjects) {
		for (const k of Object.keys(countsObj)) {
			counts[k] ??= 0;
			counts[k] += countsObj[k];
		}
	}
	return counts;
}

client.on('ready', async () => {
	console.log(`Logged in as ${client.user.tag}!`);

	try {
		// set icon
		const guild = await client.guilds.fetch(guildID);
		const iconsDir = path.join(path.dirname(import.meta.url.substr(7)), 'icons');
		const files = await fs.promises.readdir(iconsDir);
		const newIcon = files[Math.floor(Math.random() * files.length)];
		await guild.edit({ icon: path.join(iconsDir, newIcon) });

		// fetch messages
		for (const c of channels) {
			const channel = await client.channels.fetch(c);
			data[c] = {
				counts: {},
				mostRecentCounted: 0,
				...data[c],
				name: channel.name,
			};
			const newCounts = [];
			const [counts, newMostRecent] = await fetchChannel(channel, data[c].mostRecentCounted, Infinity);
			newCounts.push(counts);
			const activeThreads = await channel.threads.fetchActive();
			console.log(`got ${activeThreads.threads.size} active threads`);
			for (const thread of activeThreads.threads.values()) {
				newCounts.push((await fetchChannel(thread, data[c].mostRecentCounted, newMostRecent))[0]);
			}
			// ew archived threads
			let before = undefined, hasMore = true;
			while (hasMore) {
				const archivedThreads = await channel.threads.fetchArchived({ before });
				hasMore = archivedThreads.hasMore;
				console.log(`got ${archivedThreads.threads.size} archived threads`);
				for (const thread of archivedThreads.threads.values()) {
					if (before === undefined || thread.createdTimestamp < before.createdTimestamp) {
						before = thread;
					}
					const threadCounts = (await fetchChannel(thread, data[c].mostRecentCounted, newMostRecent))[0];
					// don't double count the message that started the thread
					const block = Math.floor((thread.createdTimestamp - startTimestamp) / interval);
					if (threadCounts.hasOwnProperty(block) && thread.createdTimestamp > data[c].mostRecentCounted) {
						console.log('thread decrement');
						threadCounts[block] -= 1;
					}
					newCounts.push(threadCounts);
				}
			}

			data[c].counts = merge(data[c].counts, ...newCounts);
			data[c].mostRecentCounted = newMostRecent;
		}
	} catch (e) {
		console.error(e);
	} finally {
		client.destroy();
		console.log(JSON.stringify(data));
		for (const k in data) {
			let sum = 0;
			for (const b in data[k].counts) {
				sum += data[k].counts[b] ?? 0;
			}
			console.log(sum);
		}
	}
});

client.login(token);
