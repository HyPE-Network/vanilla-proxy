## How to start a proxy server
1. Download and install [Go](https://go.dev/dl/)
2. Download and run Minecraft [Bedrock Dedicated Server](https://www.minecraft.net/en-us/download/server/bedrock).

```You must set online-mode=false in BDS server.properties```
```You must set client-side-chunk-generation-enabled=false in BDS server.properties```

3. Run in the console:
```
git clone https://github.com/HyPE-Network/vanilla-proxy
cd vanilla-proxy
go run main.go
```

>You can also customize config.yml for yourself

## Overview
1. You have all the advantages of the original BDS.
2. You can modify and cancel any packets passing through the proxy using the handler system.
3. You can organize the protection of your server from ddos attacks.
4. You can set the boards of the world. Blocks outside the border will not be displayed to the player. The coordinates of the corner must be divided by 16 so that the borders are displayed correctly.
5. Support for a whitelist that can store the player's xbox id(xuid).
6. You can use a special bot for Windows with operator capabilities or send commands directly to screen on your Linux server. Screen must have the name of the original server port (for example 19132).
7. You can use rcon to execute commands automatically (let's say /whitelist add "nickname").
8. A convenient API for creating forms, commands and fake inventories that are processed on the proxy side.

## Examples
Example of a simple server with survival - [vanilla-survival](https://github.com/HyPE-Network/vanilla-survival)
> The world has borders, there is a ban on placing obsidian blocks in nether (players can teleport to distant coordinates using nether portals).
> This prohibition is made by processing the breaking of the block.

## Recommendations
You can block the BDS server port using a firewall so that unauthorized players cannot enter the main server by bypassing the proxy.
