package slack

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

// RawPostMessageArguments chat.postMessage payload
type RawPostMessageArguments struct {
	Endpoint       string        `json:"-"`
	Token          string        `json:"token,omitempty"`
	Channel        string        `json:"channel,omitempty"`
	Text           string        `json:"text,omitempty"`
	AsUser         bool          `json:"as_user,omitempty"`
	Attachments    []interface{} `json:"attachments,omitempty"`
	Blocks         []interface{} `json:"blocks,omitempty"`
	IconEmoji      string        `json:"icon_emoji,omitempty"`
	IconURL        string        `json:"icon_url,omitempty"`
	LinkNames      bool          `json:"link_names,omitempty"`
	MarkDown       bool          `json:"mrkdwn,omitempty"`
	Parse          string        `json:"parse,omitempty"`
	ReplyBroadcast bool          `json:"reply_broadcast,omitempty"`
	ThreadTS       string        `json:"thread_ts,omitempty"`
	UnfurlLinks    bool          `json:"unfurl_links,omitempty"`
	UnfurlMedia    bool          `json:"unfurl_media,omitempty"`
	UserName       string        `json:"username,omitempty"`
}

// RawPostMessage do chat.postMessage with highly customizable
func (api *Client) RawPostMessage(args RawPostMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	return api.RawPostMessageContext(context.Background(), args)
}

// RawPostMessageContext do chat.postMessage with highly customizable with context
func (api *Client) RawPostMessageContext(ctx context.Context, args RawPostMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	var (
		req      *http.Request
		parser   func(*chatResponseFull) responseParser
		response chatResponseFull

		endpoint = api.endpoint + string(chatPostMessage)
	)

	switch {
	case args.Endpoint != "":
		args.Token = ""
		req, err = jsonReq(args.Endpoint, args)
		if err != nil {
			return _channel, _timestamp, _text, err
		}

		parser = func(resp *chatResponseFull) responseParser {
			return newContentTypeParser(resp)
		}

	default:
		args.Token = api.token
		req, err = jsonReq(endpoint, args)
		if err != nil {
			return _channel, _timestamp, _text, err
		}

		parser = func(resp *chatResponseFull) responseParser {
			return newJSONParser(resp)
		}
	}

	if api.Debug() {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return _channel, _timestamp, _text, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		api.Debugf("Sending request: %s", string(reqBody))
	}

	if err = doPost(ctx, api.httpclient, req, parser(&response), api); err != nil {
		return _channel, _timestamp, _text, err
	}

	_channel = response.Channel
	_timestamp = response.getMessageTimestamp()
	_text = response.Text
	err = response.Err()

	return _channel, _timestamp, _text, err
}
