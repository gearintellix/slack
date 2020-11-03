package slack

import (
	"context"
	"encoding/json"
)

// OpenRawViewArguments views.open payload
type OpenRawViewArguments struct {
	TriggerID       string        `json:"-"`
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

// OpenRawViewContext do views.open with highly customizable
func (api *Client) OpenRawView(args OpenRawViewArguments) (resp *ViewResponse, err error) {
	return api.OpenRawViewContext(context.Background(), args)
}

// OpenRawViewContext do views.open with highly customizable and context
func (api *Client) OpenRawViewContext(ctx context.Context, args OpenRawViewArguments) (resp *ViewResponse, err error) {
	var encoded []byte
	encoded, err = json.Marshal(map[string]interface{}{
		"trigger_id": args.TriggerID,
		"view":       args,
	})
	if err != nil {
		return resp, err
	}

	resp = &ViewResponse{}

	endpoint := api.endpoint + "views.open"
	err = postJSON(ctx, api.httpclient, endpoint, api.token, encoded, &resp, api)
	if err != nil {
		return resp, err
	}

	return resp, resp.Err()
}
