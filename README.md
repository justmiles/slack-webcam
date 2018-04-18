# slack-webcam
Turn your Slack avatar into an IP camera

## Installation

This package requires OpenCV version 3.4 be installed on your system, along with GoCV, which is the Go programming language wrapper used by Gobot. The best way is to follow the installation instructions on the GoCV website at https://gocv.io.

- ### macOS
  To install OpenCV on macOS follow the instructions here:

  https://gocv.io/getting-started/macos/

- ### Ubuntu
  To install on Ubuntu follow the instructions here:

  https://gocv.io/getting-started/linux/

- ### Windows
  To install on Windows follow the instructions here:

  https://gocv.io/getting-started/windows/


Now you can install the slack-webcam binary

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