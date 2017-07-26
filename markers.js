// Add your own markers to this file.

var MAPCRAFTER_MARKERS = [

	// just one example marker group
	{
		// id of the marker group, without spaces/other special chars
		"id" : "markers",
		// name of the marker group, displayed in the webinterface
		"name" : "Markers",
		// size of that icon
		"iconSize" : [32, 32],
		// whether this marker group is shown by default (optional)
		"showDefault" : true,
		// markers of this marker group...
		"markers" : {
			// ...in the world "world"
			"exmaple" : [
				// example marker, pretty format:
				{
					// position ([x, z, y])
                    "pos":[-242,59,64],
					// title when you hover over the marker
                    "title":"Exmaple",
					// text in the marker popup window
                    "text":"<b>Last update:</b> Wed Jul 26 19:50 UTC 2017<br /><b>Mapcrafter version:</b> v.2.4<br /></a>"
				},
			],
		},
	},

];
