package relay

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hdzattain/smart-gateway/constant"
	"github.com/hdzattain/smart-gateway/relay/channel"
	"github.com/hdzattain/smart-gateway/relay/channel/ali"
	"github.com/hdzattain/smart-gateway/relay/channel/aws"
	"github.com/hdzattain/smart-gateway/relay/channel/baidu"
	"github.com/hdzattain/smart-gateway/relay/channel/baidu_v2"
	"github.com/hdzattain/smart-gateway/relay/channel/claude"
	"github.com/hdzattain/smart-gateway/relay/channel/cloudflare"
	"github.com/hdzattain/smart-gateway/relay/channel/codex"
	"github.com/hdzattain/smart-gateway/relay/channel/cohere"
	"github.com/hdzattain/smart-gateway/relay/channel/coze"
	"github.com/hdzattain/smart-gateway/relay/channel/deepseek"
	"github.com/hdzattain/smart-gateway/relay/channel/dify"
	"github.com/hdzattain/smart-gateway/relay/channel/gemini"
	"github.com/hdzattain/smart-gateway/relay/channel/jimeng"
	"github.com/hdzattain/smart-gateway/relay/channel/jina"
	"github.com/hdzattain/smart-gateway/relay/channel/minimax"
	"github.com/hdzattain/smart-gateway/relay/channel/mistral"
	"github.com/hdzattain/smart-gateway/relay/channel/mokaai"
	"github.com/hdzattain/smart-gateway/relay/channel/moonshot"
	"github.com/hdzattain/smart-gateway/relay/channel/ollama"
	"github.com/hdzattain/smart-gateway/relay/channel/openai"
	"github.com/hdzattain/smart-gateway/relay/channel/palm"
	"github.com/hdzattain/smart-gateway/relay/channel/perplexity"
	"github.com/hdzattain/smart-gateway/relay/channel/replicate"
	"github.com/hdzattain/smart-gateway/relay/channel/siliconflow"
	"github.com/hdzattain/smart-gateway/relay/channel/submodel"
	taskali "github.com/hdzattain/smart-gateway/relay/channel/task/ali"
	taskdoubao "github.com/hdzattain/smart-gateway/relay/channel/task/doubao"
	taskGemini "github.com/hdzattain/smart-gateway/relay/channel/task/gemini"
	"github.com/hdzattain/smart-gateway/relay/channel/task/hailuo"
	taskjimeng "github.com/hdzattain/smart-gateway/relay/channel/task/jimeng"
	"github.com/hdzattain/smart-gateway/relay/channel/task/kling"
	tasksora "github.com/hdzattain/smart-gateway/relay/channel/task/sora"
	"github.com/hdzattain/smart-gateway/relay/channel/task/suno"
	taskvertex "github.com/hdzattain/smart-gateway/relay/channel/task/vertex"
	taskVidu "github.com/hdzattain/smart-gateway/relay/channel/task/vidu"
	"github.com/hdzattain/smart-gateway/relay/channel/tencent"
	"github.com/hdzattain/smart-gateway/relay/channel/vertex"
	"github.com/hdzattain/smart-gateway/relay/channel/volcengine"
	"github.com/hdzattain/smart-gateway/relay/channel/xai"
	"github.com/hdzattain/smart-gateway/relay/channel/xunfei"
	"github.com/hdzattain/smart-gateway/relay/channel/zhipu"
	"github.com/hdzattain/smart-gateway/relay/channel/zhipu_4v"
)

func GetAdaptor(apiType int) channel.Adaptor {
	switch apiType {
	case constant.APITypeAli:
		return &ali.Adaptor{}
	case constant.APITypeAnthropic:
		return &claude.Adaptor{}
	case constant.APITypeBaidu:
		return &baidu.Adaptor{}
	case constant.APITypeGemini:
		return &gemini.Adaptor{}
	case constant.APITypeOpenAI:
		return &openai.Adaptor{}
	case constant.APITypePaLM:
		return &palm.Adaptor{}
	case constant.APITypeTencent:
		return &tencent.Adaptor{}
	case constant.APITypeXunfei:
		return &xunfei.Adaptor{}
	case constant.APITypeZhipu:
		return &zhipu.Adaptor{}
	case constant.APITypeZhipuV4:
		return &zhipu_4v.Adaptor{}
	case constant.APITypeOllama:
		return &ollama.Adaptor{}
	case constant.APITypePerplexity:
		return &perplexity.Adaptor{}
	case constant.APITypeAws:
		return &aws.Adaptor{}
	case constant.APITypeCohere:
		return &cohere.Adaptor{}
	case constant.APITypeDify:
		return &dify.Adaptor{}
	case constant.APITypeJina:
		return &jina.Adaptor{}
	case constant.APITypeCloudflare:
		return &cloudflare.Adaptor{}
	case constant.APITypeSiliconFlow:
		return &siliconflow.Adaptor{}
	case constant.APITypeVertexAi:
		return &vertex.Adaptor{}
	case constant.APITypeMistral:
		return &mistral.Adaptor{}
	case constant.APITypeDeepSeek:
		return &deepseek.Adaptor{}
	case constant.APITypeMokaAI:
		return &mokaai.Adaptor{}
	case constant.APITypeVolcEngine:
		return &volcengine.Adaptor{}
	case constant.APITypeBaiduV2:
		return &baidu_v2.Adaptor{}
	case constant.APITypeOpenRouter:
		return &openai.Adaptor{}
	case constant.APITypeXinference:
		return &openai.Adaptor{}
	case constant.APITypeXai:
		return &xai.Adaptor{}
	case constant.APITypeCoze:
		return &coze.Adaptor{}
	case constant.APITypeJimeng:
		return &jimeng.Adaptor{}
	case constant.APITypeMoonshot:
		return &moonshot.Adaptor{} // Moonshot uses Claude API
	case constant.APITypeSubmodel:
		return &submodel.Adaptor{}
	case constant.APITypeMiniMax:
		return &minimax.Adaptor{}
	case constant.APITypeReplicate:
		return &replicate.Adaptor{}
	case constant.APITypeCodex:
		return &codex.Adaptor{}
	}
	return nil
}

func GetTaskPlatform(c *gin.Context) constant.TaskPlatform {
	channelType := c.GetInt("channel_type")
	if channelType > 0 {
		return constant.TaskPlatform(strconv.Itoa(channelType))
	}
	return constant.TaskPlatform(c.GetString("platform"))
}

func GetTaskAdaptor(platform constant.TaskPlatform) channel.TaskAdaptor {
	switch platform {
	//case constant.APITypeAIProxyLibrary:
	//	return &aiproxy.Adaptor{}
	case constant.TaskPlatformSuno:
		return &suno.TaskAdaptor{}
	}
	if channelType, err := strconv.ParseInt(string(platform), 10, 64); err == nil {
		switch channelType {
		case constant.ChannelTypeAli:
			return &taskali.TaskAdaptor{}
		case constant.ChannelTypeKling:
			return &kling.TaskAdaptor{}
		case constant.ChannelTypeJimeng:
			return &taskjimeng.TaskAdaptor{}
		case constant.ChannelTypeVertexAi:
			return &taskvertex.TaskAdaptor{}
		case constant.ChannelTypeVidu:
			return &taskVidu.TaskAdaptor{}
		case constant.ChannelTypeDoubaoVideo, constant.ChannelTypeVolcEngine:
			return &taskdoubao.TaskAdaptor{}
		case constant.ChannelTypeSora, constant.ChannelTypeOpenAI:
			return &tasksora.TaskAdaptor{}
		case constant.ChannelTypeGemini:
			return &taskGemini.TaskAdaptor{}
		case constant.ChannelTypeMiniMax:
			return &hailuo.TaskAdaptor{}
		}
	}
	return nil
}
