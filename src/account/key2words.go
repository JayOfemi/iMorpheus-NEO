package main

import (
	"crypto/sha256"
	"math/big"
)

const addressChecksum8bitsLen = 1 //8*addressChecksum8bitsLen bits checksum
const wordsCount = 24
const wordPerBit = 11
const dictionaryWordsCnt = 1 << wordPerBit - 1

var privKeysDictionary = [dictionaryWordsCnt]string{
	"ability", "able", "about", "above", "abroad", "absence", "absent", "accent", "accept", "accident",
	"according", "account", "ache", "achieve", "across", "action", "active", "activity", "actor", "advertisement",
	"advance", "advantage", "adventure", "advertise", "actual", "advice", "advise", "aeroplane", "affair", "affect",
	"afford", "afraid", "africa", "african", "after", "afternoon", "afterwards", "again", "against", "age",
	"aggression", "aggressive", "ago", "agree", "agreement", "agricultural", "agriculture", "ahead", "aid", "aids",
	"aim", "air", "aircraft", "airline", "airmail", "airplane", "airport", "alarm", "alive", "all",
	"allow", "almost", "alone", "along", "aloud", "already", "also", "although", "altogether", "always",
	"amaze", "ambulance", "among", "amuse", "amusement", "ancestor", "ancient", "and", "anger", "angry",
	"animal", "announce", "announcement", "annoy", "another", "answer", "ant", "antarctic", "antique", "anxious",
	"anybody", "anyhow", "anyone", "anything", "anyway", "anywhere", "apartment", "apologize", "apology", "appear",
	"appearance", "apple", "application", "apply", "appointment", "arctic", "area", "argue", "argument", "arise",
	"arithmetic", "arm", "armchair", "army", "arose", "around", "arrange", "art", "arrow", "arrive",
	"article", "artist", "aside", "ask", "asleep", "assistant", "astonish", "astronaut", "astronomy", "athlete",
	"atmosphere", "atom", "attack", "attempt", "attend", "attention", "attitude", "attract", "attractive",
	"baby", "back", "backache", "background", "backward", "bacon", "bacteria", "bacterium", "bad",
	"baggage", "bake", "bakery", "balance", "balcony", "ball", "ballet", "balloon", "ballpoint", "bamboo",
	"barbecue", "barber", "scholarship", "trousers", "truck", "true", "wipe", "wish", "with", "within",
	"women", "wonderful", "wood", "wooden", "wool", "woollen", "word", "wore", "work", "workday",
	"worker", "workforce", "workmate", "workplace", "works", "world", "worn", "worried", "worry", "worse",
	"worst", "worth", "worthless", "worthwhile", "would", "wound", "wounded", "write", "writing",
	"written", "wrong", "wrote", "year", "yellow", "yes", "yesterday", "yet", "you", "young",
	"troublesome", "barbershop", "bargain", "bark", "base", "baseball", "basement", "basic", "basin", "basket",
	"basketball", "bat", "bathe", "bathrobe", "bathroom", "bathtub", "battery", "battle", "battleground", "bay",
	"beam", "bean", "bear", "beat", "beaten", "beautiful", "beauty", "became", "because",
	"become", "bed", "bedclothes", "bedroom", "bee", "beef", "beehive", "beer", "before", "beg",
	"began", "begin", "beginning", "begun", "behave", "behaviour", "behind", "beijing", "being", "belgium",
	"belief", "believe", "bell", "below", "belt", "bench", "bend", "beneath", "bent", "beside",
	"besides", "best", "better", "between", "beyond", "bicycle", "big", "bike", "bill", "billion",
	"biology", "bird", "birdcage", "birth", "birthday", "birthplace", "biscuit", "bit", "bite", "bitten",
	"bitter", "blame", "blank", "blanket", "bleed", "bless", "blew", "blind", "block", "blood",
	"blouse", "blow", "blown", "blue", "board", "boat", "boating", "body", "bodybuilding",
	"boil", "bomb", "bone", "book", "bookcase", "bookmark", "bookshelf", "bookshop", "bookstore", "boot",
	"boring", "born", "borrow", "boss", "botany", "both", "bottle", "bottom", "bought", "bound",
	"bow", "bowl", "box", "boxing", "boy", "brain", "brake", "branch", "brave", "bravery",
	"bread", "break", "breakfast", "breath", "breathe", "brick", "bride", "bridegroom", "bring", "britain",
	"british", "broad", "broadcast", "broke", "broken", "broom", "brother", "brotherhood", "brought",
	"brown", "brunch", "brush", "bucket", "buddhism", "buddhist", "build", "building", "built", "bun",
	"burn", "burnt", "burst", "bury", "bus", "businessman", "businesswoman", "but", "butter", "butterfly",
	"button", "buy", "caf", "cage", "cake", "call", "calm", "came", "camel", "camera",
	"camp", "can", "canadian", "canal", "cancel", "candy", "cannot", "cant", "canteen", "cap",
	"capital", "captain", "car", "carbon", "card", "care", "careful", "careless", "carpet", "carriage",
	"carrier", "carrot", "carry", "cartoon", "carve", "case", "cash", "cast", "castle", "cat",
	"catch", "cathedral", "cattle", "caught", "celebrate", "celebration", "cell", "cent", "centigrade", "centimeter",
	"centimetre", "central", "centre", "century", "certain", "certainly", "certificate", "chain", "chair",
	"chalk", "challenging", "champion", "chance", "change", "chapter", "character", "charge", "chart", "chat",
	"cheap", "cheat", "check", "cheek", "cheer", "cheerful", "cheers", "cheese", "chemical", "chemist",
	"chemistry", "cheque", "chess", "chest", "chew", "chick", "chicken", "chief", "child", "childhood",
	"children", "chimney", "china", "chinese", "choice", "choke", "choose", "chopsticks", "chose", "chosen",
	"christian", "christmas", "church", "cigar", "cigarette", "cinema", "circle", "circus", "citizen", "city",
	"clap", "classical", "classmate", "classroom", "clean", "cleaner", "clear", "clearly", "clerk", "clever",
	"click", "clinic", "clock", "clone", "close", "cloth", "clothes", "clothing", "cloud", "cloudy",
	"club", "coach", "coal", "coast", "coat", "cock", "cocoa", "coffee", "coin", "coke", "compete",
	"cold", "colleague", "collect", "collection", "college", "color", "colour", "comb", "comfortable", "comma",
	"comment", "common", "communication", "communism", "communist", "compact", "companion", "company", "compare",
	"competition", "competitor", "complete", "composition", "compressed", "computer", "comrade",
	"conclude", "conclusion", "condition", "conference", "congratulate", "congratulation", "conj", "connection",
	"consider", "considerate", "consideration", "consist", "constant", "construct", "construction", "contain",
	"continent", "continue", "contrary", "contribution", "control", "convenience", "convenient", "conversation",
	"corn", "correct", "correction", "correspond", "cost", "cotton", "cough", "could", "count", "counter",
	"country", "countryside", "couple", "courage", "course", "court", "courtyard", "cousin", "cover", "cow",
	"crayon", "crazy", "credit", "crew", "crime", "criminal", "crop", "cross", "crossing", "crossroads",
	"crowd", "crowded", "cruel", "cry", "cubic", "culture", "cup", "cure", "curious", "currency", "cook", "cordless",
	"curtain", "cushion", "custom", "customer", "customs", "cut", "cyclist", "daily", "dam", "damp", "concert",
	"dance", "danger", "dangerous", "dare", "dark", "darkness", "dash", "data", "database", "date", "conceited",
	"daughter", "dawn", "day", "dead", "deadline", "deaf", "deal", "dear", "death", "debate", "container", "content",
	"debt", "december", "decide", "decorate", "decoration", "deed", "deer", "defence", "defend", "defense",
	"degree", "delay", "delete", "delicious", "delight", "delighted", "deliver", "demand", "dentist", "department",
	"departure", "depend", "dept", "describe", "description", "desert", "design", "desire", "destroy", "detective",
	"determination", "determine", "develop", "development", "devote", "devotion", "diagram", "dial", "dialogue",
	"diary", "dictation", "dictionary", "did", "difference", "different", "difficult", "difficulty", "dig", "digital",
	"dine", "direct", "direction", "director", "directory", "dirt", "disabled", "disadvantage", "disagree",
	"disappear", "disappoint", "disappointment", "disaster", "disc", "discourage", "discover", "discovery", "diamond",
	"discussion", "disease", "dish", "dismiss", "distant", "district", "disturb", "dive", "divide", "division",
	"dizzy", "document", "does", "dog", "dollar", "done", "door", "discrimination",
	"dormitory", "dot", "down", "download", "downstairs", "downtown", "downward", "drank", "draw", "drawer",
	"drawing", "drawn", "dream", "dreamt", "dress", "drew", "drier", "drill", "drink", "drive", "discuss",
	"driven", "driver", "drove", "drown", "drug", "drum", "drunk", "dry", "dryer",
	"duck", "duckling", "dull", "dumpling", "during", "dusk", "dust", "dustbin", "dusty", "duty", "disagreement",
	"dvd", "eager", "eagle", "ear", "early", "earn", "earth", "earthquake", "ease", "easily", "conservation",
	"east", "easter", "eastern", "eastwards", "easy", "eat", "eaten", "editor", "educate", "education", "conservative",
	"educator", "effect", "effort", "egg", "egypt", "egyptian", "eight", "eighteen", "eighth", "eighty",
	"either", "elder", "elect", "electric", "electrical", "electricity", "electronic", "elephant", "eleven", "else",
	"emergency", "emperor", "empty", "encourage", "encouragement", "end", "endless", "enemy", "energetic", "energy",
	"engine", "engineer", "england", "english", "enjoyable", "enlarge", "enough", "enquiry", "enter", "entertainment",
	"entire", "entrance", "entry", "envelope", "environment", "envy", "equality", "eraser", "error", "escape",
	"especially", "essay", "europe", "european", "eve", "event", "eventually", "ever", "everybody", "everyday",
	"everyone", "everything", "everywhere", "exact", "exactly", "exam", "examine", "example", "excellent", "except",
	"exchange", "excite", "exhibition", "exist", "exit", "expectation", "expedition", "expense", "expensive",
	"experiment", "expert", "explain", "explanation", "explode", "exploit", "explorer", "expose", "express",
	"extra", "extremely", "eye", "eyewitness", "fade", "fail", "failure", "fair", "fairly", "fairness", "experience",
	"faith", "fall", "fallen", "false", "familiar", "family", "famous", "fan", "fancy", "fantastic",
	"fantasy", "far", "fare", "farm", "farmer", "farther", "fasten", "fat", "father",
	"fault", "favor", "favourite", "fax", "fear", "feather", "february", "fed", "federal", "fee",
	"feed", "feel", "feeling", "feet", "fell", "fellow", "felt", "female", "fence", "ferry", "expression",
	"festival", "fetch", "fever", "few", "fiber", "fibre", "field", "fierce", "fifteen", "fifth",
	"fifty", "figure", "file", "find", "fingernail", "finish", "fire", "fireplace", "firewood", "fireworks",
	"firm", "first", "fish", "fist", "five", "flaming", "flash", "fled", "flee", "flesh",
	"flew", "flight", "float", "flood", "floor", "flour", "flow", "flower", "flown", "flu",
	"fly", "fog", "foggy", "fold", "folk", "follow", "following", "fond", "food", "fool",
	"foolish", "foot", "football", "for", "forbade", "forbid", "forecast", "forehead", "foreign", "foreigner",
	"foresaw", "foresee", "foreseen", "forest", "forever", "forgave", "forget", "forgetful", "forgive", "forgiven",
	"forgot", "forgotten", "fork", "fortnight", "fortunate", "fortune", "forty", "forward", "fought", "found",
	"founding", "fountain", "four", "fox", "france", "free", "freedom", "freeway", "freeze", "freezing",
	"french", "frenchman", "frenchmen", "frequent", "fresh", "friday", "fridge", "fried", "friend", "friendly",
	"friendship", "fright", "frighten", "frog", "from", "frontier", "frost", "froze", "frozen",
	"fun", "funny", "fur", "furniture", "furthest", "future", "gale", "gallery", "gallon", "game",
	"games", "garage", "garbage", "garden", "gardening", "gas", "gather", "gave", "gay", "general",
	"generation", "gentle", "geometry", "german", "germany", "gesture", "get", "gifted", "giraffe", "girl",
	"give", "given", "glad", "glance", "glare", "glass", "glasshouse", "globe", "glory", "glove",
	"glue", "goat", "god", "golden", "goldfish", "goods", "government", "gown", "graduate",
	"graduation", "grain", "grand", "granddaughter", "grandfather", "grandma", "grandmother", "grandpa",
	"granny", "grape", "great", "greece", "greedy", "greek", "green", "greengrocer", "greet", "greeting",
	"grew", "grey", "grill", "grocer", "ground", "group", "grow", "grown", "growth", "gruel", "grandson",
	"guard", "guess", "guest", "guidance", "guide", "guilty", "guitar", "gun", "gym", "gymnasium",
	"gymnastics", "half", "hammer", "hand", "handful", "handkerchief", "handle", "handsome", "handwriting",
	"handy", "hang", "happen", "happily", "happiness", "happy", "harbor", "harbour", "hard", "hardship",
	"harmful", "harmless", "harvest", "has", "hate", "have", "hawk", "hay", "headache", "headline",
	"headmaster", "headmistress", "headteacher", "health", "healthy", "heap", "hear", "heard", "hearing", "heart",
	"heat", "heaven", "heavily", "heavy", "heel", "helicopter", "hello", "helmet", "help", "helpful",
	"hen", "here", "hero", "heroine", "hers", "herself", "hey", "hibernation", "hid", "hide",
	"high", "hill", "hillside", "hilly", "him", "hire", "his", "hobby", "hold",
	"hole", "holiday", "holy", "home", "homeland", "hometown", "homework", "honest", "honey",
	"honor", "honour", "hook", "hooray", "hope", "hopeless", "horrible", "horse", "hospital", "host",
	"hostess", "hot", "hotel", "hour", "house", "how", "howl", "huge", "human", "india",
	"humor", "humorous", "humour", "hundred", "hung", "hunger", "hungry", "hunt", "hunter", "hurricane",
	"hurry", "hurt", "husband", "hydrogen", "idea", "idiom", "immediately", "immigration", "import", "important",
	"impossible", "impress", "impression", "improve", "include", "income", "increase", "indeed", "independent",
	"influence", "inform", "initial", "injure", "inspire", "instant", "instead", "institute",
	"instruct", "instruction", "instrument", "insurance", "insure", "interesting", "international", "internet",
	"interval", "interview", "into", "introduction", "inventor", "invitation", "invite", "ireland", "irish", "iron",
	"irrigate", "irrigation", "island", "italy", "its", "jam", "japan", "japanese", "jar", "jeans", "institution",
	"jeep", "jet", "jewelry", "job", "joke", "journalist", "journey", "joy", "judgement", "juice",
	"juicy", "july", "jump", "june", "jungle", "junior", "junk", "just",
	"justice", "keep", "keeper", "kept", "kettle", "key", "kick", "kid", "kilometer", "kind", "interpreter",
	"kindness", "king", "kiss", "kite", "knives", "knock", "know", "known", "laboratory", "laborer",
	"labour", "labourer", "lack", "lady", "lamb", "lame", "lamp", "land", "language", "lantern",
	"lap", "large", "laser", "last", "late", "lately", "later", "latest", "latter", "laugh", "interrupt",
	"laughter", "laundry", "lavatory", "law", "lawyer", "lay", "lazy", "lead", "leader", "leading",
	"leaf", "league", "leak", "learnt", "least", "leather", "leave", "leaves", "lecture", "led",
	"lemonade", "lend", "lent", "let", "level", "liberate", "liberation", "librarian", "library", "license",
	"lid", "lifetime", "lift", "light", "lightning", "like", "likely", "limit", "line", "link",
	"lion", "lip", "liquid", "list", "listen", "liter", "literary", "literature", "litter", "little",
	"live", "lively", "lives", "living", "load", "loaf", "local", "lock", "locust", "long",
	"look", "loose", "lorry", "lose", "loss", "lost", "lot", "loud", "loudly", "loudspeaker",
	"love", "lovely", "low", "luck", "lucky", "luggage", "lunch", "lung", "mad", "madam",
	"madame", "made", "magazine", "magic", "maid", "main", "mainland", "major", "majority", "make",
	"male", "man", "manager", "mankind", "manner", "many", "map", "marathon",
	"marble", "march", "mark", "market", "marriage", "married", "marry", "mask", "mass", "master",
	"math", "mathematics", "maths", "matter", "maximum", "may", "mean", "meaning", "means", "meant",
	"meanwhile", "measure", "meat", "medal", "medical", "medicine", "medium", "meet", "meeting", "melon",
	"member", "memorial", "memory", "mentally", "mention", "menu", "merciful", "mercy", "merely",
	"merry", "mess", "messy", "met", "meter", "method", "metre", "mexican", "mexico", "mice",
	"microwave", "midday", "middle", "midnight", "might", "mild", "million", "millionaire", "mind",
	"minibus", "minimum", "miniskirt", "minister", "minority", "minus", "minute", "mirror", "miss", "mist",
	"mistaken", "mistook", "misunderstand", "modal", "model", "modern", "monday", "money",
	"monitor", "monkey", "month", "monument", "moon", "mop", "more", "moscow", "mosquito",
	"most", "motherland", "motor", "motorbike", "motorcycle", "motto", "mountain", "mountainous", "mourn", "mouse",
	"moustache", "move", "movement", "movie", "mud", "mum", "murder", "museum", "mushroom",
	"music", "musical", "musician", "must", "mustard", "mutton", "name", "narrow", "nation", "national",
	"native", "natural", "nature", "navy", "nearly", "neat", "neck", "necktie", "need", "neighbor",
	"neighborhood", "neighbourhood", "neither", "nephew", "nervous", "nest", "never", "new",
	"next", "night", "nine", "ninety", "ninth", "nobody", "nod", "noisily", "noisy", "none",
	"noon", "north", "northeast", "northern", "northwards", "northwest", "nose", "notebook", "nothing", "notice",
	"novel", "novelist", "november", "now", "nowhere", "nuclear", "num", "nurse", "nursery", "nursing",
	"nut", "object", "observe", "obtain", "obvious", "occupation", "occur", "ocean", "oceania", "oclock",
	"october", "office", "officer", "official", "offshore", "often", "one", "onion",
	"only", "onto", "open", "opener", "opening", "opera", "operate", "operation", "operator",
	"opposite", "orbit", "order", "ordinary", "organise", "organiser", "organization", "organize", "organizer",
	"other", "otherwise", "ottawa", "ouch", "our", "out", "outer", "outing", "outside", "outward", "origin",
	"oval", "owe", "own", "owner", "ownership", "oxygen", "pacific", "pack", "package", "packet",
	"paddle", "page", "paid", "pain", "painful", "paint", "painter", "painting", "pair", "palace",
	"pale", "pan", "pancake", "panda", "paper", "paperwork", "pardon", "parent", "paris", "park",
	"parrot", "part", "partly", "partner", "party", "pass", "passenger", "passive", "passport",
	"past", "patient", "pattern", "pause", "pavement", "pay", "pedestrian", "pen", "penny",
	"people", "pepper", "per", "percentage", "perfect", "perform", "performance", "performer", "perhaps", "period",
	"permission", "permit", "person", "personal", "personally", "persuade", "pest", "phone", "photo", "photograph",
	"photographer", "phrase", "physical", "physicist", "physics", "pianist", "piano", "pick", "picture",
	"pie", "pig", "pilot", "pin", "pink", "pity", "place", "plain", "plan", "plane",
	"planet", "plant", "plastic", "plate", "platform", "play", "player", "playground", "playmate", "playroom",
	"pleasant", "please", "pleased", "pleasure", "plenty", "plug", "plus", "pocket", "poem", "poet",
	"point", "poison", "poisonous", "pole", "policewoman", "polite", "political", "politician", "politics", "pollute",
	"pollution", "pond", "popular", "population", "pork", "port", "position", "possess", "possession", "possibility",
	"possible", "possibly", "post", "postbox", "postcard", "postcode", "postpone", "pot", "pound", "pour",
	"power", "powerful", "practical", "practice", "practise", "prairie", "praise", "pray", "precious", "prefer",
	"preference", "prep", "prepare", "prescription", "present", "presentation", "president", "press", "pressure",
	"pretty", "price", "pride", "primary", "print", "printer", "printing", "prison", "prisoner", "pretend",
	"private", "prize", "probable", "probably", "problem", "produce", "production", "profession", "professor",
	"programme", "progress", "project", "promise", "pron", "pronunciation", "proper", "properly", "protect", "prove",
	"provide", "province", "pub", "publicly", "publish", "pull", "pulse", "pump", "punctual", "punctuate", "program",
	"punctuation", "punish", "punishment", "pupil", "pure", "purpose", "purse", "push", "put", "puzzled",
	"pyramid", "quality", "quantity", "quarrel", "question", "queue", "quick", "quiet", "quilt", "quite",
	"quiz", "race", "racial", "radiation", "radio", "radioactive", "radium", "rag", "railway", "rain",
	"rainbow", "raincoat", "rainfall", "rainy", "raise", "ran", "rank", "rapid", "rare", "raw",
	"ray", "reach", "read", "real", "reality", "realize", "really", "reason", "reasonable", "rebuild",
	"receipt", "receive", "receiver", "recent", "reception", "receptionist", "recite", "recognise", "recognize",
	"recommend", "record", "recorder", "rectangle", "recycle", "red", "reduce", "refer", "refreshments",
	"refusal", "refuse", "regard", "regards", "register", "regret", "regular", "regulation", "reject", "relate",
	"relation", "relationship", "relative", "relax", "relay", "religion", "religious", "remain", "remark", "remember",
	"remind", "repair", "repairs", "repeat", "replace", "reply", "report", "reporter", "represent", "republic",
	"request", "require", "requirement", "rescue", "research", "reservation", "reserve", "resist", "respect", "rest",
	"restaurant", "restrict", "result", "retell", "retire", "return", "reuse", "review", "reviewer", "revision",
	"rewind", "rewrite", "rice", "rich", "rid", "ridden", "riddle", "ride", "right", "ring", "ripe",
	"ripen", "rise", "risen", "river", "road", "roast", "rob", "robot", "rock", "rocket",
	"rode", "role", "roll", "roller", "rooster", "root", "rope", "rose", "rot",
	"rough", "round", "roundabout", "row", "rubber", "rubbish", "rude", "rugby", "ruin", "rule",
	"ruler", "run", "rung", "runner", "running", "rush", "russia", "russian", "sad", "sadness",
	"safe", "safety", "sailing", "sailor", "salad", "salary", "sale", "salesgirl", "salesman", "saleswoman",
	"salt", "salty", "salute", "same", "sand", "sandwich", "sang", "sank", "sat", "satisfaction",
	"satisfy", "saturday", "sauce", "saucer", "sausage", "savage", "save", "saw", "scholar",
}

// Generates a 8bits checksum for input
func checksum8(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksum8bitsLen]
}

// Verify the checksum is correct
func verifyChecksum(code []byte) bool {
	checksum := checksum8(code[addressChecksum8bitsLen:])
	if len(checksum) != addressChecksum8bitsLen {
		return false
	}

	for k, v := range checksum {
		if v != code[k] {
			return false
		}
	}

	return true
}

//Encode code to wordValue
func encode2WordValue(code []byte) [wordsCount]int64 {
	var wordValue [wordsCount]int64

	for i := 0; i < wordsCount; i++ {
		valSrc := big.NewInt(0).SetBytes(code)
		valTemp := valSrc.Rsh(valSrc, wordPerBit * uint(i))
		valAnd:= big.NewInt(dictionaryWordsCnt)
		val := valTemp.And(valTemp, valAnd)
		wordValue[i] = val.Int64()
	}

	return wordValue
}

//Decode wordValue to code
func decode2Code(wordValue [wordsCount]int64) []byte {
	var code []byte

	val := big.NewInt(0)
	for i := 0; i < wordsCount; i++ {
		valAdd := big.NewInt(wordValue[wordsCount-i-1])
		val = val.Lsh(val, wordPerBit)
		val = val.Add(val, valAdd)
	}

	code = val.Bytes()
	return code
}

//Encode private key to keywords
func PrivKey2Words(privKey *big.Int) []string {
	checksum := checksum8(privKey.Bytes())
	code := append(checksum, privKey.Bytes()...)

	wordValue := encode2WordValue(code)

	words:= make([]string, dictionaryWordsCnt)
	for i := 0; i < wordsCount; i++ {
		words[i] = privKeysDictionary[wordValue[i]]
	}

	return words
}

//Decode keywords to private key
func Words2PrivKey(words []string) *big.Int {
	wordsDictionary := make(map[string]int64, dictionaryWordsCnt)
	for i := int64(0); i < dictionaryWordsCnt; i++ {
		wordsDictionary[privKeysDictionary[i]] = i
	}

	///test-------------------------------------------------------------------------
	//for j := int64(0); j < dictionaryWordsCnt; j++ {
	//	fmt.Println(wordsDictionary[privKeysDictionary[j]])
	//	if wordsDictionary[privKeysDictionary[j]] != j {
	//		fmt.Println(privKeysDictionary[j])
	//		fmt.Println("Test error!")
	//		break
	//	}
	//}
	///end test---------------------------------------------------------------------

	var wordValue [wordsCount]int64
	for k := range wordValue {
		wordValue[k] = wordsDictionary[words[k]]
	}
	code := decode2Code(wordValue)

	if verifyChecksum(code) {
		key := new(big.Int).SetBytes(code[addressChecksum8bitsLen:])
		return key
	} else {
		return big.NewInt(0)
	}
}
