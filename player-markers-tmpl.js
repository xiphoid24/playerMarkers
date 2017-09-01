var PLAYER_MARKERS = [
    // just one example marker group
    {
        // id of the marker group, without spaces/other special chars
        "id" : "players",
        // name of the marker group, displayed in the webinterface
        "name" : "Players",
        // icon of the markers belonging to that group (optional)
        "icon" : "steve.png",
        // size of that icon
        "iconSize" : [30, 48],
        // whether this marker group is shown by default (optional)
        "showDefault" : true,
        // markers of this marker group...
        "markers" : {
            // ...in the world "the_end"
            "end" : [{{ range $player := .endPlayers }}
                {
                    // position ([x, z, y])
                    "pos" : [{{ $player.X }}, {{ $player.Z }}, {{ $player.Y }}],
                    // title when you hover over the marker
                    "title" :"{{ $player.Username }}",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">{{ $player.Username }}</div><br><b>Location:</b> X: {{ $player.X }}, Y: {{ $player.Y }}, Z: {{ $player.Z }}',
                    // override the icon of a single marker (optional)
                    "icon" : "{{ $player.Uuid }}.png",
                },
            {{ end }}],
            // ...in the world "nether"
            "nether" : [{{ range $player := .netherPlayers }}
                {
                    // position ([x, z, y])
                    "pos" : [{{ $player.X }}, {{ $player.Z }}, {{ $player.Y }}],
                    // title when you hover over the marker
                    "title" :"{{ $player.Username }}",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">{{ $player.Username }}</div><br><b>Location:</b> X: {{ $player.X }}, Y: {{ $player.Y }}, Z: {{ $player.Z }}',
                    // override the icon of a single marker (optional)
                    "icon" : "{{ $player.Uuid }}.png",
                },
            {{ end }}],
            // ...in the world "overworld"
            "test" : [{{ range $player := .overworldPlayers }}
                {
                    // position ([x, z, y])
                    "pos" : [{{ $player.X }}, {{ $player.Z }}, {{ $player.Y }}],
                    // title when you hover over the marker
                    "title" :"{{ $player.Username }}",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">{{ $player.Username }}</div><br><b>Location:</b> X: {{ $player.X }}, Y: {{ $player.Y }}, Z: {{ $player.Z }}',
                    // override the icon of a single marker (optional)
                    "icon" : "{{ $player.Uuid }}.png",
                },
            {{ end }}],
        },
    },
];

if (MAPCRAFTER_MARKERS == null || MAPCRAFTER_MARKERS == undefined) {
    MAPCRAFTER_MARKERS = [];
}

MAPCRAFTER_MARKERS = MAPCRAFTER_MARKERS.concat(PLAYER_MARKERS);
