package slack

import (
	"context"
	"encoding/json"
)

// RawViewObject view payload
type RawViewObject struct {
	Type            ViewType      `json:"type"`
	Title           interface{}   `json:"title"`
	Blocks          []interface{} `json:"blocks"`
	Close           interface{}   `json:"close,omitempty"`
	Submit          interface{}   `json:"submit,omitempty"`
	PrivateMetadata string        `json:"private_metadata,omitempty"`
	CallbackID      string        `json:"callback_id,omitempty"`
	ClearOnClose    bool          `json:"clear_on_close,omitempty"`
	NotifyOnClose   bool          `json:"notify_on_close,omitempty"`
	ExternalID      string        `json:"external_id,omitempty"`
}

// OpenRawViewArguments views.open payload
type OpenRawViewArguments struct {
	TriggerID string        `json:"trigger_id"`
	View      RawViewObject `json:"view"`
}

// OpenRawViewContext do views.open with highly customizable
func (api *Client) OpenRawView(args OpenRawViewArguments) (resp *ViewResponse, err error) {
	return api.OpenRawViewContext(context.Background(), args)
}

// OpenRawViewContext do views.open with highly customizable and context
func (api *Client) OpenRawViewContext(ctx context.Context, args OpenRawViewArguments) (resp *ViewResponse, err error) {
	var encoded []byte
	encoded, err = json.Marshal(args)
	if err != nil {
		return resp, err
	}

	if api.Debug() {
		api.Debugf("Sending request views.open: %s", string(encoded))
	}

	resp = &ViewResponse{}

	endpoint := api.endpoint + "views.open"
	err = postJSON(ctx, api.httpclient, endpoint, api.token, encoded, &resp, api)
	if err != nil {
		return resp, err
	}

	return resp, resp.Err()
}

// UpdateRawViewArguments views.update payload
type UpdateRawViewArguments struct {
	ExternalID string        `json:"external_id,omitempty"`
	Hash       string        `json:"hash,omitempty"`
	ViewID     string        `json:"view_id,omitempty"`
	View       RawViewObject `json:"view"`
}

// UpdateRawView do views.update with highly customizable
func (api *Client) UpdateRawView(args UpdateRawViewArguments) (resp *ViewResponse, err error) {
	return api.UpdateRawViewContext(context.Background(), args)
}

// UpdateRawViewContext do views.update with highly customizable and context
func (api *Client) UpdateRawViewContext(ctx context.Context, args UpdateRawViewArguments) (resp *ViewResponse, err error) {
	if args.ExternalID == "" && args.ViewID == "" {
		return resp, ErrParametersMissing
	}

	var encoded []byte
	encoded, err = json.Marshal(args)
	if err != nil {
		return resp, err
	}

	if api.Debug() {
		api.Debugf("Sending request views.update: %s", string(encoded))
	}

	endpoint := api.endpoint + "views.update"
	err = postJSON(ctx, api.httpclient, endpoint, api.token, encoded, resp, api)
	if err != nil {
		return resp, err
	}

	return resp, resp.Err()
}
