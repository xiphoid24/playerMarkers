# Player Markers
Generate player markers with go for Mapcrafter

This is designed to generate player markers for a Mapcrafter map. It uses the world/playerdata/ to find the last known location of the players.

It useful in  a situation where the Mapcrafter map is not on the same server as the Minecraft server or you are unable to install mods.

It uses [minero-go](https://github.com/minero/minero-go/tree/master/proto/nbt) NBT package to read and parse the player.dat. 
Next it will call the mojang API to retreive a username from a UUID and then it calls the [Visage](https://visage.surgeplay.com/index.html) API to download the user skin.

After successfully retreiving all the player information it generates a js file to tell Mapcrafter to create the markers as seen [here](https://docs.mapcrafter.org/builds/stable/markers.html#manually-specifying-markers).

## Use
Clone or go get the this repository.

You can modify the player-markers-tmpl.js to reflect your specific Mapcrafter setup.

Create a json config file.

Finally run `playerMarkers -c=/location/of/json/config.json`

NOTE: if you have other manually set markers then you will want to generate a seperate player-markers.js 
file which merges your existing manual markers with the player markers to create one MAPCRAFTER_MARKERS array.
You will also need to modify your index.html Mapcrafter template to include this newly generated file.

Added player-markers.js under markers.js

```html
<script type="text/javascript" src="config.js"></script>
<script type="text/javascript" src="markers.js"></script>
<script type="text/javascript" src="player-markers.js"></script>
<script type="text/javascript" src="markers-generated.js"></script>
```

## Config File

This is a simple json config file to tell the program where to find your world and where to same the generated js file.

**api-url** *string*      rest API endpoint to get usernames from UUIDs. 
This should be left ommited unless you are modifying the go source code

**skin-url** *string*     rest API url to download skins.
This should be left omitted unless you are modifying the go source code

**js-path** *string*      full path including file name of where the final js file should be saved

**js-tmpl-path** *string* full path of where the js template is located. NOTE: this should be valid go template syntax

**dat-dirs** *string array* each element of the array should be a full path to a playerdata directory inside a worldyou want to render markers for

#### Example
```json
{
    "api-url": "https://sessionserver.mojang.com/session/minecraft/profile/",
    "skin-url": "https://visage.surgeplay.com/frontfull/50/",
    "skin-dir": "/home/user/minecraft-map/static/markers/",
    "js-path": "/home/user/minecraft-map/player-markers.js",
    "js-tmpl-path": "/home/user/playerMarkers/player-markers-tmpl.js",
    "dat-dirs": [
        "/home/user/world/playerdata/",
        "/home/user/world_the_end/playerdata/"
    ]
}
```
