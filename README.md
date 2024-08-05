# Vanilla Proxy

Vanilla Proxy is a simple proxy servers that enabled Packet management and protection of BDS server attacks through a end to end proxy system.

## Overview

1. You have all the advantages of the original BDS.
2. You can modify and cancel any packets passing through the proxy using the handler system.
3. You can organize the protection of your server from DDOS attacks.
4. You can set the boards of the world. Blocks outside the border will not be displayed to the player. The coordinates of the corner must be divided by 16 so that the borders are displayed correctly.
5. Support for a whitelist that can store the player's xbox id(xuid).
6. A Player List manager that stores a players `SelfSignedID` ensuring as long as a player is signed into a XBL account, they wont get there player reset.

## Getting Started

> [!IMPORTANT]
> You must configure BDS `server.properties` to disable `online-mode` and `client-side-chunk-generation-enabled`

1. Download and install [Go](https://go.dev/dl/)
2. Download and run Minecraft [Bedrock Dedicated Server](https://www.minecraft.net/en-us/download/server/bedrock).
3. Run in the console:

```
git clone https://github.com/smell-of-curry/vanilla-proxy
cd vanilla-proxy
go run main.go
```

## Configuration

Configuration in the Vanilla Proxy is all managed through the [config.toml](config.toml.example). This file holds all details including database, server connection, api, world border and more.
To get started, copy the [config.toml.example](config.toml.example) file and rename it to `config.toml`. Inside this file you can first set the `Connection` properties.

### Connection

These are the most important as this is what the proxy uses to send the upstream connections too.

- `ProxyAddress` - This is the address that the THIS proxy server will run on. This is the address that you want people to connect too.
- `RemoteAddress` - This is the address that the BDS server is running on. Ensure that the BDS server is running on a different port.

### Api

This vanilla proxy uses a API to fetch and set players connection details. Details such as connection IP, name, and XUID are saved for moderation capabilities.

- `ApiHost` - This is the API host, something like `https://pokebedrock.com` would work.
- `ApiKey` - This is a authentication token which will be passed as a password.

### Database

This proxy uses a database to fetch claims, and local player details that are important for operations. This must be managed to ensure the claims work correctly.

- `Host` - This is the host of the database, usually localhost `http://127.0.0.1`
- `Key` - This is the authentication key, which is passed as a password to ensure a authenticated session.
- `Name` - This is the database name you want to connect too, it would switch between servers, something like `black`, `white`, `testing` is used.
- `Port` - This is the database port that the server is running on, which will attach to `Host` when creating a request.

### Logging

This proxy uses a bit of logging so that the discord and the staff team are updated on the server details. Details about sign changes, and failed to ping alerts get sent to discord webhooks to ensure the staff is alerted.

- `DiscordChatLogsWebhook` - this is the endpoint you want chat logs sent too.
- `DiscordCommandLogsWebhook` - this is the endpoint where when players use commands get sent too.
- `DiscordSignLogsIconURL` - This is a icon to be used when a sign log is sent.
- `DiscordSignLogsWebhook` this is the destination for where sign edit logs should be sent.
- `DiscordStaffAlertsWebhook` - this is a endpoint for staff alerts, things like failed to ping database, etc.

### Resources

This is still a work in progress feature, however this configuration can allow use of custom resource packs to be downloaded by players.

- `PackURLs` - this is a string array of URLs that the players must download and activate to play.

### Server

Server holds less essential server configuration that changes connection aspects.

- `DisableXboxAuth` - specifies if authentication of players that join is disabled. If set to true, no verification will be done to ensure that the player connecting is authenticated using their XBOX Live account.
- `Prefix` - Prefix is used to specify the current server in error logs. For example `TESTING` would be sent with a logging endpoint to tell readers this came from `TESTING` server.
- `ViewDistance` - Manages the distance players can view through the chunk handler. This is important for large servers.
- `Whitelist` - If the whitelist is turned on and limiting players from joining. Whitelist can be managed through the (whitelist.json)[whitelist.json] file.

### WorldBorder

The proxy system comes with a pre-built world border that limits chunks, entities, and ticking from happening outside the world border. This is important for large servers and servers that pregenerate chunks.

- `Enabled` - if the world border is enabled.
- `MaxX` & `MaxZ` - Holds the Max location in the positive direction for the border, example `6000`
- `MinX` & `MinZ` - Holds the Minimum location in the negative direction for the border, example `-6000`
