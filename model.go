package tweetgo

// OAuthRequestTokenInput contains the possible inputs to the request token endpoint.
type OAuthRequestTokenInput struct {
	OAuthCallback   *string `schema:"oauth_callback"`
	XAuthAccessType *string `schema:"x_auth_access_type"`
}

// OAuthRequestTokenOutput contains the results of calling request token.
type OAuthRequestTokenOutput struct {
	OAuthToken             string `schema:"oauth_token,required"`
	OAuthTokenSecret       string `schema:"oauth_token_secret,required"`
	OAuthCallbackConfirmed bool   `schema:"oauth_callback_confirmed,required"`
}

// OAuthAccessTokenInput contains the input necessary to exchange a request token for an access token
type OAuthAccessTokenInput struct {
	OAuthToken    *string `schema:"oauth_token"`
	OAuthVerifier *string `schema:"oauth_verifier"`
}

// OAuthAccessTokenOutput contains the results of calling OAuth Access Token, including a long lived token
type OAuthAccessTokenOutput struct {
	OAuthToken       string `schema:"oauth_token"`
	OAuthTokenSecret string `schema:"oauth_token_secret"`
	UserID           int64  `schema:"user_id"`
	ScreenName       string `schema:"screen_name"`
}

// ListsListInput contains the possible inputs when listing a list
type ListsListInput struct {
	UserID     *int64  `schema:"user_id"`
	ScreenName *string `schema:"screen_name"`
	Reverse    *bool   `schema:"reverse"`
}

// ListsListOutput contains the output of listing the lists
type ListsListOutput struct {
	ID              int64  `json:"id"`
	IDStr           string `json:"id_str"`
	Name            string `json:"name"`
	URI             string `json:"uri"`
	SubscriberCount int    `json:"subscriber_count"`
	MemberCount     int    `json:"member_count"`
	Mode            string `json:"mode"`
	Description     string `json:"description"`
	Slug            string `json:"slug"`
	FullName        string `json:"full_name"`
	CreatedAt       string `json:"created_at"`
	Following       bool   `json:"following"`
	User            user   `json:"user"`
}

// ListsMembersInput contains the possible inputs when listing the members of a list
type ListsMembersInput struct {
	ListID          *int64  `schema:"list_id"`
	Slug            *string `schema:"slug"`
	OwnerScreenName *string `schema:"owner_screen_name"`
	OwnerID         *int64  `schema:"owner_id"`
	Count           *int    `schema:"count"`
	Cursor          *int    `schema:"cursor"`
	IncludeEntities *bool   `schema:"include_entities"`
	SkipStatus      *bool   `schema:"skip_status"`
}

// ListsMembersOutput contians the output from listing the members of a list
type ListsMembersOutput struct {
	Users             []user `json:"users"`
	NextCursor        int    `json:"next_cursor"`
	NextCursorStr     string `json:"next_cursor_str"`
	PreviousCursor    int    `json:"previous_cursor"`
	PreviousCursorStr string `json:"previous_cursor_str"`
	TotalCount        int    `json:"total_count"`
}

// ListsMembersShowInput contains the possible inputs for the lists/members/show endpoint
type ListsMembersShowInput struct {
	ListID          *int64  `schema:"list_id"`
	Slug            *string `schema:"slug"`
	UserID          *int64  `schema:"user_id"`
	ScreenName      *string `schema:"screen_name"`
	OwnerScreenName *string `schema:"owner_screen_name"`
	OwnerID         *int64  `schema:"owner_id"`
	IncludeEntities *bool   `schema:"include_entities"`
	SkipStatus      *bool   `schema:"skip_status"`
}

// ListsMembersShowOutput contains the output for the lists/members/show endpoint
type ListsMembersShowOutput struct {
	user
}

// StatusesUpdateInput contains the possible inputs when updating a status
type StatusesUpdateInput struct {
	Status                    *string  `schema:"status"`
	InReplyToStatusID         *int64   `schema:"in_reply_to_status_id"`
	AutoPopulateReplyMetadata *bool    `schema:"auto_populate_reply_metadata"`
	ExcludeReplyUserIDs       *string  `schema:"exclude_reply_user_ids"`
	AttachmentURL             *string  `schema:"attachment_url"`
	MediaIDs                  *string  `schema:"media_ids"`
	PossiblySensitive         *bool    `schema:"possibly_sensitive"`
	Lat                       *float64 `schema:"lat"`
	Long                      *float64 `schema:"long"`
	PlaceID                   *string  `schema:"place_id"`
	DisplayCoordinates        *bool    `schema:"display_coordinates"`
	TrimUser                  *bool    `schema:"trim_user"`
	EnableDMCommands          *bool    `schema:"enable_dmcommands"`
	FailDMCommands            *bool    `schema:"fail_dmcommands"`
	CardURI                   *string  `schema:"card_uri"`
}

// StatusesUpdateOutput contains the output from posting a status update
type StatusesUpdateOutput struct {
	CreatedAt            string `json:"created_at"`
	ID                   int64  `json:"id"`
	IDStr                string `json:"id_str"`
	Text                 string `json:"text"`
	Source               string `json:"source"`
	Truncated            bool   `json:"truncated"`
	InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64  `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  string `json:"in_reply_to_screen_name"`
	User                 user   `json:"user"`
}

// StatusesFilterInput contains the input options for getting filtered statuses
type StatusesFilterInput struct {
	Follow        *string `schema:"follow"`
	Track         *string `schema:"track"`
	Locations     *string `schema:"locations"`
	Delimited     *string `schema:"delimited"`
	StallWarnings *string `schema:"stall_warnings"`
}

// StatusesFilterOutput contains the output for a single response from the filtered statuses endpoint
type StatusesFilterOutput struct {
	tweet
	QuotedStatus    tweet `json:"quoted_status"`
	RetweetedStatus tweet `json:"retweeted_status"`
	ExtendedTweet   tweet `json:"extended_tweet"`
}

// StatusesUserTimelineInput contains the input options for getting the users timeline statuses
type StatusesUserTimelineInput struct {
	UserID         *int64  `schema:"user_id"`
	ScreenName     *string `schema:"screen_name"`
	SinceID        *int64  `schema:"since_id"`
	Count          *int    `schema:"count"`
	MaxID          *int64  `schema:"max_id"`
	TrimUser       *bool   `schema:"trim_user"`
	ExcludeReplies *bool   `schema:"exclude_replies"`
	IncludeRts     *bool   `schema:"include_rts"`
}

// StatusesUserTimelineOutput contains the output for a response from the users timeline endpoint
type StatusesUserTimelineOutput struct {
	tweet
}

// ---------------------------------------------------------------------------------------------------------------------
// STRUCTS FOR COMPOSITION
// ---------------------------------------------------------------------------------------------------------------------

// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/tweet-object
type tweet struct {
	CreatedAt            string           `json:"created_at"`
	ID                   int64            `json:"id"`
	IDStr                string           `json:"id_str"`
	Text                 string           `json:"text"`
	Source               string           `json:"source"`
	Truncated            bool             `json:"truncated"`
	InReplyToStatusID    int64            `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string           `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64            `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string           `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  string           `json:"in_reply_to_screen_name"`
	User                 user             `json:"user"`
	Coordinates          coordinates      `json:"coordinates"`
	Place                place            `json:"place"`
	QuotedStatusID       int64            `json:"quoted_status_id"`
	QuotedStatusIDStr    string           `json:"quoted_status_id_str"`
	IsQuoteStatus        bool             `json:"is_quote_status"`
	QuoteCount           int              `json:"quote_count"`
	ReplyCount           int              `json:"reply_count"`
	RetweetCount         int              `json:"retweet_count"`
	FavoriteCount        int              `json:"favorite_count"`
	Entities             entities         `json:"entities"`
	ExtendedEntities     extendedEntities `json:"extended_entities"`
	Favorited            bool             `json:"favorited"`
	Retweeted            bool             `json:"retweeted"`
	PossiblySensitive    bool             `json:"possibly_sensitive"`
	FilterLevel          string           `json:"filter_level"`
	Lang                 string           `json:"lang"`
	MatchingRules        []rule           `json:"matching_rules"`
	// I did not include any "Additional Tweet attributes"
}

// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/user-object
type user struct {
	ID                   int64    `json:"id"`
	IDStr                string   `json:"id_str"`
	Name                 string   `json:"name"`
	ScreenName           string   `json:"screen_name"`
	Location             string   `json:"location"`
	Derived              derived  `json:"derived"`
	URL                  string   `json:"url"`
	Description          string   `json:"description"`
	Protected            bool     `json:"protected"`
	Verified             bool     `json:"verified"`
	FollowersCount       int      `json:"followers_count"`
	FriendsCount         int      `json:"friends_count"`
	ListedCount          int      `json:"listed_count"`
	FavouritesCount      int      `json:"favourites_count"`
	StatusesCount        int      `json:"statuses_count"`
	CreatedAt            string   `json:"created_at"`
	ProfileBannerURL     string   `json:"profile_banner_url"`
	ProfileImageURLHTTPS string   `json:"profile_image_url_https"`
	DefaultProfile       bool     `json:"default_profile"`
	DefaultProfileImage  bool     `json:"default_profile_image"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
	WithheldScope        string   `json:"withheld_scope"`
}

// https://developer.twitter.com/en/docs/tweets/enrichments/overview/profile-geo
type derived struct {
	Locations []location `json:"locations"`
}

type location struct {
	Country     string      `json:"country"`
	CountryCode string      `json:"country_code"`
	Locality    string      `json:"locality"`
	Region      string      `json:"region"`
	SubRegion   string      `json:"sub_region"`
	FullName    string      `json:"full_name"`
	Geo         coordinates `json:"geo"`
}

type coordinates struct {
	// Coordinates will be listed as [long, lat]
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type place struct {
	ID          string      `json:"id"`
	URL         string      `json:"url"`
	PlaceType   string      `json:"place_type"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	CountryCode string      `json:"country_code"`
	BoundingBox boundingBox `json:"bounding_box"`
}

type boundingBox struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

type entities struct {
	Hashtags     []hashtag      `json:"hashtags"`
	Media        []media        `json:"media"`
	URLs         []entityURL    `json:"urls"`
	UserMentions []userMentions `json:"user_mentions"`
	Symbols      []symbol       `json:"symbols"`
	Polls        []poll         `json:"polls"`
}

// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/extended-entities-object
type extendedEntities struct {
	Media []media `json:"media"`
}

type hashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type media struct {
	DisplayURL          string              `json:"display_url"`
	ExpandedURL         string              `json:"expanded_url"`
	ID                  int64               `json:"id"`
	IDStr               string              `json:"id_str"`
	Indices             []int               `json:"indices"`
	MediaURL            string              `json:"media_url"`
	MediaURLHTTPS       string              `json:"media_url_https"`
	Sizes               sizes               `json:"sizes"`
	SourceStatusID      int64               `json:"source_status_id"`
	SourceStatusIDStr   string              `json:"source_status_id_str"`
	Type                string              `json:"photo"`
	URL                 string              `json:"url"`
	VideoInfo           videoInfo           `json:"video_info"`
	AdditionalMediaInfo additionalMediaInfo `json:"additional_media_info"`
}

type sizes struct {
	Thumb  size `json:"thumb"`
	Large  size `json:"large"`
	Medium size `json:"medium"`
	Small  size `json:"small"`
}

type size struct {
	W      int    `json:"w"`
	H      int    `json:"h"`
	Resize string `json:"resize"`
}

type entityURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type userMentions struct {
	ID         int64  `json:"id"`
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type symbol struct {
	Indices []int  `json:"indicies"`
	Text    string `json:"text"`
}

type poll struct {
	Options         []option `json:"options"`
	EndDatetime     string   `json:"end_datetime"`
	DurationMinutes int      `json:"duration_minutes"`
}

type option struct {
	Position int    `json:"position"`
	Text     string `json:"text"`
}

type videoInfo struct {
	AspectRatio    []int          `json:"aspect_ratio"`
	DurationMillis int            `json:"duration_millis"`
	Variants       []videoVariant `json:"variants"`
}

type videoVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type additionalMediaInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Embeddable  bool   `json:"embeddable"`
	Monetizable bool   `json:"monetizable"`
}

type rule struct {
	Tag   string `json:"tag"`
	ID    int64  `json:"id"`
	IDStr string `json:"id_str"`
}

/*
scheduled tweet endpoints
*/

//AccountsScheduledTweetsInput ...
type AccountsScheduledTweetsInput struct {
	AccountId   *string `schema:"account_id"`
	ScheduledAt *string `schema:"scheduled_at"`
	AsUserId    *int64  `schema:"as_user_id"`
	Text        *string `schema:"text"`
	CardUri     *string `schema:"card_uri"`
	MediaKeys   *string `schema:"media_keys"`
	Nullcast    *bool   `schema:"nullcast"`
}

//AccountsScheduledTweetsOutput ...
type AccountsScheduledTweetsOutput struct {
}
