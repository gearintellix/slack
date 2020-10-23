package slack

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

// RawSendMessageArguments chat.postMessage payload
type RawSendMessageArguments struct {
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

// RawSendMessage do chat.postMessage with highly customizable
func (api *Client) RawSendMessage(args RawSendMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	return api.RawSendMessageContext(context.Background(), args)
}

// RawSendMessageContext do chat.postMessage with highly customizable and context
func (api *Client) RawSendMessageContext(ctx context.Context, args RawSendMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	var (
		req      *http.Request
		response chatResponseFull
	)

	args.Token = api.token
	req, err = jsonReq(api.endpoint+string(chatPostMessage), args)
	if err != nil {
		return _channel, _timestamp, _text, err
	}

	if api.Debug() {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return _channel, _timestamp, _text, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		api.Debugf("Sending request: %s", string(reqBody))
	}

	if err = doPost(ctx, api.httpclient, req, newJSONParser(&response), api); err != nil {
		return _channel, _timestamp, _text, err
	}

	_channel = response.Channel
	_timestamp = response.getMessageTimestamp()
	_text = response.Text
	err = response.Err()

	return _channel, _timestamp, _text, err
}

// RawResponseMessageArguments hook message payload
type RawResponseMessageArguments struct {
	Endpoint        string        `json:"-"`
	ThreadTS        string        `json:"thread_ts,omitempty"`
	Text            string        `json:"text,omitempty"`
	Attachments     []interface{} `json:"attachments,omitempty"`
	Blocks          []interface{} `json:"blocks,omitempty"`
	ResponseType    string        `json:"response_type,omitempty"`
	ReplaceOriginal bool          `json:"replace_original"`
	DeleteOriginal  bool          `json:"delete_original"`
}

// RawResponseMessage do hook message with highly customizable
func (api *Client) RawResponseMessage(args RawResponseMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	return api.RawResponseMessageContext(context.Background(), args)
}

// RawResponseMessageContext do hook message with highly customizable and context
func (api *Client) RawResponseMessageContext(ctx context.Context, args RawResponseMessageArguments) (_channel string, _timestamp string, _text string, err error) {
	var (
		req      *http.Request
		response chatResponseFull
	)

	req, err = jsonReq(args.Endpoint, args)
	if err != nil {
		return _channel, _timestamp, _text, err
	}

	if api.Debug() {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return _channel, _timestamp, _text, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		api.Debugf("Sending request: %s", string(reqBody))
	}

	if err = doPost(ctx, api.httpclient, req, newContentTypeParser(&response), api); err != nil {
		return _channel, _timestamp, _text, err
	}

	_channel = response.Channel
	_timestamp = response.getMessageTimestamp()
	_text = response.Text
	err = response.Err()

	return _channel, _timestamp, _text, err
}
