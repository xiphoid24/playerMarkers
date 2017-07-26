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
            // ...in the world "myworls"
            "exmaple" : [
                
                
                // example marker, pretty format:
                {
                    // position ([x, z, y])
                    "pos" : [-285, 138, 68],
                    // title when you hover over the marker
                    "title" :"xiphoid24",
                    // text in the marker popup window
                    "text" : '<div style="text-align: center;">xiphoid24</div><br><b>Location:</b> X: -285, Y: 68, Z: 138',
                    // override the icon of a single marker (optional)
                    "icon" : "b3bddfb444a44c349a342b0e640da150.png",
                },
                
                
            ],
        },
    },
];

if (MAPCRAFTER_MARKERS == null || MAPCRAFTER_MARKERS == undefined) {
    MAPCRAFTER_MARKERS = [];
}

MAPCRAFTER_MARKERS = MAPCRAFTER_MARKERS.concat(PLAYER_MARKERS);
