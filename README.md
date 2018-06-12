# slackstamp
Like LINE, Replace Emoji to Stamp Image with SpreadSheet.

This idea is an implementation of .

## usage

Set conf.toml or environment variables. If you use the variables, set raw json into ACCESS_TOKEN.

```
[Slack]
token = <slack regacy token>

[Server]
port = <port number>

[Spreadsheet]
token = <path to accesstoken.json>
id = <spreadsheet id>
name = <sheet name in the sheet>

```

Set OutgoingWebhooks in Slack. Trigger word are  `:` and `rmstamp`.

Set google api authentification.

Set key emoji and the URL in your spreadsheet.

## LICENSE

See LICENSE. 
