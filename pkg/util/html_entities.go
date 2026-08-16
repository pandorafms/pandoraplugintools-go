package util

import "strings"

// entityRow maps one HTML entity name to its Unicode code point. Ported
// verbatim from general.py's _ENT2CHR table (extracted programmatically from
// the Python source to avoid transcription errors).
type entityRow struct {
	Entity    string
	CodePoint rune
}

var entityTable = []entityRow{
	{Entity: "#x00", CodePoint: 0},
	{Entity: "#x01", CodePoint: 1},
	{Entity: "#x02", CodePoint: 2},
	{Entity: "#x03", CodePoint: 3},
	{Entity: "#x04", CodePoint: 4},
	{Entity: "#x05", CodePoint: 5},
	{Entity: "#x06", CodePoint: 6},
	{Entity: "#x07", CodePoint: 7},
	{Entity: "#x08", CodePoint: 8},
	{Entity: "#x09", CodePoint: 9},
	{Entity: "#x0a", CodePoint: 10},
	{Entity: "#x0b", CodePoint: 11},
	{Entity: "#x0c", CodePoint: 12},
	{Entity: "#x0d", CodePoint: 13},
	{Entity: "#x0e", CodePoint: 14},
	{Entity: "#x0f", CodePoint: 15},
	{Entity: "#x10", CodePoint: 16},
	{Entity: "#x11", CodePoint: 17},
	{Entity: "#x12", CodePoint: 18},
	{Entity: "#x13", CodePoint: 19},
	{Entity: "#x14", CodePoint: 20},
	{Entity: "#x15", CodePoint: 21},
	{Entity: "#x16", CodePoint: 22},
	{Entity: "#x17", CodePoint: 23},
	{Entity: "#x18", CodePoint: 24},
	{Entity: "#x19", CodePoint: 25},
	{Entity: "#x1a", CodePoint: 26},
	{Entity: "#x1b", CodePoint: 27},
	{Entity: "#x1c", CodePoint: 28},
	{Entity: "#x1d", CodePoint: 29},
	{Entity: "#x1e", CodePoint: 30},
	{Entity: "#x1f", CodePoint: 31},
	{Entity: "#x20", CodePoint: 32},
	{Entity: "quot", CodePoint: 34},
	{Entity: "amp", CodePoint: 38},
	{Entity: "#039", CodePoint: 39},
	{Entity: "#40", CodePoint: 40},
	{Entity: "#41", CodePoint: 41},
	{Entity: "lt", CodePoint: 60},
	{Entity: "gt", CodePoint: 62},
	{Entity: "#92", CodePoint: 92},
	{Entity: "#x80", CodePoint: 128},
	{Entity: "#x81", CodePoint: 129},
	{Entity: "#x82", CodePoint: 130},
	{Entity: "#x83", CodePoint: 131},
	{Entity: "#x84", CodePoint: 132},
	{Entity: "#x85", CodePoint: 133},
	{Entity: "#x86", CodePoint: 134},
	{Entity: "#x87", CodePoint: 135},
	{Entity: "#x88", CodePoint: 136},
	{Entity: "#x89", CodePoint: 137},
	{Entity: "#x8a", CodePoint: 138},
	{Entity: "#x8b", CodePoint: 139},
	{Entity: "#x8c", CodePoint: 140},
	{Entity: "#x8d", CodePoint: 141},
	{Entity: "#x8e", CodePoint: 142},
	{Entity: "#x8f", CodePoint: 143},
	{Entity: "#x90", CodePoint: 144},
	{Entity: "#x91", CodePoint: 145},
	{Entity: "#x92", CodePoint: 146},
	{Entity: "#x93", CodePoint: 147},
	{Entity: "#x94", CodePoint: 148},
	{Entity: "#x95", CodePoint: 149},
	{Entity: "#x96", CodePoint: 150},
	{Entity: "#x97", CodePoint: 151},
	{Entity: "#x98", CodePoint: 152},
	{Entity: "#x99", CodePoint: 153},
	{Entity: "#x9a", CodePoint: 154},
	{Entity: "#x9b", CodePoint: 155},
	{Entity: "#x9c", CodePoint: 156},
	{Entity: "#x9d", CodePoint: 157},
	{Entity: "#x9e", CodePoint: 158},
	{Entity: "#x9f", CodePoint: 159},
	{Entity: "#xa0", CodePoint: 160},
	{Entity: "#xa1", CodePoint: 161},
	{Entity: "#xa2", CodePoint: 162},
	{Entity: "#xa3", CodePoint: 163},
	{Entity: "#xa4", CodePoint: 164},
	{Entity: "#xa5", CodePoint: 165},
	{Entity: "#xa6", CodePoint: 166},
	{Entity: "#xa7", CodePoint: 167},
	{Entity: "#xa8", CodePoint: 168},
	{Entity: "#xa9", CodePoint: 169},
	{Entity: "#xaa", CodePoint: 170},
	{Entity: "#xab", CodePoint: 171},
	{Entity: "#xac", CodePoint: 172},
	{Entity: "#xad", CodePoint: 173},
	{Entity: "#xae", CodePoint: 174},
	{Entity: "#xaf", CodePoint: 175},
	{Entity: "#xb0", CodePoint: 176},
	{Entity: "#xb1", CodePoint: 177},
	{Entity: "#xb2", CodePoint: 178},
	{Entity: "#xb3", CodePoint: 179},
	{Entity: "#xb4", CodePoint: 180},
	{Entity: "#xb5", CodePoint: 181},
	{Entity: "#xb6", CodePoint: 182},
	{Entity: "#xb7", CodePoint: 183},
	{Entity: "#xb8", CodePoint: 184},
	{Entity: "#xb9", CodePoint: 185},
	{Entity: "#xba", CodePoint: 186},
	{Entity: "#xbb", CodePoint: 187},
	{Entity: "#xbc", CodePoint: 188},
	{Entity: "#xbd", CodePoint: 189},
	{Entity: "#xbe", CodePoint: 190},
	{Entity: "Aacute", CodePoint: 193},
	{Entity: "Auml", CodePoint: 196},
	{Entity: "Eacute", CodePoint: 201},
	{Entity: "Euml", CodePoint: 203},
	{Entity: "Iacute", CodePoint: 205},
	{Entity: "Iuml", CodePoint: 207},
	{Entity: "Ntilde", CodePoint: 209},
	{Entity: "Oacute", CodePoint: 211},
	{Entity: "Ouml", CodePoint: 214},
	{Entity: "Uacute", CodePoint: 218},
	{Entity: "Uuml", CodePoint: 220},
	{Entity: "aacute", CodePoint: 225},
	{Entity: "auml", CodePoint: 228},
	{Entity: "eacute", CodePoint: 233},
	{Entity: "euml", CodePoint: 235},
	{Entity: "iacute", CodePoint: 237},
	{Entity: "iuml", CodePoint: 239},
	{Entity: "ntilde", CodePoint: 241},
	{Entity: "oacute", CodePoint: 243},
	{Entity: "ouml", CodePoint: 246},
	{Entity: "uacute", CodePoint: 250},
	{Entity: "uuml", CodePoint: 252},
	{Entity: "OElig", CodePoint: 338},
	{Entity: "oelig", CodePoint: 339},
	{Entity: "Scaron", CodePoint: 352},
	{Entity: "scaron", CodePoint: 353},
	{Entity: "Yuml", CodePoint: 376},
	{Entity: "fnof", CodePoint: 402},
	{Entity: "circ", CodePoint: 710},
	{Entity: "tilde", CodePoint: 732},
	{Entity: "Alpha", CodePoint: 913},
	{Entity: "Beta", CodePoint: 914},
	{Entity: "Gamma", CodePoint: 915},
	{Entity: "Delta", CodePoint: 916},
	{Entity: "Epsilon", CodePoint: 917},
	{Entity: "Zeta", CodePoint: 918},
	{Entity: "Eta", CodePoint: 919},
	{Entity: "Theta", CodePoint: 920},
	{Entity: "Iota", CodePoint: 921},
	{Entity: "Kappa", CodePoint: 922},
	{Entity: "Lambda", CodePoint: 923},
	{Entity: "Mu", CodePoint: 924},
	{Entity: "Nu", CodePoint: 925},
	{Entity: "Xi", CodePoint: 926},
	{Entity: "Omicron", CodePoint: 927},
	{Entity: "Pi", CodePoint: 928},
	{Entity: "Rho", CodePoint: 929},
	{Entity: "Sigma", CodePoint: 931},
	{Entity: "Tau", CodePoint: 932},
	{Entity: "Upsilon", CodePoint: 933},
	{Entity: "Phi", CodePoint: 934},
	{Entity: "Chi", CodePoint: 935},
	{Entity: "Psi", CodePoint: 936},
	{Entity: "Omega", CodePoint: 937},
	{Entity: "alpha", CodePoint: 945},
	{Entity: "beta", CodePoint: 946},
	{Entity: "gamma", CodePoint: 947},
	{Entity: "delta", CodePoint: 948},
	{Entity: "epsilon", CodePoint: 949},
	{Entity: "zeta", CodePoint: 950},
	{Entity: "eta", CodePoint: 951},
	{Entity: "theta", CodePoint: 952},
	{Entity: "iota", CodePoint: 953},
	{Entity: "kappa", CodePoint: 954},
	{Entity: "lambda", CodePoint: 955},
	{Entity: "mu", CodePoint: 956},
	{Entity: "nu", CodePoint: 957},
	{Entity: "xi", CodePoint: 958},
	{Entity: "omicron", CodePoint: 959},
	{Entity: "pi", CodePoint: 960},
	{Entity: "rho", CodePoint: 961},
	{Entity: "sigmaf", CodePoint: 962},
	{Entity: "sigma", CodePoint: 963},
	{Entity: "tau", CodePoint: 964},
	{Entity: "upsilon", CodePoint: 965},
	{Entity: "phi", CodePoint: 966},
	{Entity: "chi", CodePoint: 967},
	{Entity: "psi", CodePoint: 968},
	{Entity: "omega", CodePoint: 969},
	{Entity: "thetasym", CodePoint: 977},
	{Entity: "upsih", CodePoint: 978},
	{Entity: "piv", CodePoint: 982},
	{Entity: "ensp", CodePoint: 8194},
	{Entity: "emsp", CodePoint: 8195},
	{Entity: "thinsp", CodePoint: 8201},
	{Entity: "zwnj", CodePoint: 8204},
	{Entity: "zwj", CodePoint: 8205},
	{Entity: "lrm", CodePoint: 8206},
	{Entity: "rlm", CodePoint: 8207},
	{Entity: "ndash", CodePoint: 8211},
	{Entity: "mdash", CodePoint: 8212},
	{Entity: "lsquo", CodePoint: 8216},
	{Entity: "rsquo", CodePoint: 8217},
	{Entity: "sbquo", CodePoint: 8218},
	{Entity: "ldquo", CodePoint: 8220},
	{Entity: "rdquo", CodePoint: 8221},
	{Entity: "bdquo", CodePoint: 8222},
	{Entity: "dagger", CodePoint: 8224},
	{Entity: "Dagger", CodePoint: 8225},
	{Entity: "bull", CodePoint: 8226},
	{Entity: "hellip", CodePoint: 8230},
	{Entity: "permil", CodePoint: 8240},
	{Entity: "prime", CodePoint: 8242},
	{Entity: "Prime", CodePoint: 8243},
	{Entity: "lsaquo", CodePoint: 8249},
	{Entity: "rsaquo", CodePoint: 8250},
	{Entity: "oline", CodePoint: 8254},
	{Entity: "frasl", CodePoint: 8260},
	{Entity: "euro", CodePoint: 8364},
	{Entity: "image", CodePoint: 8465},
	{Entity: "weierp", CodePoint: 8472},
	{Entity: "real", CodePoint: 8476},
	{Entity: "trade", CodePoint: 8482},
	{Entity: "alefsym", CodePoint: 8501},
	{Entity: "larr", CodePoint: 8592},
	{Entity: "uarr", CodePoint: 8593},
	{Entity: "rarr", CodePoint: 8594},
	{Entity: "darr", CodePoint: 8595},
	{Entity: "harr", CodePoint: 8596},
	{Entity: "crarr", CodePoint: 8629},
	{Entity: "lArr", CodePoint: 8656},
	{Entity: "uArr", CodePoint: 8657},
	{Entity: "rArr", CodePoint: 8658},
	{Entity: "dArr", CodePoint: 8659},
	{Entity: "hArr", CodePoint: 8660},
	{Entity: "forall", CodePoint: 8704},
	{Entity: "part", CodePoint: 8706},
	{Entity: "exist", CodePoint: 8707},
	{Entity: "empty", CodePoint: 8709},
	{Entity: "nabla", CodePoint: 8711},
	{Entity: "isin", CodePoint: 8712},
	{Entity: "notin", CodePoint: 8713},
	{Entity: "ni", CodePoint: 8715},
	{Entity: "prod", CodePoint: 8719},
	{Entity: "sum", CodePoint: 8721},
	{Entity: "minus", CodePoint: 8722},
	{Entity: "lowast", CodePoint: 8727},
	{Entity: "radic", CodePoint: 8730},
	{Entity: "prop", CodePoint: 8733},
	{Entity: "infin", CodePoint: 8734},
	{Entity: "ang", CodePoint: 8736},
	{Entity: "and", CodePoint: 8743},
	{Entity: "or", CodePoint: 8744},
	{Entity: "cap", CodePoint: 8745},
	{Entity: "cup", CodePoint: 8746},
	{Entity: "int", CodePoint: 8747},
	{Entity: "there4", CodePoint: 8756},
	{Entity: "sim", CodePoint: 8764},
	{Entity: "cong", CodePoint: 8773},
	{Entity: "asymp", CodePoint: 8776},
	{Entity: "ne", CodePoint: 8800},
	{Entity: "equiv", CodePoint: 8801},
	{Entity: "le", CodePoint: 8804},
	{Entity: "ge", CodePoint: 8805},
	{Entity: "sub", CodePoint: 8834},
	{Entity: "sup", CodePoint: 8835},
	{Entity: "nsub", CodePoint: 8836},
	{Entity: "sube", CodePoint: 8838},
	{Entity: "supe", CodePoint: 8839},
	{Entity: "oplus", CodePoint: 8853},
	{Entity: "otimes", CodePoint: 8855},
	{Entity: "perp", CodePoint: 8869},
	{Entity: "sdot", CodePoint: 8901},
	{Entity: "lceil", CodePoint: 8968},
	{Entity: "rceil", CodePoint: 8969},
	{Entity: "lfloor", CodePoint: 8970},
	{Entity: "rfloor", CodePoint: 8971},
	{Entity: "lang", CodePoint: 9001},
	{Entity: "rang", CodePoint: 9002},
	{Entity: "loz", CodePoint: 9674},
	{Entity: "spades", CodePoint: 9824},
	{Entity: "clubs", CodePoint: 9827},
	{Entity: "hearts", CodePoint: 9829},
	{Entity: "diams", CodePoint: 9830},
}

// charToEntity maps a raw code point to its "&entity;" representation, used
// by SafeInput. entityReplacer maps "&entity;" back to the raw character,
// used by SafeOutput.
var (
	charToEntity   map[rune]string
	entityReplacer *strings.Replacer
)

func init() {
	charToEntity = make(map[rune]string, len(entityTable))

	pairs := make([]string, 0, len(entityTable)*2)
	for _, row := range entityTable {
		entity := "&" + row.Entity + ";"
		charToEntity[row.CodePoint] = entity
		pairs = append(pairs, entity, string(row.CodePoint))
	}

	entityReplacer = strings.NewReplacer(pairs...)
}

// SafeInput encodes characters in s that have an HTML entity representation
// (control characters, quotes, accented letters, etc.) into that entity form,
// leaving unrecognized characters unchanged. Ports general.py's safe_input:
// note the Python source's own docstring on that function is misleading
// (it says "decode entities to a clear string"), but its actual
// implementation encodes raw characters into entities, which is what this
// mirrors.
func SafeInput(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if entity, ok := charToEntity[r]; ok {
			b.WriteString(entity)
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// SafeOutput decodes every "&entity;" occurrence in s back into its raw
// character. Ports general.py's safe_output.
func SafeOutput(s string) string {
	if s == "" {
		return ""
	}

	return entityReplacer.Replace(s)
}
