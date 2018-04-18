# slack-webcam
Turn your Slack avatar into an IP camera

## Installation

    sudo curl -L https://github.com/justmiles/slack-webcam/releases/download/v0.0.1/slack-webcam.linux-amd64 -o /usr/local/bin/slack-webcam
    sudo chmod +x /usr/local/bin/slack-webcam

## Usage

    Usage of slack-webcam:
      -debug
        	enable debug logging
      -device string
        	camera device index
      -token string
        	slack token for auth
          
[Generate a Slack token](https://api.slack.com/custom-integrations/legacy-tokens) for your personal user and pass that to slack-webcam

    slack-webcam -token "XYXYXYX"

If you have more than one camera on your machine, you can reference the device's index using the `-device` flag. This defaults to `0`.

Currently, this uploads one image and exits. It's up to you to determine how often it runs.