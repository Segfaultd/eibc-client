package main

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type slacker struct {
	*slack.Client // TODO: abstract this out
	channelID     string
	enabled       bool
}

func (oc *orderClient) begOnSlack(ctx context.Context, coin sdk.Coin, chainID, node string) (string, error) {
	if !oc.slack.enabled {
		oc.logger.Debug("Slack is disabled")
		return "", nil
	}
	address := oc.account.GetAddress().String()
	oc.logger.With(
		zap.String("denom", coin.Denom),
		zap.String("amount", coin.Amount.String()),
		zap.String("address", address),
	).Debug("Slack post @poor-bots")

	message := fmt.Sprintf("please sir, send some crypto: address: %s, amount: %s, denom: %s, chainID: %s, node: %s",
		address, coin.Amount.String(), coin.Denom, chainID, node)

	respChannel, respTimestamp, err := oc.slack.PostMessageContext(
		ctx,
		oc.slack.channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return "", err
	}

	oc.logger.With(
		zap.String("channel", respChannel),
		zap.String("timestamp", respTimestamp),
	).Debug("Slack message successfully sent")

	time.Sleep(time.Second * 5)

	_, _, _ = oc.slack.PostMessageContext(
		ctx,
		oc.slack.channelID,
		slack.MsgOptionText("I can wait...", false),
	)
	return respTimestamp, nil
}
