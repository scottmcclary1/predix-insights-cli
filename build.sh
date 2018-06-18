#!/bin/bash

# Check Operating System
uname -s | grep Darwin > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo 'Error: the Predix Insights CLI is only supported on MacOS.'
    exit 1
fi

# Install Dependencies
which go > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo 'Error: Go is not installed. Please install Go manaully and then rebuild PI CLI.'
    exit 1
fi

rm -rf vendor

if [ ! -d vendor ]; then
    glide install --strip-vendor
fi

## Build Source Code
which go > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo 'Error: Go is not installed. Please install Go manaully and then rebuild PI CLI.'
    exit 1
fi
go build -ldflags "-X github.build.ge.com/predix-data-services/predix-insights-cli/cmd.Version=`cat version` -X github.build.ge.com/predix-data-services/predix-insights-cli/cmd.GitHash=`git rev-parse HEAD`" -o pi
if [ ! $? -eq 0 ]; then
    echo 'Error: PI CLI source code could not be compiled.'
    exit 1
fi

## Install PI CLI in GOBIN
if [ -z $GOBIN ]
then
    echo 'Error: GOBIN environment variable is not defined. Please set manually and then rebuild PI CLI.'
    exit 1
fi
if [ $? -eq 0 ]
then
    #go install
    mv pi $GOBIN/
fi

## Install Bash Completions using homebrew
which brew > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo 'Error: brew is not installed. Please install brew manaully (see command below) and then rebuild PI CLI:'
    echo '$ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"'
    exit 1
fi
brew ls --versions bash-completion > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    brew install bash-completion
fi

str='if [ -f $(brew --prefix)/etc/bash_completion ]; then . $(brew --prefix)/etc/bash_completion; fi'
grep 'then . $(brew --prefix)/etc/bash_completion; fi' ~/.bash_profile > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo $str >> ~/.bash_profile
fi

## Generate BASH Completion File
if [ ! -z $ADMIN ]
then
    GENERATE_BASH_COMPLETION_FILE=true pi -h > /dev/null 2>&1
fi

## Install BASH Completion File 
bashCompletionDir=$(brew --prefix)/etc/bash_completion.d
if [ ! -d $bashCompletionDir ]; then
    mkdir -p $bashCompletionDir
fi
cp scripts/pi_completion.sh $bashCompletionDir/

source ~/.bash_profile

pi -h > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
    echo "PI CLI installation failed..."
    exit 1
fi

#echo "PI CLI installation complete!"
echo
echo "************ATTENTION************"
echo "Execute the following command to complete installation:"
echo "$ source ~/.bash_profile"
echo
exit 0
