This example presumes:

1. You're using mogrus as logger.

2. You want to log messages from something to Mattermost.

3. You have ``SendMessage(channelID string, msg string)`` in your Mattermost's
client library for sending messages to desired channel.

4. ``cl.DebugChannelID is`` a channel ID to which mogrus will log.

Here's io.Writer-compatible struct for using with mogrus:

```
var logMessage = `Level: %s
Timestamp: %s
 
    %s`
 
// MattermostChannelLogger provides io.Writer compatible interface for
// writing logs to Mattermost channel.
type MattermostChannelLogger struct{}
 
func (mcl MattermostChannelLogger) Write(p []byte) (n int, err error) {
    // Generate eye-candy message.
    pstr := string(p)
 
    timer := regexp.MustCompile("time=\"([0-9. -:]+)\"")
    timeFound := timer.FindAllStringSubmatch(pstr, -1)
    time := timeFound[0][1]
 
    levelr := regexp.MustCompile("level=([a-zA-Z]+)")
    levelFound := levelr.FindAllStringSubmatch(pstr, -1)
    level := levelFound[0][1]
 
    msgr := regexp.MustCompile("msg=\"(.+)\"")
    msgFound := msgr.FindAllStringSubmatch(pstr, -1)
    logmsg := msgFound[0][1]
 
    msg := fmt.Sprintf(logMessage, level, time, logmsg)
 
    cl.SendMessage(cl.DebugChannelID, msg)
 
    return len(p), nil
}
```

Initialize it after connecting to Mattermost server like:

```
// Register output for logging into Mattermost's channel.
c.Log.CreateOutput("Mattermost debug channel", MattermostChannelLogger{}, false, cfg.System.Logging.Level)
```

Where ``c.Log`` is a mogrus instance, ``cfg.System.Logging.Level`` is one of logging levels ("debug", "info", "warn" or "error").