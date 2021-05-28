package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	actionName = Flag{
		Name:      "Action Name",
		LongForm:  "action",
		ShortForm: "a",
		Help:      "Name of the action.",
	}

	actionTrigger = Flag{
		Name:      "Trigger Type",
		LongForm:  "trigger",
		ShortForm: "t",
		Help:      "Trigger of the action.",
	}
)

func actionsCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Manage resources for actions",
		Long:  "Manage resources for applications.",
	}

	cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(actionsLogsCmd(cli))

	return cmd
}

func actionsLogsCmd(cli *cli) *cobra.Command {
	var inputs struct {
		ActionName string
		Trigger    string
	}

	cmd := &cobra.Command{
		Use:   "logs",
		Args:  cobra.MaximumNArgs(1),
		Short: "View and follow the actions logs",
		Long:  "View and follow the actions logs filtering by action or trigger type.",
		Example: `auth0 actions logs tail
auth0 actions logs tail --action <action-name> --trigger <trigger type>
auth0 actions logs tail -a <action-name> -t <trigger type>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			triggerOpts, err := pickerOptionsWithNone(cli.triggerPickerOptions())
			if err != nil {
				return err
			}

			// Prompt for app type
			if err := actionTrigger.Select(cmd, &inputs.Trigger, triggerOpts.labels(), nil); err != nil {
				return err
			}

			actionOpts, err := pickerOptionsWithNone(cli.actionPickerOptions())
			if err != nil {
				return err
			}

			if err := actionName.Select(cmd, &inputs.ActionName, actionOpts.labels(), nil); err != nil {
				return err
			}

			// We're covering both cases of action being supplied
			// via --action flag, or being chosen.
			var actionID string
			if inputs.ActionName != "" && inputs.ActionName != optNone {
				a, err := cli.getActionByName(cmd.Context(), inputs.ActionName)
				if err != nil {
					return err
				}
				actionID = a.ID
			}

			if inputs.Trigger == optNone {
				inputs.Trigger = ""
			}

			log.Printf("chosen trigger: %q, action: %q", inputs.Trigger, actionID)
			return nil
		},
	}

	actionName.RegisterString(cmd, &inputs.ActionName, "")
	actionTrigger.RegisterString(cmd, &inputs.Trigger, "")
	return cmd
}

func (c *cli) triggerPickerOptions() (pickerOptions, error) {
	triggers, err := c.listActionsTriggers(context.TODO())
	if err != nil {
		return nil, err
	}

	var opts pickerOptions

	for _, t := range triggers {
		opts = append(opts, pickerOption{
			value: t.ID,
			label: t.ID,
		})
	}

	return opts, nil
}

func (c *cli) actionPickerOptions() (pickerOptions, error) {
	actions, err := c.listActions(context.TODO())
	if err != nil {
		return nil, err
	}

	var opts pickerOptions

	for _, a := range actions {
		opts = append(opts, pickerOption{
			value: a.Name,
			label: a.Name,
		})
	}

	return opts, nil
}

type action struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *cli) listActions(ctx context.Context) ([]action, error) {
	const path = "/api/v2/actions/actions"

	var payload struct {
		Actions []action `json:"actions"`
	}

	if err := c.doReq(ctx, http.MethodGet, path, nil, &payload); err != nil {
		return nil, err
	}

	return payload.Actions, nil
}

func (c *cli) getActionByName(ctx context.Context, name string) (action, error) {
	// NOTE(cyx): there's no API call to get an action by name so we list.
	actions, err := c.listActions(ctx)
	if err != nil {
		return action{}, err
	}

	for _, a := range actions {
		if a.Name == name {
			return a, nil
		}
	}

	return action{}, fmt.Errorf("Unable to find action with name: %s", name)
}

type trigger struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

const statusCurrent = "CURRENT"

func (c *cli) listActionsTriggers(ctx context.Context) ([]trigger, error) {
	var payload struct {
		Triggers []trigger `json:"triggers"`
	}

	const path = "/api/v2/actions/triggers"
	if err := c.doReq(ctx, http.MethodGet, path, nil, &payload); err != nil {
		return nil, err
	}

	result := make([]trigger, 0, len(payload.Triggers))
	for _, tt := range payload.Triggers {
		if tt.Status == statusCurrent {
			result = append(result, tt)
		}
	}
	return result, nil
}

func (c *cli) doReq(ctx context.Context, method, path string, in, out interface{}) error {
	tenant, err := c.getTenant()
	if err != nil {
		return err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   tenant.Domain,
		Path:   path,
	}

	var body io.Reader
	if in != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(in); err != nil {
			return errors.Wrap(err, "json Encode failed")
		}
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return errors.Wrap(err, "NewRequest failed")
	}
	req.Header.Set("Authorization", "Bearer "+tenant.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "http request failed: %s %s", method, path)
	}

	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return errors.Wrap(err, "json Decode failed")
	}

	return nil
}

const optNone = "<none>"

func pickerOptionsWithNone(opts pickerOptions, err error) (pickerOptions, error) {
	if err != nil {
		return nil, err
	}

	result := append(pickerOptions{}, pickerOption{label: optNone, value: optNone})
	result = append(result, opts...)
	return result, nil
}
