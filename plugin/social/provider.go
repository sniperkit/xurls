package social

var SocialRegexProvider = map[string]string{

	// RegEx for public URLs. Prefixed by `/in/`, `/profile/`, or `/pub/`. Alphanumerical chacters, hyphens and underscores.
	"linkedin": `^(http(s)?:\/\/)?(www\.)?linkedin\.com\/(in|profile|pub)\/([A-z 0-9 _ -]+)\/?$`,

	// RegEx for usernames. 39 alphanumerical characters. Hyphens are allowed, but usernames cannot start or end with one.
	"github": `^(http(s)?:\/\/)?(www\.)?github\.com\/(?!-)[A-z 0-9 -]{1,39}[^-]\/?$`,

	// RegEx for repository names.
	"github-repo": `((git|ssh|http(s)?)|(git@[\w\.]+))(:(\/\/)?)([\w\.@\:/\-~]+)(\.git)(\/)?`,

	// RegEx for usernames. 15 alphanumerical characters and underscores but no other special characters.
	"twitter": `^(http(s)?:\/\/)?(www\.)?twitter\.com\/[A-z 0-9 _]{1,15}\/?$`,

	// RegEx for usernames. 3-20 alphanumberical characters, dashes, and underscores. Prefix of `/user/`.
	"reddit": `^(http(s)?:\/\/)?(www\.)?reddit\.com\/user\/[A-z 0-9 _ -]{3,20}\/?$`,

	/*
		// TODO:

		// Facebook
		"facebook": ``,

		// Instagram
		"instagram": ``,

		// Dribble
		"dribbble": ``,

		// Spotify
		"spotify": ``,

		// Spkype
		"skype": ``,

		// Soundcloud
		"soundcloud": ``,

		// Google plus
		"google-plus": ``,

		// Telegram
		"telegram": ``,

		// Medium
		"medium": ``,

		// Youtube
		"youtube": ``,

		// StackOverflow
		"stackoverflow": ``,

		// Vimeo
		"vimeo": ``,

		// Keybase
		"keybase": ``,

		// Pinterest
		"pinterest": ``,

	*/

	//-- End
}
