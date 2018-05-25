package memo

const (
	CodePrefix = 0x6d

	CodeTest              = 0x00
	CodeSetName           = 0x01
	CodePost              = 0x02
	CodeReply             = 0x03
	CodeLike              = 0x04
	CodeSetProfile        = 0x05
	CodeFollow            = 0x06
	CodeUnfollow          = 0x07
	CodeSetImageBaseUrl   = 0x08
	CodeAttachPicture     = 0x09
	CodeSetProfilePicture = 0x0A
	CodeRepost            = 0x0B
	CodeTopicMessage      = 0x0C

	CodePollCreate = 0x10
	CodePollOption = 0x13
	CodePollVote   = 0x14
)

const (
	CodePollTypeSingle = 0x01
	CodePollTypeMulti  = 0x02
	CodePollTypeRank   = 0x03
)

func GetAllCodes() [][]byte {
	return [][]byte{
		{CodePrefix, CodeTest},
		{CodePrefix, CodeSetName},
		{CodePrefix, CodePost},
		{CodePrefix, CodeReply},
		{CodePrefix, CodeLike},
		{CodePrefix, CodeSetProfile},
		{CodePrefix, CodeFollow},
		{CodePrefix, CodeUnfollow},
		{CodePrefix, CodeSetImageBaseUrl},
		{CodePrefix, CodeAttachPicture},
		{CodePrefix, CodeSetProfilePicture},
		{CodePrefix, CodeRepost},
		{CodePrefix, CodeTopicMessage},
		{CodePrefix, CodePollCreate},
		{CodePrefix, CodePollOption},
		{CodePrefix, CodePollVote},
	}
}
