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
            "the_end" : [
                {{ range $user := . }}
                {{ if eq $user.Dimension 1 }}
                // example marker, pretty format:
                {
                    // position ([x, z, y])
                    "pos" : [{{ $user.X }}, {{ $user.Z }}, {{ $user.Y }}],
                    // title when you hover over the marker
                    "title" :"{{ $user.Username }}",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">{{ $user.Username }}</div><br><b>Location:</b> X: {{ $user.X }}, Y: {{ $user.Y }}, Z: {{ $user.Z }}',
                    // override the icon of a single marker (optional)
                    "icon" : "{{ $user.Uuid }}.png",
                },
                {{ end }}
                {{ end }}
            ],
            // ...in the world "the_end"
            "myworld" : [
                {{ range $user := . }}
                {{ if ne $user.Dimension 1 }}
                // example marker, pretty format:
                {
                    // position ([x, z, y])
                    "pos" : [{{ $user.X }}, {{ $user.Z }}, {{ $user.Y }}],
                    // title when you hover over the marker
                    "title" :"{{ $user.Username }}",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">{{ $user.Username }}</div><br><b>Location:</b> X: {{ $user.X }}, Y: {{ $user.Y }}, Z: {{ $user.Z }}',
                    // override the icon of a single marker (optional)
                    "icon" : "{{ $user.Uuid }}.png",
                },
                {{ end }}
                {{ end }}
            ],
        },
    },
];

if (MAPCRAFTER_MARKERS == null || MAPCRAFTER_MARKERS == undefined) {
    MAPCRAFTER_MARKERS = [];
}

MAPCRAFTER_MARKERS = MAPCRAFTER_MARKERS.concat(PLAYER_MARKERS);
