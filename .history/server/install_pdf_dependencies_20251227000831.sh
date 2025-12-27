#!/bin/bash

# Installation script for PDF generation dependencies

echo "==================================================="
echo "PDF Generation System - Dependency Installation"
echo "==================================================="
echo ""

# Detect OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Detected macOS"
    echo "Installing wkhtmltopdf using Homebrew..."
    
    # Check if Homebrew is installed
    if ! command -v brew &> /dev/null; then
        echo "ERROR: Homebrew is not installed."
        echo "Please install Homebrew first: https://brew.sh/"
        exit 1
    fi
    
    brew install wkhtmltopdf
    
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Detected Linux"
    
    # Detect Linux distribution
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        
        if [[ "$OS" == *"Ubuntu"* ]] || [[ "$OS" == *"Debian"* ]]; then
            echo "Installing wkhtmltopdf for Ubuntu/Debian..."
            sudo apt-get update
            sudo apt-get install -y wkhtmltopdf
            
        elif [[ "$OS" == *"CentOS"* ]] || [[ "$OS" == *"Red Hat"* ]]; then
            echo "Installing wkhtmltopdf for CentOS/RHEL..."
            sudo yum install -y wkhtmltopdf
            
        else
            echo "Unsupported Linux distribution: $OS"
            echo "Please install wkhtmltopdf manually."
            exit 1
        fi
    else
        echo "Cannot detect Linux distribution"
        echo "Please install wkhtmltopdf manually."
        exit 1
    fi
    
else
    echo "Unsupported operating system: $OSTYPE"
    echo "Please install wkhtmltopdf manually from:"
    echo "https://wkhtmltopdf.org/downloads.html"
    exit 1
fi

echo ""
echo "==================================================="
echo "Verifying installation..."
echo "==================================================="

if command -v wkhtmltopdf &> /dev/null; then
    echo "✓ wkhtmltopdf installed successfully!"
    echo ""
    wkhtmltopdf --version
    echo ""
    echo "==================================================="
    echo "Installation complete!"
    echo "You can now generate PDFs from the application."
    echo "==================================================="
else
    echo "✗ Installation failed. wkhtmltopdf not found."
    echo "Please install manually from:"
    echo "https://wkhtmltopdf.org/downloads.html"
    exit 1
fi
