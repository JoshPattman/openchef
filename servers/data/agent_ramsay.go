package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"utils"

	"github.com/JoshPattman/jpf"
)

type agentRamsaySearchRequest struct {
	Query string `json:"query"`
}

type agentRamsaySearchResponse struct {
	ResponseText string `json:"response_text"`
}

func agentRamsaySearch(req agentRamsaySearchRequest) (agentRamsaySearchResponse, error) {
	model := jpf.NewStandardOpenAIModel(os.Getenv("OPENAI_KEY"), "gpt-4o-mini", 0, 0, 0.5)
	agent := &agentRamsay{}
	state := agentRamsayState{
		userQuery: req.Query,
	}
	for step := range jpf.RunAgent(model, agent, state) {
		if step.Error != nil {
			return agentRamsaySearchResponse{}, step.Error
		}
		logger.Debug("agent action taken", "action_type", fmt.Sprintf("%T", step.Action), "action args", step.Action)
		if ans, ok := step.Action.(agentRamsayFinalAnswerAction); ok {
			return agentRamsaySearchResponse{
				ResponseText: ans.FinalAnswer,
			}, nil
		}
	}
	return agentRamsaySearchResponse{}, fmt.Errorf("no final answer")
}

var _ jpf.Agent[agentRamsayState] = &agentRamsay{}
var _ jpf.Action[agentRamsayState] = agentRamsayFinalAnswerAction{}
var _ jpf.Action[agentRamsayState] = agentRamsayThinkingAction{}
var _ jpf.Action[agentRamsayState] = agentRamsaySearchRecipesAction{}

type agentRamsayState struct {
	userQuery     string
	history       []jpf.Message
	finalResponse string
}
type agentRamsay struct{}

var arSys = `
- You are Agent Ramsay, a profesional Chef AI.
- The user will provide you with a query of what they want to cook, and you (the 'assistant') will help them come up with a recipe.
- Each time you respond, you should always respond with some plaintext:
	- This is a useful way to collect your thoughts before making any calls to tools or final answers.
	- You can write anything you think would be helpful to yourself in this, the user will never see it though.
- You may optionall also include AT MOST ONE tool call with one of the follwoing formats:
	- The following json format (with no extra text or backticks):
		{
			"tool": "search_recipes",
			"query": "<What to search for>"
		}
		- This can be used to search a state-of-the-art dataset of recipes for anything that is semantically similar.
		- You can make as many recipe searches as you would like to inform your final answer.
	- The following json format (with no extra text or backticks):
		{
			"tool": "final_answer",
			"final_answer": "<Your final answer that will be shown directly to the user>"
		}
		- This will terminate your process, and will be used to make your final answer to the user.
		- You should not call this without first at least thinking in a plaintext response, but you should also usually call other tools first.
		- The only text that the user can see is the text inside of the "final_answer" key.
- You may need to come up with new recipes by combining ones from the dataset with user preferences.
- In your final answer, you may tell the user which recipes you were inspired by by citing the IDs of the recipes in the format [1], [2].
`

// BuildInputMessages implements jpf.Function.
func (a *agentRamsay) BuildInputMessages(state agentRamsayState) ([]jpf.Message, error) {
	return append([]jpf.Message{
		{
			Role:    jpf.SystemRole,
			Content: arSys,
		},
		{
			Role:    jpf.UserRole,
			Content: state.userQuery,
		},
	}, state.history...), nil
}

// ParseResponseText implements jpf.Function.
func (a *agentRamsay) ParseResponseText(resp string) (jpf.Action[agentRamsayState], error) {
	jsObj, ok := extractJsonObject(resp)
	if !ok {
		return agentRamsayThinkingAction{
			RawResponse: resp,
		}, nil
	}
	respObj := struct {
		Tool        string `json:"tool"`
		Query       string `json:"query"`
		FinalAnswer string `json:"final_answer"`
	}{}
	err := json.Unmarshal([]byte(jsObj), &respObj)
	if err != nil {
		return nil, err
	}
	switch respObj.Tool {
	case "search_recipes":
		if respObj.Query == "" {
			return nil, fmt.Errorf("must specify a query to search recipes")
		}
		return agentRamsaySearchRecipesAction{
			Query:       respObj.Query,
			RawResponse: resp,
		}, nil
	case "final_answer":
		if respObj.FinalAnswer == "" {
			return nil, fmt.Errorf("must specify a final answer to answer")
		}
		return agentRamsayFinalAnswerAction{
			FinalAnswer: respObj.FinalAnswer,
		}, nil
	}
	return nil, fmt.Errorf("unrecognised type %s", respObj.Tool)
}

type agentRamsayFinalAnswerAction struct {
	FinalAnswer string
}

// DoAction implements jpf.Action.
func (a agentRamsayFinalAnswerAction) DoAction(state agentRamsayState) (agentRamsayState, bool, error) {
	state.finalResponse = a.FinalAnswer
	return state, true, nil
}

type agentRamsaySearchRecipesAction struct {
	Query       string
	RawResponse string
}

// DoAction implements jpf.Action.
func (a agentRamsaySearchRecipesAction) DoAction(state agentRamsayState) (agentRamsayState, bool, error) {
	req := utils.SemanticSearchRequest{
		Search: a.Query,
		MaxN:   5,
	}
	resp, err := semanticSearch(req)
	if err != nil {
		return agentRamsayState{}, false, err
	}
	bs, err := json.Marshal(resp)
	if err != nil {
		return agentRamsayState{}, false, err
	}
	state.history = append(state.history, jpf.Message{
		Role:    jpf.AssistantRole,
		Content: a.RawResponse,
	})
	state.history = append(state.history, jpf.Message{
		Role:    jpf.SystemRole,
		Content: string(bs),
	})
	return state, false, nil
}

type agentRamsayThinkingAction struct {
	RawResponse string
}

// DoAction implements jpf.Action.
func (a agentRamsayThinkingAction) DoAction(state agentRamsayState) (agentRamsayState, bool, error) {
	state.history = append(state.history, jpf.Message{
		Role:    jpf.AssistantRole,
		Content: a.RawResponse,
	})
	state.history = append(state.history, jpf.Message{
		Role:    jpf.SystemRole,
		Content: "Thank you, your last response was understood as plaintext thinking (no tool call). Please continue...",
	})
	return state, false, nil
}

func extractJsonObject(text string) (string, bool) {
	re := regexp.MustCompile("(?s)({.*?})")
	js := re.FindStringSubmatch(text)
	if js == nil {
		logger.Debug("extracted json object", "obj", js)
		return "", false
	}
	return js[1], true
}
