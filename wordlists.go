package namegen

// DefaultAdjectives and DefaultNouns back New(). They're deliberately plain
// and short so generated names stay readable and easy to say out loud.
var DefaultAdjectives = []string{
	"amber", "bold", "brave", "breezy", "bright", "calm", "cheerful", "clever",
	"cosmic", "crisp", "curious", "daring", "dusty", "eager", "fierce", "frosty",
	"gentle", "glowing", "golden", "graceful", "hidden", "honest", "humble", "jaunty",
	"jolly", "keen", "lively", "lonely", "lucky", "mellow", "merry", "misty",
	"nimble", "noble", "patient", "quiet", "radiant", "rapid", "rugged", "rustic",
	"sharp", "silent", "steady", "sunny", "swift", "vivid", "wandering", "wry",
}

var DefaultNouns = []string{
	"alder", "badger", "brook", "canyon", "cavern", "cedar", "cliff", "comet",
	"cove", "coyote", "crag", "delta", "dune", "egret", "ember", "falcon",
	"fern", "fjord", "forest", "glacier", "grove", "harbor", "heron", "hollow",
	"ibis", "island", "juniper", "lagoon", "lantern", "marsh", "meadow", "oasis",
	"orchard", "otter", "pebble", "plateau", "prairie", "quail", "raven", "reef",
	"ridge", "river", "sparrow", "summit", "thicket", "tundra", "willow", "wren",
}

// SpaceAdjectives and SpaceNouns back NewWithTheme(ThemeSpace), for names
// like "lunar-comet" or "distant-nebula".
var SpaceAdjectives = []string{
	"astral", "celestial", "cosmic", "distant", "galactic", "gravitational",
	"interstellar", "lunar", "orbital", "planetary", "solar", "stellar",
}

var SpaceNouns = []string{
	"asteroid", "comet", "eclipse", "galaxy", "meteor", "nebula",
	"nova", "orbit", "pulsar", "quasar", "satellite", "supernova",
}

// TechAdjectives and TechNouns back NewWithTheme(ThemeTech), for names like
// "cached-daemon" or "recursive-kernel".
var TechAdjectives = []string{
	"binary", "cached", "cloud", "compiled", "digital", "encrypted",
	"modular", "parallel", "recursive", "scalable", "virtual", "wireless",
}

var TechNouns = []string{
	"array", "buffer", "cache", "daemon", "kernel", "packet",
	"pipeline", "protocol", "register", "server", "socket", "thread",
}
