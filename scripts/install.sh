#!/bin/bash

echo "Setting up Go dependencies..."

packages=(
    "go.uber.org/zap"
    "github.com/spf13/viper"
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/json-iterator/go"
    "github.com/google/wire"
)

for pkg in "${packages[@]}"; do
    echo "Fetching $pkg..."
    go get "$pkg"
done

echo "Setup complete!"
