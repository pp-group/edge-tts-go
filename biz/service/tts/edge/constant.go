package edge

type MessageType int

const (
	TrustedClientToken = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WssURL             = "wss://speech.platform.bing.com/consumer/speech/synthesize/" + "readaloud/edge/v1?TrustedClientToken=" + TrustedClientToken

	defaultVoice = "zh-CN-XiaoxiaoNeural"
)
