#!/bin/bash

# Installation script for PDF generation dependencies

echo "==================================================="
echo "PDF Generation System - Dependency Installation"
echo "==================================================="
echo ""

# Detect OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Detected macOS"
    echo "Checking for Chrome/Chromium..."
    
    if [ -d "/Applications/Google Chrome.app" ]; then
        echo "✓ Google Chrome is already installed!"
    elif [ -d "/Applications/Chromium.app" ]; then
        echo "✓ Chromium is already installed!"
    else
        echo "Installing Google Chrome using Homebrew..."
        
        # Check if Homebrew is installed
        if ! command -v brew &> /dev/null; then
            echo "ERROR: Homebrew is not installed."
            echo "Please install Homebrew first: https://brew.sh/"
            exit 1
        fi
        
        brew install --cask google-chrome
    fi
    
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Detected Linux"
    
    # Detect Linux distribution
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        
        if [[ "$OS" == *"Ubuntu"* ]] || [[ "$OS" == *"Debian"* ]]; then
            echo "Installing Chromium for Ubuntu/Debian..."
            sudo apt-get update
            sudo apt-get install -y chromium-browser
            
        elif [[ "$OS" == *"CentOS"* ]] || [[ "$OS" == *"Red Hat"* ]]; then
            echo "Installing Chromium for CentOS/RHEL..."
            sudo yum install -y chromium
            
        else
            echo "Unsupported Linux distribution: $OS"
            echo "Please install Chrome or Chromium manually."
            exit 1
        fi
    else
        echo "Cannot detect Linux distribution"
        echo "Please install Chrome or Chromium manually."
        exit 1
    fi
    
else
    echo "Unsupported operating system: $OSTYPE"
    echo "Please install Chrome or Chromium manually from:"
    echo "https://www.google.com/chrome/"
    exit 1
fi

echo ""
echo "==================================================="
echo "Installing Go dependencies..."
echo "==================================================="

cd "$(dirname "$0")"
go mod download

echo ""
echo "==================================================="
echo "Verifying installation..."
echo "==================================================="

if command -v google-chrome &> /dev/null || command -v chromium &> /dev/null || command -v chromium-browser &> /dev/null; then
    echo "✓ Chrome/Chromium installed successfully!"
    echo ""
    echo "==================================================="
    echo "Installation complete!"
    echo "You can now generate PDFs from the application."
    echo "==================================================="
else
    echo "✗ Installation verification failed."
    echo "Please install Chrome or Chromium manually from:"
    echo "https://www.google.com/chrome/"
    exit 1
fi
